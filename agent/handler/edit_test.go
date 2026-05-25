package handler

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestEditHandler(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "foo.txt")
	os.WriteFile(filePath, []byte("old content"), 0644)

	h := EditHandler(dir)
	body := bytes.NewBufferString("new content")
	r := httptest.NewRequestWithContext(context.Background(), "PUT", "/api/edit?path=foo.txt", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if string(data) != "new content" {
		t.Fatalf("bad file content: %s", string(data))
	}
}
