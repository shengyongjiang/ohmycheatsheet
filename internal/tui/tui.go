package tui

import (
	"fmt"
	"strings"
	"time"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
	"github.com/shengyongjiang/ohmycheatsheet/internal/resolver"
	"github.com/shengyongjiang/ohmycheatsheet/internal/shuffle"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
)

type displayEntry struct {
	pageName   string
	entry      model.Entry
	isBackfill bool
}

type Model struct {
	page          *model.Page
	store         store.StateStore
	resolver      *resolver.Resolver
	entries       []displayEntry
	visible       []int
	cursor        int
	scrollOffset  int
	showAll       bool
	originalCount int
	parsedPages   map[string]*model.Page
	width         int
	height        int
	helpOpen      bool
	dirty         bool
}

func New(page *model.Page, states map[int]model.EntryState, st store.StateStore, res *resolver.Resolver, seed int64) Model {
	m := Model{
		page:          page,
		store:         st,
		resolver:      res,
		originalCount: min(10, len(page.Entries)),
		parsedPages:   make(map[string]*model.Page),
	}

	shuffled := shuffle.ShuffleEntries(page.Entries, seed)
	for _, e := range shuffled {
		m.entries = append(m.entries, displayEntry{
			pageName:   page.Name,
			entry:      e,
			isBackfill: false,
		})
	}

	m.parsedPages[page.Name] = page

	m.rebuildVisible()
	return m
}

func (m *Model) getState(de displayEntry) model.EntryState {
	es, _ := m.store.GetEntryState(de.pageName, de.entry.Index)
	return es
}

func (m *Model) cycleStateForward() {
	if len(m.visible) == 0 {
		return
	}
	de := m.entries[m.visible[m.cursor]]
	es := m.getState(de)
	es.State = es.State.Next()
	es.LastReviewed = time.Now()
	es.ReviewCount++
	es.Fingerprint = de.entry.Fingerprint
	m.store.SetEntryState(de.pageName, de.entry.Index, es)
	m.dirty = true
}

func (m *Model) cycleStateBackward() {
	if len(m.visible) == 0 {
		return
	}
	de := m.entries[m.visible[m.cursor]]
	es := m.getState(de)
	es.State = es.State.Prev()
	es.LastReviewed = time.Now()
	es.ReviewCount++
	es.Fingerprint = de.entry.Fingerprint
	m.store.SetEntryState(de.pageName, de.entry.Index, es)
	m.dirty = true
}

func (m *Model) setState(target model.MemoryState) {
	if len(m.visible) == 0 {
		return
	}
	de := m.entries[m.visible[m.cursor]]
	es := m.getState(de)
	es.State = target
	es.LastReviewed = time.Now()
	es.ReviewCount++
	es.Fingerprint = de.entry.Fingerprint
	m.store.SetEntryState(de.pageName, de.entry.Index, es)
	m.dirty = true
}

func (m *Model) resetAll() {
	for _, idx := range m.visible {
		de := m.entries[idx]
		m.store.SetEntryState(de.pageName, de.entry.Index, model.EntryState{})
	}
	m.entries = m.entries[:0]
	for _, e := range m.page.Entries {
		m.entries = append(m.entries, displayEntry{
			pageName:   m.page.Name,
			entry:      e,
			isBackfill: false,
		})
	}
	m.dirty = true
	m.cursor = 0
	m.scrollOffset = 0
	m.rebuildVisible()
}

func (m *Model) rebuildVisible() {
	m.visible = m.visible[:0]
	for i, de := range m.entries {
		es := m.getState(de)
		if es.State == model.StateRemembered && !m.showAll {
			continue
		}
		m.visible = append(m.visible, i)
	}
	if len(m.visible) < m.originalCount && !m.showAll {
		m.backfill()
	}
	if m.cursor >= len(m.visible) {
		m.cursor = max(0, len(m.visible)-1)
	}
}

func (m *Model) backfill() {
	needed := m.originalCount - len(m.visible)
	if needed <= 0 {
		return
	}

	for _, cmd := range m.collectRelatedCommands() {
		if needed <= 0 {
			break
		}
		page := m.loadPage(cmd)
		if page == nil {
			continue
		}
		for _, entry := range page.Entries {
			if needed <= 0 {
				break
			}
			de := displayEntry{pageName: page.Name, entry: entry, isBackfill: true}
			if m.isAlreadyInEntries(de) {
				continue
			}
			es := m.getState(de)
			if es.State == model.StateRemembered {
				continue
			}
			m.entries = append(m.entries, de)
			m.visible = append(m.visible, len(m.entries)-1)
			needed--
		}
	}
}

func (m *Model) collectRelatedCommands() []string {
	if m.resolver == nil {
		return nil
	}
	seen := map[string]bool{m.page.Name: true}
	var result []string

	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}

	if prefixed, err := m.resolver.ListRelatedCommands(m.page.Name); err == nil {
		for _, p := range prefixed {
			add(p)
		}
	}

	return result
}

func (m *Model) loadPage(command string) *model.Page {
	if p, ok := m.parsedPages[command]; ok {
		return p
	}
	if m.resolver == nil {
		return nil
	}
	page, err := m.resolver.Resolve(command)
	if err != nil {
		return nil
	}
	m.parsedPages[command] = page
	return page
}

func (m *Model) isAlreadyInEntries(de displayEntry) bool {
	for _, existing := range m.entries {
		if existing.pageName == de.pageName && existing.entry.Index == de.entry.Index {
			return true
		}
	}
	return false
}

func (m *Model) ensureCursorVisible() {
	vh := m.viewportHeight()
	if vh <= 0 {
		return
	}
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+vh {
		m.scrollOffset = m.cursor - vh + 1
	}
}

func (m *Model) viewportHeight() int {
	available := m.height - 6
	if available < 1 {
		return 1
	}
	vh := available / 2
	if vh > 15 {
		vh = 15
	}
	return vh
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
				m.ensureCursorVisible()
			}

		case matchKey(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.ensureCursorVisible()
			}

		case matchKey(msg, keys.Right):
			m.cycleStateForward()

		case matchKey(msg, keys.Left):
			m.cycleStateBackward()

		case matchKey(msg, keys.MarkRemember):
			m.setState(model.StateRemembered)

		case matchKey(msg, keys.MarkReview):
			m.setState(model.StateNeedsReview)

		case matchKey(msg, keys.TogAll):
			m.showAll = !m.showAll
			m.rebuildVisible()

		case matchKey(msg, keys.Reset):
			m.resetAll()

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
	b.WriteString(header + "\n")

	if len(m.visible) == 0 {
		b.WriteString(dimStyle.Render("  All entries remembered! Press 'a' to show all.") + "\n")
	} else {
		vh := m.viewportHeight()
		end := m.scrollOffset + vh
		if end > len(m.visible) {
			end = len(m.visible)
		}

		for vi := m.scrollOffset; vi < end; vi++ {
			idx := m.visible[vi]
			de := m.entries[idx]
			es := m.getState(de)
			isCursor := vi == m.cursor

			if de.isBackfill && vi > m.scrollOffset {
				prevIdx := m.visible[vi-1]
				if !m.entries[prevIdx].isBackfill {
					b.WriteString(separatorStyle.Render("  ── backfill ──") + "\n")
				}
			}

			var stateIndicator string
			switch es.State {
			case model.StateNeedsReview:
				stateIndicator = reviewTagStyle.Render("+")
			case model.StateRemembered:
				stateIndicator = dimStyle.Render("x")
			default:
				stateIndicator = dimStyle.Render("o")
			}

			desc := de.entry.Description
			if de.isBackfill {
				desc = backfillPageStyle.Render("["+de.pageName+"]") + " " + desc
			}

			var dStyle, cStyle lipgloss.Style
			switch es.State {
			case model.StateRemembered:
				dStyle = rememberedDescStyle
				cStyle = rememberedCmdStyle
			case model.StateNeedsReview:
				dStyle = needsReviewDescStyle
				cStyle = needsReviewCmdStyle
			default:
				dStyle = descStyle
				cStyle = cmdStyle
			}

			if isCursor {
				b.WriteString(fmt.Sprintf("  %s  %s\n", stateIndicator, selectedStyle.Render(desc)))
				b.WriteString(fmt.Sprintf("       %s\n", selectedCmdStyle.Render(de.entry.Command)))
			} else {
				b.WriteString(fmt.Sprintf("  %s  %s\n", stateIndicator, dStyle.Render(desc)))
				b.WriteString(fmt.Sprintf("       %s\n", cStyle.Render(de.entry.Command)))
			}
		}
	}

	b.WriteString("\n")

	hiddenCount := 0
	for _, de := range m.entries {
		es := m.getState(de)
		if es.State == model.StateRemembered && !m.showAll {
			hiddenCount++
		}
	}
	if hiddenCount > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  (%d remembered entries hidden)", hiddenCount)) + "\n")
	}

	allToggle := "a:show all"
	if m.showAll {
		allToggle = "a:filter"
	}
	status := statusBarStyle.Width(m.width).Render(
		helpStyle.Render(fmt.Sprintf("←/→:cycle  x:remember  Enter:review  j/k:nav  %s  r:reset  q:quit  ?:help", allToggle)),
	)
	b.WriteString(status)

	return lipgloss.NewStyle().Width(m.width).Render(b.String())
}

func (m Model) viewHelp() string {
	help := `
  Key Bindings:

  ← Left     Cycle state backward (not remembered -> needs review -> remembered)
  → Right    Cycle state forward (not remembered -> remembered -> needs review)
  x / X      Mark as remembered
  Enter      Mark as needs review
  j / Down   Next entry
  k / Up     Previous entry
  a          Toggle show all / filter remembered
  r          Reset all states and refresh list
  q / Esc    Quit (auto-saves)
  ?          Toggle this help

  Memory States:

  o not remembered   Default, always shown
  x remembered       Hidden by default
  + needs review     Highlighted, shown with priority

  Entries from related or random pages are shown as backfill
  when you've memorized entries from the current page.

  Press any key to close help.
`
	return headerStyle.Render(titleStyle.Render("omcs help")) + help
}
