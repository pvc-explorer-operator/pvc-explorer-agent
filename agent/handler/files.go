// Package handler provides HTTP handlers for PVC file operations.
package handler

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileEntry represents a single file or directory entry in a listing.
type FileEntry struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

// ListResponse is the JSON response for a directory listing.
type ListResponse struct {
	Path    string      `json:"path"`
	Entries []FileEntry `json:"entries"`
}

// DeleteResponse is the JSON response for a delete operation.
type DeleteResponse struct {
	Deleted string `json:"deleted"`
}

// ClearResponse is the JSON response for a clear operation.
type ClearResponse struct {
	Cleared bool `json:"cleared"`
}

// FilesHandler returns an HTTP handler that lists files and handles
// GET (list) and DELETE (remove) requests.
func FilesHandler(root string, isReadonly func(*http.Request) bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		p := r.URL.Query().Get("path")
		abs, err := safeJoin(root, p)
		if err != nil {
			http.Error(w, `{"error":"bad path"}`, http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			entries, err := os.ReadDir(abs)
			if err != nil {
				if os.IsNotExist(err) {
					http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				} else {
					http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
				}
				return
			}
			var resp ListResponse
			resp.Path = p
			for _, entry := range entries {
				info, err := entry.Info()
				if err != nil {
					// skip entries we cannot stat
					continue
				}
				resp.Entries = append(resp.Entries, FileEntry{
					Name:    entry.Name(),
					Size:    info.Size(),
					IsDir:   entry.IsDir(),
					ModTime: info.ModTime().UTC().Format("2006-01-02T15:04:05Z07:00"),
				})
			}
			_ = json.NewEncoder(w).Encode(resp)
		case http.MethodDelete:
			if isReadonly(r) {
				http.Error(w, `{"error":"read-only mode"}`, http.StatusForbidden)
				return
			}
			err := os.RemoveAll(abs)
			if err != nil {
				if os.IsNotExist(err) {
					http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				} else {
					log.Printf("delete %q: %v", abs, err)
					http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
				}
				return
			}
			_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: p})
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
}

// ClearHandler returns an HTTP handler that removes all files under root.
func ClearHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		for _, entry := range entries {
			// best-effort removal; ignore errors per-file
			_ = os.RemoveAll(filepath.Join(root, entry.Name()))
		}
		_ = json.NewEncoder(w).Encode(ClearResponse{Cleared: true})
	})
}

func safeJoin(root, rel string) (string, error) {
	clean := path.Clean("/" + rel)
	abs := filepath.Join(root, clean)
	if !strings.HasPrefix(abs, root) {
		return "", fs.ErrInvalid
	}
	return abs, nil
}

func recover500(w http.ResponseWriter) {
	if r := recover(); r != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
	}
}
