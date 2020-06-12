// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package decode_cef

import (
	"github.com/codragonzuo/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "decode_cef", asset.ModuleFieldsPri, AssetDecodeCef); err != nil {
		panic(err)
	}
}

// AssetDecodeCef returns asset data.
// This is the base64 encoded gzipped contents of processors/decode_cef.
func AssetDecodeCef() string {
	return "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"
}
