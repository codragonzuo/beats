package cmd

import (
	"{beat_path}/beater"

	cmd "github.com/codragonzuo/beats/libbeat/cmd"
	"github.com/codragonzuo/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "{beat}"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})
