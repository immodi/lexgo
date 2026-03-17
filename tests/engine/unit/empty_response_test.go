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

func TestEngine_EmptyResponse(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, res.StatusCode)
	}
}
