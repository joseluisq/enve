package fs

import (
	"fmt"
	"os"
)

func FileExists(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path was empty or not provided")
	}
	if info, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("error: cannot access file '%s'.\n%v", filePath, err)
	} else {
		if info.IsDir() {
			return fmt.Errorf("error: file path '%s' is a directory", filePath)
		}
		return nil
	}
}

func DirExists(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("error: directory path was empty or not provided")
	}
	if info, err := os.Stat(dirPath); err != nil {
		return fmt.Errorf("error: cannot access directory '%s'.\n%v", dirPath, err)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("error: directory path '%s' is a file", dirPath)
		}
		return nil
	}
}
