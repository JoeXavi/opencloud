package config

import (
	"context"

	"github.com/opencloud-eu/opencloud/pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	LogLevel string `yaml:"loglevel" env:"OC_LOG_LEVEL;FILELOGGER_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"1.0.0"`

	LogFilePath string `yaml:"log_file_path" env:"FILELOGGER_LOG_FILE_PATH" desc:"The path to the JSON log file where events will be recorded." introductionVersion:"1.0.0"`

	Events Events `yaml:"events"`

	Context context.Context `yaml:"-"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OC_EVENTS_ENDPOINT" desc:"The address of the event system." introductionVersion:"1.0.0"`
	Cluster              string `yaml:"cluster" env:"OC_EVENTS_CLUSTER" desc:"The clusterID of the event system." introductionVersion:"1.0.0"`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OC_INSECURE;OC_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates." introductionVersion:"1.0.0"`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OC_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate." introductionVersion:"1.0.0"`
	EnableTLS            bool   `yaml:"enable_tls" env:"OC_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker." introductionVersion:"1.0.0"`
	AuthUsername         string `yaml:"username" env:"OC_EVENTS_AUTH_USERNAME" desc:"The username to authenticate with the events broker." introductionVersion:"1.0.0"`
	AuthPassword         string `yaml:"password" env:"OC_EVENTS_AUTH_PASSWORD" desc:"The password to authenticate with the events broker." introductionVersion:"1.0.0"`
}
