package unit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
)

func TestRouterBasic(t *testing.T) {
	luaVm := vm.MakeLuaVm()
	r, _ := router.MakeRouter(luaVm)

	r.ConstructRouterNode(&router.Handler{
		Pattern: "/hello",
		Handler: nil,
		HijackHandler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello, world"))
		},
		Params: map[string]string{},
		Method: framework.GET,
		Mws:    nil,
	})

	r.ConstructRouterNode(&router.Handler{
		Pattern: "/post-hello",
		Handler: nil,
		HijackHandler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("post"))
		},
		Params: map[string]string{},
		Method: framework.POST,
		Mws:    nil,
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()

	s_req := httptest.NewRequest("POST", "/post-hello", nil)
	s_rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	r.ServeHTTP(s_rec, s_req)

	if rec.Code != 200 {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if s_rec.Body.Len() != 4 {
		t.Fatalf("expected body length, got %d", s_rec.Body.Len())
	}
}
