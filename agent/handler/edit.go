package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type EditResponse struct {
	Written string `json:"written"`
}

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
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		tmp, err := os.CreateTemp(dir, ".edit-*")
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		defer func() { tmp.Close(); os.Remove(tmp.Name()) }()
		if _, err = io.Copy(tmp, r.Body); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		if err = os.Rename(tmp.Name(), abs); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(EditResponse{Written: p})
	})
}
