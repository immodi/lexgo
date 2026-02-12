package vm

import (
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type LVm interface {
	LoadMainLuaFile(mainFilePath string) error
	RunFunction(fn *LuaFunction, args ...LuaValue) error
	GetFunction(name string) (*LuaFunction, error)
	GetGlobal(name string) (LuaValue, error)
	SetGlobal(name string, value LuaValue)
	NewTable() *LuaTable
	NewFunction(fn func(LVm) LuaValue) *LuaFunction
	Return(values ...LuaValue) int
	NewMultiReturnFunction(fn func(LVm) int) *LuaFunction
	Error(mes string)
	Close()

	// Check functions - validate and retrieve typed values from stack
	CheckString(index int) (string, error)
	CheckNumber(index int) (float64, error)
	CheckBool(index int) (bool, error)
	CheckTable(index int) (*LuaTable, error)
	CheckFunction(index int) (*LuaFunction, error)
}

type LuaVm struct {
	L  *lua.LState
	mu sync.Mutex
}

func MakeLuaVm() *LuaVm {
	L := lua.NewState()
	return &LuaVm{L: L}
}

// Public API methods (alphabetical)

func (luaVm *LuaVm) Close() {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()
	luaVm.L.Close()
}

func (luaVm *LuaVm) DoString(code string) error {
	return luaVm.L.DoString(code)
}

func (luaVm *LuaVm) GetFunction(name string) (*LuaFunction, error) {
	var fn *LuaFunction
	err := luaVm.withLock(func(L *lua.LState) error {
		lv := L.GetGlobal(name)
		if lv.Type() != lua.LTFunction {
			return fmt.Errorf("'%s' is not a function", name)
		}
		lfn, ok := lv.(*lua.LFunction)
		if !ok {
			return fmt.Errorf("failed to convert to LFunction")
		}
		fn = &LuaFunction{LFunction: lfn}
		return nil
	})
	return fn, err
}

func (luaVm *LuaVm) GetGlobal(name string) (LuaValue, error) {
	var result lua.LValue
	err := luaVm.withLock(func(L *lua.LState) error {
		result = L.GetGlobal(name)
		return nil
	})
	out := FromLua(result)
	return out, err
}

func (luaVm *LuaVm) SetGlobal(name string, value LuaValue) {
	luaVm.L.SetGlobal(name, ToLua(value))
}

func (luaVm *LuaVm) LoadMainLuaFile(mainFilePath string) error {
	return luaVm.L.DoFile(mainFilePath)
}

func (luaVm *LuaVm) NewTable() *LuaTable {
	var tbl *lua.LTable
	luaVm.withLock(func(L *lua.LState) error {
		tbl = L.NewTable()
		return nil
	})
	return &LuaTable{LTable: tbl}
}

func (luaVm *LuaVm) NewFunction(fn func(LVm) LuaValue) *LuaFunction {
	var lfn *lua.LFunction
	luaVm.withLock(func(L *lua.LState) error {
		lfn = L.NewFunction(func(l *lua.LState) int {
			result := fn(luaVm)
			if result != nil && result.Type() != LTNil {
				luaVm.push(result)
				return 1
			}
			return 0
		})
		return nil
	})
	return &LuaFunction{LFunction: lfn}
}

func (luaVm *LuaVm) NewMultiReturnFunction(fn func(LVm) int) *LuaFunction {
	var lfn *lua.LFunction
	luaVm.withLock(func(L *lua.LState) error {
		lfn = L.NewFunction(func(l *lua.LState) int {
			return fn(luaVm)
		})
		return nil
	})
	return &LuaFunction{LFunction: lfn}
}

func (luaVm *LuaVm) Return(values ...LuaValue) int {
	for _, v := range values {
		luaVm.push(v)
	}
	return len(values)
}

func (luaVm *LuaVm) RunFunction(fn *LuaFunction, args ...LuaValue) error {
	top := luaVm.L.GetTop()
	defer luaVm.L.SetTop(top)

	luaVm.push(fn)
	for _, arg := range args {
		luaVm.push(arg)
	}

	return luaVm.pcall(len(args), 0)
}

func (luaVm *LuaVm) Error(mes string) {
	// No lock - called from within Lua execution
	luaVm.L.RaiseError(mes)
}

// Check functions - NO LOCKS (called from within Lua execution context)
func (luaVm *LuaVm) CheckString(index int) (string, error) {
	lv := luaVm.L.Get(index)
	if lv.Type() != lua.LTString {
		return "", fmt.Errorf("expected string at index %d, got %s", index, lv.Type())
	}
	return string(lv.(lua.LString)), nil
}

func (luaVm *LuaVm) CheckNumber(index int) (float64, error) {
	lv := luaVm.L.Get(index)
	if lv.Type() != lua.LTNumber {
		return 0, fmt.Errorf("expected number at index %d, got %s", index, lv.Type())
	}
	return float64(lv.(lua.LNumber)), nil
}

func (luaVm *LuaVm) CheckBool(index int) (bool, error) {
	lv := luaVm.L.Get(index)
	if lv.Type() != lua.LTBool {
		return false, fmt.Errorf("expected bool at index %d, got %s", index, lv.Type())
	}
	return bool(lv.(lua.LBool)), nil
}

func (luaVm *LuaVm) CheckTable(index int) (*LuaTable, error) {
	lv := luaVm.L.Get(index)
	if lv.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected table at index %d, got %s", index, lv.Type())
	}
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("failed to convert to table")
	}
	return &LuaTable{LTable: tbl}, nil
}

func (luaVm *LuaVm) CheckFunction(index int) (*LuaFunction, error) {
	lv := luaVm.L.Get(index)
	if lv.Type() != lua.LTFunction {
		return nil, fmt.Errorf("expected function at index %d, got %s", index, lv.Type())
	}
	fn, ok := lv.(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("failed to convert to function")
	}
	return &LuaFunction{LFunction: fn}, nil
}

// Internal/low-level methods - assume lock is held

func (luaVm *LuaVm) pcall(nargs int, nret int) error {
	return luaVm.L.PCall(nargs, nret, nil)
}

func (luaVm *LuaVm) push(value LuaValue) {
	luaVm.L.Push(ToLua(value))
}

func (luaVm *LuaVm) withLock(fn func(L *lua.LState) error) error {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()
	return fn(luaVm.L)
}
