//go:build windows
// +build windows

package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func TestPlainEnv(t *testing.T) {
	expected := []string{
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
	}

	basePath := path.Dir("./../")

	envFile := basePath + "/fixtures/plain.env"
	psFile := basePath + "/fixtures/test.ps1"

	cmd := exec.Command("go", "run", basePath+"/main.go", "-f", envFile, "powershell", "-ExecutionPolicy", "Bypass", "-File", psFile)

	var out bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		t.Errorf("error trying to read the .env file.\n %s", err)
	}

	actual := strings.Split(out.String(), "\n")
	for i, exp := range expected {
		act := strings.TrimRight(actual[i], "\r")
		if exp != act {
			t.Errorf("actual: [%s] expected: [%s]", act, exp)
		}
	}
}
