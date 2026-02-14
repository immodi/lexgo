package vm

import lua "github.com/yuin/gopher-lua"

// Public interface used across the VM and modules.
type LuaValue interface {
	String() string
	Type() LValueType
}

// --- Concrete value types ---------------------------------------------------

type LuaNil struct{}

type LuaBool lua.LBool
type LuaNumber lua.LNumber
type LuaString lua.LString

// wrap pointer types so we can hold and forward the underlying lua objects
type LuaFunction struct{ *lua.LFunction }
type LuaTable struct{ *lua.LTable }

// --- LValueType + names (kept as you had them) ------------------------------

type LValueType int

const (
	LTNil LValueType = iota
	LTBool
	LTNumber
	LTString
	LTFunction
	LTUserData
	LTThread
	LTTable
	LTChannel
)

var lValueNames = [9]string{
	"nil", "boolean", "number", "string", "function", "userdata", "thread", "table", "channel",
}

func (vt LValueType) String() string { return lValueNames[int(vt)] }

// --- Implementations of String() and Type() --------------------------------

// Nil
func (LuaNil) String() string   { return "nil" }
func (LuaNil) Type() LValueType { return LTNil }

// Bool
func (v LuaBool) String() string   { return lua.LBool(v).String() }
func (v LuaBool) Type() LValueType { return LTBool }

// Number
func (v LuaNumber) String() string   { return lua.LNumber(v).String() }
func (v LuaNumber) Type() LValueType { return LTNumber }

// String
func (v LuaString) String() string   { return string(v) }
func (v LuaString) Type() LValueType { return LTString }

// Function
func (v *LuaFunction) String() string {
	if v == nil || v.LFunction == nil {
		return "function:nil"
	}
	return v.LFunction.String()
}
func (v *LuaFunction) Type() LValueType { return LTFunction }

// Table
func (t *LuaTable) String() string {
	if t == nil || t.LTable == nil {
		return "table:nil"
	}
	return t.LTable.String()
}
func (t *LuaTable) Type() LValueType { return LTTable }

// --- Converters between lua.LValue and your LuaValue ------------------------

// FromLua converts a gopher-lua value into your LuaValue wrapper.
func FromLua(lv lua.LValue) LuaValue {
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
		// Unknown/unsupported types map to nil (adjust if you want to support userdata/thread etc.)
		return LuaNil{}
	}
}

// ToLua converts your LuaValue into a gopher-lua lua.LValue for VM operations.
func ToLua(v LuaValue) lua.LValue {
	if v == nil {
		return lua.LNil
	}

	switch vv := v.(type) {
	case LuaNil:
		return lua.LNil
	case LuaBool:
		return lua.LBool(vv)
	case LuaNumber:
		return lua.LNumber(vv)
	case LuaString:
		return lua.LString(vv)
	case *LuaFunction:
		if vv == nil || vv.LFunction == nil {
			return lua.LNil
		}
		return vv.LFunction
	case *LuaTable:
		if vv == nil || vv.LTable == nil {
			return lua.LNil
		}
		return vv.LTable
	default:
		// If some external implementation of LuaValue is provided, attempt best-effort conversion:
		// fallback to String() as a safe lua.LString.
		return lua.LString(v.String())
	}
}

// --- LuaTable helper methods using LuaValue --------------------------------

// SetField sets a string-keyed field on the table.
func (t *LuaTable) SetField(key string, value LuaValue) {
	if t == nil || t.LTable == nil {
		return
	}
	t.LTable.RawSetString(key, ToLua(value))
}

// Set sets a key/value pair where key and value are arbitrary LuaValue.
func (t *LuaTable) Set(key LuaValue, value LuaValue) {
	if t == nil || t.LTable == nil {
		return
	}
	t.LTable.RawSet(ToLua(key), ToLua(value))
}

// GetField returns the value at the string key as a LuaValue.
func (t *LuaTable) GetField(key string) LuaValue {
	if t == nil || t.LTable == nil {
		return LuaNil{}
	}
	return FromLua(t.LTable.RawGetString(key))
}

func GenericGetField[T LuaValue](field *LuaTable, key string) (T, bool) {
	val, ok := field.GetField(key).(T)
	return val, ok
}

// Get returns the value for an arbitrary key.
func (t *LuaTable) Get(key LuaValue) LuaValue {
	if t == nil || t.LTable == nil {
		return LuaNil{}
	}
	return FromLua(t.LTable.RawGet(ToLua(key)))
}

// Append pushes a value onto the array portion of the table.
func (t *LuaTable) Append(value LuaValue) {
	if t == nil || t.LTable == nil {
		return
	}
	t.LTable.Append(ToLua(value))
}

// Len returns the length of the table (as Lua's # operator would).
func (t *LuaTable) Len() int {
	if t == nil || t.LTable == nil {
		return 0
	}
	return t.LTable.Len()
}

// ForEach iterates over the table and calls fn with wrapped LuaValue key/values.
func (t *LuaTable) ForEach(fn func(key, value LuaValue)) {
	if t == nil || t.LTable == nil {
		return
	}
	t.LTable.ForEach(func(k, v lua.LValue) {
		fn(FromLua(k), FromLua(v))
	})
}
