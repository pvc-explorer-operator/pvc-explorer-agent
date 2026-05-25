// Package handler provides HTTP handlers for PVC file operations.
package handler

import (
	"encoding/json"
	"net/http"
	"syscall"
)

// SpaceResponse is the JSON response for disk space usage.
type SpaceResponse struct {
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

// SpaceHandler returns an HTTP handler that reports disk space usage.
func SpaceHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var stat syscall.Statfs_t
		if err := syscall.Statfs(root, &stat); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bavail * uint64(stat.Bsize)
		used := total - free
		_ = json.NewEncoder(w).Encode(SpaceResponse{
			Used:  used,
			Total: total,
			Free:  free,
		})
	})
}
