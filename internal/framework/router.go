package framework

import (
	"immodi/lexgo/internal/vm"
	"net/http"
)

type HTTPRoute struct {
	Path   string
	Method HTTPMethod
}

type RouterDriver interface {
	RegisterLuaMethodHandler(fn *vm.LuaFunction, path string, method string)
	ResgisterLuaErrorHandler(fn *vm.LuaFunction)
	ResgisterLuaNotFoundHandler(fn *vm.LuaFunction)
	RegisterLuaMiddleware(fn *vm.LuaFunction)

	GetAllRegistredRoutes() map[string][]string
}

func ResgisterRouter(luaVm vm.LVm, routerDriver RouterDriver) *vm.LuaTable {
	appTbl := luaVm.NewTable()

	registerMethod := func(method string) *vm.LuaFunction {
		return luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
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

	appTbl.SetField("notFound", luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error(err.Error())
			return nil
		}

		routerDriver.ResgisterLuaNotFoundHandler(fn)
		return nil
	}))

	appTbl.SetField("error", luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error(err.Error())
			return nil
		}

		routerDriver.ResgisterLuaErrorHandler(fn)
		return nil
	}))

	appTbl.SetField("use", luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
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

func ExecuteLuaHandler(luaVm vm.LVm, errFn *vm.LuaFunction, fn *vm.LuaFunction, luaReq *LuaRequest, luaRes *LuaResponse) {
	if fn == nil {
		luaRes.Flush()
		return
	}

	if err := luaVm.RunFunction(
		fn,
		luaReq.MakeLuaRequest(),
		luaRes.MakeLuaResponse(),
	); err != nil {
		HandleServerError(luaVm, errFn, err.Error(), luaRes)
	}

	luaRes.Flush()
}

func HandleServerError(luaVm vm.LVm, errFn *vm.LuaFunction, errMsg string, luaRes *LuaResponse) {
	luaRes.Reset()

	if errFn == nil {
		http.Error(luaRes.HttpWriter, errMsg, http.StatusInternalServerError)
		return
	}

	err := luaVm.RunFunction(errFn, vm.LuaString(errMsg), luaRes.MakeLuaResponse())
	if err != nil {
		luaRes.Reset()
		http.Error(luaRes.HttpWriter, err.Error(), http.StatusInternalServerError)
	} else {
		luaRes.Flush()
	}
}
