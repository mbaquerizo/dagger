package mcp

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewMCPHandler(pool poolIface, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				JSONRPC: "2.0",
				Error:   &Error{Code: ErrCodeParse, Message: "failed to read request body"},
			})

			return
		}

		var req Request

		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				JSONRPC: "2.0",
				Error:   &Error{Code: ErrCodeParse, Message: "Invalid JSON"},
			})
			return
		}

		server := NewServer(NewDBService(pool, baseURL))

		resp := server.HandleRequest(r.Context(), req)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
