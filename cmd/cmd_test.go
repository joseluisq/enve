//go:build !windows
// +build !windows

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainEnv(t *testing.T) {
	expected := strings.Join([]string{
		"DB_PROTOCOL=tcp",
		"DB_HOST=127.0.0.1",
		"DB_PORT=3306",
		"DB_DEFAULT_CHARACTER_SET=utf8",
		"DB_EXPORT_GZIP=true",
		"DB_EXPORT_FILE_PATH=dbname.sql.gz",
		"DB_NAME=dbname",
		"DB_USERNAME=username",
		"DB_PASSWORD=passwd",
		"DB_ARGS=",
	}, "\n")

	t.Run("should read .env file", func(t *testing.T) {
		basePath := path.Dir("./../")
		envFile := basePath + "/fixtures/cmd/devel.env"
		bashFile := basePath + "/fixtures/cmd/test.sh"

		cmd := exec.Command("go", "run", basePath+"/main.go", "-f", envFile, bashFile)

		var out bytes.Buffer
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			assert.Error(t, err, "error trying to read the .env file.")
		}

		actual := strings.Trim(out.String(), "\n")
		assert.Equal(t, expected, actual, "one or more env keys have wrong values")
	})
}

func TestOverwriteDisabledPlainEnv(t *testing.T) {
	expected := strings.Join([]string{
		"DB_PROTOCOL=udp",
		"DB_HOST=127.0.0.1",
		"DB_PORT=3306",
		"DB_DEFAULT_CHARACTER_SET=utf8",
		"DB_EXPORT_GZIP=true",
		"DB_EXPORT_FILE_PATH=dbname.sql.gz",
		"DB_NAME=dbname",
		"DB_USERNAME=username",
		"DB_PASSWORD=passwd",
		"DB_ARGS=",
	}, "\n")

	t.Run("should not overwrite env vars", func(t *testing.T) {
		basePath := path.Dir("./../")
		envFile := basePath + "/fixtures/cmd/devel.env"
		bashFile := basePath + "/fixtures/cmd/test.sh"

		// Set DB_PROTOCOL as UDP before running the script
		if err := os.Setenv("DB_PROTOCOL", "udp"); err != nil {
			assert.Error(t, err, "error setting DB_PROTOCOL environment variable")
		}

		cmd := exec.Command("go", "run", basePath+"/main.go", "-f", envFile, bashFile)

		var out bytes.Buffer
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			assert.Error(t, err, "error trying to read the .env file.")
		}

		actual := strings.Trim(out.String(), "\n")
		assert.Equal(t, expected, actual, "one or more env keys have wrong values")
	})
}

func TestOverwriteEnabledPlainEnv(t *testing.T) {
	expected := strings.Join([]string{
		"DB_PROTOCOL=tcp",
		"DB_HOST=127.0.0.1",
		"DB_PORT=3306",
		"DB_DEFAULT_CHARACTER_SET=utf8",
		"DB_EXPORT_GZIP=true",
		"DB_EXPORT_FILE_PATH=dbname.sql.gz",
		"DB_NAME=dbname",
		"DB_USERNAME=username",
		"DB_PASSWORD=passwd",
		"DB_ARGS=",
	}, "\n")

	t.Run("should overwrite env vars", func(t *testing.T) {
		basePath := path.Dir("./../")
		envFile := basePath + "/fixtures/cmd/devel.env"
		bashFile := basePath + "/fixtures/cmd/test.sh"

		// Set DB_PROTOCOL as UDP before running the script
		if err := os.Setenv("DB_PROTOCOL", "udp"); err != nil {
			assert.Error(t, err, "error setting DB_PROTOCOL environment variable")
		}

		cmd := exec.Command("go", "run", basePath+"/main.go", "-w", "-f", envFile, bashFile)

		var out bytes.Buffer
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			assert.Error(t, err, "error trying to read the .env file.")
		}

		actual := strings.Trim(out.String(), "\n")
		assert.Equal(t, expected, actual, "one or more env keys have wrong values")
	})
}

const maxArgsCount = 128

func TestExecute(t *testing.T) {
	basePath := path.Dir("./../")
	envFile := basePath + "/fixtures/cmd/devel.env"

	tests := []struct {
		name        string
		vargs       []string
		expectedErr error
	}{
		{
			name: "should return error for too many arguments",
			vargs: func() []string {
				// Create a slice with more arguments than the allowed maximum
				args := make([]string, maxArgsCount+2)
				args[0] = "app"
				for i := 1; i < len(args); i++ {
					args[i] = "arg"
				}
				return args
			}(),
			expectedErr: fmt.Errorf("error: number of arguments exceeds the limit of %d", maxArgsCount),
		},
		{
			name:        "should return error for non-existent file",
			expectedErr: errors.New("error: cannot access file '.env'.\nstat .env: no such file or directory"),
		},
		{
			name:        "should return error for non-existent command",
			vargs:       []string{"app", "--file", envFile, "notfoundcmd"},
			expectedErr: errors.New("error: executable 'notfoundcmd' was not found.\nexec: \"notfoundcmd\": executable file not found in $PATH"),
		},
		{
			name:  "should execute command successfully",
			vargs: []string{"app", "--no-file", "pwd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Execute(tt.vargs); tt.expectedErr != nil {
				assert.Error(t, err, "error was not expected but got one")
				assert.Equal(t, err.Error(), tt.expectedErr.Error(), "Error message does not match the expected one")
			} else {
				assert.NoError(t, err, "unexpected error but got none")
			}
		})
	}
}
