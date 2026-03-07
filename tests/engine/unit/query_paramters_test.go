package engine_unit_test

import (
	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_QueryParams(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			q := req.URL.Query().Get("name")
			w.Write([]byte(q))
		}),
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/?name=lexgo", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	if string(data) != "lexgo" {
		t.Fatalf("expected lexgo got %s", data)
	}
}
