package tui

import bubbletea "github.com/charmbracelet/bubbletea"

type keyMap struct {
	Up     []string
	Down   []string
	Tab    []string
	TogAll []string
	Quit   []string
	Help   []string
}

var keys = keyMap{
	Up:     []string{"up", "k"},
	Down:   []string{"down", "j"},
	Tab:    []string{"tab"},
	TogAll: []string{"a"},
	Quit:   []string{"q", "esc", "ctrl+c"},
	Help:   []string{"?"},
}

func matchKey(msg bubbletea.KeyMsg, bindings []string) bool {
	for _, b := range bindings {
		if msg.String() == b {
			return true
		}
	}
	return false
}
