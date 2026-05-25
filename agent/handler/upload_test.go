package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadHandler(t *testing.T) {
	dir := t.TempDir()
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("file", "foo.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	io.WriteString(fw, "hello upload")
	w.Close()

	h := UploadHandler(dir)
	r := httptest.NewRequestWithContext(context.Background(), "POST", "/api/upload?path=", buf)
	r.Header.Set("Content-Type", w.FormDataContentType())
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, r)
	if resp.Code != 200 {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	data, err := os.ReadFile(filepath.Join(dir, "foo.txt"))
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if string(data) != "hello upload" {
		t.Fatalf("bad file content: %s", string(data))
	}
}
