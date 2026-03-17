package userlibs

import (
	"fmt"
	"immodi/lexgo/internal/vm"
	lx_embed "immodi/lexgo/lua/libs/lx"
)

func MakeLX(vm vm.LVm) (vm.LuaValue, error) {
	source := lx_embed.LxSource
	fn, err := vm.LoadLuaString(source)
	if err != nil {
		return nil, err
	}

	out, err := vm.RunFunctionMultiReturn(fn)
	if err != nil {
		return nil, err
	}

	if len(out) < 1 {
		return nil, fmt.Errorf("failed to get the 'lx' object from library source")
	}

	tbl := out[0]
	return tbl, nil
}
