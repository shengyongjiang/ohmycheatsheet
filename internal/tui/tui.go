package tui

import (
	"fmt"
	"strings"
	"time"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shengyongjiang/ocheetsheet/internal/model"
	"github.com/shengyongjiang/ocheetsheet/internal/store"
)

type Model struct {
	page     *model.Page
	states   map[int]model.EntryState
	store    store.StateStore
	cursor   int
	showAll  bool
	visible  []int
	width    int
	height   int
	helpOpen bool
	dirty    bool
}

func New(page *model.Page, states map[int]model.EntryState, st store.StateStore) Model {
	m := Model{
		page:   page,
		states: states,
		store:  st,
	}
	m.rebuildVisible()
	return m
}

func (m *Model) rebuildVisible() {
	m.visible = m.visible[:0]
	for _, e := range m.page.Entries {
		es, ok := m.states[e.Index]
		if !m.showAll && ok && es.State == model.StateRemembered {
			continue
		}
		m.visible = append(m.visible, e.Index)
	}
	if m.cursor >= len(m.visible) {
		m.cursor = max(0, len(m.visible)-1)
	}
}

func (m *Model) currentState() model.MemoryState {
	if len(m.visible) == 0 {
		return model.StateNotRemembered
	}
	idx := m.visible[m.cursor]
	es, ok := m.states[idx]
	if !ok {
		return model.StateNotRemembered
	}
	return es.State
}

func (m *Model) cycleState() {
	if len(m.visible) == 0 {
		return
	}
	idx := m.visible[m.cursor]
	es := m.states[idx]
	es.State = es.State.Next()
	es.LastReviewed = time.Now()
	es.ReviewCount++
	entry := m.page.Entries[idx]
	es.Fingerprint = entry.Fingerprint
	m.states[idx] = es
	m.store.SetEntryState(m.page.Name, idx, es)
	m.dirty = true
	m.rebuildVisible()
}

func (m Model) Init() bubbletea.Cmd {
	return nil
}

func (m Model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case bubbletea.KeyMsg:
		if m.helpOpen {
			m.helpOpen = false
			return m, nil
		}

		switch {
		case matchKey(msg, keys.Quit):
			if m.dirty {
				m.store.Save()
			}
			return m, bubbletea.Quit

		case matchKey(msg, keys.Down):
			if m.cursor < len(m.visible)-1 {
				m.cursor++
			}

		case matchKey(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case matchKey(msg, keys.Tab):
			m.cycleState()

		case matchKey(msg, keys.TogAll):
			m.showAll = !m.showAll
			m.rebuildVisible()

		case matchKey(msg, keys.Help):
			m.helpOpen = true
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.helpOpen {
		return m.viewHelp()
	}

	var b strings.Builder

	header := headerStyle.Render(
		titleStyle.Render(m.page.Name) + dimStyle.Render(" - "+m.page.Description),
	)
	b.WriteString(header + "\n\n")

	if len(m.visible) == 0 {
		b.WriteString(dimStyle.Render("  All entries remembered! Press 'a' to show all.") + "\n")
	} else {
		idx := m.visible[m.cursor]
		entry := m.page.Entries[idx]
		state := m.currentState()

		counter := counterStyle.Render(fmt.Sprintf("[%d/%d]", m.cursor+1, len(m.visible)))

		var stateTag string
		switch state {
		case model.StateNeedsReview:
			stateTag = reviewTagStyle.Render("  * needs review")
		case model.StateRemembered:
			stateTag = rememberedTagStyle.Render("  v remembered")
		default:
			stateTag = dimStyle.Render("  o not remembered")
		}

		b.WriteString(fmt.Sprintf("  %s%s\n\n", counter, stateTag))
		b.WriteString(fmt.Sprintf("  %s:\n", descStyle.Render(entry.Description)))
		b.WriteString(fmt.Sprintf("  %s\n", cmdStyle.Render(entry.Command)))
	}

	b.WriteString("\n")

	hiddenCount := len(m.page.Entries) - len(m.visible)
	if hiddenCount > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  (%d remembered entries hidden)", hiddenCount)) + "\n")
	}

	allToggle := "a:show all"
	if m.showAll {
		allToggle = "a:filter"
	}
	status := statusBarStyle.Width(m.width).Render(
		helpStyle.Render(fmt.Sprintf("Tab:cycle state  j/k:navigate  %s  q:quit  ?:help", allToggle)),
	)
	b.WriteString(status)

	return lipgloss.NewStyle().Width(m.width).Render(b.String())
}

func (m Model) viewHelp() string {
	help := `
  Key Bindings:

  Tab        Cycle memory state (not remembered -> remembered -> needs review)
  j / Down   Next entry
  k / Up     Previous entry
  a          Toggle show all / filter remembered
  q / Esc    Quit (auto-saves)
  ?          Toggle this help

  Memory States:

  o not remembered   Default, always shown
  v remembered       Hidden by default
  * needs review     Highlighted, shown with priority

  Press any key to close help.
`
	return headerStyle.Render(titleStyle.Render("ocs help")) + help
}
