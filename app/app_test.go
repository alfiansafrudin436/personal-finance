package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// newTestEcho builds the real route table. The repositories are constructed
// with a nil database handle, which is fine because these tests only check
// which requests reach a handler at all, never what a handler reads.
func newTestEcho(t *testing.T) *echo.Echo {
	t.Helper()

	config.Application = &config.App{
		Version: "test",
		JWT:     config.JWTConfig{SecretKey: "test-secret"},
	}

	e := echo.New()
	registerRoutes(e)
	return e
}

func do(t *testing.T, e *echo.Echo, method, target, token string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// Health is public, and is also the canary for the catch-all route that Echo
// registers for a group with middleware: if that ever shadowed real routes,
// this would come back 401.
func TestHealthIsPublic(t *testing.T) {
	e := newTestEcho(t)

	rec := do(t, e, http.MethodGet, "/api/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/health = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body)
	}
}

// The auth endpoints must stay reachable without a token, or nobody could ever
// obtain one. They are registered on a sibling group, so the JWT middleware on
// the protected group must not apply to them.
func TestAuthRoutesDoNotRequireToken(t *testing.T) {
	e := newTestEcho(t)

	for _, path := range []string{"/api/auth/login", "/api/auth/register", "/api/auth/forgot-password"} {
		rec := do(t, e, http.MethodPost, path, "")
		if rec.Code == http.StatusUnauthorized {
			t.Errorf("POST %s = 401, want the handler to be reached without a token", path)
		}
	}
}

func TestProtectedRoutesRequireToken(t *testing.T) {
	e := newTestEcho(t)

	paths := []string{
		"/api/users/me",
		"/api/accounts",
		"/api/categories",
		"/api/transactions",
		"/api/budgets",
		"/api/reports/dashboard",
	}

	for _, path := range paths {
		rec := do(t, e, http.MethodGet, path, "")
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("GET %s without a token = %d, want %d", path, rec.Code, http.StatusUnauthorized)
		}
	}
}

func TestProtectedRoutesRejectInvalidToken(t *testing.T) {
	e := newTestEcho(t)

	rec := do(t, e, http.MethodGet, "/api/accounts", "not-a-jwt")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET /api/accounts with a malformed token = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// A token signed with a different secret must not pass, which is what stops a
// caller from minting their own.
func TestProtectedRoutesRejectForeignSignature(t *testing.T) {
	e := newTestEcho(t)

	config.Application.JWT.SecretKey = "a-different-secret"
	token, err := utils.ParseToken(utils.TokenParams{ID: uuid.NewString()})
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	config.Application.JWT.SecretKey = "test-secret"

	rec := do(t, e, http.MethodGet, "/api/accounts", token)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET /api/accounts with a foreign-signed token = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// A token whose subject is not a UUID cannot identify a user, so it must be
// refused rather than reaching a handler that would scope a query by it.
func TestProtectedRoutesRejectNonUUIDSubject(t *testing.T) {
	e := newTestEcho(t)

	token, err := utils.ParseToken(utils.TokenParams{ID: "not-a-uuid"})
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	rec := do(t, e, http.MethodGet, "/api/accounts", token)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET /api/accounts with a non-UUID subject = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
