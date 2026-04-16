package team

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebUIProxyHandlerForwardsOnboardingRoutes(t *testing.T) {
	var gotPath string
	var gotAuth string
	var gotQuery string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer upstream.Close()

	b := NewBroker()
	req := httptest.NewRequest(http.MethodGet, "/onboarding/state?step=providers", nil)
	rec := httptest.NewRecorder()

	b.webUIProxyHandler(upstream.URL, "").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if gotPath != "/onboarding/state" {
		t.Fatalf("expected proxied onboarding path, got %q", gotPath)
	}
	if gotQuery != "step=providers" {
		t.Fatalf("expected query to be forwarded, got %q", gotQuery)
	}
	if gotAuth != "Bearer "+b.Token() {
		t.Fatalf("expected broker auth header, got %q", gotAuth)
	}
	if body := strings.TrimSpace(rec.Body.String()); body != `{"ok":true}` {
		t.Fatalf("unexpected proxied body %q", body)
	}
}
