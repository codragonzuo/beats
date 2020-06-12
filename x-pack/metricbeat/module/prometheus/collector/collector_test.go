// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build !integration

package collector

import (
	"testing"

	mbtest "github.com/codragonzuo/beats/metricbeat/mb/testing"

	_ "github.com/codragonzuo/beats/x-pack/metricbeat/module/prometheus"

	// Import common fields for validation
	_ "github.com/codragonzuo/beats/metricbeat/module/prometheus"
)

func TestData(t *testing.T) {
	mbtest.TestDataFiles(t, "prometheus", "collector")
}
