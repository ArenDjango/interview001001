//nolint:errcheck
package app

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	bundebug "github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/urfave/cli/v2"
	"gitlab.b2broker.tech/pbsr/pbsr/backend/go/libs/nats"
)

type appCtxKey struct{}

func AppFromContext(ctx context.Context) *App {
	return ctx.Value(appCtxKey{}).(*App)
}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

type App struct {
	ctx context.Context
	cfg *Config

	stopping uint32
	stopCh   chan struct{}

	onStop      appHooks
	onAfterStop appHooks

	// lazy init
	dbOnce sync.Once
	db     *bun.DB

	natsOnce sync.Once
	nats     *nats.Connection
}

func New(ctx context.Context, cfg *Config) *App {
	app := &App{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
	app.ctx = ContextWithApp(ctx, app)
	return app
}

func StartCLI(c *cli.Context) (context.Context, *App, error) {
	return Start(c.Context, c.Command.Name, c.String("env"))
}

func Start(ctx context.Context, service, envName string) (context.Context, *App, error) {
	cfg := LoadConfig(ctx)

	return StartConfig(ctx, cfg)
}

func StartConfig(ctx context.Context, cfg *Config) (context.Context, *App, error) {
	app := New(ctx, cfg)
	if err := onStart.Run(ctx, app); err != nil {
		return nil, nil, err
	}
	return app.ctx, app, nil
}

func (app *App) Stop() {
	_ = app.onStop.Run(app.ctx, app)
	_ = app.onAfterStop.Run(app.ctx, app)
}

func (app *App) OnStop(name string, fn HookFunc) {
	app.onStop.Add(newHook(name, fn))
}

func (app *App) OnAfterStop(name string, fn HookFunc) {
	app.onAfterStop.Add(newHook(name, fn))
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *Config {
	return app.cfg
}

func (app *App) Running() bool {
	return !app.Stopping()
}

func (app *App) Stopping() bool {
	return atomic.LoadUint32(&app.stopping) == 1
}

func (app *App) IsDebug() bool {
	return app.cfg.Debug
}

func (app *App) GetNatsConnection() *nats.Connection {
	log.Info().Msg("Getting nats connection")
	app.natsOnce.Do(func() {
		log.Info().Msgf("creating connection... host: %s, port: %d, JwtFilePath: %s",
			app.cfg.Nats.Host, app.cfg.Nats.Port, app.cfg.Nats.JwtCredentialFilePath)
		natsConfig := &nats.Config{
			Host:        app.cfg.Nats.Host,
			Port:        app.cfg.Nats.Port,
			JwtFilePath: app.cfg.Nats.JwtCredentialFilePath,
		}

		natsConnection, err := nats.NewNatsConnection(natsConfig)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to init NATS connection, original error: %s", err)
		}

		app.nats = natsConnection
	})

	return app.nats
}

func (app *App) DB() *bun.DB {
	app.dbOnce.Do(func() {
		postgresURL := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(app.cfg.DB.User, app.cfg.DB.Password),
			Host:   fmt.Sprintf("%s:%s", app.cfg.DB.Host, app.cfg.DB.Port),
			Path:   app.cfg.DB.Database,
		}
		connString := postgresURL.String()

		config, err := pgx.ParseConfig(connString)
		if err != nil {
			panic(err)
		}

		sqldb := stdlib.OpenDB(*config)

		err = sqldb.Ping()
		if err != nil {
			panic(err)
		}

		db := bun.NewDB(sqldb, pgdialect.New())
		db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(app.cfg.DB.Database)))

		if app.cfg.Debug {
			db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}

		app.OnStop("db.Close", func(ctx context.Context, _ *App) error {
			return db.Close()
		})

		app.db = db
	})

	return app.db
}

//------------------------------------------------------------------------------

func (app *App) WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3) //nolint:mnd
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}
