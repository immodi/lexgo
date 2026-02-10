package vm

import (
	"net/http"
)

type HTTPRoute struct {
	Path   string
	Method HTTPMethod
}

type RouterDriver interface {
	RegisterLuaMethodHandler(fn *LuaFunction, path string, method string)
	ResgisterLuaErrorHandler(fn *LuaFunction)
	ResgisterLuaNotFoundHandler(fn *LuaFunction)
	RegisterLuaMiddleware(fn *LuaFunction)
}

func ResgisterRouter(luaVm LVm, routerDriver RouterDriver) *LuaTable {
	appTbl := luaVm.NewTable()

	registerMethod := func(method string) *LuaFunction {
		return luaVm.NewFunction(func(l LVm) LuaValue {
			path, err := l.CheckString(1)
			if err != nil {
				l.Error(err.Error())
				return nil
			}

			fn, err := l.CheckFunction(2)
			if err != nil {
				l.Error(err.Error())
				return nil
			}

			routerDriver.RegisterLuaMethodHandler(fn, path, method)
			return nil
		})
	}

	appTbl.SetField("get", registerMethod("GET"))
	appTbl.SetField("post", registerMethod("POST"))
	appTbl.SetField("put", registerMethod("PUT"))
	appTbl.SetField("delete", registerMethod("DELETE"))
	appTbl.SetField("patch", registerMethod("PATCH"))
	appTbl.SetField("options", registerMethod("OPTIONS"))

	appTbl.SetField("notFound", luaVm.NewFunction(func(l LVm) LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error(err.Error())
			return nil
		}

		routerDriver.ResgisterLuaNotFoundHandler(fn)
		return nil
	}))

	appTbl.SetField("error", luaVm.NewFunction(func(l LVm) LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error(err.Error())
			return nil
		}

		routerDriver.ResgisterLuaErrorHandler(fn)
		return nil
	}))

	appTbl.SetField("use", luaVm.NewFunction(func(l LVm) LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error(err.Error())
			return nil
		}

		routerDriver.RegisterLuaMiddleware(fn)
		return nil
	}))

	return appTbl
}

func ExecuteLuaHandler(luaVm LVm, errFn *LuaFunction, fn *LuaFunction, luaReq *LuaRequest, luaRes *LuaResponse) {
	if fn == nil {
		http.Error(luaRes.HttpWriter, "404 Not Found", http.StatusNotFound)
		return
	}

	err := luaVm.RunFunction(fn, luaReq.MakeLuaRequest(), luaRes.MakeLuaResponse())
	if err != nil {
		HandleServerError(luaVm, errFn, err.Error(), luaRes)
	} else {
		luaRes.Flush()
	}
}

func HandleServerError(luaVm LVm, errFn *LuaFunction, errMsg string, luaRes *LuaResponse) {
	luaRes.Reset()

	if errFn == nil {
		http.Error(luaRes.HttpWriter, errMsg, http.StatusInternalServerError)
		return
	}

	err := luaVm.RunFunction(errFn, LuaString(errMsg), luaRes.MakeLuaResponse())
	if err != nil {
		luaRes.Reset()
		http.Error(luaRes.HttpWriter, err.Error(), http.StatusInternalServerError)
	} else {
		luaRes.Flush()
	}
}
