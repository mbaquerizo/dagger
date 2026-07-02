package issues

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mbaquerizo/dagger/internal/auth"
)

func NewGetIssueHandler(pool poolIface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceID, ok := auth.WorkspaceIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		displayID := chi.URLParam(r, "displayId")

		if displayID == "" {
			http.Error(w, "Missing display ID", http.StatusBadRequest)
			return
		}

		projectID, ok := auth.ProjectIDFromContext(r.Context())

		var projectIDPtr *int

		if ok {
			projectIDPtr = &projectID
		}

		issueContext, err := GetIssueContext(r.Context(), pool, displayID, workspaceID, projectIDPtr)

		if err != nil {
			if errors.Is(err, ErrIssueNotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			if errors.Is(err, ErrProjectIDMismatch) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		markdown := RenderIssueContext(issueContext)

		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(markdown))
	}
}

func NewListIssuesHandler(pool poolIface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceID, ok := auth.WorkspaceIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		projectID, ok := auth.ProjectIDFromContext(r.Context())

		var projectIDPtr *int

		if ok {
			projectIDPtr = &projectID
		}

		status := r.URL.Query().Get("status")

		if status == "" {
			status = "open"
		}

		issues, err := ListIssues(r.Context(), pool, status, workspaceID, projectIDPtr)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(issues)
	}
}

func NewUpdateIssueStatusHandler(pool poolIface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceID, ok := auth.WorkspaceIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		projectID, ok := auth.ProjectIDFromContext(r.Context())
		var projectIDPtr *int

		if ok {
			projectIDPtr = &projectID
		}

		var req UpdateStatusRequest

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		displayID := chi.URLParam(r, "displayId")

		if displayID == "" {
			http.Error(w, "Missing display ID", http.StatusBadRequest)
			return
		}

		err = ValidateStatus(req.Status)

		if err != nil {
			http.Error(w, "Invalid status value", http.StatusUnprocessableEntity)
			return
		}

		err = UpdateIssueStatus(r.Context(), pool, req, displayID, workspaceID, projectIDPtr)

		if err != nil {
			if errors.Is(err, ErrIssueNotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
