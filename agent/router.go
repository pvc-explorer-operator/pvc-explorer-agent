package agent

import (
	"net/http"

	"github.com/pvc-explorer-operator/pvc-explorer-agent/agent/handler"
)

func RegisterRoutes(mux *http.ServeMux, root string, isReadonly func(*http.Request) bool) {
	mux.Handle("/api/space", handler.SpaceHandler(root))
	mux.Handle("/api/files", handler.FilesHandler(root, isReadonly))
	mux.Handle("/api/download", handler.DownloadHandler(root))
	mux.Handle("/api/upload", guardWrite(handler.UploadHandler(root), isReadonly))
	mux.Handle("/api/edit", guardWrite(handler.EditHandler(root), isReadonly))
	mux.Handle("/api/rename", guardWrite(handler.RenameHandler(root), isReadonly))
	mux.Handle("/api/clear", guardWrite(handler.ClearHandler(root), isReadonly))
}

func guardWrite(h http.Handler, isReadonly func(*http.Request) bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isReadonly(r) {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"read-only mode"}`, http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
