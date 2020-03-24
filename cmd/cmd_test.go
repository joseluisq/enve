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

	cwd, err := os.Getwd()

	if err != nil {
		t.Error(err)
	}

	basePath := path.Dir(cwd)

	envFile := basePath + "/fixtures/plain.env"
	bashFile := basePath + "/fixtures/test.sh"

	cmd := exec.Command("go", "run", basePath+"/main.go", "-f", envFile, bashFile)

	var out bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		t.Error("error trying to read the .env file")
	}

	actual := strings.Trim(string(out.Bytes()), "\n")

	if expected != actual {
		t.Error("one or more env keys have wrong values")
	}
}
