package handler

import (
	"encoding/json"
	"net/http"
	"os"
)

type RenameRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type RenameResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func RenameHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var req RenameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		absFrom, err := safeJoin(root, req.From)
		if err != nil {
			http.Error(w, `{"error":"bad source path"}`, http.StatusBadRequest)
			return
		}
		absTo, err := safeJoin(root, req.To)
		if err != nil {
			http.Error(w, `{"error":"bad destination path"}`, http.StatusBadRequest)
			return
		}
		if _, err := os.Lstat(absFrom); err != nil {
			if os.IsNotExist(err) {
				http.Error(w, `{"error":"source not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		if err := os.Rename(absFrom, absTo); err != nil {
			http.Error(w, `{"error":"rename failed"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(RenameResponse{From: req.From, To: req.To})
	})
}
