package router

import (
	"log"
	"net/http"

	lua "github.com/yuin/gopher-lua"
	"immodi/lexgo/internal/vm"
)

type Router struct {
	LuaVm        *vm.LuaVm
	Routes       map[Hanlder]*lua.LFunction
	NotFoundFunc *lua.LFunction
}

type Hanlder struct {
	Path   string
	Method HTTPMethod
}

func MakeRouter(vm *vm.LuaVm) *Router {
	router := &Router{
		LuaVm:  vm,
		Routes: make(map[Hanlder]*lua.LFunction),
	}

	router.RegisterFramework()
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fn, ok := router.Routes[Hanlder{Path: req.URL.Path, Method: HTTPMethod(req.Method)}]
	luaReq := &LuaRequest{HttpRequest: req, LuaVm: router.LuaVm}
	luaRes := &LuaResponse{HttpWriter: w, LuaVm: router.LuaVm}

	if !ok {
		if router.NotFoundFunc != nil {
			router.LuaVm.L.Push(router.NotFoundFunc)
			router.LuaVm.L.Push(luaReq.MakeLuaRequest())
			router.LuaVm.L.Push(luaRes.MakeLuaResponse())
			if err := router.LuaVm.L.PCall(2, 0, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.NotFound(w, req)
		return
	}

	L := router.LuaVm.L
	L.Push(fn)
	L.Push(luaReq.MakeLuaRequest())
	L.Push(luaRes.MakeLuaResponse())

	if err := L.PCall(2, 0, nil); err != nil {
		log.Printf("Lua error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (router *Router) RegisterFramework() {
	L := router.LuaVm.L

	tbl := L.NewTable()
	L.SetGlobal("lexgo", tbl)

	L.SetField(tbl, "new", L.NewFunction(func(l *lua.LState) int {
		app := router.MakeLuaAppTable()
		l.Push(app)
		return 1
	}))
}

func (router *Router) MakeLuaAppTable() *lua.LTable {
	L := router.LuaVm.L

	appTbl := L.NewTable()

	L.SetField(appTbl, "get", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Hanlder{Path: path, Method: GET}] = fn
		return 0
	}))

	L.SetField(appTbl, "post", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Hanlder{Path: path, Method: POST}] = fn
		return 0
	}))

	L.SetField(appTbl, "notFound", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		router.NotFoundFunc = fn
		return 0
	}))

	return appTbl
}
