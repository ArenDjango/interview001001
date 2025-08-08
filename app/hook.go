package app

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var onStart appHooks

func OnStart(name string, fn HookFunc) {
	onStart.Add(newHook(name, fn))
}

//------------------------------------------------------------------------------

type HookFunc func(ctx context.Context, app *App) error

type appHooks struct {
	mu    sync.Mutex
	hooks []appHook
}

func (hs *appHooks) Add(hook appHook) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.hooks = append(hs.hooks, hook)
}

func (hs *appHooks) Run(ctx context.Context, app *App) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, h := range hs.hooks {
		h := h //nolint:copyloopvar
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.run(ctx, app); err != nil {
				mu.Lock()
				defer mu.Unlock()
				fmt.Printf("hook=%q failed: %s\n", h.name, err)
				errors = append(errors, err)
			}
		}()
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errors)
	}

	return nil
}

type appHook struct {
	name string
	fn   HookFunc
}

func newHook(name string, fn HookFunc) appHook {
	return appHook{
		name: name,
		fn:   fn,
	}
}

func (h appHook) run(ctx context.Context, app *App) error {
	const timeout = 30 * time.Second

	done := make(chan struct{})
	errc := make(chan error)

	go func() {
		start := time.Now()
		if err := h.fn(ctx, app); err != nil {
			errc <- err
			return
		}
		if d := time.Since(start); d > time.Second {
			fmt.Printf("hook=%q took %s\n", h.name, d)
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case err := <-errc:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("hook=%q timed out after %s", h.name, timeout)
	}
}
