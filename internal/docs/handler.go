package docs

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mbaquerizo/dagger/internal/auth"
)

func NewGetDocHandler(pool poolIface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		displayID := chi.URLParam(r, "displayId")

		if displayID == "" {
			http.Error(w, "Missing display ID", http.StatusBadRequest)
			return
		}

		workspaceID, ok := auth.WorkspaceIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		doc, err := GetDoc(r.Context(), pool, displayID, workspaceID)

		if err != nil {
			if errors.Is(err, ErrDocNotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
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
