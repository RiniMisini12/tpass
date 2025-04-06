package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Security SecurityConfig `toml:"security"`
}

type SecurityConfig struct {
	UID          string `toml:"uid"`
	GID          string `toml:"gid"`
	Username     string `toml:"username"`
	MainPassword string `toml:"password"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}

func GetOSSecurity() (SecurityConfig, error) {
	current, err := user.Current()
	if err != nil {
		return SecurityConfig{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return SecurityConfig{
		UID:      current.Uid,
		GID:      current.Gid,
		Username: current.Username,
	}, nil
}
