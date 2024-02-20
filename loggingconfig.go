package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"sync"
	"time"
)

var (
	globalConfig *logConfig
)

type logConfig struct {
	isFile         bool
	fileDir        string
	maxSize        int64
	deleteDuration time.Duration
	maxSecond      int64
	splitDuration  time.Duration
	//prefix used to generate log file names
	prefix string
	//Using the observer pattern
	observers   []*logging
	observersMu sync.Mutex
	//The current output file of the log. Since multiple goroutines share this resource, a read-write lock needs to be added.
	file *os.File
	//Used to protect mutually exclusive resources file
	fileMu sync.RWMutex
	//stop channel
	exitChan chan struct{}
}

// the fileDir is the log save path, default value is current path
//
// the prefix is generated log filename prefix, the default value is log
//
// the maxSize default is 20,unit is MB
//
// the maxSecond is seconds for log retention, after which it will be automatically deleted
func BuildLogConfig(isFile bool, fileDir, prefix, splitDurationStr, deleteDurationStr string, maxSize, maxSecond int) error {

	if fileDir != "" {
		globalConfig.fileDir = fileDir
	}

	if prefix != "" {
		globalConfig.prefix = prefix
	}
	if maxSize > 20 {
		globalConfig.maxSize = int64(maxSize)
	}
	if maxSecond > 5 {
		globalConfig.maxSecond = int64(maxSecond)
	}

	globalConfig.maxSize = globalConfig.maxSize << 20 //MB

	globalConfig.splitDuration = time.Second * 5
	if splitDurationStr != "" {
		splitDuration, err := time.ParseDuration(splitDurationStr)
		if err != nil {
			return err
		}
		globalConfig.splitDuration = splitDuration
	}

	globalConfig.deleteDuration = time.Second * 5
	if deleteDurationStr != "" {
		deleteDuration, err := time.ParseDuration(deleteDurationStr)
		if err != nil {
			return err
		}
		globalConfig.deleteDuration = deleteDuration
	}

	globalConfig.isFile = isFile
	path, err := getCurrentPath()
	if err != nil {
		return err
	}
	globalConfig.fileDir = path
	if isFile {
		if err := globalConfig.initLogFile(); err != nil {
			return err
		}
	}
	globalConfig.exitChan = make(chan struct{})

	return nil
}

// attach observer
func (f *logConfig) Attach(observer *logging) {
	f.observersMu.Lock()
	defer f.observersMu.Unlock()
	f.observers = append(f.observers, observer)
}

// detach observer
func (f *logConfig) Detach(observer *logging) {
	f.observersMu.Lock()
	defer f.observersMu.Unlock()
	for i := 0; i < len(f.observers); {
		//Compare if it is the same object
		if reflect.DeepEqual(f.observers[i], observer) {
			f.observers = append(f.observers[:i], f.observers[i+1:]...)
			continue
		}
		i++
	}
}

// It is recommended to call this method on the last line of the main function, in order to notify all observers
func Notify() {
	globalConfig.Notify()
}

// notify observer,this method is not thread safe
// and needs to be explicitly locked when called
func (f *logConfig) Notify() {
	if !f.isFile {
		return
	}

	for _, observer := range f.observers {
		observer.SetOutPut(f.file)
	}
}

// Enable log write to file mode
func Start() {
	globalConfig.Start()
}

func Stop() {
	globalConfig.Stop()
}

func (f *logConfig) Start() {
	if !f.isFile {
		return
	}

	go f.start()
}

func (f *logConfig) Stop() {
	if !f.isFile {
		return
	}

	f.exitChan <- struct{}{}
}

func (f *logConfig) initLogFile() error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("init log file  panic: %v\n", err)
		}
	}()

	fileName := generateFileName(globalConfig.prefix)

	logFile := joinFilePath(globalConfig.fileDir, fileName)
	if !isExist(globalConfig.fileDir) {
		if err := os.Mkdir(globalConfig.fileDir, 0755); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("open file '%s' fail: %v", logFile, err)
		return fmt.Errorf("open file '%s' fail: %v", logFile, err)
	}
	f.fileMu.Lock()
	f.file = file
	//Used to modify the log object during initialization to modify its path
	f.Notify()
	f.fileMu.Unlock()
	return nil
}

func (f *logConfig) start() {
	splitTicker := time.NewTicker(f.splitDuration)
	defer splitTicker.Stop()

	deleteTicker := time.NewTicker(f.deleteDuration)
	defer deleteTicker.Stop()
	express := fmt.Sprintf(`%s_\d{8}_\d{4}`, globalConfig.prefix)
	reg := regexp.MustCompile(express)

	for {
		select {
		case <-splitTicker.C:
			if err := f.splitOnce(); err != nil {
				log.Fatalf("split log failed, %v", err)
			}
		case <-deleteTicker.C:
			if err := f.deleteOnce(reg); err != nil {
				log.Fatalf("delete log failed, %v", err)
			}
		case <-f.exitChan:
			f.observersMu.Lock()

			for _, observer := range f.observers {
				observer.Close()
			}

			f.observersMu.Unlock()
			return
		}
	}
}

func (f *logConfig) splitOnce() error {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("log split panic: %v\n", err)
		}
	}()

	fi, err := f.file.Stat()
	if err != nil {
		return err
	}

	diff := fi.Size() - globalConfig.maxSize
	if diff < 0 {
		return nil
	}

	fileName := generateFileName(globalConfig.prefix)
	logFile := joinFilePath(globalConfig.fileDir, fileName)
	if !isExist(globalConfig.fileDir) {
		if err := os.Mkdir(globalConfig.fileDir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	oldFile := f.file
	f.file = file
	f.Notify()

	return oldFile.Close()
}

func (f *logConfig) deleteOnce(reg *regexp.Regexp) error {
	f.fileMu.RLock()
	defer f.fileMu.RUnlock()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("delete file once panic: %v\n", err)
		}
	}()

	fis, err := ioutil.ReadDir(globalConfig.fileDir)
	if err != nil {
		return err
	}

	fi, err := f.file.Stat()
	if err != nil {
		return err
	}
	currFileName := fi.Name()
	for _, fi := range fis {
		if fi.Name() == currFileName {
			continue
		}
		fileName := fi.Name()
		if !reg.Match([]byte(fileName)) {
			continue
		}
		now := time.Now().Unix()
		modify := fi.ModTime().Unix()
		diff := now - modify
		if diff < globalConfig.maxSecond {
			continue
		}
		filename := joinFilePath(globalConfig.fileDir, fi.Name())
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("delete file '%s' err:%v", fi.Name(), err)
		}

	}
	return nil
}
