package config

import (
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

// HealthCheck
type HealthCheck struct {
	Path     string        `yaml:"path"`
	Interval time.Duration `yaml:"interval"`
}

// Backend
type Backend struct {
	Host   string `yaml:"host"`
	Weight int32  `yaml:"weight"`
}

// Config
type Config struct {
	Port    string `yaml:"port"`
	Retries int32  `yaml:"retries"`

	// ignored if `health` is in backend spec
	Health HealthCheck `yaml:"health"`

	Backends []Backend `yaml:"backends"`
}

// Load
func Load(r io.Reader) (*Config, error) {

	// default configurations
	cfg := &Config{
		Port:    "3000",
		Retries: 1,
		Health: HealthCheck{
			Path:     "/",
			Interval: 30 * time.Second,
		},
	}

	data, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)

	return cfg, err
}
