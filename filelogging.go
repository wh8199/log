package log

import (
	"os"
	"sync"
)

type FileLogging interface {
	LoggingInterface
}
type fileLogging struct {
	mu sync.RWMutex
	logging
	isLogFile bool
	fileDir   string
	maxSize   int64
	prefix    string
	fileName  string
}

func NewFileLogging(name string, level LoggingLevel, callerLevel int, isLogFile bool, fileDir, prefix string, maxSize int64) FileLogging {
	if fileDir == "" {
		fileDir = getCurrentPath()
	}
	if prefix == "" {
		prefix = "log"
	}
	if maxSize < 20 {
		maxSize = 20 //MB
	}
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
		mu:      sync.RWMutex{},
		fileDir: fileDir,
		maxSize: maxSize,
		prefix:  prefix,
	}
	if isLogFile {
		logging.Start()
	}
	return logging
}
func (f *fileLogging) Start() {
	if err := f.initLogFile(); err != nil {
		f.Fatalf("log file is err:%s", err.Error())
	}
}

func (f *fileLogging) initLogFile() error {
	f.fileName = generateFileName(f.prefix)
	logFile := joinFilePath(f.fileDir, f.fileName)
	if !isExist(logFile) {
		if err := os.Mkdir(f.fileDir, 0755); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.SetOutPut(file)
	return nil
}

func (f *fileLogging) logWriter() error {
	logFile := joinFilePath(f.fileDir, f.fileName)
	if fileSize(logFile) >= f.maxSize*1024*1024 {
		f.fileName = generateFileName(f.prefix)
		logFile := joinFilePath(f.fileDir, f.fileName)
		if !isExist(logFile) {
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

}

// // init filelogger split by fileSize
// func (f *fileLogging) initLoggerBySize() {

// 	f.mu.Lock()
// 	defer f.mu.Unlock()

// 	logFile := joinFilePath(f.fileDir, f.fileName)
// 	for i := 1; i <= f.fileCount; i++ {
// 		if !isExist(logFile + "." + strconv.Itoa(i)) {
// 			break
// 		}

// 		f.suffix = i
// 	}

// 	if !f.isMustSplit() {
// 		if !isExist(f.fileDir) {
// 			os.Mkdir(f.fileDir, 0755)
// 		}
// 		f.logFile, _ = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
// 		f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
// 	} else {
// 		f.split()
// 	}

// 	go f.logWriter()
// 	go f.fileMonitor()
// }

// // used for determine the fileLogger f is time to split.
// // size: once the current fileLogger's fileSize >= config.fileSize need to split
// // daily: once the current fileLogger stands for yesterday need to split
// func (f *fileLogging) isMustSplit() bool {

// 	logFile := joinFilePath(f.fileDir, f.fileName)
// 	if f.fileCount > 1 {
// 		if fileSize(logFile) >= f.fileSize {
// 			return true
// 		}
// 	}

// 	return false
// }

// // Split fileLogger
// func (f *fileLogging) split() {

// 	logFile := joinFilePath(f.fileDir, f.fileName)

// 	f.suffix = int(f.suffix%f.fileCount + 1)
// 	if f.logFile != nil {
// 		f.logFile.Close()
// 	}

// 	logFileBak := logFile + "." + strconv.Itoa(f.suffix)
// 	if isExist(logFileBak) {
// 		os.Remove(logFileBak)
// 	}
// 	os.Rename(logFile, logFileBak)

// 	f.logFile, _ = os.Create(logFile)
// 	f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)

// }

// // After some interval time, goto check the current fileLogger's size or date
// func (f *fileLogging) fileMonitor() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			f.lg.Printf("FileLogger's FileMonitor() catch panic: %v\n", err)
// 		}
// 	}()

// 	//TODO  load logScan interval from config file
// 	logScan := DEFAULT_LOG_SCAN

// 	timer := time.NewTicker(time.Duration(logScan) * time.Second)
// 	for {
// 		select {
// 		case <-timer.C:
// 			f.fileCheck()
// 		}
// 	}
// }

// // If the current fileLogger need to split, just split
// func (f *fileLogging) fileCheck() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			f.lg.Printf("FileLogger's FileCheck() catch panic: %v\n", err)
// 		}
// 	}()

// 	if f.isMustSplit() {
// 		f.mu.Lock()
// 		defer f.mu.Unlock()

// 		f.split()
// 	}
// }
