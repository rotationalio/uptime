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
