package config

import (
	"time"
)

// HealthCheck
type HealthCheck struct {
	Path     string        `yaml:"path"`
	Interval time.Duration `yaml:"interval"`
}

// Backend
type Backend struct {
	Host   string      `yaml:"host"`
	Weight uint8       `yaml:"weight"`
	Health HealthCheck `yaml:"health"`
}

// Config
type Config struct {
	Port    uint16 `yaml:"port"`
	Retries uint8  `yaml:"retries"`

	// ignored if `health` is in backend spec
	Health HealthCheck `yaml:"health"`

	Backends []Backend `yaml:"backends"`
}

// Load
func Load() (*Config, error) {

	// default configurations
	cfg := &Config{
		Port:    3000,
		Retries: 2,
		Health: HealthCheck{
			Path: "/",
		},
	}

	return cfg, nil
}
