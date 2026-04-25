package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
)

func TestJsonStore_SetAndGet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := NewJSONStore(path)
	if err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}

	now := time.Now()
	err = s.SetEntryState("tmux", 0, model.EntryState{
		State:        model.StateRemembered,
		LastReviewed: now,
		ReviewCount:  1,
		Fingerprint:  "abc123",
	})
	if err != nil {
		t.Fatalf("SetEntryState: %v", err)
	}

	es, ok := s.GetEntryState("tmux", 0)
	if !ok {
		t.Fatal("expected entry state to exist")
	}
	if es.State != model.StateRemembered {
		t.Errorf("state = %q, want %q", es.State, model.StateRemembered)
	}
	if es.Fingerprint != "abc123" {
		t.Errorf("fingerprint = %q, want %q", es.Fingerprint, "abc123")
	}
}

func TestJsonStore_GetPageStates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, _ := NewJSONStore(path)

	s.SetEntryState("tmux", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "a"})
	s.SetEntryState("tmux", 2, model.EntryState{State: model.StateNeedsReview, Fingerprint: "b"})
	s.SetEntryState("curl", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "c"})

	states := s.GetPageStates("tmux")
	if len(states) != 2 {
		t.Fatalf("tmux states count = %d, want 2", len(states))
	}
	if states[0].State != model.StateRemembered {
		t.Errorf("tmux/0 state = %q", states[0].State)
	}
	if states[2].State != model.StateNeedsReview {
		t.Errorf("tmux/2 state = %q", states[2].State)
	}
}

func TestJsonStore_Persistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s1, _ := NewJSONStore(path)
	s1.SetEntryState("tmux", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "x"})
	if err := s1.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2, err := NewJSONStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	es, ok := s2.GetEntryState("tmux", 0)
	if !ok {
		t.Fatal("expected state to persist")
	}
	if es.State != model.StateRemembered {
		t.Errorf("state = %q after reload", es.State)
	}
}

func TestJsonStore_ListTrackedPages(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, _ := NewJSONStore(path)
	s.SetEntryState("tmux", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "a"})
	s.SetEntryState("curl", 1, model.EntryState{State: model.StateNeedsReview, Fingerprint: "b"})

	pages := s.ListTrackedPages()
	if len(pages) != 2 {
		t.Fatalf("pages count = %d, want 2", len(pages))
	}
}

func TestJsonStore_ResetPage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, _ := NewJSONStore(path)
	s.SetEntryState("tmux", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "a"})
	s.SetEntryState("tmux", 1, model.EntryState{State: model.StateNeedsReview, Fingerprint: "b"})
	s.SetEntryState("curl", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "c"})

	s.ResetPage("tmux")

	states := s.GetPageStates("tmux")
	if len(states) != 0 {
		t.Errorf("tmux states after reset = %d, want 0", len(states))
	}
	curlStates := s.GetPageStates("curl")
	if len(curlStates) != 1 {
		t.Errorf("curl states should be untouched, got %d", len(curlStates))
	}
}

func TestJsonStore_ResetAll(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, _ := NewJSONStore(path)
	s.SetEntryState("tmux", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "a"})
	s.SetEntryState("curl", 0, model.EntryState{State: model.StateRemembered, Fingerprint: "b"})

	s.ResetAll()

	if len(s.ListTrackedPages()) != 0 {
		t.Error("expected no tracked pages after ResetAll")
	}
}

func TestJsonStore_NotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, _ := NewJSONStore(path)

	_, ok := s.GetEntryState("tmux", 99)
	if ok {
		t.Error("expected no state for non-existent entry")
	}
}
