package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	CacheDir     string `json:"cache_dir"`
	StateFile    string `json:"state_file"`
	ColorEnabled bool   `json:"color_enabled"`
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
	cacheDir, _ := os.UserCacheDir()
	if cacheDir == "" {
		cacheDir = filepath.Join(home, ".cache")
	}
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = filepath.Join(home, ".config")
	}
	return &Config{
		CacheDir:     filepath.Join(cacheDir, "omcs", "cheatsh"),
		StateFile:    filepath.Join(configDir, "omcs", "state.json"),
		ColorEnabled: true,
	}
}

func DefaultConfigPath() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "omcs", "config.json")
}
