package engine_unit_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"immodi/lexgo/internal/engine"
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
)

func TestEngine_Get(t *testing.T) {
	r, rd := router.New()
	luaVm := vm.New()

	mes := []byte("hello world!")

	r.AppendRoute(&router.Handler{
		Pattern: "/",
		Method:  framework.GET,
		Handler: rd.MakeGoHandler(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(mes)
		}),
	})

	engine := engine.New(r, luaVm)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, mes) {
		t.Fatalf("expected %q got %q", mes, data)
	}
}
