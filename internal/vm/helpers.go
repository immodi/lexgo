package vm

import (
	"fmt"
)

func luaTableToMap(tbl *LuaTable) map[string]any {
	result := make(map[string]any)
	tbl.ForEach(func(key, value LuaValue) {
		var k string
		switch key := key.(type) {
		case LuaString:
			k = string(key)
		case LuaNumber:
			k = fmt.Sprintf("%v", float64(key))
		default:
			k = ToLua(key).String()
		}

		switch v := value.(type) {
		case LuaString:
			result[k] = string(v)
		case LuaNumber:
			result[k] = float64(v)
		case LuaBool:
			result[k] = bool(v)
		case *LuaTable:
			result[k] = luaTableToMap(v)
		default:
			result[k] = ToLua(v).String()
		}
	})
	return result
}

func mapToLuaTable(LVm LVm, data map[string]any) *LuaTable {
	tbl := LVm.NewTable()
	for k, v := range data {
		switch val := v.(type) {
		case string:
			tbl.SetField(k, LuaString(val))
		case float64:
			tbl.SetField(k, LuaNumber(val))
		case bool:
			tbl.SetField(k, LuaBool(val))
		case map[string]any:
			tbl.SetField(k, mapToLuaTable(LVm, val))
		default:
			tbl.SetField(k, LuaString(fmt.Sprintf("%v", val)))
		}
	}
	return tbl
}
