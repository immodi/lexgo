package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_StaticRoute(t *testing.T) {
	r, rd := router.MakeRouter()

	r.AppendRoute(&router.Handler{
		Pattern: "/hello",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	route := &framework.HTTPRoute{
		Path:   "/hello",
		Method: framework.GET,
	}

	h, _, _ := r.Match(route)

	if h.Pattern != "/hello" {
		t.Fatalf("expected /hello, got %s", h.Pattern)
	}
}
