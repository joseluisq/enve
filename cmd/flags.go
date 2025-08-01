package cmd

import (
	cli "github.com/joseluisq/cline"
)

var Flags = []cli.Flag{
	cli.FlagString{
		Name:    "file",
		Aliases: []string{"f"},
		Value:   ".env",
		Summary: "Load environment variables from a file path (optional)",
	},
	cli.FlagString{
		Name:    "output",
		Aliases: []string{"o"},
		Value:   "text",
		Summary: "Output environment variables using text, json or xml format",
	},
	cli.FlagBool{
		Name:    "overwrite",
		Aliases: []string{"w"},
		Value:   false,
		Summary: "Overwrite environment variables if already set",
	},
	cli.FlagString{
		Name:    "chdir",
		Aliases: []string{"c"},
		Summary: "Change currrent working directory",
	},
	cli.FlagBool{
		Name:    "new-environment",
		Aliases: []string{"n"},
		Value:   false,
		Summary: "Start a new environment with only variables from the .env file or stdin",
	},
	cli.FlagBool{
		Name:    "ignore-environment",
		Aliases: []string{"i"},
		Value:   false,
		Summary: "Starts with an empty environment, ignoring any existing environment variables",
	},
	cli.FlagBool{
		Name:    "no-file",
		Aliases: []string{"z"},
		Value:   false,
		Summary: "Do not load a .env file",
	},
	cli.FlagBool{
		Name:    "stdin",
		Aliases: []string{"s"},
		Value:   false,
		Summary: "Read only environment variables from stdin and ignore the .env file",
	},
}
