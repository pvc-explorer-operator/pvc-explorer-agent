package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadHandler_File(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "foo.txt")
	os.WriteFile(filePath, []byte("hello world"), 0644)

	h := DownloadHandler(dir)
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/api/download?path=foo.txt", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Disposition"); ct == "" {
		t.Fatalf("missing Content-Disposition")
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("hello world")) {
		t.Fatalf("file content missing")
	}
}

func TestDownloadHandler_ZipDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "foo.txt"), []byte("abc"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "bar.txt"), []byte("def"), 0644)

	h := DownloadHandler(dir)
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/api/download?path=", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	zr, err := zip.NewReader(bytes.NewReader(w.Body.Bytes()), int64(w.Body.Len()))
	if err != nil {
		t.Fatalf("bad zip: %v", err)
	}
	found := false
	for _, f := range zr.File {
		if f.Name == "foo.txt" {
			found = true
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			rc.Close()
			if string(data) != "abc" {
				t.Fatalf("bad file content in zip")
			}
		}
	}
	if !found {
		t.Fatalf("foo.txt not found in zip")
	}
}
