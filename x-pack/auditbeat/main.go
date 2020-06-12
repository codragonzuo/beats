// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"os"

	"github.com/codragonzuo/beats/x-pack/auditbeat/cmd"

	// Register modules.
	_ "github.com/codragonzuo/beats/auditbeat/module/auditd"
	_ "github.com/codragonzuo/beats/auditbeat/module/file_integrity"

	// Register includes.
	_ "github.com/codragonzuo/beats/auditbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
