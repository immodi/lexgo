package router

import (
	"log"
	"net/http"

	"immodi/lexgo/internal/middlewares"
	"immodi/lexgo/internal/vm"

	lua "github.com/yuin/gopher-lua"
)

type Router struct {
	LuaVm           *vm.LuaVm
	Routes          map[Handler]*lua.LFunction
	MiddleWares     []*lua.LFunction
	NotFoundFunc    *lua.LFunction
	ServerErrorFunc *lua.LFunction
}

type Handler struct {
	Path   string
	Method HTTPMethod
}

func MakeRouter(vm *vm.LuaVm) *Router {
	router := &Router{
		LuaVm:       vm,
		Routes:      make(map[Handler]*lua.LFunction),
		MiddleWares: make([]*lua.LFunction, 0),
	}
	router.RegisterFramework()
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := Handler{Path: req.URL.Path, Method: HTTPMethod(req.Method)}
	fn, ok := router.Routes[handler]

	luaReq := &LuaRequest{HttpRequest: req, LuaVm: router.LuaVm}
	luaRes := &LuaResponse{HttpWriter: w, LuaVm: router.LuaVm, Written: false}

	if !ok {
		router.handleNotFound(luaReq, luaRes, req)
		return
	}

	// Execute middlewares first, then the handler
	if len(router.MiddleWares) > 0 {
		ctx := middlewares.NewContext(
			router,
			luaReq,
			luaRes,
			fn,
		)

		middlewares.Execute(ctx, router.MiddleWares)
	} else {
		router.executeLuaHandler(fn, luaReq, luaRes)
	}
}

func (router *Router) executeLuaHandler(fn *lua.LFunction, luaReq *LuaRequest, luaRes *LuaResponse) {
	L := router.LuaVm.L
	L.Push(fn)
	L.Push(luaReq.MakeLuaRequest())
	L.Push(luaRes.MakeLuaResponse())

	if err := L.PCall(2, 0, nil); err != nil {
		router.handleServerError(err.Error(), luaRes)
	}
}

func (router *Router) handleNotFound(luaReq *LuaRequest, luaRes *LuaResponse, req *http.Request) {
	if router.NotFoundFunc == nil {
		log.Println("Lua error: no registered `app.notFound()` handler")
		http.NotFound(luaRes.HttpWriter, req)
		return
	}

	L := router.LuaVm.L
	L.Push(router.NotFoundFunc)
	L.Push(luaReq.MakeLuaRequest())
	L.Push(luaRes.MakeLuaResponse())

	if err := L.PCall(2, 0, nil); err != nil {
		router.handleServerError(err.Error(), luaRes)
	}
}

func (router *Router) handleServerError(errMsg string, luaRes *LuaResponse) {
	if router.ServerErrorFunc == nil {
		log.Printf("Lua error: no registered `app.error()` handler: %s", errMsg)
		http.Error(luaRes.HttpWriter, errMsg, http.StatusInternalServerError)
		return
	}

	L := router.LuaVm.L
	L.Push(router.ServerErrorFunc)
	L.Push(lua.LString(errMsg))
	L.Push(luaRes.MakeLuaResponse())

	if err := L.PCall(2, 0, nil); err != nil {
		log.Printf("Lua error in error handler: %s", err)
		http.Error(luaRes.HttpWriter, err.Error(), http.StatusInternalServerError)
	}
}
