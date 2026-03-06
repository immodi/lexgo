package router_unit_test

import (
	"net/http"
	"testing"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
)

func TestRouter_MethodIsolation(t *testing.T) {
	r, rd := router.New()

	r.AppendRoute(&router.Handler{
		Pattern: "/hello",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {}),
	})

	route := &framework.HTTPRoute{
		Path:   "/hello",
		Method: framework.POST,
	}

	h, _, _ := r.Match(route)

	if h.Method == framework.GET {
		t.Fatal("GET handler should not match POST")
	}
}
