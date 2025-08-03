package cmd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	cli "github.com/joseluisq/cline"
	"github.com/joseluisq/enve/env"
	"github.com/stretchr/testify/assert"
)

const defaultEnvFile = "devel.env"

var basePath = path.Dir("./../")

func newArgsWithFile(filename string, args []string) []string {
	base := append([]string{"enve-test"}, "-f", filepath.Join(basePath, "fixtures", "handler", filename))
	return append(base, args...)
}

func newArgs(args []string) []string {
	return newArgsWithFile(defaultEnvFile, args)
}

func TestAppHandler_Output(t *testing.T) {
	tests := []struct {
		// Input
		name string
		args []string

		// Output
		err        error
		expectText []string // []string{"HOST=127.0.0.1"}
		expectJSON *env.Environment
		expectXML  *env.Environment
	}{
		{
			name:       "should output nothing with no args",
			args:       newArgs([]string{}),
			expectText: []string{""},
		},
		{
			name: "should output help with available flags",
			args: newArgs([]string{"--help"}),
			expectText: []string{
				"enve",
				"Run a program in a modified environment",
				"v1.0.0-beta.1",
				"-f --file",
				"-o --output",
				"-w --overwrite",
				"-c --chdir",
				"-n --new-environment",
				"-i --ignore-environment",
				"-z --no-file",
				"-s --stdin",
				"-h --help",
				"-v --version",
			},
		},
		{
			name:       "should output with --new-environment as text",
			args:       newArgs([]string{"--new-environment"}),
			expectText: []string{""},
		},
		{
			name: "should output with --new-environment as json",
			args: newArgs([]string{"--new-environment", "--output", "json"}),
			expectJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name: "should output with --new-environment as xml",
			args: newArgs([]string{"--new-environment", "--output", "xml"}),
			expectXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup app
			app := cli.New()
			app.Name = "enve"
			app.Summary = "Run a program in a modified environment"
			app.Version = "v1.0.0-beta.1"
			app.Flags = Flags
			app.Handler = appHandler

			// Capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				assert.Fail(t, "Failed to create pipe: %v", err)
			}
			os.Stdout = w

			if err := app.Run(tt.args); tt.err != nil {
				assert.EqualError(
					t, err, tt.err.Error(), "app.Run() with args %v failed: %v", tt.args, err,
				)
				return
			} else {
				assert.NoError(t, err, "app.Run() with args %v failed: %v", tt.args, err)
			}

			// close writer and restore stdout
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				assert.Fail(t, "Failed to copy output: %v", err)
			}

			output := buf.Bytes()

			if tt.expectJSON != nil {
				var vars env.Environment
				if err := json.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal JSON output: %v", err)
				}
				assert.ElementsMatch(
					t, vars.Env, tt.expectJSON.Env, "JSON output should match to %#v", tt.expectJSON,
				)
			}

			if tt.expectXML != nil {
				var vars env.Environment
				if err := xml.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal XML output: %v", err)
				}
				assert.ElementsMatch(
					t, vars.Env, tt.expectXML.Env, "XML output should match to %#v", tt.expectXML,
				)
			}

			for _, s := range tt.expectText {
				assert.Contains(t, string(output), s, "Output should contain %q", s)
			}
		})
	}
}
