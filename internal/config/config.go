package config

import (
	_ "embed"
	"fmt"
	"os/user"

	"github.com/pelletier/go-toml/v2"
)

//go:embed config.toml
var configFileData []byte

type Config struct {
	Security SecurityConfig `toml:"security"`
}

type SecurityConfig struct {
	UID          string `toml:"uid"`
	GID          string `toml:"gid"`
	Username     string `toml:"username"`
	MainPassword string `toml:"password"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := toml.Unmarshal(configFileData, &cfg); err != nil {
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
