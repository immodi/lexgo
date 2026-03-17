package engine_unit_test

import (
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_Middleware_ModifiesResponse(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	// Middleware adds a prefix
	mw := rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("mw-"))
	})

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("ok"))
		}),
		Mws: []framework.RouterHandler{mw},
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	expected := "mw-ok"
	if string(data) != expected {
		t.Fatalf("expected %q got %q", expected, data)
	}
}
