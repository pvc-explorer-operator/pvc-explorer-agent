package main

import (
	"path/filepath"
	"testing"
)

func TestResolveRootReturnsAbsolutePath(t *testing.T) {
	root, err := resolveRoot("./testdata/demo")
	if err != nil {
		t.Fatalf("resolveRoot returned error: %v", err)
	}
	if !filepath.IsAbs(root) {
		t.Fatalf("expected absolute root, got %q", root)
	}
	if filepath.Base(root) != "demo" {
		t.Fatalf("expected root to end with demo, got %q", root)
	}
}