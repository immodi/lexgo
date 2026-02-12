package defaultmiddlewares

import (
	"immodi/lexgo/internal/vm"
	"strings"
)

func DefaultLuaCORS(LVm vm.LVm, getRoutes func() map[string][]string) *vm.LuaFunction {
	return LVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		req, err := l.CheckTable(1)
		if err != nil {
			return nil
		}
		res, err := l.CheckTable(2)
		if err != nil {
			return nil
		}
		next, err := l.CheckFunction(3)
		if err != nil {
			return nil
		}

		setHeaderFn, ok := vm.GenericGetField[*vm.LuaFunction](res, "setHeader")
		if !ok {
			return nil
		}

		l.RunFunction(
			setHeaderFn,
			vm.LuaString("Access-Control-Allow-Origin"),
			vm.LuaString("*"),
		)
		url := req.GetField("url").String()
		registerdRotues := getRoutes()
		allowedMethods, ok := registerdRotues[url]
		allowedMethodsString := "GET, POST, PUT, DELETE, OPTIONS"
		// TODO: when in production mode, make it ""

		if ok {
			allowedMethodsString = strings.Join(allowedMethods, ", ")
		}
		l.RunFunction(
			setHeaderFn,
			vm.LuaString("Access-Control-Allow-Methods"),
			vm.LuaString(allowedMethodsString),
		)
		l.RunFunction(
			setHeaderFn,
			vm.LuaString("Access-Control-Allow-Headers"),
			vm.LuaString("Content-Type, Authorization"),
		)

		method := req.GetField("method").String()
		if method == "OPTIONS" {
			statusFn, ok := vm.GenericGetField[*vm.LuaFunction](res, "status")
			if !ok {
				return nil
			}
			l.RunFunction(statusFn, vm.LuaNumber(204))
		}

		l.RunFunction(next)
		return nil
	})
}
