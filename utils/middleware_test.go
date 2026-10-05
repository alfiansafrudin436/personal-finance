package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-finance/config"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func withTestConfig(t *testing.T) {
	t.Helper()
	config.Application = &config.App{
		JWT: config.JWTConfig{SecretKey: "test-secret"},
	}
}

// The positive path: a valid token reaches the handler and GetUserID returns
// the subject, which is what every scoped query relies on.
func TestJWTMiddlewarePassesUserIDToHandler(t *testing.T) {
	withTestConfig(t)

	userID := uuid.New()
	token, err := ParseToken(TokenParams{ID: userID.String()})
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	e := echo.New()
	var seen uuid.UUID
	var ok bool

	handler := JWTMiddleware()(func(c echo.Context) error {
		seen, ok = GetUserID(c)
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	if err := handler(e.NewContext(req, rec)); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !ok {
		t.Fatal("GetUserID reported no authenticated user")
	}
	if seen != userID {
		t.Errorf("GetUserID = %s, want %s", seen, userID)
	}
}

func TestJWTMiddlewareRejectsMissingAndMalformedHeaders(t *testing.T) {
	withTestConfig(t)

	cases := map[string]string{
		"no header":        "",
		"wrong scheme":     "Basic abc",
		"bearer with none": "Bearer ",
		"not a jwt":        "Bearer not-a-jwt",
	}

	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			reached := false

			handler := JWTMiddleware()(func(c echo.Context) error {
				reached = true
				return c.NoContent(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}
			rec := httptest.NewRecorder()

			if err := handler(e.NewContext(req, rec)); err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if reached {
				t.Error("handler was reached despite an unusable token")
			}
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

// Without the middleware there is no user on the context, so GetUserID has to
// report that rather than hand back a zero UUID that would look like an owner.
func TestGetUserIDWithoutMiddleware(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	if _, ok := GetUserID(c); ok {
		t.Error("GetUserID reported an authenticated user on a bare context")
	}
}
