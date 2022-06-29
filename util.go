package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Get current path
func getCurrentPath() (string, error) {
	absPath, error := filepath.Abs(filepath.Dir(os.Args[0]))
	return absPath, error
}

// Determine a file or a path exists in the os
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// joinFilePath joins path & file into a single string
func joinFilePath(path, file string) string {
	return filepath.Join(path, file)
}

func generateFileName(prefix string) string {
	current := time.Now()
	format := fmt.Sprintf("%s_20060102_150405.log", prefix)
	return current.Format(format)
}
