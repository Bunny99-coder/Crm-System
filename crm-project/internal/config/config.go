// File: internal/config/config.go
package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/jackc/pgx/v5/stdlib"
)

func init() {
	sql.Register("postgres", &stdlib.Driver{})
}

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
	Roles struct {
		SalesAgentID int `yaml:"-"` // Not from YAML, populated from DB
		ReceptionID  int `yaml:"-"` // Not from YAML, populated from DB
	} `yaml:"-"`
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

	logger.Info("Database URL from config", "url", cfg.Database.URL)
	// Establish database connection to fetch role IDs
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Fetch Sales_Agent role ID
	err = db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", "Sales_Agent").Scan(&cfg.Roles.SalesAgentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Sales_Agent role ID: %w", err)
	}

	// Fetch Reception role ID
	err = db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", "Reception").Scan(&cfg.Roles.ReceptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Reception role ID: %w", err)
	}

	logger.Info("configuration loaded successfully",
		"SalesAgentRoleID", cfg.Roles.SalesAgentID,
		"ReceptionRoleID", cfg.Roles.ReceptionID)

	return &cfg, nil
}