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
		return fmt.Errorf("cannot access file: %s; %v", filePath, err)
	} else {
		if info.IsDir() {
			return fmt.Errorf("file path is a directory: %s", filePath)
		}
		return nil
	}
}

func DirExists(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path was empty or not provided")
	}
	if info, err := os.Stat(dirPath); err != nil {
		return fmt.Errorf("cannot access directory: %s; %v", dirPath, err)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("directory path is a file: %s", dirPath)
		}
		return nil
	}
}
