package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_NotFound(t *testing.T) {
	r, rd := router.MakeRouter()

	r.NotFoundFunc = rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {})

	route := &framework.HTTPRoute{
		Path:   "/missing",
		Method: framework.GET,
	}

	h, _, _ := r.Match(route)

	if h.Pattern != "/missing" {
		t.Fatal("expected fallback handler")
	}
}
