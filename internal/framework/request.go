package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"immodi/lexgo/internal/vm"
	"io"
	"net/http"
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
	LuaVm       vm.LVm
	Params      map[string]string
}

func ConstructRequest(req *http.Request, vm vm.LVm, params map[string]string) *LuaRequest {
	return &LuaRequest{HttpRequest: req, LuaVm: vm, Params: params}
}

func (req *LuaRequest) MakeLuaRequest() *vm.LuaTable {
	luaReq := req.LuaVm.NewTable()
	req.setBasicFields(luaReq)
	req.setBodyField(luaReq)
	return luaReq
}

func (req *LuaRequest) setBasicFields(luaReq *vm.LuaTable) {
	luaReq.SetField("method", vm.LuaString(req.HttpRequest.Method))
	luaReq.SetField("url", vm.LuaString(req.HttpRequest.URL.Path))
	luaReq.SetField("origin", vm.LuaString(req.HttpRequest.Header.Get("Origin")))

	req.setQueryParameters(luaReq)
	req.setRequestParameters(luaReq)
}

func (req *LuaRequest) setBodyField(luaReq *vm.LuaTable) {
	if !req.shouldParseBody() {
		return
	}
	contentType := req.HttpRequest.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		req.parseJSONBody(luaReq)
	case "application/x-www-form-urlencoded":
		req.parseFormBody(luaReq)
	}
}

func (req *LuaRequest) shouldParseBody() bool {
	method := req.HttpRequest.Method
	return method == string(POST) || method == string(PUT) || method == string(PATCH)
}

func (req *LuaRequest) parseJSONBody(luaReq *vm.LuaTable) {
	bodyData := map[string]interface{}{}
	data, err := io.ReadAll(req.HttpRequest.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}
	req.HttpRequest.Body = io.NopCloser(bytes.NewBuffer(data))
	if len(data) > 0 {
		err = json.Unmarshal(data, &bodyData)
		if err != nil {
			fmt.Println("Warning: failed to decode JSON body:", err)
		}
	}
	luaReq.SetField("body", mapToLuaTable(req.LuaVm, bodyData))
}

func (req *LuaRequest) parseFormBody(luaReq *vm.LuaTable) {
	err := req.HttpRequest.ParseForm()
	if err != nil {
		return
	}
	formTable := req.LuaVm.NewTable()
	for key, values := range req.HttpRequest.PostForm {
		req.setFormField(formTable, key, values)
	}
	luaReq.SetField("body", formTable)
}

func (req *LuaRequest) setFormField(formTable *vm.LuaTable, key string, values []string) {
	if len(values) == 1 {
		formTable.SetField(key, vm.LuaString(values[0]))
	} else {
		arr := req.LuaVm.NewTable()
		for _, v := range values {
			arr.Append(vm.LuaString(v))
		}
		formTable.SetField(key, arr)
	}
}

func (req *LuaRequest) setQueryParameters(luaReq *vm.LuaTable) {
	queryTbl := req.LuaVm.NewTable()
	for key, values := range req.HttpRequest.URL.Query() {
		valTbl := req.LuaVm.NewTable()
		for _, v := range values {
			valTbl.Append(vm.LuaString(v))
		}
		queryTbl.SetField(key, valTbl)
	}
	luaReq.SetField("query", queryTbl)
}

func (req *LuaRequest) setRequestParameters(luaReq *vm.LuaTable) {
	paramsTbl := req.LuaVm.NewTable()
	for key, value := range req.Params {
		paramsTbl.SetField(key, vm.LuaString(value))
	}
	luaReq.SetField("params", paramsTbl)
}
