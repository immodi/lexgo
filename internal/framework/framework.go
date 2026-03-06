package framework

import (
	"fmt"
	"immodi/lexgo/internal/vm"
)

func RegisterFramework(
	luaVm vm.LVm,
	routerDriver RouterDriver,
	restartServerChannel chan struct{},
) (*AppData, error) {

	tbl := luaVm.NewTable()
	data := &AppData{}

	luaVm.SetGlobal("lexgo", tbl)

	newFn := createNewAppFunction(luaVm, data, routerDriver, restartServerChannel)
	tbl.SetField("new", newFn)

	if err := registerLibraries(luaVm, tbl); err != nil {
		return nil, err
	}

	registerMiddlewares(luaVm, tbl, data, routerDriver)

	return data, nil
}

func registerLibraries(luaVm vm.LVm, tbl *vm.LuaTable) error {
	lx, err := LxTable(luaVm)
	if err != nil {
		return fmt.Errorf("failed to load 'lx' library into runtime")
	}

	tbl.SetField("lx", lx)
	return nil
}

func registerMiddlewares(luaVm vm.LVm, tbl *vm.LuaTable, data *AppData, routerDriver RouterDriver) {
	RegisterDefaultMiddlewares(luaVm, tbl, &CORSRuntime{
		appData:      data,
		routerDriver: routerDriver,
	})
}
