package vm

import lua "github.com/yuin/gopher-lua"

func RegisterFramework(L *lua.LState, routerDriver RouterDriver) {
	tbl := L.NewTable()
	L.SetGlobal("lexgo", tbl)

	L.SetField(tbl, "new", L.NewFunction(func(l *lua.LState) int {
		app := ResgisterRouter(L, routerDriver)
		l.Push(app)
		return 1
	}))

	RegisterDefaultMiddlewares(L, tbl)
}
