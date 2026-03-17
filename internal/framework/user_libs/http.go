package userlibs

import (
	"bytes"
	"encoding/json"
	"immodi/lexgo/internal/vm"
	"io"
	"net/http"
)

const MaxBodySize = 5 << 20 // 5MB

type HTTP struct {
	LVm vm.LVm
}

func MakeHTTP(lvm vm.LVm) (vm.LuaValue, error) {
	tbl := lvm.NewTable()
	httpObj := HTTP{LVm: lvm}
	tbl.SetField("get", lvm.NewMultiReturnFunction(func(l vm.LVm) int {
		url, err := l.CheckString(1)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		res, err := httpObj.get(url)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		return l.Return(res, vm.LuaNil{})
	}))

	tbl.SetField("post", lvm.NewMultiReturnFunction(func(l vm.LVm) int {
		url, err := l.CheckString(1)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		contentType, err := l.CheckString(2)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		body, err := l.CheckTable(3)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		bodyMap := make(map[string]any)
		var luaValueToGo func(v vm.LuaValue) any
		luaValueToGo = func(v vm.LuaValue) any {
			switch val := v.(type) {
			case vm.LuaNil:
				return nil
			case vm.LuaBool:
				return bool(val)
			case vm.LuaNumber:
				return float64(val)
			case vm.LuaString:
				return string(val)
			case *vm.LuaTable:
				nested := make(map[string]any)
				val.ForEach(func(k, v vm.LuaValue) {
					nested[k.String()] = luaValueToGo(v)
				})
				return nested
			default:
				// fallback to string representation
				return val.String()
			}
		}

		body.ForEach(func(key, value vm.LuaValue) {
			bodyMap[key.String()] = luaValueToGo(value)
		})

		res, err := httpObj.post(url, contentType, bodyMap)
		if err != nil {
			return l.Return(vm.LuaNil{}, vm.LuaString(err.Error()))
		}

		return l.Return(res, vm.LuaNil{})
	}))

	return tbl, nil
}

func (h *HTTP) get(url string) (*vm.LuaTable, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resTbl, err := h.packageResponseTbl(res)
	if err != nil {
		return nil, err
	}

	return resTbl, nil
}

func (h *HTTP) post(url, contentType string, body map[string]any) (*vm.LuaTable, error) {
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	responseBody := bytes.NewBuffer(postBody)

	res, err := http.Post(url, contentType, responseBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resTbl, err := h.packageResponseTbl(res)
	if err != nil {
		return nil, err
	}

	return resTbl, nil
}

func (h *HTTP) packageResponseTbl(res *http.Response) (*vm.LuaTable, error) {
	resp := h.LVm.NewTable()

	statusTbl := vm.LuaNumber(res.StatusCode)

	body, err := io.ReadAll(io.LimitReader(res.Body, MaxBodySize))
	if err != nil {
		return nil, err
	}

	bodyTbl := vm.LuaString(string(body))
	headersTbl := h.LVm.NewTable()
	for key, values := range res.Header {
		vs := h.LVm.NewTable()
		for i, v := range values {
			vs.Set(vm.LuaNumber(i+1), vm.LuaString(v))
		}

		headersTbl.SetField(key, vs)
	}

	resp.SetField("status", statusTbl)
	resp.SetField("body", bodyTbl)
	resp.SetField("headers", headersTbl)

	return resp, nil
}
