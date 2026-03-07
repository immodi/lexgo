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

func TestEngine_Middleware_ModifiesRequest(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	mw := rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
		req.Header.Set("X-Middleware", "yes")
	})

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(req.Header.Get("X-Middleware")))
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

	if string(data) != "yes" {
		t.Fatalf("expected header modification by middleware, got %q", data)
	}
}
