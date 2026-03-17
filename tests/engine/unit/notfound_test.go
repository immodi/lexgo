package engine_unit_test

import (
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_NotFound(t *testing.T) {
	r, _ := router.New()
	luaVm := vm.New()

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d got %d", http.StatusNotFound, res.StatusCode)
	}
}
