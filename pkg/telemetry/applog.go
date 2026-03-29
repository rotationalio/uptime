package telemetry

import (
	"log/slog"
	"os"

	otelslog "go.opentelemetry.io/contrib/bridges/otelslog"

	"go.rtnl.ai/uptime/pkg/config"
	"go.rtnl.ai/x/rlog"
	"go.rtnl.ai/x/rlog/console"
)

// initializeLogging configures the process root logger: stdout (JSON or console) with rlog
// level names, optionally fan-out to OpenTelemetry logs via otelslog when telemetry is on.
func initializeLogging() error {
	// Get the configured log level, set it as the global log level, and merge
	// custom level string replacer.
	conf := config.MustGet()
	rlog.SetLevel(slog.Level(conf.LogLevel))
	opts := rlog.MergeWithCustomLevels(rlog.WithGlobalLevel(nil))

	// Create the stdout handler: text or JSON.
	var stdout slog.Handler
	if conf.ConsoleLog {
		stdout = console.New(os.Stdout, &console.Options{HandlerOptions: opts})
	} else {
		stdout = slog.NewJSONHandler(os.Stdout, opts)
	}

	// If telemetry is enabled, fan-out to OpenTelemetry logs via otelslog.
	var root slog.Handler = stdout
	if !disabled && upLoggerProvider != nil {
		otelHandler := otelslog.NewHandler(ServiceName())
		root = slog.NewMultiHandler(stdout, otelHandler)
	}

	// Create the root logger and set it as the default rlog logger.
	rlog.SetDefault(rlog.New(slog.New(root)))
	return nil
}
