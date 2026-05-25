// Package handler provides HTTP handlers for PVC file operations.
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// EditResponse is the JSON response for an edit operation.
type EditResponse struct {
	Written string `json:"written"`
}

// EditHandler returns an HTTP handler that writes (creates or overwrites)
// a file in the PVC.
func EditHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		p := r.URL.Query().Get("path")
		abs, err := safeJoin(root, p)
		if err != nil {
			http.Error(w, `{"error":"bad path"}`, http.StatusBadRequest)
			return
		}
		dir := filepath.Dir(abs)
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		tmp, err := os.CreateTemp(dir, ".edit-*")
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		defer func() { _ = tmp.Close(); _ = os.Remove(tmp.Name()) }()
		if _, err = io.Copy(tmp, r.Body); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		if err = os.Rename(tmp.Name(), abs); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(EditResponse{Written: p})
	})
}
