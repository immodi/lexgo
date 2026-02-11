package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"immodi/lexgo/internal/vm"
	"net/http"
)

type LuaResponse struct {
	HttpWriter http.ResponseWriter
	LuaVm      vm.LVm
	buf        bytes.Buffer
	headerBuf  http.Header
	statusCode int
}

func ConstructResponse(w http.ResponseWriter, vm vm.LVm) *LuaResponse {
	return &LuaResponse{
		HttpWriter: w,
		LuaVm:      vm,
		statusCode: 0,
		buf:        *bytes.NewBuffer([]byte{}),
		headerBuf:  make(http.Header),
	}
}

func (res *LuaResponse) MakeLuaResponse() *vm.LuaTable {
	luaRes := res.LuaVm.NewTable()

	statusFn := res.LuaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		res.handleStatus()
		return nil
	})
	setHeaderFn := res.LuaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		res.handleSetHeader()
		return nil
	})
	htmlFn := res.LuaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		if res.ensureNotWritten() != nil {
			return nil
		}
		res.handleHTML()
		return nil
	})
	jsonFn := res.LuaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		if res.ensureNotWritten() != nil {
			return nil
		}
		res.handleJSON()
		return nil
	})
	rawFn := res.LuaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		if res.ensureNotWritten() != nil {
			return nil
		}
		res.handleRaw()
		return nil
	})

	luaRes.SetField("status", statusFn)
	luaRes.SetField("setHeader", setHeaderFn)
	luaRes.SetField("html", htmlFn)
	luaRes.SetField("json", jsonFn)
	luaRes.SetField("raw", rawFn)

	return luaRes
}

func (res *LuaResponse) ensureNotWritten() error {
	contentHeader := res.headerBuf.Get("Content-Type")
	if res.buf.Len() > 0 || contentHeader != "" {
		errMsg := "response already sent - cannot call multiple response methods"
		res.LuaVm.Error(errMsg)
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (res *LuaResponse) Flush() {
	res.ensureStatus()

	for k, vals := range res.headerBuf {
		for _, v := range vals {
			res.HttpWriter.Header().Add(k, v)
		}
	}

	res.HttpWriter.WriteHeader(res.statusCode)
	res.HttpWriter.Write(res.buf.Bytes())
}

func (res *LuaResponse) Reset() {
	res.buf.Reset()
	res.headerBuf = make(http.Header)
	res.statusCode = 0
}

func (res *LuaResponse) handleStatus() {
	code, err := res.LuaVm.CheckNumber(1)
	if err != nil {
		code = float64(http.StatusOK)
	}
	if code == 0 {
		code = float64(http.StatusOK)
	}
	res.statusCode = int(code)
}

func (res *LuaResponse) handleHTML() {
	msg, err := res.LuaVm.CheckString(1)
	if err != nil {
		res.LuaVm.Error(err.Error())
		return
	}

	res.ensureStatus()
	res.headerBuf.Set("Content-Type", "text/html")
	res.buf.WriteString(msg)
}

func (res *LuaResponse) handleRaw() {
	msg, err := res.LuaVm.CheckString(1)
	if err != nil {
		res.LuaVm.Error(err.Error())
		return
	}

	res.ensureStatus()
	res.buf.WriteString(msg)
}

func (res *LuaResponse) handleJSON() {
	tbl, err := res.LuaVm.CheckTable(1)
	if err != nil {
		res.LuaVm.Error(err.Error())
		return
	}

	goMap := luaTableToMap(tbl)
	data, err := json.Marshal(goMap)
	if err != nil {
		res.LuaVm.Error(fmt.Sprintf("json marshal failed: %v", err))
		return
	}

	res.ensureStatus()
	res.headerBuf.Set("Content-Type", "application/json")
	res.buf.Write(data)
}

func (res *LuaResponse) handleSetHeader() {
	key, err := res.LuaVm.CheckString(1)
	if err != nil {
		res.LuaVm.Error(err.Error())
		return
	}

	val, err := res.LuaVm.CheckString(2)
	if err != nil {
		res.LuaVm.Error(err.Error())
		return
	}

	res.headerBuf.Set(key, val)
}

func (res *LuaResponse) ensureStatus() {
	if res.statusCode == 0 {
		res.statusCode = http.StatusOK
	}
}
