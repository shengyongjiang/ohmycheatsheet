package source

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const cacheTTL = 7 * 24 * time.Hour

type CheatshSource struct {
	cacheDir string
}

func NewCheatshSource(cacheDir string) *CheatshSource {
	return &CheatshSource{cacheDir: cacheDir}
}

func (s *CheatshSource) Fetch(command string) (string, error) {
	cachePath := filepath.Join(s.cacheDir, command+".txt")

	if info, err := os.Stat(cachePath); err == nil {
		if time.Since(info.ModTime()) < cacheTTL {
			data, err := os.ReadFile(cachePath)
			if err == nil {
				return string(data), nil
			}
		}
	}

	url := fmt.Sprintf("https://cheat.sh/%s?T", command)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "omcs/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cheat.sh returned %d for %q", resp.StatusCode, command)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	content := string(body)

	if err := os.MkdirAll(s.cacheDir, 0o755); err != nil {
		return content, nil
	}
	os.WriteFile(cachePath, body, 0o644)

	return content, nil
}

func (s *CheatshSource) ListCachedCommands() []string {
	entries, err := os.ReadDir(s.cacheDir)
	if err != nil {
		return nil
	}
	var commands []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".txt")
		if !strings.Contains(name, "/:list") {
			commands = append(commands, name)
		}
	}
	return commands
}
