package vm

import (
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type HTTPHandler struct {
	Path   string
	Method HTTPMethod
}

type RouterDriver interface {
	RegisterLuaMethodHandler(fn *lua.LFunction, path string, method string)
	ResgisterLuaErrorHandler(fn *lua.LFunction)
	ResgisterLuaNotFoundHandler(fn *lua.LFunction)
	RegisterLuaMiddleware(fn *lua.LFunction)
}

func ResgisterRouter(L *lua.LState, routerDriver RouterDriver) *lua.LTable {
	appTbl := L.NewTable()

	L.SetField(appTbl, "get", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: GET}] = fn
		routerDriver.RegisterLuaMethodHandler(fn, path, "GET")
		return 0
	}))

	L.SetField(appTbl, "post", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: POST}] = fn

		routerDriver.RegisterLuaMethodHandler(fn, path, "POST")
		return 0
	}))

	L.SetField(appTbl, "put", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: PUT}] = fn

		routerDriver.RegisterLuaMethodHandler(fn, path, "PUT")
		return 0
	}))

	L.SetField(appTbl, "delete", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: DELETE}] = fn

		routerDriver.RegisterLuaMethodHandler(fn, path, "DELETE")
		return 0
	}))

	L.SetField(appTbl, "patch", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: PATCH}] = fn

		routerDriver.RegisterLuaMethodHandler(fn, path, "PATCH")
		return 0
	}))

	L.SetField(appTbl, "options", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		// router.Routes[Handler{Path: path, Method: OPTIONS}] = fn

		routerDriver.RegisterLuaMethodHandler(fn, path, "OPTIONS")
		return 0
	}))

	L.SetField(appTbl, "notFound", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		// router.NotFoundFunc = fn
		routerDriver.ResgisterLuaNotFoundHandler(fn)
		return 0
	}))

	L.SetField(appTbl, "error", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		// router.ServerErrorFunc = fn
		routerDriver.ResgisterLuaErrorHandler(fn)
		return 0
	}))

	L.SetField(appTbl, "use", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		// router.MiddleWares = append(router.MiddleWares, fn)
		routerDriver.RegisterLuaMiddleware(fn)
		return 0
	}))

	return appTbl
}

func ExecuteLuaHandler(L *lua.LState, errFn *lua.LFunction, fn *lua.LFunction, luaReq *LuaRequest, luaRes *LuaResponse) {
	L.Push(fn)
	L.Push(luaReq.MakeLuaRequest())
	L.Push(luaRes.MakeLuaResponse())
	defer L.SetTop(0)

	if err := L.PCall(2, 0, nil); err != nil {
		// log.Printf("Lua error in handler: %s", err)
		HandleServerError(L, errFn, err.Error(), luaRes)
	}
}

func HandleServerError(L *lua.LState, errFn *lua.LFunction, errMsg string, luaRes *LuaResponse) {
	if errFn == nil {
		// log.Printf("Lua error: no registered `app.error()` handler: %s", errMsg)
		http.Error(luaRes.HttpWriter, errMsg, http.StatusInternalServerError)
		return
	}

	defer L.SetTop(0)

	L.Push(errFn)
	L.Push(lua.LString(errMsg))
	L.Push(luaRes.MakeLuaResponse())

	if err := L.PCall(2, 0, nil); err != nil {
		// log.Printf("Lua error in error handler: %s", err)
		http.Error(luaRes.HttpWriter, err.Error(), http.StatusInternalServerError)
	}
}
