package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestSpaceHandler(t *testing.T) {
	dir := t.TempDir()
	h := SpaceHandler(dir)
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/api/space", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp SpaceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if resp.Total == 0 {
		t.Fatalf("expected nonzero total")
	}
}
