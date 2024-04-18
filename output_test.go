package log

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func cleanLogs() {
	files, err := filepath.Glob("*.log")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		os.Remove(file)
	}
}

func TestGenerateLog(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	currentTime := time.Now()
	fileName := generateFileName(fileOutput.Prefix, currentTime)
	defer cleanLogs()

	if err := fileOutput.generateFileWithTime(currentTime); err != nil {
		t.Error(err)
		return
	}

	if _, err := os.Stat(fileName); err != nil {
		t.Error(err)
	}
}

func TestParseFileTime(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	currentTime := time.Now()
	fileName := generateFileName(fileOutput.Prefix, currentTime)
	defer cleanLogs()

	if err := fileOutput.generateFileWithTime(currentTime); err != nil {
		t.Error(err)
		return
	}

	ts, err := fileOutput.parseFileTime(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	if ts != currentTime.Unix() {
		t.Error("parse file time failed")
		return
	}
}

func TestParseFileName(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	currentTime := time.Now()
	fileName := generateFileName(fileOutput.Prefix, currentTime)

	ok, ts, err := fileOutput.parseFileName(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("parse log file name failed")
		return
	}

	if ts != currentTime.Unix() {
		t.Error("parse log file time failed")
		return
	}

	ok, _, err = fileOutput.parseFileName("test.data")
	if err == nil || ok {
		t.Error("test parsing log name failed")
		return
	}
}

func TestCheckLogSize(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	currentTime := time.Now()
	if err := fileOutput.generateFileWithTime(currentTime); err != nil {
		t.Error(err)
		return
	}
	defer cleanLogs()

	bigLog, err := fileOutput.checkLogFileSize()
	if err != nil {
		t.Error(err)
		return
	}

	if bigLog {
		t.Error("check log file size failed")
		return
	}

	for i := 0; i < 1024*1024*2; i++ {
		fileOutput.Write([]byte{'a'})
	}

	bigLog, err = fileOutput.checkLogFileSize()
	if err != nil {
		t.Error(err)
		return
	}

	if !bigLog {
		t.Error("check log file size failed")
		return
	}
}

func TestNotExistLogSize(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
		fileName: "test.log",
	}

	needNewLog, err := fileOutput.checkLogFileSize()
	if err != nil {
		t.Error(err)
		return
	}

	if !needNewLog {
		t.Error("check the size of none-exist file failed")
	}
}

func TestGetAllLogs(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	defer cleanLogs()
	logs := []int{30, 45, 60, 75, 80}
	currentTime := time.Now()

	newLogs := map[string]int64{}

	for _, log := range logs {
		expireTime := currentTime.Add(time.Second * time.Duration(-1*log))
		if err := fileOutput.generateFileWithTime(expireTime); err != nil {
			t.Error(err)
			return
		}

		newLogs[fileOutput.fileName] = expireTime.Unix()
	}

	expireLogs, err := fileOutput.getAllLogs()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(newLogs, expireLogs) {
		t.Error("get all log file failed")
		return
	}
}

func TestCleanExpiredLogs(t *testing.T) {
	fileOutput := FileOutput{
		LogRotateConfig: LogRotateConfig{
			EnableLogFile: true,
			Prefix:        "log",
			FileDir:       ".",
			MaxSize:       1,
			MaxLogLife:    60,
		},
	}

	defer cleanLogs()
	logs := []int{30, 45, 60, 75, 80}
	currentTime := time.Now()

	newLogs := map[string]int64{}

	for i := len(logs) - 1; i >= 0; i-- {
		expireTime := currentTime.Add(time.Second * time.Duration(-1*logs[i]))
		if err := fileOutput.generateFileWithTime(expireTime); err != nil {
			t.Error(err)
			return
		}

		if currentTime.Unix()-expireTime.Unix() <= fileOutput.MaxLogLife {
			newLogs[fileOutput.fileName] = expireTime.Unix()
		}
	}

	if err := fileOutput.cleanExpiredLogs(currentTime.Unix()); err != nil {
		t.Error(err)
		return
	}

	expireLogs, err := fileOutput.getAllLogs()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(newLogs, expireLogs) {
		t.Error("get all log file failed")
		return
	}
}
