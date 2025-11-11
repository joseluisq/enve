package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromReader(t *testing.T) {
	t.Run("should return an EnvReader from an io.Reader", func(t *testing.T) {
		reader := strings.NewReader("KEY=VALUE")
		envReader := FromReader(reader)
		assert.NotNil(t, envReader, "should not be nil")

		env, ok := envReader.(*Env)
		assert.True(t, ok, "should be of type *Env")
		assert.Equal(t, reader, env.r, "should contain the provided reader")
	})
}

func TestFromPath(t *testing.T) {
	tests := []struct {
		name       string
		expected   func() (path string, data []byte, err error)
		createFile bool
	}{
		{
			name: "should return an error for a non-existent file",
			expected: func() (path string, data []byte, err error) {
				path = "non-existent-file.env"
				err = fmt.Errorf("error: cannot access file '%s'.", path)
				return
			},
		},
		{
			name: "should return an error for a directory",
			expected: func() (path string, data []byte, err error) {
				path = t.TempDir()
				err = fmt.Errorf("error: file path '%s' is a directory", path)
				return
			},
		},
		{
			name: "should return an EnvFile for an existing file",
			expected: func() (path string, data []byte, err error) {
				path = filepath.Join(t.TempDir(), "test.env")
				data = []byte("KEY=VALUE")
				return
			},
			createFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedData, expectedErr := tt.expected()

			if tt.createFile {
				if err := os.WriteFile(expectedPath, expectedData, 0644); err != nil {
					assert.NoError(t, err, "should write temp file without error")
				}
			}

			if envFile, err := FromPath(expectedPath); expectedErr != nil {
				assert.Error(t, err, "should return an error for a non-existent file")
				assert.Contains(t, err.Error(), expectedErr.Error(), "error message should indicate file access issue")
				assert.Nil(t, envFile, "should return a nil EnvFile")
			} else {
				assert.NoError(t, err, "should not return an error for an existing file")
				assert.NotNil(t, envFile, "should return a non-nil EnvFile")

				err = envFile.Close()
				assert.NoError(t, err, "should close the file without error")
			}
		})
	}
}

func TestEnv_Parse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMap Map
		expectedErr bool
	}{
		{
			name:        "should parse valid key-value pairs",
			input:       "KEY1=VALUE1\nKEY2=VALUE2",
			expectedMap: Map{"KEY1": "VALUE1", "KEY2": "VALUE2"},
		},
		{
			name:        "should return an empty map for an empty reader",
			expectedMap: Map{},
		},
		{
			name:        "should handle various valid formats",
			input:       "  KEY1 = VALUE1 #comment\nexport KEY2=VALUE2",
			expectedMap: Map{"KEY1": "VALUE1", "KEY2": "VALUE2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			env := &Env{r: reader}
			envMap, err := env.Parse()

			if tt.expectedErr {
				assert.Error(t, err, "should return an error")
			} else {
				assert.NoError(t, err, "should not return an error")
				assert.Equal(t, tt.expectedMap, envMap, "parsed map should match expected")
			}
		})
	}
}

func TestEnv_Load(t *testing.T) {
	t.Run("should load variables when overload is false", func(t *testing.T) {
		t.Setenv("EXISTING_KEY", "initial_value")

		reader := strings.NewReader("NEW_KEY=new_value\nEXISTING_KEY=new_value_overwritten")
		env := &Env{r: reader}
		err := env.Load(false)

		assert.NoError(t, err, "should load without error")
		assert.Equal(t, "new_value", os.Getenv("NEW_KEY"), "should set new environment variable")
		assert.Equal(t, "initial_value", os.Getenv("EXISTING_KEY"), "should not overwrite existing environment variable")
	})

	t.Run("should load and overwrite variables when overload is true", func(t *testing.T) {
		t.Setenv("EXISTING_KEY", "initial_value")

		reader := strings.NewReader("NEW_KEY=new_value\nEXISTING_KEY=new_value_overwritten")
		env := &Env{r: reader}
		err := env.Load(true)

		assert.NoError(t, err, "should load without error")
		assert.Equal(t, "new_value", os.Getenv("NEW_KEY"), "should set new environment variable")
		assert.Equal(t, "new_value_overwritten", os.Getenv("EXISTING_KEY"), "should overwrite existing environment variable")
	})

	t.Run("should return error on parse failure", func(t *testing.T) {
		reader := strings.NewReader("INVALID-INPUT")
		env := &Env{r: reader}
		err := env.Load(false)

		assert.Error(t, err, "should return an error on parse failure")
	})
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestEnv_Parse_Error(t *testing.T) {
	t.Run("should return error when reader fails", func(t *testing.T) {
		env := &Env{r: &errorReader{}}
		_, err := env.Parse()
		assert.Error(t, err, "should return an error if reading fails")
	})
}

func TestEnv_Close(t *testing.T) {
	t.Run("should close the file if it's an os.File", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.env")
		err := os.WriteFile(filePath, []byte("KEY=VALUE"), 0644)
		assert.NoError(t, err)

		file, err := os.Open(filePath)
		assert.NoError(t, err)

		env := &Env{r: file}
		err = env.Close()
		assert.NoError(t, err, "should close the file without error")

		_, err = file.Read(make([]byte, 1))
		assert.Error(t, err, "should be an error reading from a closed file")
		assert.True(t, env.closed, "closed flag should be true")
	})

	t.Run("should not return an error if reader is not an os.File", func(t *testing.T) {
		reader := strings.NewReader("KEY=VALUE")
		env := &Env{r: reader}
		err := env.Close()
		assert.NoError(t, err, "should not return an error for non-file readers")
		assert.False(t, env.closed, "closed flag should be false")
	})

	t.Run("should do nothing on second close", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.env")
		err := os.WriteFile(filePath, []byte("KEY=VALUE"), 0644)
		assert.NoError(t, err)

		file, err := os.Open(filePath)
		assert.NoError(t, err)

		env := &Env{r: file}
		err = env.Close()
		assert.NoError(t, err, "first close should be successful")
		assert.True(t, env.closed, "closed flag should be true after first close")

		err = env.Close()
		assert.NoError(t, err, "second close should also be successful (no-op)")
		assert.True(t, env.closed, "closed flag should remain true")
	})

	t.Run("should not return an error if reader is nil", func(t *testing.T) {
		env := &Env{r: nil}
		err := env.Close()
		assert.NoError(t, err, "should not return an error for a nil reader")
	})
}
