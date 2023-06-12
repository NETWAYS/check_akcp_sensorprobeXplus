package main

import (
	"fmt"
	"os"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

const readme = `Check plugin for the AKCP Sensorprobe X plus`

const license = `
Copyright (C) 2023 NETWAYS GmbH <info@netways.de>
`

// nolint: gochecknoglobals
var (
	// These get filled at build time with the proper vaules
	version = "development"
	commit  = "HEAD"
	date    = "latest"
)

//goland:noinspection GoBoolExpressions
func buildVersion() string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\ndate: %s", result, date)
	}

	result += "\n" + license

	return result
}

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	version := buildVersion()
	plugin.Name = "check_akcp_sensorprobeXplus"
	plugin.Version = version
	plugin.Readme = readme
	plugin.Timeout = 30

	config := &Config{}
	config.BindArguments(plugin.FlagSet)

	plugin.ParseArguments()

	if len(os.Args) <= 1 {
		plugin.FlagSet.Usage()
		check.Exit(check.Unknown, "No arguments given")
	}

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	plugin.SetupLogging()

	var overall result.Overall

	err = config.Run(&overall)
	if err != nil {
		check.ExitError(err)
	}

	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
}
