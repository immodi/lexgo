package router

import (
	"encoding/json"
	"fmt"
	"immodi/lexgo/internal/vm"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type LuaResponse struct {
	HttpWriter http.ResponseWriter
	LuaVm      *vm.LuaVm
	statusCode int
}

func (res *LuaResponse) MakeLuaResponse() *lua.LTable {
	var luaRes *lua.LTable

	res.LuaVm.WithLock(func(L *lua.LState) error {
		luaRes = L.NewTable()

		L.SetField(luaRes, "status", L.NewFunction(res.handleStatus))
		L.SetField(luaRes, "html", L.NewFunction(res.handleHTML))
		L.SetField(luaRes, "raw", L.NewFunction(res.handleRaw))
		L.SetField(luaRes, "json", L.NewFunction(res.handleJSON))
		L.SetField(luaRes, "setHeader", L.NewFunction(res.handleSetHeader))

		return nil
	})

	return luaRes
}

func (res *LuaResponse) handleStatus(L *lua.LState) int {
	code := L.CheckInt(1)
	if code == 0 {
		code = http.StatusOK
	}
	res.statusCode = code
	return 0
}

func (res *LuaResponse) handleHTML(L *lua.LState) int {
	msg := L.CheckString(1)
	res.ensureStatusCode()

	res.HttpWriter.Header().Set("Content-Type", "text/html")
	res.HttpWriter.WriteHeader(res.statusCode)
	fmt.Fprint(res.HttpWriter, msg)

	return 0
}

func (res *LuaResponse) handleRaw(L *lua.LState) int {
	msg := L.CheckString(1)
	res.ensureStatusCode()

	res.HttpWriter.WriteHeader(res.statusCode)
	fmt.Fprint(res.HttpWriter, msg)

	return 0
}

func (res *LuaResponse) handleJSON(L *lua.LState) int {
	tbl := L.CheckTable(1)
	goMap := luaTableToMap(tbl)

	data, err := json.Marshal(goMap)
	if err != nil {
		L.RaiseError("failed to marshal JSON: %v", err)
		return 0
	}

	res.ensureStatusCode()
	res.HttpWriter.Header().Set("Content-Type", "application/json")
	res.HttpWriter.WriteHeader(res.statusCode)
	fmt.Fprint(res.HttpWriter, string(data))

	return 0
}

func (res *LuaResponse) handleSetHeader(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckString(2)
	res.HttpWriter.Header().Set(key, val)
	return 0
}

func (res *LuaResponse) ensureStatusCode() {
	if res.statusCode == 0 {
		res.statusCode = http.StatusOK
	}
}
