package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
)

func newFixture(t *testing.T) (pgxmock.PgxPoolIface, func(http.Handler) http.Handler) {
	t.Helper()

	pool, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))

	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { pool.Close() })

	return pool, NewMiddleware(pool)
}

func TestMiddleware_HealthCheckBypass(t *testing.T) {
	_, middleware := newFixture(t)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestMiddleware_MissingAuthHeader(t *testing.T) {
	_, middleware := newFixture(t)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when auth is missing")
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_WrongPrefix(t *testing.T) {
	_, middleware := newFixture(t)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when auth is missing")
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	req.Header.Set("Authorization", "Bearer abc_123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_ValidKey(t *testing.T) {
	pool, middleware := newFixture(t)

	pool.ExpectQuery("").
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "workspace_id", "project_id"}).
			AddRow(1, 5, nil))

	var capturedWorkspaceID int

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedWorkspaceID, _ = WorkspaceIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	req.Header.Set("Authorization", "Bearer dgr_iliketurtles")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if capturedWorkspaceID != 5 {
		t.Errorf("expected workspace ID 5, got %d", capturedWorkspaceID)
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestMiddleware_ValidKeyWithProjectID(t *testing.T) {
	pool, middleware := newFixture(t)

	projectID := 33

	pool.ExpectQuery("").
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "workspace_id", "project_id"}).
			AddRow(1, 5, &projectID))

	var capturedProjectID int

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedProjectID, _ = ProjectIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	req.Header.Set("Authorization", "Bearer dgr_iliketurtles")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if capturedProjectID != 33 {
		t.Errorf("expected project ID 33, got %d", capturedProjectID)
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestMiddleware_InvalidKey(t *testing.T) {
	pool, middleware := newFixture(t)

	pool.ExpectQuery("").
		WithArgs(pgxmock.AnyArg()).
		WillReturnError(pgx.ErrNoRows)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for invalid key")
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	req.Header.Set("Authorization", "Bearer dgr_iliketurtles")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestMiddleware_DBError(t *testing.T) {
	pool, middleware := newFixture(t)

	pool.ExpectQuery("").
		WithArgs(pgxmock.AnyArg()).
		WillReturnError(errors.New("failed to make query"))

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for db error")
	}))

	req := httptest.NewRequest("GET", "/api/v1/someEndpoint", nil)
	req.Header.Set("Authorization", "Bearer dgr_iliketurtles")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}

	if err := pool.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
