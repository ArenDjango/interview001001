package main

import (
	"context"
	"online-registration/app"
	"online-registration/cmd/migrations"
	"online-registration/internal/interview/domain/handler"
	"online-registration/internal/interview/domain/usecase"
	repository2 "online-registration/internal/interview/infrastructure/db/repository"

	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

type appServicesAndDependencies struct {
	ctx              context.Context
	cancel           context.CancelFunc
	app              *app.App
	gracefulShutdown func()
}

func main() {
	app := &cli.App{
		Name: "service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			httpCommand,
			newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

var httpCommand = &cli.Command{
	Name:  "http",
	Usage: "start http proxy server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: ":8086",
			Usage: "serve address",
		},
	},
	Action: func(c *cli.Context) error {
		servicesAndDependencies, err := startAppAndServices(
			c.Context,
			c.String("env"),
		)
		if err != nil {
			return err
		}
		defer servicesAndDependencies.gracefulShutdown()

		repository := repository2.NewDBEventRepository(servicesAndDependencies.app.DB())
		usecase := usecase.NewCreateEventUseCase(
			repository,
		)

		proxyHandlerInstance := handler.NewHandler(
			usecase,
		)

		router := gin.Default()
		v1 := router.Group("/api/v1")
		{
			v1.POST("/events", proxyHandlerInstance.CreateEvent)
		}

		srv := &http.Server{
			Addr:    c.String("addr"),
			Handler: router,
		}

		go func() {
			log.Info().
				Str("addr", c.String("addr")).
				Str("external_api", c.String("external-api-url")).
				Msg("Starting HTTP proxy server...")

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msg(fmt.Sprintf("HTTP server error: %v", err))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info().Msg("Shutting down HTTP server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("HTTP server shutdown error")
		}

		servicesAndDependencies.cancel()
		return nil
	},
}

func startAppAndServices(ctx context.Context, env string) (
	*appServicesAndDependencies,
	error,
) {
	ctx, cancel := context.WithCancel(ctx)

	_, appInstance, err := app.Start(ctx, "service", env)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start app: %w", err)
	}

	gracefulShutdown := func() {
		cancel()
		appInstance.Stop()
	}

	return &appServicesAndDependencies{
		ctx,
		cancel,
		appInstance,
		gracefulShutdown,
	}, nil
}

//nolint:funlen
func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}
					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}
					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}
					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}
					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(ctx, name)
					if err != nil {
						return err
					}
					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())
					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					migrator := migrate.NewMigrator(app.DB(), migrations)
					group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
					if err != nil {
						return err
					}
					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}
					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}
