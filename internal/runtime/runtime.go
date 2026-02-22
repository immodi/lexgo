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
	Port   int32
}

func New(luaFile string) (*Runtime, error) {
	luaVm := vm.MakeLuaVm()
	router, routerDriver := router.MakeRouter()
	engine := engine.MakeEngine(router, luaVm)
	app, err := framework.RegisterFramework(luaVm, routerDriver)
	if err != nil {
		return nil, err
	}

	err = luaVm.LoadMainLuaFile(luaFile)
	if err != nil {
		return nil, err
	}

	if app.Port == 0 {
		return nil, fmt.Errorf("invalid application port, please use 'app.listen()'")
	}

	return &Runtime{
		LuaVm:  luaVm,
		Engine: engine,
		Port:   app.Port,
	}, nil
}

func (r *Runtime) Listen(addr string) error {
	log.Printf("HTTP server starting at http://%s...\n", addr)
	err := http.ListenAndServe(addr, r.Engine)
	if err != nil {
		log.Println("failed to start server:", err)
		return err
	}

	return nil
}

func (r *Runtime) Close() {
	if r.LuaVm != nil {
		r.LuaVm.Close()
	}
}
