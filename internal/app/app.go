package app

import (
	"fmt"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"log"
	"net/http"
)

type App struct {
	LuaVm  vm.LVm
	Router *router.Router
	Port   int32
}

func New(luaFile string) (*App, error) {
	luaVm := vm.MakeLuaVm()
	router, routerDriver := router.MakeRouter(luaVm)
	data, err := vm.RegisterFramework(router.LuaVm, routerDriver)
	if err != nil {
		return nil, err
	}

	err = luaVm.LoadMainLuaFile(luaFile)
	if err != nil {
		return nil, err
	}

	if data.Port == 0 {
		return nil, fmt.Errorf("invalid application port, please use 'app.listen()'")
	}

	return &App{
		LuaVm:  luaVm,
		Router: router,
		Port:   data.Port,
	}, nil
}

func (a *App) Listen(addr string) error {
	log.Printf("HTTP server starting at http://%s...\n", addr)
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Println("failed to start server:", err)
		return err
	}

	return nil
}

func (a *App) Close() {
	if a.LuaVm != nil {
		a.LuaVm.Close()
	}
}
