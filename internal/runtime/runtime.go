package runtime

import (
	"context"
	"fmt"
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"log"
	"net/http"
	"sync"
	"time"
)

type Runtime struct {
	LuaVm     vm.LVm
	Engine    http.Handler
	Router    framework.RouterDriver
	Port      func() int32
	RestartCh chan struct{}
	ErrorCh   chan error
}

func New() (*Runtime, error) {
	luaVm := vm.MakeLuaVm()
	router, routerDriver := router.New()
	engine := engine.MakeEngine(router, luaVm)
	restartChannel := make(chan struct{})
	errorCh := make(chan error, 1)

	app, err := framework.RegisterFramework(luaVm, routerDriver, restartChannel)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		LuaVm:     luaVm,
		Engine:    engine,
		Router:    routerDriver,
		RestartCh: restartChannel,
		ErrorCh:   errorCh,
		Port: func() int32 {
			return app.Port
		},
	}, nil
}

func (r *Runtime) DoFile(path string) error {
	return r.LuaVm.DoLuaFile(path)
}

func (r *Runtime) DoString(code string) error {
	return r.LuaVm.DoLuaString(code)
}

func (r *Runtime) listen(server *http.Server, ctx context.Context) error {
	log.Printf("HTTP server starting at http://%s...\n", server.Addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println("failed to start server:", err)
		return err
	}

	return nil

}
func (r *Runtime) Start(ctx context.Context) error {
	if r.Port() == 0 {
		return fmt.Errorf("invalid application port, please use 'app.listen()'")
	}

	for {
		var wg sync.WaitGroup
		serverCtx, severCancel := context.WithCancel(ctx)
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", r.Port()),
			Handler: r.Engine,
		}

		wg.Go(func() {
			err := r.listen(server, serverCtx)
			if err != nil {
				r.ErrorCh <- err
				return
			}
		})

		select {
		case <-r.RestartCh:
			severCancel()
			wg.Wait()
			continue
		case err := <-r.ErrorCh:
			severCancel()
			wg.Wait()
			return err
		case <-ctx.Done():
			severCancel()
			wg.Wait()
			return nil
		}

	}

}

func (r *Runtime) Close() {
	if r.LuaVm != nil {
		r.LuaVm.Close()
	}
}
