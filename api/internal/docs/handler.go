package docs

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mbaquerizo/dagger/internal/auth"
)

func NewGetDocHandler(pool poolIface) http.HandlerFunc {
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

		doc, err := GetDoc(r.Context(), pool, displayID, workspaceID, projectIDPtr)

		if err != nil {
			if errors.Is(err, ErrDocNotFound) {
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

		markdown := RenderDoc(doc)

		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(markdown))
	}
}

func NewUpdateDocHandler(pool poolIface) http.HandlerFunc {
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

		var req UpdateDocRequest

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := ValidateUpdateDocRequest(req); err != nil {
			if errors.Is(err, ErrNoFieldsToUpdate) {
				http.Error(w, "At least one of title, body, or status must be provided", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		displayID := chi.URLParam(r, "displayId")

		if displayID == "" {
			http.Error(w, "Missing display ID", http.StatusBadRequest)
			return
		}

		err = UpdateDoc(r.Context(), pool, req, displayID, workspaceID, projectIDPtr)

		if err != nil {
			if errors.Is(err, ErrDocNotFound) {
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
	}
}
