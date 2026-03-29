package config

import (
	"sync"

	"github.com/gin-gonic/gin"
	"go.rtnl.ai/confire"
	"go.rtnl.ai/x/rlog"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be UPTIME_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "uptime"

// Config contains all of the configuration parameters for the Uptime server which
// are loaded from the environment and should be validated before use.
type Config struct {
	Maintenance    bool              `default:"false" desc:"if true the server will start in maintenance mode"`
	Mode           string            `default:"release" desc:"specify the mode of the gin server (release, debug, testing)"`
	LogLevel       rlog.LevelDecoder `default:"info" split_words:"true" desc:"set the log level for the server"`
	ConsoleLog     bool              `default:"false" split_words:"true" desc:"if true the server will log to the console in text format"`
	BindAddr       string            `default:":8000" split_words:"true" desc:"the ip address and port to bind the server to"`
	AllowedOrigins []string          `split_words:"true" default:"http://localhost:8000" desc:"a list of allowed origins for CORS"`
	Telemetry      TelemetryConfig
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

// New creates a new Config instance and loads the configuration from the environment,
// validating the configuration and returning an error if the configuration is invalid
// or could not be parsed from environment variables.
//
// NOTE: New should only be used for testing, for module access to the config use Get().
func New() (conf *Config, err error) {
	// NOTE: confire.Process calls Validate() internally.
	conf = &Config{}
	if err = confire.Process(Prefix, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c Config) Validate() (err error) {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return confire.Invalid("", "mode", "gin mode must be one of: release, debug, test")
	}
	return nil
}

//============================================================================
// Config Package Management
//============================================================================

var (
	mu   sync.RWMutex
	err  error
	load sync.Once
	conf *Config
)

func Get() (Config, error) {
	load.Do(func() {
		mu.Lock()
		defer mu.Unlock()

		if conf == nil {
			conf, err = New()
		}
	})
	mu.RLock()
	defer mu.RUnlock()
	if conf != nil {
		return *conf, err
	}
	return Config{}, err
}

func MustGet() Config {
	conf, err := Get()
	if err != nil {
		panic(err)
	}
	return conf
}

func Set(c Config) {
	mu.Lock()
	defer mu.Unlock()
	conf = &c
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	conf = nil
	err = nil
	load = sync.Once{}
}
