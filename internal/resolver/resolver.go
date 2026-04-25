package resolver

import (
	"fmt"
	"strings"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
	"github.com/shengyongjiang/ohmycheatsheet/internal/parser"
	"github.com/shengyongjiang/ohmycheatsheet/internal/source"
)

type Resolver struct {
	source *source.CheatshSource
}

func New(src *source.CheatshSource) *Resolver {
	return &Resolver{source: src}
}

func (r *Resolver) Resolve(command string) (*model.Page, error) {
	content, err := r.source.Fetch(command)
	if err != nil {
		return nil, fmt.Errorf("fetch %q: %w", command, err)
	}
	page, err := parser.ParseCheatsh(content, command)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", command, err)
	}
	if len(page.Entries) == 0 {
		return nil, fmt.Errorf("command %q not found on cheat.sh", command)
	}
	return page, nil
}

func (r *Resolver) ListRelatedCommands(prefix string) ([]string, error) {
	cached := r.source.ListCachedCommands()
	var related []string
	dashPrefix := prefix + "-"
	for _, cmd := range cached {
		if cmd != prefix && strings.HasPrefix(cmd, dashPrefix) {
			related = append(related, cmd)
		}
	}
	return related, nil
}

func (r *Resolver) ListAllCommands() ([]string, error) {
	return r.source.ListCachedCommands(), nil
}
