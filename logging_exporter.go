package fluentdexporter

import (
	"context"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/senomas/otel_fluentd_exporter/internal/otlptext"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type fluentdExporter struct {
	tagTrace         string
	tagMetric        string
	tagLog           string
	fluent           *fluent.Fluent
	logsMarshaler    plog.Marshaler
	metricsMarshaler pmetric.Marshaler
	tracesMarshaler  ptrace.Marshaler
}

func (s *fluentdExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	data := map[string]interface{}{}
	data["resource spans"] = td.ResourceSpans().Len()
	data["spans"] = td.SpanCount()

	buf, err := s.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	data["traces"] = string(buf)
	err = s.fluent.Post(s.tagTrace, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *fluentdExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	data := map[string]interface{}{}
	data["resource metrics"] = md.ResourceMetrics().Len()
	data["metrics"] = md.MetricCount()
	data["data points"] = md.DataPointCount()
	buf, err := s.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	data["traces"] = string(buf)
	err = s.fluent.Post(s.tagMetric, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *fluentdExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	data := map[string]interface{}{}
	data["resource logs"] = ld.ResourceLogs().Len()
	data["log records"] = ld.LogRecordCount()
	buf, err := s.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	data["traces"] = string(buf)
	err = s.fluent.Post(s.tagLog, data)
	if err != nil {
		return err
	}
	return nil
}

func newFluentdExporter(tag string, fluent *fluent.Fluent) *fluentdExporter {
	return &fluentdExporter{
		tagTrace:         tag + ".trace",
		tagMetric:        tag + ".metric",
		tagLog:           tag + ".log",
		fluent:           fluent,
		logsMarshaler:    otlptext.NewTextLogsMarshaler(),
		metricsMarshaler: otlptext.NewTextMetricsMarshaler(),
		tracesMarshaler:  otlptext.NewTextTracesMarshaler(),
	}
}

func fluentdSync(fluent *fluent.Fluent) func(context.Context) error {
	return func(context.Context) error {
		return nil
	}
}
