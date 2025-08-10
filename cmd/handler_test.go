package cmd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	cli "github.com/joseluisq/cline"
	"github.com/joseluisq/enve/env"
	"github.com/stretchr/testify/assert"
)

const defaultEnvFile = "devel.env"

var baseDirPath = filepath.Join(path.Dir("./../"))
var fixturePath = filepath.Join(baseDirPath, "fixtures", "handler")

func newArgs(args []string) []string {
	return append([]string{"enve-test"}, args...)
}

func newArgsWithFile(filename string, args []string) []string {
	return newArgs(append(
		[]string{"-f", filepath.Join(fixturePath, filename)},
		args...,
	))
}

func newArgsDefault(args []string) []string {
	return newArgsWithFile(defaultEnvFile, args)
}

// ElementsContain asserts that all elements in listB are contained in listA.
func ElementsContain(t assert.TestingT, listA any, listB any, msgAndArgs ...any) (ok bool) {
	aVal := reflect.ValueOf(listA)
	bVal := reflect.ValueOf(listB)

	if aVal.Kind() != reflect.Slice || bVal.Kind() != reflect.Slice {
		return assert.Fail(t, "ElementsContain only accepts slice arguments", msgAndArgs...)
	}

	// Build multiset for listA
	counts := make(map[any]int)
	for i := 0; i < aVal.Len(); i++ {
		val := aVal.Index(i).Interface()
		counts[val]++
	}

	// Check that each element in listB is present in listA
	for i := 0; i < bVal.Len(); i++ {
		val := bVal.Index(i).Interface()
		if counts[val] == 0 {
			return assert.Fail(
				t, fmt.Sprintf("Expected element %+v not found in listA: %+v", val, listA), msgAndArgs...,
			)
		}
		counts[val]--
	}

	return true
}

func TestAppHandler_Output(t *testing.T) {
	tests := []struct {
		// Input
		name string
		args []string

		// Output
		globalEnvs   []string
		expectedText []string // []string{"HOST=127.0.0.1"}
		expectedJSON *env.Environment
		expectedXML  *env.Environment
		expectedErr  error
	}{
		{
			name:         "should output nothing with no args provided",
			args:         newArgsDefault([]string{}),
			expectedText: []string{""},
		},
		{
			name: "should output help with available flags",
			args: newArgsDefault([]string{"--help"}),
			expectedText: []string{
				"enve-test",
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
			name: "should output variables as text",
			args: newArgsDefault([]string{"--output", "text"}),
			globalEnvs: []string{
				"API_URL=http://localhost:3000",
			},
			expectedText: []string{
				"API_URL=http://localhost:3000",
				"HOST=127.0.0.1",
				"PORT=8080",
				"DEBUG=true",
				"LOG_LEVEL=info",
			},
		},
		{
			name: "should output variables as json",
			args: newArgsDefault([]string{"--output", "json"}),
			globalEnvs: []string{
				"SERVER_IP=192.168.1.1",
			},
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "SERVER_IP", Value: "192.168.1.1"},
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name: "should output variables as xml",
			args: newArgsDefault([]string{"--output", "xml"}),
			globalEnvs: []string{
				"SERVER2_IP=192.168.1.1",
			},
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "SERVER2_IP", Value: "192.168.1.1"},
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name:         "should output variables with --new-environment as text",
			args:         newArgsDefault([]string{"--new-environment"}),
			expectedText: []string{""},
		},
		{
			name: "should output variables with --new-environment as json",
			args: newArgsDefault([]string{"--new-environment", "--output", "json"}),
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name: "should output variables with --new-environment as xml",
			args: newArgsDefault([]string{"--new-environment", "--output", "xml"}),
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "HOST", Value: "127.0.0.1"},
					{Name: "PORT", Value: "8080"},
					{Name: "DEBUG", Value: "true"},
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name: "should output variables with --no-file as text",
			args: newArgsDefault([]string{"--no-file", "--output", "text"}),
			globalEnvs: []string{
				"HOST=0.0.0.0",
			},
			expectedText: []string{"HOST=0.0.0.0"},
		},
		{
			name: "should output variables with --no-file as json",
			args: newArgsDefault([]string{"--no-file", "--new-environment", "--output", "json"}),
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{},
			},
		},
		{
			name: "should output variables with --no-file as xml",
			args: newArgsDefault([]string{"--no-file", "--ignore-environment", "--output", "xml"}),
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{},
			},
		},
		{
			name: "should overwrite variables and output as text",
			args: newArgsDefault([]string{"--overwrite", "--output", "text"}),
			globalEnvs: []string{
				"HOST=192.168.1.1",
			},
			expectedText: []string{"HOST=127.0.0.1"},
		},
		{
			name: "should overwrite variables and output as xml",
			args: newArgsDefault([]string{"--overwrite", "--output", "xml"}),
			globalEnvs: []string{
				"HOST=192.168.1.1",
			},
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "HOST", Value: "127.0.0.1"},
				},
			},
		},
		{
			name: "should overwrite variables and output as json",
			args: newArgsDefault([]string{"--overwrite", "--output", "json"}),
			globalEnvs: []string{
				"LOG_LEVEL=error",
			},
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		},
		{
			name:        "should return error if env file does not exist in new working dir",
			args:        newArgs([]string{"--chdir", "./cmd", "--output", "text"}),
			expectedErr: fmt.Errorf("error: cannot access directory './cmd'."),
		},
		{
			name: "should output variables if env file exist in new working dir",
			args: newArgs([]string{"--chdir", fixturePath}),
			expectedText: []string{
				"SERVER=localhost",
				"IP=192.168.1.120",
				"LEVEL=info",
			},
		},
		{
			name: "should output variables as xml if env file exist in new working dir",
			args: newArgs([]string{"--chdir", fixturePath, "-o", "xml"}),
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "SERVER", Value: "localhost"},
					{Name: "IP", Value: "192.168.1.120"},
					{Name: "LEVEL", Value: "info"},
				},
			},
		},
		{
			name: "should output variables as json if env file exist in new working dir",
			args: newArgs([]string{"--chdir", fixturePath, "-o", "json"}),
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "SERVER", Value: "localhost"},
					{Name: "IP", Value: "192.168.1.120"},
					{Name: "LEVEL", Value: "info"},
				},
			},
		},
	}

	CWD, err := os.Getwd()
	if err != nil {
		assert.Fail(t, "Failed to get current working directory for tests: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset working directory for tests that will change it
			if slices.Contains(tt.args, "--chdir") || slices.Contains(tt.args, "-c") {
				if err := os.Chdir(CWD); err != nil {
					assert.Fail(t, "Failed to reset working directory before test: %v", err)
				}
			}

			// Setup app
			app := cli.New()
			app.Name = "enve-test"
			app.Summary = "Run a program in a modified environment"
			app.Version = "v1.0.0-beta.1"
			app.Flags = Flags
			app.Handler = appHandler

			if tt.globalEnvs != nil {
				for _, envVar := range tt.globalEnvs {
					parts := bytes.SplitN([]byte(envVar), []byte{'='}, 2)
					if len(parts) == 2 {
						if err := os.Setenv(string(parts[0]), string(parts[1])); err != nil {
							assert.Fail(t, "Failed to set environment variable %s: %v", envVar, err)
						}
					} else {
						assert.Fail(t, "Invalid environment variable format: %s", envVar)
					}
				}
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				assert.Fail(t, "Failed to create pipe for stdout capture: %v", err)
			}
			os.Stdout = w

			// Ensure stdout is restored even if the test panics
			defer func() { os.Stdout = oldStdout }()

			var outCopiedChan = make(chan struct{})
			var buf bytes.Buffer

			go func() {
				defer close(outCopiedChan)
				// NOTE: `io.Copy` will block here until the writer (w) is closed
				_, err := io.Copy(&buf, r)
				assert.NoError(t, err, "Failed to copy output from pipe reader")
			}()

			t.Logf("  Running app as '%v'", strings.Join(tt.args, " "))
			runErr := app.Run(tt.args)

			// Close the pipe's writer end to unblock the `io.Copy` in the goroutine above
			_ = w.Close()
			<-outCopiedChan

			output := buf.Bytes()

			if tt.expectedErr != nil {
				assert.Error(t, runErr, "Expected error but got none for args %v", tt.args)
				assert.Contains(
					t, runErr.Error(), tt.expectedErr.Error(), "app.Run() with args %v failed: %v", tt.args, runErr,
				)
			} else {
				assert.NoError(t, runErr, "app.Run() with args %v", tt.args)
			}

			if tt.expectedJSON != nil {
				var vars env.Environment
				if err := json.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal JSON output: %v", err)
				}

				ElementsContain(
					t, vars.Env, tt.expectedJSON.Env, "JSON output should match to %#v", tt.expectedJSON,
				)
			}

			if tt.expectedXML != nil {
				var vars env.Environment
				if err := xml.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal XML output: %v", err)
				}
				ElementsContain(
					t, vars.Env, tt.expectedXML.Env, "XML output should match to %#v", tt.expectedXML,
				)
			}

			for _, s := range tt.expectedText {
				assert.Contains(t, string(output), s, "Output should contain %q", s)
			}

		})
	}
}
