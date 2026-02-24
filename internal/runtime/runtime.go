package runtime

import (
	"fmt"
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"log"
	"net/http"
)

type Runtime struct {
	LuaVm  vm.LVm
	Engine http.Handler
	Router framework.RouterDriver
	Port   func() int32
}

func New() (*Runtime, error) {
	luaVm := vm.MakeLuaVm()
	router, routerDriver := router.New()
	engine := engine.MakeEngine(router, luaVm)

	app, err := framework.RegisterFramework(luaVm, routerDriver)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		LuaVm:  luaVm,
		Engine: engine,
		Router: routerDriver,
		Port: func() int32 {
			return app.Port
		},
	}, nil
}

func (r *Runtime) LoadFile(path string) error {
	return r.LuaVm.LoadMainLuaFile(path)
}

func (r *Runtime) LoadString(code string) error {
	return r.LuaVm.LoadLuaString(code)
}

func (r *Runtime) listen() error {
	addr := fmt.Sprintf(":%d", r.Port())
	log.Printf("HTTP server starting at http://%s...\n", addr)

	err := http.ListenAndServe(addr, r.Engine)
	if err != nil {
		log.Println("failed to start server:", err)
		return err
	}

	return nil
}

func (r *Runtime) Start() error {
	if r.Port() == 0 {
		return fmt.Errorf("invalid application port, please use 'app.listen()'")
	}

	return r.listen()
}

func (r *Runtime) Close() {
	if r.LuaVm != nil {
		r.LuaVm.Close()
	}
}
