package render

import (
	"fmt"
	"strings"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
	"github.com/shengyongjiang/ohmycheatsheet/internal/shuffle"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	italic = "\033[3m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	white  = "\033[37m"
)

func RenderText(page *model.Page, states map[int]model.EntryState, showAll bool, noColor bool, seed int64) string {
	var b strings.Builder

	title := page.Name
	desc := page.Description
	if !noColor {
		title = bold + white + title + reset
	}
	b.WriteString(fmt.Sprintf("\n  %s\n\n", title))
	if desc != "" {
		b.WriteString(fmt.Sprintf("  %s\n", desc))
	}

	if page.URL != "" {
		b.WriteString(fmt.Sprintf("  %s\n", page.URL))
	}
	b.WriteString("\n")

	shuffled := shuffle.ShuffleEntries(page.Entries, seed)

	displayLimit := 10
	displayed := 0
	hiddenCount := 0
	for _, entry := range shuffled {
		es, hasState := states[entry.Index]
		state := model.StateNotRemembered
		if hasState {
			state = es.State
		}
		if state == model.StateRemembered && !showAll {
			hiddenCount++
			continue
		}
		if !showAll && displayed >= displayLimit {
			continue
		}
		renderEntry(&b, entry.Description, entry.Command, state, "", noColor)
		displayed++
	}

	if hiddenCount > 0 {
		notice := fmt.Sprintf("  (%d remembered entries hidden. Use --all to show.)", hiddenCount)
		if !noColor {
			notice = fmt.Sprintf("  %s%s(%d remembered entries hidden. Use --all to show.)%s", dim, italic, hiddenCount, reset)
		}
		b.WriteString(notice + "\n")
	}

	remaining := len(page.Entries) - hiddenCount - displayed
	if !showAll && remaining > 0 {
		more := fmt.Sprintf("  (%d more entries available. Use -i for interactive mode.)", remaining)
		if !noColor {
			more = fmt.Sprintf("  %s%s(%d more entries available. Use -i for interactive mode.)%s", dim, italic, remaining, reset)
		}
		b.WriteString(more + "\n")
	}

	return b.String()
}

func renderEntry(b *strings.Builder, desc, cmd string, state model.MemoryState, fromPage string, noColor bool) {
	descLine := desc
	cmdLine := cmd

	if fromPage != "" {
		prefix := fmt.Sprintf("[%s] ", fromPage)
		if noColor {
			descLine = prefix + descLine
		} else {
			descLine = italic + prefix + reset + descLine
		}
	}

	switch state {
	case model.StateRemembered:
		if noColor {
			descLine = fmt.Sprintf("  %s  [remembered]", descLine)
		} else {
			descLine = fmt.Sprintf("  %s%s  [remembered]%s", dim, descLine, reset)
			cmdLine = fmt.Sprintf("%s%s%s", dim, cmdLine, reset)
		}
	case model.StateNeedsReview:
		if noColor {
			descLine = fmt.Sprintf("  + %s  [needs review]", descLine)
		} else {
			descLine = fmt.Sprintf("  %s%s+ %s  [needs review]%s", bold, yellow, descLine, reset)
			cmdLine = fmt.Sprintf("%s%s%s", cyan, cmdLine, reset)
		}
	default:
		descLine = fmt.Sprintf("  %s", descLine)
	}

	b.WriteString(fmt.Sprintf("%s:\n", descLine))
	b.WriteString(fmt.Sprintf("    %s\n\n", cmdLine))
}
