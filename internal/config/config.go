// Package config handles loading and parsing the gateway configuration from YAML files.
package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main gateway configuration structure.
type Config struct {
	Server ServerConfig `yaml:"server"`
	CORS   CORSConfig   `yaml:"cors"`
	Routes []Route      `yaml:"routes"`
}

// ServerConfig defines server-specific settings like host, port, and timeouts.
type ServerConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// CORSConfig defines Cross-Origin Resource Sharing settings.
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// Route defines a single routing rule for proxying requests.
type Route struct {
	Path    string   `yaml:"path"`
	Prefix  bool     `yaml:"prefix"`
	Target  string   `yaml:"target"`
	Methods []string `yaml:"methods"`
}

// Load reads and parses the configuration file from the given path.
func Load(configPath string) (*Config, error) {
	// #nosec G304 - configPath is controlled by application environment variables
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
