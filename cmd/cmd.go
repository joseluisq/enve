package cmd

import (
	"fmt"
	"os"

	"github.com/joseluisq/cline/app"
	"github.com/joseluisq/cline/handler"
)

// Build-time application values
var (
	versionNumber string = "devel"
	buildTime     string
	buildCommit   string
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	ap := app.New()
	ap.Name = "enve"
	ap.Summary = "Run a program in a modified environment providing an optional .env file or variables from stdin"
	ap.Version = versionNumber
	ap.BuildTime = buildTime
	ap.BuildCommit = buildCommit
	ap.Flags = Flags
	ap.Handler = appHandler

	if err := handler.New(ap).Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
