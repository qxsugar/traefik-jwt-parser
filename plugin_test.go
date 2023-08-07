package traefik_jwt_parser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCase(t *testing.T) {
	cfg := CreateConfig()
	cfg.TrustKeys = []string{"Uid", "CompanyId"}
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM3MzExODgsIlVpZCI6MTIzLCJDb21wYW55SWQiOjIyMn0.JXKTxWufbGWnPQmuAXnTb8iOY_TgTAHsaaL3BZuueDs"
	url := "https://localhost"

	handler, err := New(ctx, next, cfg, "traefik-jwt-parser")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test authorization", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Authorization", token)

		handler.ServeHTTP(writer, request)

		for key, value := range request.Header {
			t.Log(key, strings.Join(value, ""))
		}
	})

	t.Run("test query", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://localhost?Authorization="+token, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(writer, request)

		for key, value := range request.Header {
			t.Log(key, strings.Join(value, ""))
		}
	})
}
