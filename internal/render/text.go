package render

import (
	"fmt"
	"strings"

	"github.com/shengyongjiang/ocheetsheet/internal/model"
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

func RenderText(page *model.Page, states map[int]model.EntryState, showAll bool, noColor bool) string {
	var b strings.Builder

	title := page.Name
	desc := page.Description
	if !noColor {
		title = bold + white + title + reset
	}
	b.WriteString(fmt.Sprintf("\n  %s\n\n", title))
	b.WriteString(fmt.Sprintf("  %s\n", desc))

	if len(page.SeeAlso) > 0 {
		b.WriteString(fmt.Sprintf("  See also: %s\n", strings.Join(page.SeeAlso, ", ")))
	}
	if page.URL != "" {
		b.WriteString(fmt.Sprintf("  %s\n", page.URL))
	}
	b.WriteString("\n")

	hiddenCount := 0
	for _, entry := range page.Entries {
		es, hasState := states[entry.Index]
		state := model.StateNotRemembered
		if hasState {
			state = es.State
		}

		if state == model.StateRemembered && !showAll {
			hiddenCount++
			continue
		}

		descLine := entry.Description
		cmdLine := entry.Command

		switch {
		case state == model.StateRemembered && showAll:
			if noColor {
				descLine = fmt.Sprintf("  %s  [remembered]", descLine)
			} else {
				descLine = fmt.Sprintf("  %s%s  [remembered]%s", dim, descLine, reset)
				cmdLine = fmt.Sprintf("%s%s%s", dim, cmdLine, reset)
			}
		case state == model.StateNeedsReview:
			if noColor {
				descLine = fmt.Sprintf("  ★ %s  [needs review]", descLine)
			} else {
				descLine = fmt.Sprintf("  %s%s★ %s  [needs review]%s", bold, yellow, descLine, reset)
				cmdLine = fmt.Sprintf("%s%s%s", cyan, cmdLine, reset)
			}
		default:
			descLine = fmt.Sprintf("  %s", descLine)
		}

		b.WriteString(fmt.Sprintf("%s:\n", descLine))
		b.WriteString(fmt.Sprintf("    %s\n\n", cmdLine))
	}

	if hiddenCount > 0 {
		notice := fmt.Sprintf("  (%d remembered entries hidden. Use --all to show.)", hiddenCount)
		if !noColor {
			notice = fmt.Sprintf("  %s%s(%d remembered entries hidden. Use --all to show.)%s", dim, italic, hiddenCount, reset)
		}
		b.WriteString(notice + "\n")
	}

	return b.String()
}
