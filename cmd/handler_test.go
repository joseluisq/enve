package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	cli "github.com/joseluisq/cline"
	"github.com/stretchr/testify/assert"
)

func fixturePath(filename string) string {
	return filepath.Join("..", "fixtures", "handler", filename)
}

func TestAppHandler_Output(t *testing.T) {
	// Setup app
	app := cli.New()
	app.Name = "enve"
	app.Summary = "Run a program in a modified environment"
	app.Version = "test"
	app.Flags = Flags
	app.Handler = appHandler

	tests := []struct {
		// Input
		name string
		args []string

		// Output
		err         error
		contain     string
		containList []string
	}{
		{
			name: "should show help with available flags",
			args: []string{"enve", "-h"},
			containList: []string{
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
			name: "should run with new environment",
			args: []string{"enve", "--new-environment", "-f", fixturePath(".env")},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		assert.Fail(t, "Failed to create pipe: %v", err)
	}
	os.Stdout = w

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run(tt.args)

			// close writer and restore stdout
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				assert.Fail(t, "Failed to copy output: %v", err)
			}

			if tt.err != nil {
				assert.EqualError(
					t, err, tt.err.Error(), "app.Run() with args %v failed: %v", tt.args, err,
				)
				return
			}

			assert.NoError(t, err, "app.Run() with args %v failed: %v", tt.args, err)

			if tt.contain != "" {
				assert.Contains(t, buf.String(), tt.contain, "Output should contain %q", tt.contain)
			}
			for _, s := range tt.containList {
				assert.Contains(t, buf.String(), s, "Output should contain %q", s)
			}
		})
	}
}
