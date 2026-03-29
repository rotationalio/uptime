package telemetry

import (
	"errors"

	"go.opentelemetry.io/otel/metric"
)

var (
	GroupsMonitored   metric.Int64Gauge     // The number of groups that are being actively monitored and their statuses (operational, degraded, partial, major, maintenance)
	ServicesMonitored metric.Int64Gauge     // The number of services that are being actively monitored and their statuses (operational, degraded, partial, major, maintenance)
	ServiceLatency    metric.Int64Histogram // The latency of the uptime checks for monitored services (in milliseconds)
)

func initializeMetrics(meter metric.Meter) (err error) {
	var merr error

	if GroupsMonitored, merr = meter.Int64Gauge("groups.monitored",
		metric.WithDescription("The number of groups that are being actively monitored"),
		metric.WithUnit("groups"),
	); merr != nil {
		err = errors.Join(err, merr)
	}

	if ServicesMonitored, merr = meter.Int64Gauge("services.monitored",
		metric.WithDescription("The number of services that are being actively monitored"),
		metric.WithUnit("services"),
	); merr != nil {
		err = errors.Join(err, merr)
	}

	if ServiceLatency, merr = meter.Int64Histogram("services.latency",
		metric.WithDescription("The latency of the uptime checks for monitored services (in milliseconds)"),
		metric.WithUnit("ms"),
	); merr != nil {
		err = errors.Join(err, merr)
	}

	return err
}
