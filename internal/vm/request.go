package vm

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	LuaVm       LVm
	Params      map[string]string
}

func (req *LuaRequest) MakeLuaRequest() *LuaTable {
	luaReq := req.LuaVm.NewTable()
	req.setBasicFields(&luaReq)
	req.setBodyField(&luaReq)
	return &luaReq
}

func (req *LuaRequest) setBasicFields(luaReq *LuaTable) {
	luaReq.SetField("method", LuaString(req.HttpRequest.Method))
	luaReq.SetField("url", LuaString(req.HttpRequest.URL.Path))
	req.setQueryParameters(luaReq)
	req.setRequestParameters(luaReq)
}

func (req *LuaRequest) setBodyField(luaReq *LuaTable) {
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

func (req *LuaRequest) parseJSONBody(luaReq *LuaTable) {
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

func (req *LuaRequest) parseFormBody(luaReq *LuaTable) {
	err := req.HttpRequest.ParseForm()
	if err != nil {
		return
	}
	formTable := req.LuaVm.NewTable()
	for key, values := range req.HttpRequest.PostForm {
		req.setFormField(&formTable, key, values)
	}
	luaReq.SetField("body", formTable)
}

func (req *LuaRequest) setFormField(formTable *LuaTable, key string, values []string) {
	if len(values) == 1 {
		formTable.SetField(key, LuaString(values[0]))
	} else {
		arr := req.LuaVm.NewTable()
		for _, v := range values {
			arr.Append(LuaString(v))
		}
		formTable.SetField(key, arr)
	}
}

func (req *LuaRequest) setQueryParameters(luaReq *LuaTable) {
	queryTbl := req.LuaVm.NewTable()
	for key, values := range req.HttpRequest.URL.Query() {
		valTbl := req.LuaVm.NewTable()
		for _, v := range values {
			valTbl.Append(LuaString(v))
		}
		queryTbl.SetField(key, valTbl)
	}
	luaReq.SetField("query", queryTbl)
}

func (req *LuaRequest) setRequestParameters(luaReq *LuaTable) {
	paramsTbl := req.LuaVm.NewTable()
	for key, value := range req.Params {
		paramsTbl.SetField(key, LuaString(value))
	}
	luaReq.SetField("params", paramsTbl)
}
