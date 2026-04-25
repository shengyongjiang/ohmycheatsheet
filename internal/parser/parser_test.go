package parser

import (
	"testing"
)

const cheatshGit = `#[cheat.sheets:git]
# git
# Set your identity.
git config --global user.name "John Doe"
git config --global user.email johndoe@example.com

# Stage all changes for commit.
git add [--all|-A]

# Stash changes locally. This will keep the changes in a separate changelist, -
# called 'stash', and the working directory is cleaned. You can apply changes
# from the stash at any time.
git stash

#[cheat:git]
# To start a new branch:
git checkout -b <branch_name>

#[tldr:git]
# git
# Distributed version control system.

# Show help:
git --help
`

func TestParseCheatsh_Basic(t *testing.T) {
	page, err := ParseCheatsh(cheatshGit, "git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Name != "git" {
		t.Errorf("name = %q, want %q", page.Name, "git")
	}

	if page.Description != "git Set your identity." {
		t.Errorf("description = %q, want %q", page.Description, "git Set your identity.")
	}

	if len(page.Entries) < 4 {
		t.Fatalf("entries count = %d, want at least 4", len(page.Entries))
	}

	e0 := page.Entries[0]
	if e0.Description != "Stage all changes for commit." {
		t.Errorf("entry[0].Description = %q", e0.Description)
	}
	if e0.Command != "git add [--all|-A]" {
		t.Errorf("entry[0].Command = %q", e0.Command)
	}

	e1 := page.Entries[1]
	if e1.Description != "Stash changes locally. This will keep the changes in a separate changelist, - called 'stash', and the working directory is cleaned. You can apply changes from the stash at any time." {
		t.Errorf("entry[1].Description = %q", e1.Description)
	}
	if e1.Command != "git stash" {
		t.Errorf("entry[1].Command = %q", e1.Command)
	}
}

func TestParseCheatsh_MultipleSourceSections(t *testing.T) {
	page, err := ParseCheatsh(cheatshGit, "git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, e := range page.Entries {
		if e.Description == "To start a new branch" {
			found = true
			if e.Command != "git checkout -b <branch_name>" {
				t.Errorf("command = %q", e.Command)
			}
		}
	}
	if !found {
		t.Error("entry from cheat section not found")
	}

	found = false
	for _, e := range page.Entries {
		if e.Description == "Show help" {
			found = true
			if e.Command != "git --help" {
				t.Errorf("command = %q", e.Command)
			}
		}
	}
	if !found {
		t.Error("entry from tldr section not found")
	}
}

func TestParseCheatsh_EmptyContent(t *testing.T) {
	page, err := ParseCheatsh("", "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Entries) != 0 {
		t.Errorf("entries count = %d, want 0", len(page.Entries))
	}
}

func TestParseCheatsh_Fingerprint(t *testing.T) {
	page, err := ParseCheatsh(cheatshGit, "git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range page.Entries {
		if e.Fingerprint == "" {
			t.Errorf("entry[%d].Fingerprint is empty", i)
		}
	}
}

func TestParseCheatsh_IndexSequential(t *testing.T) {
	page, err := ParseCheatsh(cheatshGit, "git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range page.Entries {
		if e.Index != i {
			t.Errorf("entry[%d].Index = %d, want %d", i, e.Index, i)
		}
	}
}

func TestParseCheatsh_NoDescription(t *testing.T) {
	content := `#[cheat.sheets:curl]
# Download a URL.
curl https://example.com
`
	page, err := ParseCheatsh(content, "curl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Entries) != 1 {
		t.Fatalf("entries count = %d, want 1", len(page.Entries))
	}
	if page.Entries[0].Description != "Download a URL." {
		t.Errorf("description = %q", page.Entries[0].Description)
	}
}
