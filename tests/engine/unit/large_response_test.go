package engine_unit_test

import (
	"bytes"
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_LargeResponse(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	large := bytes.Repeat([]byte("a"), 4096)

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Write(large)
		}),
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	if !bytes.Equal(data, large) {
		t.Fatalf("large response mismatch")
	}
}
