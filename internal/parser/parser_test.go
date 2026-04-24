package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const tmuxMd = `# tmux

> Terminal multiplexer.
> It allows multiple sessions with windows, panes, and more.
> See also: ` + "`zellij`" + `, ` + "`screen`" + `.
> More information: <https://github.com/tmux/tmux>.

- Start a new session:

` + "`tmux`" + `

- Start a new named [s]ession:

` + "`tmux {{[new|new-session]}} -s {{name}}`" + `

- Kill a session by [t]arget name:

` + "`tmux kill-session -t {{name}}`" + `
`

func TestParseString(t *testing.T) {
	page, err := ParseString(tmuxMd, "tmux", "common")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Name != "tmux" {
		t.Errorf("name = %q, want %q", page.Name, "tmux")
	}
	if page.Platform != "common" {
		t.Errorf("platform = %q, want %q", page.Platform, "common")
	}
	if page.Description != "Terminal multiplexer. It allows multiple sessions with windows, panes, and more." {
		t.Errorf("description = %q", page.Description)
	}
	if page.URL != "https://github.com/tmux/tmux" {
		t.Errorf("url = %q", page.URL)
	}
	if len(page.SeeAlso) != 2 || page.SeeAlso[0] != "zellij" || page.SeeAlso[1] != "screen" {
		t.Errorf("see_also = %v", page.SeeAlso)
	}
	if len(page.Entries) != 3 {
		t.Fatalf("entries count = %d, want 3", len(page.Entries))
	}

	e0 := page.Entries[0]
	if e0.Index != 0 {
		t.Errorf("entry[0].Index = %d", e0.Index)
	}
	if e0.Description != "Start a new session" {
		t.Errorf("entry[0].Description = %q", e0.Description)
	}
	if e0.Command != "tmux" {
		t.Errorf("entry[0].Command = %q", e0.Command)
	}
	if e0.Fingerprint == "" {
		t.Error("entry[0].Fingerprint is empty")
	}

	e1 := page.Entries[1]
	if e1.Description != "Start a new named [s]ession" {
		t.Errorf("entry[1].Description = %q", e1.Description)
	}
	if e1.Command != "tmux {{[new|new-session]}} -s {{name}}" {
		t.Errorf("entry[1].Command = %q", e1.Command)
	}
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "common", "tmux.md")
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, []byte(tmuxMd), 0o644)

	page, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Name != "tmux" {
		t.Errorf("name = %q, want %q", page.Name, "tmux")
	}
	if page.Platform != "common" {
		t.Errorf("platform = %q, want %q", page.Platform, "common")
	}
}

func TestParseEmptyPage(t *testing.T) {
	input := `# empty

> An empty command.
`
	page, err := ParseString(input, "empty", "common")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Entries) != 0 {
		t.Errorf("entries count = %d, want 0", len(page.Entries))
	}
}

func TestParseDescription_NoSeeAlso(t *testing.T) {
	input := `# curl

> Transfer data from or to a server.
> More information: <https://curl.se/docs/manpage.html>.

- Download a URL:

` + "`curl {{https://example.com}}`" + `
`
	page, err := ParseString(input, "curl", "common")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Description != "Transfer data from or to a server." {
		t.Errorf("description = %q", page.Description)
	}
	if len(page.SeeAlso) != 0 {
		t.Errorf("see_also = %v, want empty", page.SeeAlso)
	}
	if page.URL != "https://curl.se/docs/manpage.html" {
		t.Errorf("url = %q", page.URL)
	}
	if len(page.Entries) != 1 {
		t.Fatalf("entries count = %d, want 1", len(page.Entries))
	}
}
