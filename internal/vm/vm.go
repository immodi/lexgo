package vm

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type LuaVm struct {
	L  *lua.LState
	mu sync.Mutex
}

func MakeLuaVm() *LuaVm {
	L := lua.NewState()

	return &LuaVm{L: L}
}

func (luaVm *LuaVm) LoadMainLuaFile(mainFilePath string) error {
	return luaVm.WithLock(func(L *lua.LState) error {
		return L.DoFile(mainFilePath)
	})
}

func (luaVm *LuaVm) Push(value lua.LValue) {
	luaVm.WithLock(func(L *lua.LState) error {
		L.Push(value)
		return nil
	})
}

func (luaVm *LuaVm) Call(fn *lua.LFunction, nargs int, nret int) {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()

	luaVm.L.Call(nargs, nret)
}

func (luaVm *LuaVm) PCall(nargs int, nret int, fn *lua.LFunction) error {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()

	return luaVm.L.PCall(nargs, nret, nil)
}

func (luaVm *LuaVm) WithLock(fn func(L *lua.LState) error) error {
	luaVm.mu.Lock()
	defer luaVm.mu.Unlock()
	return fn(luaVm.L)
}

func (luaVm *LuaVm) Close() {
	luaVm.L.Close()
}
