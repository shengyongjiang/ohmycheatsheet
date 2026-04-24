package resolver

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestCache(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, sub := range []string{"common", "osx", "linux"} {
		os.MkdirAll(filepath.Join(dir, sub), 0o755)
	}
	os.WriteFile(filepath.Join(dir, "common", "tmux.md"), []byte("# tmux\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "common", "curl.md"), []byte("# curl\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "osx", "brew.md"), []byte("# brew\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "osx", "tmux.md"), []byte("# tmux\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "linux", "apt.md"), []byte("# apt\n"), 0o644)
	return dir
}

func TestResolve_CommonCommand(t *testing.T) {
	dir := setupTestCache(t)
	r := New(dir, []string{"osx", "common"})

	path, err := r.Resolve("curl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(dir, "common", "curl.md")
	if path != expected {
		t.Errorf("path = %q, want %q", path, expected)
	}
}

func TestResolve_PlatformPriority(t *testing.T) {
	dir := setupTestCache(t)
	r := New(dir, []string{"osx", "common"})

	path, err := r.Resolve("tmux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(dir, "osx", "tmux.md")
	if path != expected {
		t.Errorf("path = %q, want %q (should prefer osx)", path, expected)
	}
}

func TestResolve_PlatformOnly(t *testing.T) {
	dir := setupTestCache(t)
	r := New(dir, []string{"osx", "common"})

	path, err := r.Resolve("brew")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(dir, "osx", "brew.md")
	if path != expected {
		t.Errorf("path = %q, want %q", path, expected)
	}
}

func TestResolve_NotFound(t *testing.T) {
	dir := setupTestCache(t)
	r := New(dir, []string{"osx", "common"})

	_, err := r.Resolve("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}

func TestResolve_CrossPlatformNotVisible(t *testing.T) {
	dir := setupTestCache(t)
	r := New(dir, []string{"osx", "common"})

	_, err := r.Resolve("apt")
	if err == nil {
		t.Fatal("expected error: apt is linux-only, shouldn't resolve on osx")
	}
}
