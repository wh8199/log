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
	isFile      bool
	fileDir     string
	maxSize     int64
	maxHour     int64
	prefix      string             //prefix used to generate log file names
	observers   []LoggingInterface //Using the observer pattern
	observersMu sync.Mutex
	file        *os.File     //The current output file of the log. Since multiple goroutines share this resource, a read-write lock needs to be added.
	fileMu      sync.RWMutex //Used to protect mutually exclusive resources file
}

func init() {
	globalConfig = &logConfig{
		isFile:      false,
		observers:   make([]LoggingInterface, 0, 1000),
		observersMu: sync.Mutex{},
		fileMu:      sync.RWMutex{},
	}
}

//the fileDir is the log save path, default value is current path
//
//the prefix is generated log filename prefix, the default value is log
//
//the maxSize default is 20,unit is MB
//
//the maxHour is hours for log retention, after which it will be automatically deleted
func SetLogConfig(isFile bool, fileDir, prefix string, maxSize, maxHour int) {
	globalConfig.isFile = true
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

//notify observer,this method is not thread safe and needs to be explicitly locked when called
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
	for {
		if err := f.initLogFile(); err != nil {
			log.Print(err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
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
	f.Notify()
	f.fileMu.Unlock()
	return nil
}

func (f *logConfig) logSplit() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f.splitOnce(); err != nil {
				log.Printf("%v", err)
			}
		}
	}
}
func (f *logConfig) splitOnce() error {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("log writer panic: %v\n", err)
		}
	}()

	fi, err := f.file.Stat()
	if err != nil {
		return err
	}
	if fi.Size() >= globalConfig.maxSize {
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
	return nil
}

func (f *logConfig) logDelete() {
	ticker := time.NewTicker(5 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f.deleteOnce(); err != nil {
				log.Printf("%v", err)
			}
		}
	}

}

func (f *logConfig) deleteOnce() error {
	f.fileMu.RLock()
	defer f.fileMu.Unlock()
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
		now := time.Now().Unix()
		modify := fi.ModTime().Unix()
		if now-modify > globalConfig.maxHour*60*60 {
			continue
		}
		express := fmt.Sprintf(`%s_\d{8}_\d{4}`, globalConfig.prefix)
		reg := regexp.MustCompile(express)
		fileName := fi.Name()
		if reg.Match([]byte(fileName)) {
			if err := os.Remove(joinFilePath(globalConfig.fileDir, fi.Name())); err != nil {
				return fmt.Errorf("delete file '%s' err:%v", fi.Name(), err)
			}
		}
	}
	return nil
}
