package main

import (
	//"fmt"
	"os"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

const readme = `Check plugin for the AKCP Sensorprobe X plus`

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "check_akcp_sensorprobeXplus"
	plugin.Version = "1.0.1"
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

	// Modify output here until https://github.com/Icinga/icinga2/issues/9379 is fixed
	// TODO: Remove this dirty hack
	var output string
	var perfdata []string

	for _, partial := range overall.Outputs {
		tmp := strings.Split(partial, "|")
		output += "\n" + tmp[0]
		if len(tmp) > 1 {
			perfdata = append(perfdata, strings.TrimSpace(tmp[1]))
		}
	}

	output = overall.GetSummary() + "\n" + output + " | " + strings.Join(perfdata, " ")

	check.ExitRaw(overall.GetStatus(), output)
}
