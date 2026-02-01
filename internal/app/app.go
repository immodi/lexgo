package app

import (
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"log"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type App struct {
	LuaVm  *lua.LState
	Router *router.Router
}

func New(luaFile string) (*App, error) {
	luaVm := vm.MakeLuaVm()
	router, routerDriver := router.MakeRouter(luaVm)
	vm.RegisterFramework(router.LuaVm.L, routerDriver)

	err := luaVm.LoadMainLuaFile(luaFile)
	if err != nil {
		return nil, err
	}

	return &App{
		LuaVm:  luaVm.L,
		Router: router,
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
