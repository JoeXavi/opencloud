package defaults

import (
	"github.com/opencloud-eu/opencloud/services/filepostprocessing/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := &config.Config{}
	EnsureDefaults(cfg)
	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "filepostprocessing",
		},
		Events: config.Events{
			Endpoint: "127.0.0.1:9233",
			Cluster:  "opencloud-cluster",
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if cfg.LogFilePath == "" {
		cfg.LogFilePath = "/tmp/file_events.json"
	}

	if cfg.Service.Name == "" {
		cfg.Service.Name = "filepostprocessing"
	}

	if cfg.Events.Endpoint == "" {
		cfg.Events.Endpoint = "127.0.0.1:9233"
	}
	if cfg.Events.Cluster == "" {
		cfg.Events.Cluster = "opencloud-cluster"
	}
}

func Sanitize(cfg *config.Config) {
	// Add sanitization if needed
}
