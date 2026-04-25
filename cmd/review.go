package cmd

import (
	"fmt"
	"strings"
	"time"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shengyongjiang/ohmycheatsheet/internal/config"
	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
	"github.com/shengyongjiang/ohmycheatsheet/internal/resolver"
	"github.com/shengyongjiang/ohmycheatsheet/internal/source"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review [command]",
	Short: "Review entries marked as needs-review",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}

type reviewItem struct {
	pageKey string
	entry   model.Entry
	state   model.EntryState
}

type reviewModel struct {
	items    []reviewItem
	store    store.StateStore
	cursor   int
	revealed bool
	width    int
	height   int
	dirty    bool
}

func (m reviewModel) Init() bubbletea.Cmd { return nil }

func (m reviewModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case bubbletea.KeyMsg:
		switch {
		case matchReviewKey(msg, "q", "esc", "ctrl+c"):
			if m.dirty {
				m.store.Save()
			}
			return m, bubbletea.Quit

		case matchReviewKey(msg, " "):
			m.revealed = true

		case matchReviewKey(msg, "tab"):
			if m.revealed {
				item := &m.items[m.cursor]
				item.state.State = item.state.State.Next()
				item.state.LastReviewed = time.Now()
				item.state.ReviewCount++
				m.store.SetEntryState(item.pageKey, item.entry.Index, item.state)
				m.dirty = true
			}

		case matchReviewKey(msg, "n", "j", "down"):
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.revealed = false
			}

		case matchReviewKey(msg, "p", "k", "up"):
			if m.cursor > 0 {
				m.cursor--
				m.revealed = false
			}
		}
	}
	return m, nil
}

func matchReviewKey(msg bubbletea.KeyMsg, keys ...string) bool {
	for _, k := range keys {
		if msg.String() == k {
			return true
		}
	}
	return false
}

var (
	reviewDimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	reviewCmdStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Padding(0, 2)
	reviewDescStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	reviewBarStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(lipgloss.Color("240")).Foreground(lipgloss.Color("240")).Padding(0, 1)
	reviewTagStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	reviewRemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	reviewCountStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

func (m reviewModel) View() string {
	var b strings.Builder

	if len(m.items) == 0 {
		return reviewDimStyle.Render("  No entries to review.\n")
	}

	item := m.items[m.cursor]
	counter := reviewCountStyle.Render(fmt.Sprintf("Review (%d of %d)", m.cursor+1, len(m.items)))

	header := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Render(counter)
	b.WriteString(header + "\n\n")

	b.WriteString(fmt.Sprintf("  %s:\n", reviewDimStyle.Render(item.pageKey)))
	b.WriteString(fmt.Sprintf("  %s\n\n", reviewDescStyle.Render(item.entry.Description)))

	if m.revealed {
		b.WriteString(fmt.Sprintf("  %s\n\n", reviewCmdStyle.Render(item.entry.Command)))

		var stateTag string
		switch item.state.State {
		case model.StateNeedsReview:
			stateTag = reviewTagStyle.Render("★ needs review")
		case model.StateRemembered:
			stateTag = reviewRemStyle.Render("✓ remembered")
		default:
			stateTag = reviewDimStyle.Render("○ not remembered")
		}
		b.WriteString(fmt.Sprintf("  %s\n", stateTag))
	} else {
		b.WriteString(reviewDimStyle.Render("  [press space to reveal]") + "\n")
	}
	b.WriteString("\n")

	var statusText string
	if m.revealed {
		statusText = "Tab:cycle state  j/k:navigate  q:quit"
	} else {
		statusText = "space:reveal  j/k:navigate  q:quit"
	}
	status := reviewBarStyle.Width(m.width).Render(reviewDimStyle.Render(statusText))
	b.WriteString(status)

	return b.String()
}

func runReview(cmd *cobra.Command, args []string) error {
	cfgPath := flagConfigPath
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	st, err := store.NewJSONStore(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	src := source.NewCheatshSource(cfg.CacheDir)
	res := resolver.New(src)
	var items []reviewItem

	var pagesToReview []string
	if len(args) > 0 {
		pagesToReview = []string{args[0]}
	} else {
		pagesToReview = st.ListTrackedPages()
	}

	for _, pageKey := range pagesToReview {
		states := st.GetPageStates(pageKey)
		page, err := res.Resolve(pageKey)
		if err != nil {
			continue
		}
		for _, entry := range page.Entries {
			es, ok := states[entry.Index]
			if ok && es.State == model.StateNeedsReview {
				items = append(items, reviewItem{
					pageKey: pageKey,
					entry:   entry,
					state:   es,
				})
			}
		}
	}

	if len(items) == 0 {
		fmt.Println("No entries to review.")
		return nil
	}

	m := reviewModel{items: items, store: st}
	p := bubbletea.NewProgram(m, bubbletea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	return nil
}
