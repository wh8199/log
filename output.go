package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type OutPut interface {
	io.Writer
	Rotate() error
}

var _ OutPut = &BufferOutput{}

// BufferOutput is use for unit test
type BufferOutput struct {
	bytes.Buffer
}

func (b *BufferOutput) Rotate() error {
	return nil
}

var _ OutPut = &FileOutput{}

type FileOutput struct {
	*os.File
	fileName string
	LogRotateConfig
}

func NewFileOutput(rotateConfig LogRotateConfig) (io.Writer, error) {
	fileOutput := &FileOutput{
		LogRotateConfig: rotateConfig,
	}

	if err := fileOutput.generateFile(); err != nil {
		return nil, err
	}

	return fileOutput, nil
}

func (f *FileOutput) parseFileTime(fileName string) (int64, error) {
	t, err := time.ParseInLocation(fmt.Sprintf("%s_20060102_150405.log", f.Prefix), fileName, time.Local)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func (f *FileOutput) parseFileName(fileName string) (bool, int64, error) {
	if !(strings.HasSuffix(fileName, ".log") && strings.HasPrefix(fileName, f.Prefix)) {
		return false, 0, fmt.Errorf("invalid log file")
	}

	ts, err := f.parseFileTime(fileName)
	if err != nil {
		return false, 0, err
	}

	return true, ts, nil
}

func (f *FileOutput) getAllLogs() (map[string]int64, error) {
	fileInfos, err := ioutil.ReadDir(f.FileDir)
	if err != nil {
		return nil, err
	}

	ret := map[string]int64{}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()

		isLogFile, ts, err := f.parseFileName(name)
		if err != nil || !isLogFile {
			continue
		}

		ret[name] = ts
	}

	return ret, nil
}

func (f *FileOutput) cleanExpiredLogs(now int64) error {
	logs, err := f.getAllLogs()
	if err != nil {
		return err
	}

	for logName, ts := range logs {
		if logName == f.fileName {
			continue
		}

		if ts+f.MaxLogLife < now {
			os.Remove(logName)
		}
	}

	return nil
}

func (f *FileOutput) generateFileWithTime(t time.Time) error {
	fileName := generateFileName(f.Prefix, t)

	logFile := joinFilePath(f.FileDir, fileName)
	if !isExist(f.FileDir) {
		if err := os.Mkdir(f.FileDir, 0755); err != nil {
			return err
		}
	}

	if f.File != nil {
		f.File.Close()
	}

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	f.fileName = logFile
	f.File = file

	return nil
}

func (f *FileOutput) generateFile() error {
	return f.generateFileWithTime(time.Now().Local())
}

func (f *FileOutput) checkLogFileSize() (bool, error) {
	fileInfo, err := os.Stat(f.fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}

		return false, err
	}

	return fileInfo.Size() > f.MaxSize*1024*1024, nil
}

func (f *FileOutput) Rotate() error {
	needNewFile, err := f.checkLogFileSize()
	if err != nil {
		return err
	}

	t := time.Now().Local()

	if needNewFile {
		if err := f.generateFileWithTime(t); err != nil {
			return err
		}
	}

	return f.cleanExpiredLogs(t.Unix())
}
