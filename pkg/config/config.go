package config

import (
	"github.com/rotationalio/confire"
	"go.rtnl.ai/gimlet/logger"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be UPTIME_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "uptime"

// Config contains all of the configuration parameters for the Uptime server which
// are loaded from the environment and should be validated before use.
type Config struct {
	Maintenance bool                `default:"false" desc:"if true the server will start in maintenance mode"`
	Mode        string              `default:"release" desc:"specify the mode of the gin server (release, debug, testing)"`
	LogLevel    logger.LevelDecoder `default:"info" desc:"set the log level for the server"`
	BindAddr    string              `default:":8000" desc:"the ip address and port to bind the server to"`
	Telemetry   TelemetryConfig
}

// Telemetry is primarily configured via the open telemetry sdk environment variables.
// As such there is no need to specify OTel specific configuration here. This config
// is used primarily to enable/disable telemetry and to set values for custom telemetry.
//
// See: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
// For the environment variables that can be used to configure telemetry.
//
// See Also: https://oneuptime.com/blog/post/2026-02-06-opentelemetry-environment-variables-zero-code/view
// For OpenTelemetry configuration best practices.
type TelemetryConfig struct {
	Enabled     bool   `default:"true" desc:"disable telemetry by setting this environment variable to false"`
	ServiceName string `split_words:"true" env:"OTEL_SERVICE_NAME" desc:"override the default name of the service, used for logging and telemetry"`
	ServiceAddr string `split_words:"true" env:"GIMLET_OTEL_SERVICE_ADDR" desc:"the primary server name if it is known. E.g. the server name directive in an Nginx config. Should include host identifier and port if it is used; empty if not known."`
}

func New() (conf *Config, err error) {
	conf = &Config{}
	if err = confire.Process(Prefix, conf); err != nil {
		return nil, err
	}

	if err = conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) Validate() (err error) {
	return nil
}
