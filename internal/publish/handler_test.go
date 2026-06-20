package publish

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/pashagolub/pgxmock/v5"
)

func TestHandler_PublishDoc(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`INSERT INTO docs`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
	}

	body, err := json.Marshal(req)

	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/v1/publish", bytes.NewReader(body))
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := NewHandler(mockPool, "http://localhost:8080")
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code 201, but got %d", w.Code)
	}

	var resp PublishResponse

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected DisplayID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestHandler_ValidationError(t *testing.T) {
	req := PublishRequest{
		Type:      "adr",
		Title:     "",
		Body:      "Test Body",
		ProjectID: 1,
	}

	body, err := json.Marshal(req)

	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/v1/publish", bytes.NewReader(body))
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := NewHandler(nil, "http://localhost:8080")
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, but got %d", w.Code)
	}

	var errs []ValidationError
	err = json.Unmarshal(w.Body.Bytes(), &errs)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(errs) == 0 {
		t.Error("expected at least one validation error")
	}
}
