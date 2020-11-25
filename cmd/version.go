package cmd

import (
	cli "github.com/joseluisq/cline"
)

// application version values
var (
	versionNumber string = "devel"
	buildTime     string
)

// VersionFlag builds a new Version flag
func VersionFlag() *cli.FlagBool {
	return &cli.FlagBool{
		Name:    "version",
		Aliases: []string{"v"},
		Summary: "shows the current version",
	}
}
