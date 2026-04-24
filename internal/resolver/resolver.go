package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Resolver struct {
	cachePath     string
	platformOrder []string
}

func New(cachePath string, platformOrder []string) *Resolver {
	return &Resolver{
		cachePath:     cachePath,
		platformOrder: platformOrder,
	}
}

func NewDefault(cachePath string) *Resolver {
	return New(cachePath, defaultPlatformOrder())
}

func defaultPlatformOrder() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"osx", "common"}
	case "linux":
		return []string{"linux", "common"}
	case "windows":
		return []string{"windows", "common"}
	case "freebsd":
		return []string{"freebsd", "common"}
	default:
		return []string{"common"}
	}
}

func (r *Resolver) Resolve(command string) (string, error) {
	filename := command + ".md"
	for _, platform := range r.platformOrder {
		path := filepath.Join(r.cachePath, platform, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("command %q not found in tldr cache", command)
}

func (r *Resolver) ListAllCommands() ([]string, error) {
	seen := make(map[string]bool)
	var commands []string
	for _, platform := range r.platformOrder {
		dir := filepath.Join(r.cachePath, platform)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".md")
			if !seen[name] {
				seen[name] = true
				commands = append(commands, name)
			}
		}
	}
	return commands, nil
}
