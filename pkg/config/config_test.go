package config_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/confire/contest"
	"go.rtnl.ai/uptime/pkg/config"
	"go.rtnl.ai/x/rlog"
)

var testEnv = contest.Env{
	"UPTIME_MAINTENANCE":       "true",
	"UPTIME_MODE":              "debug",
	"UPTIME_LOG_LEVEL":         "debug",
	"UPTIME_CONSOLE_LOG":       "true",
	"UPTIME_BIND_ADDR":         ":8888",
	"UPTIME_ALLOWED_ORIGINS":   "http://localhost:8888",
	"UPTIME_STATIC_SERVE":      "true",
	"UPTIME_STATIC_ROOT":       "/tmp",
	"UPTIME_STATIC_URL":        "/assets",
	"UPTIME_TELEMETRY_ENABLED": "false",
	"OTEL_SERVICE_NAME":        "uptime",
	"GIMLET_OTEL_SERVICE_ADDR": "localhost:8888",
}

var validConfig = config.Config{
	Maintenance:    true,
	Mode:           "debug",
	LogLevel:       rlog.LevelDecoder(slog.LevelDebug),
	ConsoleLog:     true,
	BindAddr:       ":8888",
	AllowedOrigins: []string{"http://localhost:8888"},
	Static: config.StaticConfig{
		Serve: true,
		Root:  "/tmp",
		URL:   "/assets",
	},
	Telemetry: config.TelemetryConfig{
		Enabled:     false,
		ServiceName: "uptime",
		ServiceAddr: "localhost:8888",
	},
}

func TestConfig(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		t.Cleanup(testEnv.Set())

		conf, err := config.New()
		require.NoError(t, err)
		require.Equal(t, validConfig, *conf)
	})

	t.Run("InvalidMode", func(t *testing.T) {
		t.Cleanup(testEnv.Set())

		invalid := contest.Env{
			"UPTIME_MODE": "invalid",
		}

		errs := map[string]string{
			"UPTIME_MODE": "invalid configuration: mode gin mode must be one of: release, debug, test",
		}

		for key := range invalid {
			cleanup := invalid.Set(key)

			_, err := config.New()
			require.EqualError(t, err, errs[key])

			cleanup()
		}

	})
}

func TestStaticConfigValidate(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		require.NoError(t, validConfig.Static.Validate(), "valid static config should pass validation")
	})

	t.Run("Invalid", func(t *testing.T) {
		tests := []struct {
			name string
			conf config.StaticConfig
			err  string
		}{
			{
				"root is required",
				config.StaticConfig{
					Serve: true,
					URL:   "/static",
				},
				"invalid configuration: static.root is required but not set",
			},
			{
				"root does not exist",
				config.StaticConfig{
					Serve: true,
					Root:  "/var/lib/www/not-a-directory",
					URL:   "/static",
				},
				"invalid configuration: static.root directory does not exist",
			},
			{
				"url is required",
				config.StaticConfig{
					Serve: true,
					Root:  "../web/static",
				},
				"invalid configuration: static.url is required but not set",
			},
			{
				"url is not a valid URL",
				config.StaticConfig{
					Serve: false,
					Root:  "../web/static",
					URL:   "://not-a-url",
				},
				"invalid configuration: static.url must be a valid URL with scheme and host",
			},
			{
				"url is not an absolute path starting with a slash",
				config.StaticConfig{
					Serve: true,
					Root:  "../web/static",
					URL:   "not-a-slash",
				},
				"invalid configuration: static.url must be a valid URL or an absolute path starting with a slash",
			},
			{
				"url is a remote URL and serve is true",
				config.StaticConfig{
					Serve: true,
					Root:  "../web/static",
					URL:   "https://example.com/static",
				},
				"invalid configuration: static.url cannot use a remote URL if static files are served from the filesystem",
			},
			{
				"url is a absolute URL and serve is false",
				config.StaticConfig{
					Serve: false,
					Root:  "../web/static",
					URL:   "/static",
				},
				"invalid configuration: static.url must be a remote url if static files are not served from the filesystem",
			},
		}

		for _, test := range tests {
			err := test.conf.Validate()
			require.EqualError(t, err, test.err, "expected static config validation error on test case %q", test.name)
		}
	})
}
