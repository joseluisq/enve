//go:build !windows
// +build !windows

package cmd

import (
	"bytes"
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
