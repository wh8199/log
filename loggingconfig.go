package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	prefix         string             //prefix used to generate log file names
	observers      []LoggingInterface //Using the observer pattern
	observersMu    sync.Mutex
	file           *os.File     //The current output file of the log. Since multiple goroutines share this resource, a read-write lock needs to be added.
	fileMu         sync.RWMutex //Used to protect mutually exclusive resources file
}

//the fileDir is the log save path, default value is current path
//
//the prefix is generated log filename prefix, the default value is log
//
//the maxSize default is 20,unit is MB
//
//the maxSecond is seconds for log retention, after which it will be automatically deleted
func SetLogConfig(isFile bool, fileDir, prefix, splitDurationStr, deleteDurationStr string, maxSize, maxSecond int) error {
	globalConfig.isFile = true
	if fileDir == "" {
		globalConfig.fileDir = getCurrentPath()
	}
	if prefix == "" {
		globalConfig.prefix = "log"
	}
	if maxSize < 20 {
		globalConfig.maxSize = 20
	}
	if maxSecond < 1 {
		globalConfig.maxSecond = 1
	}
	globalConfig.maxSize = globalConfig.maxSize << 20 //MB
	if splitDurationStr == "" {
		splitDurationStr = "5s"
	}
	if deleteDurationStr == "" {
		deleteDurationStr = "5s"
	}

	splitDuration, err := time.ParseDuration(splitDurationStr)
	if err != nil {
		return err
	}
	globalConfig.splitDuration = splitDuration

	deleteDuration, err := time.ParseDuration(deleteDurationStr)
	if err != nil {
		return err
	}
	globalConfig.deleteDuration = deleteDuration

	if err := globalConfig.initLogFile(); err != nil {
		return err
	}
	return nil
}

//attach observer
func (f *logConfig) Attach(observer LoggingInterface) {
	f.observersMu.Lock()
	defer f.observersMu.Unlock()
	if !f.isFile {
		return
	}
	f.observers = append(f.observers, observer)
}

//detach observer
func (f *logConfig) Detach(observer LoggingInterface) {
	f.observersMu.Lock()
	defer f.observersMu.Unlock()
	if !f.isFile {
		return
	}
	for i := 0; i < len(f.observers); {
		if f.observers[i] == observer {
			f.observers = append(f.observers[:i], f.observers[i+1:]...)
		} else {
			i++
		}
	}
}

//notify observer,this method is not thread safe
//and needs to be explicitly locked when called
func (f *logConfig) Notify() {
	if !f.isFile {
		return
	}
	for _, observer := range f.observers {
		observer.SetOutPut(f.file)
	}
}

//Enable log write to file mode
func Start() {
	globalConfig.Start()
}

func (f *logConfig) Start() {
	if !f.isFile {
		return
	}
	go f.logSplit()
	go f.logDelete()
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
	//This method is not called during log initialization,
	//because during log initialization,
	//there is no guarantee that all log objects have been registered in the observer registry.
	//You should check the configuration during initialization to see if the file is empty
	//f.Notify()
	f.fileMu.Unlock()
	return nil
}

func (f *logConfig) logSplit() error {
	ticker := time.NewTicker(f.splitDuration)
	defer ticker.Stop()
	for {
		<-ticker.C
		if err := f.splitOnce(); err != nil {
			log.Printf("%v", err)
			return err
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
	f.file.Close()
	f.file = file
	f.Notify()
	return nil
}

func (f *logConfig) logDelete() {
	ticker := time.NewTicker(f.deleteDuration)
	defer ticker.Stop()
	express := fmt.Sprintf(`%s_\d{8}_\d{4}`, globalConfig.prefix)
	reg := regexp.MustCompile(express)
	for {
		<-ticker.C
		if err := f.deleteOnce(reg); err != nil {
			log.Printf("%v", err)
		}
	}

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
