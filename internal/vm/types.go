package vm

import lua "github.com/yuin/gopher-lua"

type LuaValue interface {
	ToLuaValue() lua.LValue
}

type (
	LuaNil      lua.LNilType
	LuaBool     lua.LBool
	LuaNumber   lua.LNumber
	LuaString   lua.LString
	LuaFunction struct{ *lua.LFunction }
	LuaTable    struct{ *lua.LTable }
)

func (v LuaNil) ToLuaValue() lua.LValue      { return lua.LNil }
func (v LuaBool) ToLuaValue() lua.LValue     { return lua.LBool(v) }
func (v LuaNumber) ToLuaValue() lua.LValue   { return lua.LNumber(v) }
func (v LuaString) ToLuaValue() lua.LValue   { return lua.LString(v) }
func (v LuaFunction) ToLuaValue() lua.LValue { return v.LFunction }
func (v LuaTable) ToLuaValue() lua.LValue    { return v.LTable }

func (t *LuaTable) SetField(key string, value LuaValue) {
	t.LTable.RawSetString(key, value.ToLuaValue())
}

func (t *LuaTable) Set(key LuaValue, value LuaValue) {
	t.LTable.RawSet(key.ToLuaValue(), value.ToLuaValue())
}

func (t *LuaTable) GetField(key string) lua.LValue {
	return t.LTable.RawGetString(key)
}

func (t *LuaTable) Get(key LuaValue) lua.LValue {
	return t.LTable.RawGet(key.ToLuaValue())
}

func (t *LuaTable) Append(value LuaValue) {
	t.LTable.Append(value.ToLuaValue())
}

func (t *LuaTable) Len() int {
	return t.LTable.Len()
}

func (t *LuaTable) ForEach(fn func(key, value LuaValue)) {
	t.LTable.ForEach(func(k, v lua.LValue) {
		key := convertToLuaValue(k)
		value := convertToLuaValue(v)
		fn(key, value)
	})
}

func convertToLuaValue(lv lua.LValue) LuaValue {
	switch v := lv.(type) {
	case *lua.LNilType:
		return LuaNil{}
	case lua.LBool:
		return LuaBool(v)
	case lua.LNumber:
		return LuaNumber(v)
	case lua.LString:
		return LuaString(v)
	case *lua.LFunction:
		return &LuaFunction{LFunction: v}
	case *lua.LTable:
		return &LuaTable{LTable: v}
	default:
		// For unsupported types, you could return nil or panic
		return LuaNil{}
	}
}
