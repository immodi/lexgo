package router

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func luaTableToMap(tbl *lua.LTable) map[string]any {
	result := make(map[string]any)
	tbl.ForEach(func(key, value lua.LValue) {
		var k string
		switch key := key.(type) {
		case lua.LString:
			k = string(key)
		case lua.LNumber:
			k = fmt.Sprintf("%v", float64(key))
		default:
			k = key.String()
		}

		switch v := value.(type) {
		case lua.LString:
			result[k] = string(v)
		case lua.LNumber:
			result[k] = float64(v)
		case *lua.LTable:
			result[k] = luaTableToMap(v)
		case lua.LBool:
			result[k] = bool(v)
		default:
			result[k] = v.String()
		}
	})
	return result
}

func mapToLuaTable(L *lua.LState, data map[string]any) *lua.LTable {
	tbl := L.NewTable()
	for k, v := range data {
		switch val := v.(type) {
		case string:
			L.SetField(tbl, k, lua.LString(val))
		case float64:
			L.SetField(tbl, k, lua.LNumber(val))
		case bool:
			L.SetField(tbl, k, lua.LBool(val))
		case map[string]any:
			L.SetField(tbl, k, mapToLuaTable(L, val))
		default:
			L.SetField(tbl, k, lua.LString(fmt.Sprintf("%v", val)))
		}
	}
	return tbl
}
