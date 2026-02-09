package vm

import (
	"encoding/json"
	"fmt"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type LuaResponse struct {
	HttpWriter http.ResponseWriter
	LuaVm      LVm
	Written    bool
	statusCode int
}

func (res *LuaResponse) MakeLuaResponse() *LuaTable {
	luaRes := res.LuaVm.NewTable()

	// Create LuaFunction wrappers for each handler
	statusFn := &LuaFunction{LFunction: res.LuaVm.(*LuaVm).L.NewFunction(res.handleStatus)}
	setHeaderFn := &LuaFunction{LFunction: res.LuaVm.(*LuaVm).L.NewFunction(res.handleSetHeader)}
	htmlFn := &LuaFunction{LFunction: res.LuaVm.(*LuaVm).L.NewFunction(res.write(res.handleHTML))}
	jsonFn := &LuaFunction{LFunction: res.LuaVm.(*LuaVm).L.NewFunction(res.write(res.handleJSON))}
	rawFn := &LuaFunction{LFunction: res.LuaVm.(*LuaVm).L.NewFunction(res.write(res.handleRaw))}

	luaRes.SetField("status", statusFn)
	luaRes.SetField("setHeader", setHeaderFn)
	luaRes.SetField("html", htmlFn)
	luaRes.SetField("json", jsonFn)
	luaRes.SetField("raw", rawFn)

	return &luaRes
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
