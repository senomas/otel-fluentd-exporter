// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fluentdexporter // import "go.opentelemetry.io/collector/exporter/loggingexporter"

import "time" // Config defines configuration for logging exporter.
type Config struct {
	Tag          string        `mapstructure:"tag,omitempty"`
	Host         string        `mapstructure:"host,omitempty"`
	Port         int           `mapstructure:"port,omitempty"`
	Timeout      time.Duration `mapstructure:"timeout,omitempty"`
	WriteTimeout time.Duration `mapstructure:"write_timeout,omitempty"`
	BufferLimit  int           `mapstructure:"buffer_limit,omitempty"`
	RetryWait    int           `mapstructure:"retry_wait,omitempty"`
	MaxRetry     int           `mapstructure:"max_retry,omitempty"`
	MaxRetryWait int           `mapstructure:"max_retry_wait,omitempty"`
	TagPrefix    string        `mapstructure:"tag_prefix,omitempty"`
	Async        bool          `mapstructure:"async,omitempty"`
}
