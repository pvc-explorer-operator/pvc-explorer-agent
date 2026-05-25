// Package handler provides HTTP handlers for PVC file operations.
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// UploadResponse is the JSON response for a file upload.
type UploadResponse struct {
	Uploaded string `json:"uploaded"`
}

// UploadHandler returns an HTTP handler that accepts file uploads for the PVC.
func UploadHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		p := r.URL.Query().Get("path")
		abs, err := safeJoin(root, p)
		if err != nil {
			http.Error(w, `{"error":"bad path"}`, http.StatusBadRequest)
			return
		}
		err = os.MkdirAll(abs, 0750)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, `{"error":"bad upload"}`, http.StatusBadRequest)
			return
		}
		defer func() { _ = file.Close() }()
		outPath := filepath.Join(abs, header.Filename)
		if !strings.HasPrefix(outPath, root) {
			http.Error(w, `{"error":"bad path"}`, http.StatusBadRequest)
			return
		}
		//nolint:gosec // path validated via strings.HasPrefix(outPath, root)
		out, err := os.Create(outPath)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		defer func() { _ = out.Close() }()
		if _, err = io.Copy(out, file); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(UploadResponse{Uploaded: path.Join(p, header.Filename)})
	})
}
