package cmd

import (
	"fmt"
	"os"

	cli "github.com/joseluisq/cline"
)

// Build-time application values
var (
	versionNumber string = "devel"
	buildTime     string
	buildCommit   string
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := cli.New()
	app.Name = "enve"
	app.Summary = "Run a program in a modified environment providing an optional .env file or variables from stdin"
	app.Version = versionNumber
	app.BuildTime = buildTime
	app.BuildCommit = buildCommit
	app.Flags = Flags
	app.Handler = appHandler

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
