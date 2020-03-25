package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"text/template"

	"github.com/urfave/cli/v2"
)

// application version values
var (
	versionNumber string = "devel"
	buildTime     string
)

var versionTemplate = `Version:      {{.Version}}
Go version:   {{.GoVersion}}
Built:        {{.BuildTime}}
OS/Arch:      {{.Os}}/{{.Arch}}`

// VersionFlag builds a new Version flag
func VersionFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "shows the current version",
	}
}

// VersionAction defines the version action function
func VersionAction(c *cli.Context) (err error) {
	if err = getVersionTemplate(os.Stdout); err != nil {
		return err
	}

	fmt.Print("\n")
	return nil
}

// getVersionTemplate write the version template
func getVersionTemplate(wr io.Writer) error {
	tmpl, err := template.New("").Parse(versionTemplate)

	if err != nil {
		return err
	}

	v := struct {
		Version   string
		GoVersion string
		BuildTime string
		Os        string
		Arch      string
	}{
		Version:   versionNumber,
		GoVersion: runtime.Version(),
		BuildTime: buildTime,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	return tmpl.Execute(wr, v)
}
