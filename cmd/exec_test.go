//go:build !windows
// +build !windows

package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_execCmd(t *testing.T) {
	basePath, err := filepath.Abs("../")
	assert.NoError(t, err)
	fixturesPath := filepath.Join(basePath, "fixtures", "cmd")
	bashFile := filepath.Join(fixturesPath, "test.sh")

	tests := []struct {
		name           string
		tailArgs       []string
		chdirPath      string
		newEnv         bool
		envVars        []string
		setupEnv       map[string]string
		expectedErr    error
		expectedOutput string
	}{
		{
			name:           "should execute a command successfully",
			tailArgs:       []string{"echo", "hello world"},
			expectedOutput: "hello world\n",
		},
		{
			name:        "should return error for non-existent command",
			tailArgs:    []string{"nonexistentcommand"},
			expectedErr: errors.New("error: executable 'nonexistentcommand' was not found."),
		},
		{
			name:     "should execute command with existing environment variables",
			tailArgs: []string{bashFile},
			expectedOutput: "" +
				"DB_PROTOCOL=udp\n" +
				"DB_HOST=127.0.0.1\n" +
				"DB_PORT=3306\n" +
				"DB_DEFAULT_CHARACTER_SET=utf8\n" +
				"DB_EXPORT_GZIP=true\n" +
				"DB_EXPORT_FILE_PATH=dbname.sql.gz\n" +
				"DB_NAME=dbname\n" +
				"DB_USERNAME=username\n" +
				"DB_PASSWORD=passwd\n" +
				"DB_ARGS=\n",
		},
		{
			name:     "should execute command with new environment variables",
			tailArgs: []string{bashFile},
			newEnv:   true,
			envVars:  []string{"DB_PROTOCOL=tcp", "DB_HOST=localhost"},
			expectedOutput: "" +
				"DB_PROTOCOL=tcp\n" +
				"DB_HOST=localhost\n" +
				"DB_PORT=\n" +
				"DB_DEFAULT_CHARACTER_SET=\n" +
				"DB_EXPORT_GZIP=\n" +
				"DB_EXPORT_FILE_PATH=\n" +
				"DB_NAME=\n" +
				"DB_USERNAME=\n" +
				"DB_PASSWORD=\n" +
				"DB_ARGS=\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setupEnv {
				t.Setenv(k, v)
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := execCmd(tt.tailArgs, tt.chdirPath, tt.newEnv, tt.envVars)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.expectedErr != nil {
				assert.Error(t, err, "expected an error but got none")
				assert.Contains(t, err.Error(), tt.expectedErr.Error(), "expected error message to match")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
			}

			assert.Equal(t, tt.expectedOutput, output, "output did not match expected")
		})
	}
}
