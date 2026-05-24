package embedui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed ui
var uiFS embed.FS

func Handler(overlayDir string) http.Handler {
	sub, err := fs.Sub(uiFS, "ui")
	if err != nil {
		panic("embedui: " + err.Error())
	}
	fsHandler := http.FileServer(http.FS(sub))
	overlayHandler := http.FileServer(http.Dir(overlayDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if overlayDir != "" {
			candidate := filepath.Join(overlayDir, filepath.Clean("/"+r.URL.Path))
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				overlayHandler.ServeHTTP(w, r)
				return
			}
		}
		fsHandler.ServeHTTP(w, r)
	})
}
