// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fluentdexporter // import "go.opentelemetry.io/collector/exporter/loggingexporter"

// Config defines configuration for logging exporter.
type Config struct {
	Tag  string `mapstructure:"tag,omitempty"`
	Host string `mapstructure:"host,omitempty"`
	Port int    `mapstructure:"port,omitempty"`
}
