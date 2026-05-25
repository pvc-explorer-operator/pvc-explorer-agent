package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzSafeJoin(f *testing.F) {
	seeds := []string{
		"/valid/path",
		"../../../etc",
		"../../etc/passwd",
		"foo/../../../bar",
		".../.../...//",
		"/",
		"",
		".",
		"..",
		"a/b/c/d/e/f",
		"a/../../../b",
		strings.Repeat("../", 50),
		"foo/..",
		"foo/./bar/../baz",
		"\x00",
		"\x00../../etc",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, rel string) {
		root := "/tmp/fuzz-root"
		abs, err := safeJoin(root, rel)
		if err != nil {
			return
		}
		if !strings.HasPrefix(abs, root) {
			t.Errorf("safeJoin(%q, %q) = %q, escaped root %q", root, rel, abs, root)
		}
	})
}

func FuzzFilesHandlerList(f *testing.F) {
	seeds := []string{
		"/",
		".",
		"foo",
		"../etc",
		"a/b/c",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, path string) {
		dir := t.TempDir()
		h := FilesHandler(dir, func(*http.Request) bool { return false })
		r := httptest.NewRequest("GET", "/api/files?path="+path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		code := w.Code
		if code != 200 && code != 400 && code != 404 && code != 500 {
			t.Errorf("unexpected status %d for path %q", code, path)
		}
	})
}

func FuzzEditHandler(f *testing.F) {
	seeds := []string{
		"test.txt",
		"../etc/passwd",
		"a/b/c/file.txt",
		"../../../tmp/foo",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, path string) {
		dir := t.TempDir()
		os.MkdirAll(filepath.Join(dir, "sub"), 0755)
		h := EditHandler(dir)
		r := httptest.NewRequest("PUT", "/api/edit?path="+path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		code := w.Code
		if code != 200 && code != 400 && code != 404 && code != 500 {
			t.Errorf("unexpected status %d for path %q", code, path)
		}
	})
}
