package router

import (
	"encoding/json"
	"immodi/lexgo/internal/vm"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	OPTIONS HTTPMethod = "OPTIONS"
)

type LuaRequest struct {
	HttpRequest *http.Request
	LuaVm       *vm.LuaVm
}

func (req *LuaRequest) MakeLuaRequest() *lua.LTable {
	var luaReq *lua.LTable

	req.LuaVm.WithLock(func(L *lua.LState) error {
		luaReq = L.NewTable()
		req.setBasicFields(L, luaReq)
		req.setBodyField(L, luaReq)
		return nil
	})

	return luaReq
}

func (req *LuaRequest) setBasicFields(L *lua.LState, luaReq *lua.LTable) {
	L.SetField(luaReq, "method", lua.LString(req.HttpRequest.Method))
	L.SetField(luaReq, "url", lua.LString(req.HttpRequest.URL.Path))

	req.setQueryParameters(L, luaReq)
}

func (req *LuaRequest) setBodyField(L *lua.LState, luaReq *lua.LTable) {
	if !req.shouldParseBody() {
		return
	}

	contentType := req.HttpRequest.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		req.parseJSONBody(L, luaReq)
	case "application/x-www-form-urlencoded":
		req.parseFormBody(L, luaReq)
	}
}

func (req *LuaRequest) shouldParseBody() bool {
	method := req.HttpRequest.Method
	return method == string(POST) || method == string(PUT) || method == string(PATCH)
}

func (req *LuaRequest) parseJSONBody(L *lua.LState, luaReq *lua.LTable) {
	var bodyData map[string]any
	err := json.NewDecoder(req.HttpRequest.Body).Decode(&bodyData)
	if err == nil {
		L.SetField(luaReq, "body", mapToLuaTable(L, bodyData))
	}
}

func (req *LuaRequest) parseFormBody(L *lua.LState, luaReq *lua.LTable) {
	err := req.HttpRequest.ParseForm()
	if err != nil {
		return
	}

	formTable := L.NewTable()
	for key, values := range req.HttpRequest.PostForm {
		req.setFormField(L, formTable, key, values)
	}
	L.SetField(luaReq, "body", formTable)
}

func (req *LuaRequest) setFormField(L *lua.LState, formTable *lua.LTable, key string, values []string) {
	if len(values) == 1 {
		L.SetField(formTable, key, lua.LString(values[0]))
	} else {
		arr := L.NewTable()
		for _, v := range values {
			arr.Append(lua.LString(v))
		}
		L.SetField(formTable, key, arr)
	}
}

func (req *LuaRequest) setQueryParameters(L *lua.LState, luaReq *lua.LTable) {
	queryTbl := L.NewTable()

	for key, values := range req.HttpRequest.URL.Query() {
		valTbl := L.NewTable()
		for _, v := range values {
			valTbl.Append(lua.LString(v))
		}
		L.SetField(queryTbl, key, valTbl)
	}

	L.SetField(luaReq, "query", queryTbl)
}

// TODO: support file uploads
