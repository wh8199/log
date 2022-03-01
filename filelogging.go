package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"time"
)

type FileLogging interface {
	LoggingInterface
	Start()
}
type fileLogging struct {
	mu sync.Mutex
	logging
	isLogFile bool
	fileDir   string
	maxSize   int64
	maxHour   int64
	prefix    string
	fileName  string
}

func NewFileLogging(name string, level LoggingLevel, callerLevel int, isLogFile bool, fileDir, prefix string, maxSize, maxHour int64) FileLogging {
	if fileDir == "" {
		fileDir = getCurrentPath()
	}
	if prefix == "" {
		prefix = "log"
	}
	if maxSize < 20 {
		maxSize = 20
	}
	if maxHour < 1 {
		maxHour = 1
	}
	maxSize = maxSize << 20 //MB
	logging := &fileLogging{
		logging: logging{
			Name:         name,
			Level:        level,
			Out:          os.Stdout,
			Pool:         NewBufferPool(),
			EnableCaller: true,
			CallerLevel:  callerLevel,
			Formater:     DefaultFormater,
		},
		mu:        sync.Mutex{},
		fileDir:   fileDir,
		isLogFile: isLogFile,
		maxSize:   maxSize,
		maxHour:   maxHour,
		prefix:    prefix,
	}
	return logging
}

func NewFileLoggingAndStart(name string, level LoggingLevel, callerLevel int, isLogFile bool, fileDir string, prefix string, maxSize int64, maxHour int64) FileLogging {
	logging := NewFileLogging(name, level, callerLevel, isLogFile, fileDir, prefix, maxSize, maxHour)
	logging.Start()
	return logging
}

func NewDefaultFileLoggingAndStart(name string, level LoggingLevel, callerLevel int, isLogFile bool) FileLogging {
	logging := NewFileLogging(name, level, callerLevel, isLogFile, "", "", 20, 7*24)
	logging.Start()
	return logging
}

func (f *fileLogging) Start() {
	if !f.isLogFile {
		return
	}
	for {
		if err := f.initLogFile(); err != nil {
			f.Error(err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	go f.logSplit()
	go f.logDelete()
}

func (f *fileLogging) initLogFile() error {
	defer func() {
		if err := recover(); err != nil {
			f.Errorf("init log file  panic: %v\n", err)
		}
	}()
	f.fileName = generateFileName(f.prefix)
	logFile := joinFilePath(f.fileDir, f.fileName)
	if !isExist(f.fileDir) {
		if err := os.Mkdir(f.fileDir, 0755); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open file '%s' fail: %v", logFile, err)
	}
	f.SetOutPut(file)
	return nil
}

func (f *fileLogging) logSplit() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f.splitOnce(); err != nil {
				fmt.Printf("%v", err)
				f.Error(err)
			}
		}
	}
}
func (f *fileLogging) splitOnce() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	defer func() {
		if err := recover(); err != nil {
			f.Errorf("log writer panic: %v\n", err)
		}
	}()
	logFile := joinFilePath(f.fileDir, f.fileName)
	fileSize := fileSize(logFile)
	if fileSize >= f.maxSize {
		f.fileName = generateFileName(f.prefix)
		logFile := joinFilePath(f.fileDir, f.fileName)
		if !isExist(f.fileDir) {
			if err := os.Mkdir(f.fileDir, 0755); err != nil {
				return err
			}
		}
		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		f.Close()
		f.SetOutPut(file)
		return nil
	}
	return nil
}

func (f *fileLogging) logDelete() {
	f.mu.Lock()
	defer f.mu.Unlock()
	ticker := time.NewTicker(5 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f.deleteOnce(); err != nil {
				fmt.Printf("%v", err)
				f.Error(err)
			}
		}
	}

}

func (f *fileLogging) deleteOnce() error {
	defer func() {
		if err := recover(); err != nil {
			f.Errorf("delete file once panic: %v\n", err)
		}
	}()

	files, err := ioutil.ReadDir(f.fileDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		now := time.Now().Unix()
		modify := file.ModTime().Unix()
		if now-modify > f.maxHour*60*60 {
			continue
		}
		express := fmt.Sprintf(`%s_\d{8}_\d{4}`, f.prefix)
		reg := regexp.MustCompile(express)
		fileName := file.Name()
		if reg.Match([]byte(fileName)) {
			if err := os.Remove(joinFilePath(f.fileDir, file.Name())); err != nil {
				return fmt.Errorf("delete file '%s' err:%v", file.Name(), err)
			}
		}
	}
	return nil
}
