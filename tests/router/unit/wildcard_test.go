package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_WildcardRoute(t *testing.T) {
	r, rd := router.MakeRouter()

	r.AppendRoute(&router.Handler{
		Pattern: "/files/*",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	route := &framework.HTTPRoute{
		Path:   "/files/a/b/c",
		Method: framework.GET,
	}

	h, _, _ := r.Match(route)

	if h.Params["*"] != "a/b/c" {
		t.Fatalf("expected wildcard a/b/c, got %s", h.Params["*"])
	}
}
