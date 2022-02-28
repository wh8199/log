package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Get the files with the prefix before the specified date of the current file
func getCurrentPathFiles() ([]string, error) {
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, err
	}
	filePaths := make([]string, 0)
	for _, file := range files {
		if file.ModTime().Unix()-time.Now().Unix() < 7*24*60*60*1e9 {
			continue
		}
		if strings.HasPrefix(file.Name(), "log") {
			filePaths = append(filePaths, file.Name())
		}
	}
	return filePaths, nil
}

// Get current path
func getCurrentPath() string {
	absPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return absPath
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

// return length in bytes for regular files
func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 0
	}

	return f.Size()
}

// return file name without dir
func shortFileName(file string) string {
	return filepath.Base(file)
}

func generateFileName(prefix string) string {
	current := time.Now()
	format := fmt.Sprintf("%s_2006-01-02_15-04.log", prefix)
	return current.Format(format)
}
