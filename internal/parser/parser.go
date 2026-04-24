package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shengyongjiang/ocheetsheet/internal/model"
)

var seeAlsoRegexp = regexp.MustCompile("`([^`]+)`")

func ParseFile(path string) (*model.Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	platform := filepath.Base(filepath.Dir(path))
	return ParseString(string(data), name, platform)
}

func ParseString(content, name, platform string) (*model.Page, error) {
	lines := strings.Split(content, "\n")
	page := &model.Page{
		Name:     name,
		Platform: platform,
	}

	var descParts []string
	var i int

	for i = 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "# ") {
			continue
		}
		if strings.HasPrefix(line, "> ") {
			text := strings.TrimPrefix(line, "> ")
			if strings.HasPrefix(text, "See also:") {
				matches := seeAlsoRegexp.FindAllStringSubmatch(text, -1)
				for _, m := range matches {
					page.SeeAlso = append(page.SeeAlso, m[1])
				}
			} else if strings.HasPrefix(text, "More information:") {
				url := strings.TrimPrefix(text, "More information: ")
				url = strings.TrimSuffix(url, ".")
				url = strings.TrimPrefix(url, "<")
				url = strings.TrimSuffix(url, ">")
				page.URL = url
			} else {
				descParts = append(descParts, text)
			}
			continue
		}
		if line == "" && len(descParts) > 0 && page.Description == "" {
			page.Description = joinDescription(descParts)
		}
		if strings.HasPrefix(line, "- ") {
			break
		}
	}
	if page.Description == "" && len(descParts) > 0 {
		page.Description = joinDescription(descParts)
	}

	entryIndex := 0
	for i < len(lines) {
		line := lines[i]
		if !strings.HasPrefix(line, "- ") {
			i++
			continue
		}

		desc := strings.TrimPrefix(line, "- ")
		desc = strings.TrimSuffix(desc, ":")

		var cmd string
		for i++; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "`") && strings.HasSuffix(trimmed, "`") {
				cmd = trimmed[1 : len(trimmed)-1]
				i++
				break
			}
			break
		}

		if cmd != "" {
			page.Entries = append(page.Entries, model.Entry{
				Index:       entryIndex,
				Description: desc,
				Command:     cmd,
				Fingerprint: model.ComputeFingerprint(desc, cmd),
			})
			entryIndex++
		}
	}

	return page, nil
}

func joinDescription(parts []string) string {
	desc := strings.Join(parts, " ")
	if !strings.HasSuffix(desc, ".") {
		desc += "."
	}
	return desc
}
