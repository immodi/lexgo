package framework

import (
	"fmt"
	"immodi/lexgo/internal/vm"
	"net/http"
)

type HTTPRoute struct {
	Path   string
	Method HTTPMethod
}

type RouterHandler interface {
	Handle(luaReq *LuaRequest, luaRes *LuaResponse, args ...vm.LuaValue) error
}

type RouterDriver interface {
	MakeLuaHandler(luaVm vm.LVm, fn *vm.LuaFunction) RouterHandler
	MakeGoHandler(fn func(w http.ResponseWriter, r *http.Request)) RouterHandler

	RegisterHandler(fn RouterHandler, path string, method string, mws []RouterHandler)
	RegisterErrorHandler(fn RouterHandler)
	RegisterNotFoundHandler(fn RouterHandler)
	RegisterMiddleware(fn RouterHandler)

	GetAllRegistredRoutes() map[string][]string
	ClearRoutes()
}

func RegisterApp(luaVm vm.LVm, routerDriver RouterDriver) *vm.LuaTable {
	appTbl := luaVm.NewTable()

	registerMethod := func(method string) *vm.LuaFunction {
		return luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
			path, err := l.CheckString(1)
			if err != nil {
				l.Error("expected a 'string' path argument")
				return nil
			}

			fn, err := l.CheckFunction(2)
			if err != nil {
				l.Error("expected a function argument with signiture 'function(req, res)'")
				return nil
			}

			mws, err := l.CheckVariadicFunctions(3)
			if err != nil {
				l.Error("expected one or more function argument with signiture 'function(req, res, next)'")
				return nil
			}

			luaHandler := routerDriver.MakeLuaHandler(l, fn)

			luaMws := make([]RouterHandler, 0)
			for _, mw := range mws {
				luaMws = append(luaMws, routerDriver.MakeLuaHandler(l, mw))
			}

			routerDriver.RegisterHandler(luaHandler, path, method, luaMws)
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
			l.Error("expected a function argument with signiture 'function(req, res)'")
			return nil
		}

		notFoundLuaHandler := routerDriver.MakeLuaHandler(l, fn)
		routerDriver.RegisterNotFoundHandler(notFoundLuaHandler)
		return nil
	}))

	appTbl.SetField("error", luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error("expected a function argument with signiture 'function(err, res)'")
			return nil
		}

		errorLuaHandler := routerDriver.MakeLuaHandler(l, fn)
		routerDriver.RegisterErrorHandler(errorLuaHandler)
		return nil
	}))

	appTbl.SetField("use", luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		fn, err := l.CheckFunction(1)
		if err != nil {
			l.Error("expected a function argument with signiture 'function(req, res, next)'")
			return nil
		}

		mw := routerDriver.MakeLuaHandler(l, fn)
		routerDriver.RegisterMiddleware(mw)
		return nil
	}))

	return appTbl
}

func ExecuteLuaHandler(luaVm vm.LVm, errFn RouterHandler, fn RouterHandler, luaReq *LuaRequest, luaRes *LuaResponse) {
	if fn == nil {
		luaRes.buf.Reset()
		luaRes.statusCode = http.StatusNotFound
		luaRes.buf.WriteString(fmt.Sprintf("Handler Not Found at => %s", luaReq.HttpRequest.URL))
		luaRes.Flush()
		return
	}

	if err := fn.Handle(luaReq, luaRes); err != nil {
		HandleServerError(luaVm, errFn, err.Error(), luaReq, luaRes)
		return
	}

	luaRes.Flush()
}

func HandleServerError(luaVm vm.LVm, errFn RouterHandler, errMsg string, luaReq *LuaRequest, luaRes *LuaResponse) {
	luaRes.Reset()

	if errFn == nil {
		http.Error(luaRes.HttpWriter, errMsg, http.StatusInternalServerError)
		return
	}

	if err := errFn.Handle(luaReq, luaRes, vm.LuaString(errMsg)); err != nil {
		luaRes.Reset()
		http.Error(luaRes.HttpWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	luaRes.Flush()
}
