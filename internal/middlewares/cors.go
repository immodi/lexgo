package middlewares

import (
	"immodi/lexgo/internal/vm"
)

type CorsMiddlewareAppProvider interface {
	GetAllowedOrigin(requestOrigin string) string
	GetRegisterdRoutes() map[string][]string
	GetAllowedMethods(url string) string
}

func DefaultLuaCORS(LVm vm.LVm, appProvider CorsMiddlewareAppProvider) *vm.LuaFunction {
	return LVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		req, err := l.CheckTable(1)
		if err != nil {
			l.Error("internal error, failed to supply the 'req' table to the default cors middleware")
			return nil
		}
		res, err := l.CheckTable(2)
		if err != nil {
			l.Error("internal error, failed to supply the 'res' table to the default cors middleware")
			return nil
		}
		next, err := l.CheckFunction(3)
		if err != nil {
			l.Error("internal error, failed to supply the 'next' function to the default cors middleware")
			return nil
		}

		setHeaderFn, ok := vm.GenericGetField[*vm.LuaFunction](res, "setHeader")
		if !ok {
			l.Error("internal error, failed to supply the 'setHeader' function to the default cors middleware")
			return nil
		}

		origin := req.GetField("origin").String()
		allowedOrigin := appProvider.GetAllowedOrigin(origin)
		if allowedOrigin != "" {
			l.RunFunction(
				setHeaderFn,
				vm.LuaString("Access-Control-Allow-Origin"),
				vm.LuaString(allowedOrigin),
			)
		}

		url := req.GetField("url").String()
		allowedMethods := appProvider.GetAllowedMethods(url)

		if allowedMethods != "" {
			l.RunFunction(
				setHeaderFn,
				vm.LuaString("Access-Control-Allow-Methods"),
				vm.LuaString(allowedMethods),
			)
		}

		l.RunFunction(
			setHeaderFn,
			vm.LuaString("Access-Control-Allow-Headers"),
			vm.LuaString("Content-Type, Authorization"),
		)

		method := req.GetField("method").String()
		if method == "OPTIONS" {
			statusFn, ok := vm.GenericGetField[*vm.LuaFunction](res, "status")
			if !ok {
				l.Error("internal error, failed to supply the 'status' function to the default cors middleware")
				return nil
			}
			l.RunFunction(statusFn, vm.LuaNumber(204))
		}

		l.RunFunction(next)
		return nil
	})
}
