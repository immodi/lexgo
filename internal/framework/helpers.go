package framework

import (
	"fmt"
	"immodi/lexgo/internal/vm"
)

func luaTableToMap(tbl *vm.LuaTable) map[string]any {
	result := make(map[string]any)
	tbl.ForEach(func(key, value vm.LuaValue) {
		var k string
		switch key := key.(type) {
		case vm.LuaString:
			k = string(key)
		case vm.LuaNumber:
			k = fmt.Sprintf("%v", float64(key))
		default:
			k = vm.ToLua(key).String()
		}

		switch v := value.(type) {
		case vm.LuaString:
			result[k] = string(v)
		case vm.LuaNumber:
			result[k] = float64(v)
		case vm.LuaBool:
			result[k] = bool(v)
		case *vm.LuaTable:
			result[k] = luaTableToMap(v)
		default:
			result[k] = vm.ToLua(v).String()
		}
	})
	return result
}

func mapToLuaTable(LVm vm.LVm, data map[string]any) *vm.LuaTable {
	tbl := LVm.NewTable()
	for k, v := range data {
		switch val := v.(type) {
		case string:
			tbl.SetField(k, vm.LuaString(val))
		case float64:
			tbl.SetField(k, vm.LuaNumber(val))
		case bool:
			tbl.SetField(k, vm.LuaBool(val))
		case map[string]any:
			tbl.SetField(k, mapToLuaTable(LVm, val))
		default:
			tbl.SetField(k, vm.LuaString(fmt.Sprintf("%v", val)))
		}
	}
	return tbl
}
