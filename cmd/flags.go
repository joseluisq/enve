package cmd

import (
	"github.com/joseluisq/cline/flag"
)

var Flags = []flag.Flag{
	flag.FlagString{
		Name:    "file",
		Aliases: []string{"f"},
		Value:   ".env",
		Summary: "Load environment variables from a file path (optional)",
	},
	flag.FlagString{
		Name:    "output",
		Aliases: []string{"o"},
		Value:   "text",
		Summary: "Output environment variables using text, json or xml format",
	},
	flag.FlagBool{
		Name:    "overwrite",
		Aliases: []string{"w"},
		Value:   false,
		Summary: "Overwrite environment variables if already set",
	},
	flag.FlagString{
		Name:    "chdir",
		Aliases: []string{"c"},
		Summary: "Change currrent working directory",
	},
	flag.FlagBool{
		Name:    "new-environment",
		Aliases: []string{"n"},
		Value:   false,
		Summary: "Start a new environment with only variables from the .env file or stdin",
	},
	flag.FlagBool{
		Name:    "ignore-environment",
		Aliases: []string{"i"},
		Value:   false,
		Summary: "Starts with an empty environment, ignoring any existing environment variables",
	},
	flag.FlagBool{
		Name:    "no-file",
		Aliases: []string{"z"},
		Value:   false,
		Summary: "Do not load a .env file",
	},
	flag.FlagBool{
		Name:    "stdin",
		Aliases: []string{"s"},
		Value:   false,
		Summary: "Read only environment variables from stdin and ignore the .env file",
	},
}
