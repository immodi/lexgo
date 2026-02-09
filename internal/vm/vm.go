package vm

import (
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type LVm interface {
	LoadMainLuaFile(mainFilePath string) error
	RunFunction(fn LuaFunction, args ...LuaValue) error
	GetFunction(name string) (*LuaFunction, error)
	GetGlobal(name string) (lua.LValue, error)
	RegisterGlobal(name string, value LuaValue)
	NewTable() LuaTable
	Close()
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
	luaVm.L.Close()
}

func (luaVm *LuaVm) DoString(code string) error {
	return luaVm.withLock(func(L *lua.LState) error {
		return L.DoString(code)
	})
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

func (luaVm *LuaVm) GetGlobal(name string) (lua.LValue, error) {
	var result lua.LValue
	err := luaVm.withLock(func(L *lua.LState) error {
		result = L.GetGlobal(name)
		return nil
	})
	return result, err
}

func (luaVm *LuaVm) LoadMainLuaFile(mainFilePath string) error {
	return luaVm.withLock(func(L *lua.LState) error {
		return L.DoFile(mainFilePath)
	})
}

func (luaVm *LuaVm) NewTable() LuaTable {
	return LuaTable{LTable: luaVm.L.NewTable()}
}

func (luaVm *LuaVm) RegisterGlobal(name string, value LuaValue) {
	luaVm.withLock(func(L *lua.LState) error {
		L.SetGlobal(name, value.ToLuaValue())
		return nil
	})
}

func (luaVm *LuaVm) RunFunction(fn LuaFunction, args ...LuaValue) error {
	return luaVm.withLock(func(L *lua.LState) error {
		luaVm.push(fn)
		for _, arg := range args {
			luaVm.push(arg)
		}
		return luaVm.pcall(len(args), 0)
	})
}

// Internal/low-level methods (alphabetical)

func (luaVm *LuaVm) call(nargs int, nret int) {
	luaVm.L.Call(nargs, nret)
}

func (luaVm *LuaVm) getTop() int {
	return luaVm.L.GetTop()
}

func (luaVm *LuaVm) pcall(nargs int, nret int) error {
	return luaVm.L.PCall(nargs, nret, nil)
}

func (luaVm *LuaVm) pop(n int) {
	luaVm.L.Pop(n)
}

func (luaVm *LuaVm) push(value LuaValue) {
	luaVm.L.Push(value.ToLuaValue())
}

func (luaVm *LuaVm) withLock(fn func(L *lua.LState) error) error {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()
	return fn(luaVm.L)
}
