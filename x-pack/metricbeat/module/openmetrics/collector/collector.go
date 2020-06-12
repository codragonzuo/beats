// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package collector

import (
	"github.com/codragonzuo/beats/metricbeat/mb"
	"github.com/codragonzuo/beats/metricbeat/mb/parse"
	"github.com/codragonzuo/beats/metricbeat/module/prometheus/collector"
)

const (
	defaultScheme = "http"
	defaultPath   = "/metrics"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   defaultPath,
	}.Build()
)

func init() {
	mb.Registry.MustAddMetricSet("openmetrics", "collector",
		collector.MetricSetBuilder("openmetrics", collector.DefaultPromEventsGeneratorFactory),
		mb.WithHostParser(hostParser))
}
