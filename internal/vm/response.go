package vm

import (
	"encoding/json"
	"fmt"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type LuaResponse struct {
	HttpWriter http.ResponseWriter
	LuaVm      *LuaVm
	Written    bool
	statusCode int
}

func (res *LuaResponse) MakeLuaResponse() *lua.LTable {
	L := res.LuaVm.L
	luaRes := L.NewTable()

	L.SetField(luaRes, "status", L.NewFunction(res.handleStatus))
	L.SetField(luaRes, "setHeader", L.NewFunction(res.handleSetHeader))

	L.SetField(luaRes, "html", L.NewFunction(res.write(res.handleHTML)))
	L.SetField(luaRes, "json", L.NewFunction(res.write(res.handleJSON)))
	L.SetField(luaRes, "raw", L.NewFunction(res.write(res.handleRaw)))

	return luaRes
}

func (res *LuaResponse) write(fn func(L *lua.LState) int) func(L *lua.LState) int {
	return func(L *lua.LState) int {

		if res.Written {
			L.RaiseError("response already sent")
			return 0
		}

		res.Written = true
		return fn(L)
	}
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
	res.ensureStatus()

	res.HttpWriter.Header().Set("Content-Type", "text/html")
	res.HttpWriter.WriteHeader(res.statusCode)
	fmt.Fprint(res.HttpWriter, msg)
	return 0
}

func (res *LuaResponse) handleRaw(L *lua.LState) int {
	msg := L.CheckString(1)
	res.ensureStatus()

	res.HttpWriter.WriteHeader(res.statusCode)
	fmt.Fprint(res.HttpWriter, msg)
	return 0
}

func (res *LuaResponse) handleJSON(L *lua.LState) int {
	tbl := L.CheckTable(1)
	goMap := luaTableToMap(tbl)

	data, err := json.Marshal(goMap)
	if err != nil {
		L.RaiseError("json marshal failed: %v", err)
		return 0
	}

	res.ensureStatus()
	res.HttpWriter.Header().Set("Content-Type", "application/json")
	res.HttpWriter.WriteHeader(res.statusCode)
	res.HttpWriter.Write(data)
	return 0
}

func (res *LuaResponse) handleSetHeader(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckString(2)
	res.HttpWriter.Header().Set(key, val)
	return 0
}

func (res *LuaResponse) ensureStatus() {
	if res.statusCode == 0 {
		res.statusCode = http.StatusOK
	}
}
