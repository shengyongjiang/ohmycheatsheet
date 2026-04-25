package parser

import (
	"strings"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
)

func ParseCheatsh(content string, command string) (*model.Page, error) {
	lines := strings.Split(content, "\n")
	page := &model.Page{
		Name: command,
	}

	var descParts []string
	var entries []model.Entry
	entryIndex := 0

	flush := func() {
		descParts = nil
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "#[") {
			flush()
			continue
		}

		if strings.TrimSpace(line) == "" {
			if len(descParts) > 0 {
				flush()
			}
			continue
		}

		if strings.HasPrefix(line, "# ") || line == "#" {
			comment := strings.TrimPrefix(line, "# ")
			comment = strings.TrimPrefix(comment, "#")
			comment = strings.TrimSpace(comment)
			if comment != "" {
				descParts = append(descParts, comment)
			}
			continue
		}

		if len(descParts) == 0 {
			continue
		}

		var cmdParts []string
		for i < len(lines) {
			l := lines[i]
			if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "#") {
				break
			}
			cmdParts = append(cmdParts, l)
			i++
		}
		i--

		desc := strings.TrimSuffix(strings.Join(descParts, " "), ":")
		cmd := strings.Join(cmdParts, "\n")

		entries = append(entries, model.Entry{
			Index:       entryIndex,
			Description: desc,
			Command:     cmd,
			Fingerprint: model.ComputeFingerprint(desc, cmd),
		})
		entryIndex++
		descParts = nil
	}

	if len(entries) > 0 {
		first := entries[0]
		if strings.EqualFold(first.Description, command) || strings.HasPrefix(strings.ToLower(first.Description), strings.ToLower(command)+" ") {
			if len(entries) > 1 {
				page.Description = first.Description
				entries = entries[1:]
				for i := range entries {
					entries[i].Index = i
				}
			}
		}
	}

	page.Entries = entries
	return page, nil
}
