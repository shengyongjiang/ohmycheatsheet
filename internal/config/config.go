package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	TldrCachePath string `json:"tldr_cache_path"`
	StateFile     string `json:"state_file"`
	ColorEnabled  bool   `json:"color_enabled"`
}

func Load(path string) (*Config, error) {
	cfg := defaults()
	if path == "" {
		return cfg, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func defaults() *Config {
	home, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = filepath.Join(home, ".config")
	}
	return &Config{
		TldrCachePath: filepath.Join(home, ".tldr", "cache", "pages"),
		StateFile:     filepath.Join(configDir, "ocs", "state.json"),
		ColorEnabled:  true,
	}
}

func DefaultConfigPath() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "ocs", "config.json")
}
