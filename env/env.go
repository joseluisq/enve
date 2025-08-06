package env

import (
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/joseluisq/enve/fs"
)

type EnvironmentVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Environment defines JSON/XML data structure
type Environment struct {
	Env []EnvironmentVar `json:"environment"`
}

type EnvFile interface {
	Load(overload bool) error
	Parse() (Map, error)
	Close() error
}

type EnvReader interface {
	Load(overload bool) error
	Parse() (Map, error)
}

type Env struct {
	r      io.Reader
	closed bool
}

func FromReader(r io.Reader) EnvReader {
	return &Env{r: r}
}

func FromPath(filePath string) (EnvFile, error) {
	if err := fs.FileExists(filePath); err != nil {
		return nil, err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &Env{r: f}, nil
}

func (e *Env) Load(overload bool) error {
	envMap, err := e.Parse()
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func (e *Env) Parse() (Map, error) {
	return godotenv.Parse(e.r)
}

func (e *Env) Close() error {
	if !e.closed {
		if f, ok := e.r.(*os.File); ok {
			e.closed = true
			return f.Close()
		}
	}
	return nil
}
