package cmd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
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
	CWD, err := os.Getwd()
	if err != nil {
		assert.Fail(t, "Failed to get current working directory for tests", err)
	}

	var baseDirPath = filepath.Join(CWD, "../")
	var fixturePath = filepath.Join(baseDirPath, "fixtures", "handler")

	var newArgs = func(args []string) []string {
		return append([]string{"enve-test"}, args...)
	}

	var newArgsWithFile = func(filename string, args []string) []string {
		return newArgs(append(
			[]string{"-f", filepath.Join(fixturePath, filename)},
			args...,
		))
	}

	var newArgsDefault = func(args []string) []string {
		return newArgsWithFile(defaultEnvFile, args)
	}

	tests := []struct {
		// Input
		name          string
		args          []string
		expectedStdin []byte
		initialEnvs   []string

		// Output
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
			initialEnvs: []string{
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
			initialEnvs: []string{
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
			initialEnvs: []string{
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
			initialEnvs: []string{
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
			initialEnvs: []string{
				"HOST=192.168.1.1",
			},
			expectedText: []string{"HOST=127.0.0.1"},
		},
		{
			name: "should overwrite variables and output as xml",
			args: newArgsDefault([]string{"--overwrite", "--output", "xml"}),
			initialEnvs: []string{
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
			initialEnvs: []string{
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
		{
			name:        "should return error if env file does not exist",
			args:        newArgs([]string{"--file", fixturePath + "-xyz", "-o", "json"}),
			expectedErr: fmt.Errorf("error: cannot access file '%s-xyz'.", fixturePath),
		},
		{
			name:        "should return error if env file cannot be parsed",
			args:        newArgsWithFile("invalid.env", []string{}),
			expectedErr: fmt.Errorf("error: cannot load env from file."),
		},
		{
			name: "should output variables as text when using stdin without initial ones",
			args: newArgs([]string{"--stdin"}),
			expectedStdin: []byte(
				"SERVER=localhost\nIP=192.168.1.120\nLEVEL=info\nAPP_URL=https://localhost",
			),
			expectedText: []string{
				"SERVER=localhost",
				"IP=192.168.1.120",
				"LEVEL=info",
				"APP_URL=https://localhost",
			},
		},
		{
			name: "should output variables as text when using stdin with initial ones",
			args: newArgs([]string{"--stdin"}),
			initialEnvs: []string{
				"SERVER=127.0.0.1",
			},
			expectedStdin: []byte(
				"SERVER=localhost\nIP=192.168.1.120\nLEVEL=info\nAPP_URL=https://localhost",
			),
			expectedText: []string{
				"SERVER=127.0.0.1",
				"IP=192.168.1.120",
				"LEVEL=info",
				"APP_URL=https://localhost",
			},
		},
		{
			name: "should output overwritten variables as json when using stdin",
			args: newArgs([]string{"--stdin", "--overwrite", "-o", "json"}),
			expectedStdin: []byte(
				"NAME=User\nEMAIL=user@example.com\nAGE=30",
			),
			expectedJSON: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "NAME", Value: "User"},
					{Name: "EMAIL", Value: "user@example.com"},
					{Name: "AGE", Value: "30"},
				},
			},
		},
		{
			name: "should output overwritten variables as xml when using stdin",
			args: newArgs([]string{"--stdin", "--overwrite", "-o", "xml"}),
			expectedStdin: []byte(
				"NAME=Gopher\nEMAIL=ghoper@example.com\nAGE=100",
			),
			expectedXML: &env.Environment{
				Env: []env.EnvironmentVar{
					{Name: "NAME", Value: "Gopher"},
					{Name: "EMAIL", Value: "ghoper@example.com"},
					{Name: "AGE", Value: "100"},
				},
			},
		},
		{
			name:        "should return an error invalid output format",
			args:        newArgs([]string{"--output", "xyz"}),
			expectedErr: fmt.Errorf("error: output format 'xyz' is not supported"),
		},
		{
			name:        "should return an error empty output value",
			args:        newArgs([]string{"--output", ""}),
			expectedErr: fmt.Errorf("error: output format was empty or not provided"),
		},
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

			if tt.initialEnvs != nil {
				for _, envVar := range tt.initialEnvs {
					parts := strings.SplitN(envVar, "=", 2)
					if len(parts) == 2 {
						t.Setenv(parts[0], parts[1])
					} else {
						assert.Fail(t, "Invalid environment variable format", envVar)
					}
				}
			}

			// Capture stdin
			if tt.expectedStdin != nil {
				oldStdin := os.Stdin
				r1, w1, err := os.Pipe()
				if err != nil {
					assert.Fail(t, "Failed to create pipe for stdin", err)
				}
				os.Stdin = r1

				defer func() { os.Stdin = oldStdin }()

				if _, err := w1.Write(tt.expectedStdin); err != nil {
					assert.Fail(t, "Failed to write to stdin pipe", err)
				}
				if err := w1.Close(); err != nil {
					assert.Fail(t, "Failed to write to stdin pipe", err)
				}
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				assert.Fail(t, "Failed to create pipe for stdout capture", err)
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
				assert.Error(t, runErr, "Expected error but got none")
				assert.Contains(t, runErr.Error(), tt.expectedErr.Error(), "Error message mismatch")
			} else {
				assert.NoError(t, runErr, "Expected no error but got: %v", runErr)
			}

			if tt.expectedJSON != nil {
				var vars env.Environment
				if err := json.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal JSON output", err)
				}

				ElementsContain(
					t, vars.Env, tt.expectedJSON.Env, "JSON output should match to %v", tt.expectedJSON,
				)
			}

			if tt.expectedXML != nil {
				var vars env.Environment
				if err := xml.Unmarshal(output, &vars); err != nil {
					assert.Fail(t, "Failed to unmarshal XML output", err)
				}
				ElementsContain(
					t, vars.Env, tt.expectedXML.Env, "XML output should match to %v", tt.expectedXML,
				)
			}

			for _, s := range tt.expectedText {
				assert.Contains(t, string(output), s, "Text output should contain %q", s)
			}
		})
	}
}
