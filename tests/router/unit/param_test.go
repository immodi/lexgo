package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_ParamRoute(t *testing.T) {
	r, rd := router.MakeRouter()

	r.AppendRoute(&router.Handler{
		Pattern: "/user/:id",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	route := &framework.HTTPRoute{
		Path:   "/user/42",
		Method: framework.GET,
	}

	h, _, _ := r.Match(route)

	if h.Params["id"] != "42" {
		t.Fatalf("expected id=42, got %s", h.Params["id"])
	}
}

