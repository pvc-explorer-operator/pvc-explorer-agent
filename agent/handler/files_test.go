package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesHandler_ListAndDelete(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "foo.txt"), []byte("abc"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "bar.txt"), []byte("def"), 0644)

	h := FilesHandler(dir, func(*http.Request) bool { return false })
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/api/files?path=", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp ListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if len(resp.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(resp.Entries))
	}

	// Test delete file
	r = httptest.NewRequestWithContext(context.Background(), "DELETE", "/api/files?path=foo.txt", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, err := os.Stat(filepath.Join(dir, "foo.txt")); !os.IsNotExist(err) {
		t.Fatalf("file not deleted")
	}

	// Test delete dir
	r = httptest.NewRequestWithContext(context.Background(), "DELETE", "/api/files?path=subdir", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, err := os.Stat(filepath.Join(dir, "subdir")); !os.IsNotExist(err) {
		t.Fatalf("dir not deleted")
	}
}

func TestClearHandler(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "foo.txt"), []byte("abc"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "bar.txt"), []byte("def"), 0644)

	h := ClearHandler(dir)
	r := httptest.NewRequestWithContext(context.Background(), "POST", "/api/clear", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	files, _ := os.ReadDir(dir)
	if len(files) != 0 {
		t.Fatalf("expected dir to be empty after clear")
	}
}
