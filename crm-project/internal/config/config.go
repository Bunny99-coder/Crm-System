// File: internal/config/config.go
package config

import (
	"os"
	"gopkg.in/yaml.v3"
	"log/slog"
)

// Config struct matches the structure of our config.yml file
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	Auth struct { // <-- ADD THIS
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"auth"`
}

// Load reads the config.yml file and returns a Config struct
func Load(path string, logger *slog.Logger) (*Config, error) {
	logger.Info("loading configuration", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}