package handler

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recover500(w)
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		p := r.URL.Query().Get("path")
		abs, err := safeJoin(root, p)
		if err != nil {
			http.Error(w, `{"error":"bad path"}`, http.StatusBadRequest)
			return
		}
		info, err := os.Stat(abs)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			}
			return
		}
		if info.IsDir() {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
			zipWriter := zip.NewWriter(w)
			defer zipWriter.Close()
			filepath.Walk(abs, func(fp string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				rel, _ := filepath.Rel(abs, fp)
				if rel == "." {
					return nil
				}
				hdr, herr := zip.FileInfoHeader(fi)
				if herr != nil {
					return nil
				}
				hdr.Name = rel
				if fi.IsDir() {
					hdr.Name += "/"
					_, err = zipWriter.CreateHeader(hdr)
					if err != nil {
						return nil
					}
					return nil
				}
				wtr, err := zipWriter.CreateHeader(hdr)
				if err != nil {
					return nil
				}
				f, err := os.Open(fp)
				if err != nil {
					return nil
				}
				_, _ = io.Copy(wtr, f)
				f.Close()
				return nil
			})
			return
		}
		f, err := os.Open(abs)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(abs))
		w.Header().Set("Content-Type", "application/octet-stream")
		io.Copy(w, f)
	})
}
