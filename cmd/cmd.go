// Package cmd implements the command line interface for the enve application.
package cmd

import (
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
func Execute(args []string) error {
	ap := app.New()
	ap.Name = "enve"
	ap.Summary = "Run a program in a modified environment providing an optional .env file or variables from stdin"
	ap.Version = versionNumber
	ap.BuildTime = buildTime
	ap.BuildCommit = buildCommit
	ap.Flags = Flags
	ap.Handler = appHandler

	return handler.New(ap).Run(args)
}
