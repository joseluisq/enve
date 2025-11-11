package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	t.Run("should return an error for an empty file path", func(t *testing.T) {
		err := FileExists("")
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "file path was empty or not provided", err.Error(), "error message should match")
	})

	t.Run("should return an error for a non-existent file", func(t *testing.T) {
		err := FileExists("non-existent-file.txt")
		assert.Error(t, err, "should return an error")
		assert.Contains(t, err.Error(), "cannot access file", "error message should indicate file access issue")
	})

	t.Run("should return an error if the path is a directory", func(t *testing.T) {
		dir := t.TempDir()
		err := FileExists(dir)
		assert.Error(t, err, "should return an error for a directory")
		assert.Contains(t, err.Error(), "is a directory", "error message should indicate it's a directory")
	})

	t.Run("should return nil for an existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		assert.NoError(t, err, "should create temp file without error")

		err = FileExists(filePath)
		assert.NoError(t, err, "should not return an error for an existing file")
	})
}

func TestDirExists(t *testing.T) {
	t.Run("should return an error for an empty directory path", func(t *testing.T) {
		err := DirExists("")
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "error: directory path was empty or not provided", err.Error(), "error message should match")
	})

	t.Run("should return an error for a non-existent directory", func(t *testing.T) {
		err := DirExists("non-existent-dir")
		assert.Error(t, err, "should return an error")
		assert.Contains(t, err.Error(), "cannot access directory", "error message should indicate directory access issue")
	})

	t.Run("should return an error if the path is a file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		assert.NoError(t, err, "should create temp file without error")

		err = DirExists(filePath)
		assert.Error(t, err, "should return an error for a file")
		assert.Contains(t, err.Error(), "is a file", "error message should indicate it's a file")
	})

	t.Run("should return nil for an existing directory", func(t *testing.T) {
		dir := t.TempDir()
		err := DirExists(dir)
		assert.NoError(t, err, "should not return an error for an existing directory")
	})
}
