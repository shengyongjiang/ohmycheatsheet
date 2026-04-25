package render

import (
	"strings"
	"testing"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
)

func testPage() *model.Page {
	return &model.Page{
		Name:        "tmux",
		Description: "Terminal multiplexer.",
		Entries: []model.Entry{
			{Index: 0, Description: "Start a new session", Command: "tmux"},
			{Index: 1, Description: "Start a new named session", Command: "tmux new-session -s name"},
			{Index: 2, Description: "Kill a session", Command: "tmux kill-session -t name"},
		},
	}
}

func TestRenderText_NoState(t *testing.T) {
	page := testPage()
	states := map[int]model.EntryState{}
	out := RenderText(page, states, false, false, 42)

	if !strings.Contains(out, "tmux") {
		t.Error("output should contain command name")
	}
	if !strings.Contains(out, "Terminal multiplexer") {
		t.Error("output should contain description")
	}
	if !strings.Contains(out, "Start a new session") {
		t.Error("output should contain entry description")
	}
}

func TestRenderText_HidesRemembered(t *testing.T) {
	page := testPage()
	states := map[int]model.EntryState{
		0: {State: model.StateRemembered},
	}
	out := RenderText(page, states, false, false, 42)

	if strings.Contains(out, "Start a new session") {
		t.Error("remembered entry should be hidden")
	}
	if !strings.Contains(out, "Start a new named session") {
		t.Error("non-remembered entry should be visible")
	}
	if !strings.Contains(out, "1 remembered") {
		t.Error("should show hidden count")
	}
}

func TestRenderText_ShowAll(t *testing.T) {
	page := testPage()
	states := map[int]model.EntryState{
		0: {State: model.StateRemembered},
	}
	out := RenderText(page, states, true, false, 42)

	if !strings.Contains(out, "Start a new session") {
		t.Error("--all should show remembered entries")
	}
}

func TestRenderText_NeedsReview(t *testing.T) {
	page := testPage()
	states := map[int]model.EntryState{
		1: {State: model.StateNeedsReview},
	}
	out := RenderText(page, states, false, false, 42)

	if !strings.Contains(out, "needs review") {
		t.Error("needs-review entry should have tag")
	}
}
