package engine_unit_test

import (
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_MethodNotAllowed(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("ok"))
		}),
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d got %d", http.StatusMethodNotAllowed, res.StatusCode)
	}
}
