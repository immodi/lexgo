package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_MiddlewareMerge(t *testing.T) {
	r, rd := router.MakeRouter()

	globalMw := rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {})
	routeMw := rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {})

	r.Mws = []framework.RouterServerHandler{globalMw}

	r.AppendRoute(&router.Handler{
		Pattern: "/hello",
		Method:  framework.GET,
		Mws:     []framework.RouterServerHandler{routeMw},
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	route := &framework.HTTPRoute{
		Path:   "/hello",
		Method: framework.GET,
	}

	h, _, _ := r.Match(route)

	if len(h.Mws) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(h.Mws))
	}
}
