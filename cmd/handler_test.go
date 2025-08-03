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

// ElementsContain asserts that all elements in listB are contained in listA.
func ElementsContain(t assert.TestingT, listA interface{}, listB interface{}, msgAndArgs ...interface{}) (ok bool) {
	aVal := reflect.ValueOf(listA)
	bVal := reflect.ValueOf(listB)

	if aVal.Kind() != reflect.Slice || bVal.Kind() != reflect.Slice {
		return assert.Fail(t, "ElementsContain only accepts slice arguments", msgAndArgs...)
	}

	// Build multiset for listA
	counts := make(map[interface{}]int)
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
		expectedErr  error
		expectedText []string // []string{"HOST=127.0.0.1"}
		expectedJSON *env.Environment
		expectedXML  *env.Environment
	}{
		{
			name:         "should output nothing with no args",
			args:         newArgs([]string{}),
			expectedText: []string{""},
		},
		{
			name: "should output as text",
			args: newArgs([]string{"--output", "text"}),
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
			name: "should output as json",
			args: newArgs([]string{"--output", "json"}),
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
			name: "should output as xml",
			args: newArgs([]string{"--output", "xml"}),
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
			name: "should output help with available flags",
			args: newArgs([]string{"--help"}),
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
			name:         "should output with --new-environment as text",
			args:         newArgs([]string{"--new-environment"}),
			expectedText: []string{""},
		},
		{
			name: "should output with --new-environment as json",
			args: newArgs([]string{"--new-environment", "--output", "json"}),
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
			name: "should output with --new-environment as xml",
			args: newArgs([]string{"--new-environment", "--output", "xml"}),
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
			name:         "should output with --no-file as text",
			args:         newArgs([]string{"--no-file", "--new-environment", "--output", "text"}),
			expectedText: []string{""},
		},
		{
			name: "should output with --no-file as json",
			args: newArgs([]string{"--no-file", "--new-environment", "--output", "json"}),
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{},
			},
		},
		{
			name: "should output with --no-file as xml",
			args: newArgs([]string{"--no-file", "--ignore-environment", "--output", "xml"}),
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				assert.Fail(t, "Failed to create pipe: %v", err)
			}
			os.Stdout = w

			if err := app.Run(tt.args); tt.expectedErr != nil {
				assert.EqualError(
					t, err, tt.expectedErr.Error(), "app.Run() with args %v failed: %v", tt.args, err,
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
