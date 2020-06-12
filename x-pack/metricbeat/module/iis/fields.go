// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package iis

import (
	"github.com/codragonzuo/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "iis", asset.ModuleFieldsPri, AssetIis); err != nil {
		panic(err)
	}
}

// AssetIis returns asset data.
// This is the base64 encoded gzipped contents of module/iis.
func AssetIis() string {
	return "eJzMkrFuhDAQRHt/xYgS6e4DXKS/Oh+ADOydNmewZfty4u8jx0AAgUUZFxSzo5m3aC940iDB7AUQOGiSKG63z0IALfnGsQ1seokPASD60Jn2pUkAjjQpTxI1BSUEcGfSrZe/zgt61dGUHF8YLEk8nHnZUdkpWIcsg5S1mhsVzZU1Rs+GveT41nh/+m5repmOLdYSLX5XgwnpScPbuHYzywBsIBAh1vFT55tqT+6b3LW8lie2TkCm/qImLOQkVGl610YdDKtOWcv9Y3QWZXHuj86Ys7q7CQf6/3twoOxBHJ7D8TFkT2HsTHk/AQAA//+gjPaq"
}
