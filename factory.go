// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fluentdexporter // import "go.opentelemetry.io/collector/exporter/loggingexporter"

import (
	"context"

	"github.com/fluent/fluent-logger-golang/fluent"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr                   = "fluentd"
	defaultSamplingInitial    = 2
	defaultSamplingThereafter = 500
)

// NewFactory creates a factory for Logging exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, component.StabilityLevelDevelopment),
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Tag:  "app",
		Host: "localhost",
		Port: 24224,
	}
}

func createTracesExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Traces, error) {
	cfg := config.(*Config)
	fluent, err := createFluent(cfg)
	if err != nil {
		return nil, err
	}
	s := newFluentdExporter(cfg.Tag, fluent)
	return exporterhelper.NewTracesExporter(ctx, set, cfg,
		s.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(fluentdSync(fluent)),
	)
}

func createMetricsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Metrics, error) {
	cfg := config.(*Config)
	fluent, err := createFluent(cfg)
	if err != nil {
		return nil, err
	}
	s := newFluentdExporter(cfg.Tag, fluent)
	return exporterhelper.NewMetricsExporter(ctx, set, cfg,
		s.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(fluentdSync(fluent)),
	)
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	fluent, err := createFluent(cfg)
	if err != nil {
		return nil, err
	}
	s := newFluentdExporter(cfg.Tag, fluent)
	return exporterhelper.NewLogsExporter(ctx, set, cfg,
		s.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(fluentdSync(fluent)),
	)
}

func createFluent(cfg *Config) (*fluent.Fluent, error) {
	return fluent.New(fluent.Config{
		FluentHost:   cfg.Host,
		FluentPort:   cfg.Port,
		Timeout:      cfg.Timeout,
		WriteTimeout: cfg.WriteTimeout,
		BufferLimit:  cfg.BufferLimit,
		RetryWait:    cfg.RetryWait,
		MaxRetry:     cfg.MaxRetry,
		MaxRetryWait: cfg.MaxRetryWait,
		TagPrefix:    cfg.TagPrefix,
		Async:        cfg.Async,
	})
}
