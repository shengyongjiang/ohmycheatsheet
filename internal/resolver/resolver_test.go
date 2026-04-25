package resolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shengyongjiang/ohmycheatsheet/internal/source"
)

func setupTestCache(t *testing.T) (string, *source.CheatshSource) {
	t.Helper()
	dir := t.TempDir()

	content := `#[cheat.sheets:tmux]
# Start a new session.
tmux

# Start a new named session.
tmux new-session -s name

# Kill a session by name.
tmux kill-session -t name
`
	os.WriteFile(filepath.Join(dir, "tmux.txt"), []byte(content), 0o644)

	gitContent := `#[cheat.sheets:git]
# Show help.
git --help
`
	os.WriteFile(filepath.Join(dir, "git.txt"), []byte(gitContent), 0o644)

	gitLogContent := `#[cheat.sheets:git-log]
# Show commit log.
git log
`
	os.WriteFile(filepath.Join(dir, "git-log.txt"), []byte(gitLogContent), 0o644)

	src := source.NewCheatshSource(dir)
	return dir, src
}

func TestResolve_CachedCommand(t *testing.T) {
	_, src := setupTestCache(t)
	r := New(src)

	page, err := r.Resolve("tmux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Name != "tmux" {
		t.Errorf("name = %q, want %q", page.Name, "tmux")
	}
	if len(page.Entries) != 3 {
		t.Errorf("entries count = %d, want 3", len(page.Entries))
	}
}

func TestResolve_ReturnsPage(t *testing.T) {
	_, src := setupTestCache(t)
	r := New(src)

	page, err := r.Resolve("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Name != "git" {
		t.Errorf("name = %q, want %q", page.Name, "git")
	}
	if len(page.Entries) < 1 {
		t.Errorf("entries count = %d, want at least 1", len(page.Entries))
	}
}

func TestListRelatedCommands(t *testing.T) {
	_, src := setupTestCache(t)
	r := New(src)

	related, err := r.ListRelatedCommands("git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(related) != 1 || related[0] != "git-log" {
		t.Errorf("related = %v, want [git-log]", related)
	}
}

func TestListAllCommands(t *testing.T) {
	_, src := setupTestCache(t)
	r := New(src)

	commands, err := r.ListAllCommands()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commands) != 3 {
		t.Errorf("commands count = %d, want 3", len(commands))
	}
}

func TestResolve_NotCached_WillFail(t *testing.T) {
	dir := t.TempDir()
	src := source.NewCheatshSource(dir)
	r := New(src)

	_, err := r.Resolve("nonexistent-command-xyz-12345")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}
