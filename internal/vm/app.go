package vm

import (
	lua "github.com/yuin/gopher-lua"
)

type AppData struct {
	Port int32
}

func RegisterFramework(L *lua.LState, routerDriver RouterDriver) (*AppData, error) {
	tbl := L.NewTable()
	data := &AppData{}

	L.SetGlobal("lexgo", tbl)
	L.SetField(tbl, "new", L.NewFunction(func(l *lua.LState) int {

		app := ResgisterRouter(L, routerDriver)

		L.SetField(app, "listen", L.NewFunction(func(L *lua.LState) int {
			port := L.CheckInt(1)
			L.SetField(app, "_port", lua.LNumber(port))
			data.Port = int32(port)
			return 0
		}))

		l.Push(app)
		return 1
	}))

	RegisterDefaultMiddlewares(L, tbl)
	return data, nil
}
