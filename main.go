package main

import (
	//"fmt"
	"os"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

const readme = `Check plugin for the AKCP Sensorprobe X plus`

func main () {
	defer check.CatchPanic()

    plugin := check.NewConfig()
    plugin.Name = "check_akcp"
    plugin.Readme = readme
    plugin.Timeout = 30
    plugin.Version = "0.1"

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

	/*
	for i, test := range overall.Outputs {
		fmt.Printf("%d: %s", i, test + "\n")
	}
	*/

	check.Exit(overall.GetStatus(), overall.GetOutput())
}
