package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type infraLogger struct {
	*log.Logger
	fileName string
}

var (
	once        sync.Once
	FolderPath  string
	infralogger *infraLogger
	Logger      = CreateLogger()
)

// check if log folder exists
func logFolderExists() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	log.Println(wd)
	parentDir := filepath.Dir(wd)
	baseDir := filepath.Dir(parentDir)
	_, err = os.Stat(filepath.Join(baseDir, "logger"))
	if err != nil {
		baseDir = filepath.Dir(baseDir)
	}
	log.Println(baseDir)
	FolderPath = baseDir + "/logger/logs/"
	logfolderErr := os.MkdirAll(FolderPath, os.ModePerm)
	if logfolderErr != nil {
		fmt.Println(logfolderErr)
	}
}

// create log file if not exist
func createLogFile(filename string) *infraLogger {
	file, _ := os.OpenFile(FolderPath+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	return &infraLogger{
		fileName: filename,
		Logger:   log.New(file, time.Now().Format("2006-01-02 15:04:05")+": Log: ", log.Lshortfile),
	}
}

func CreateLogger() *infraLogger {
	logFolderExists()
	once.Do(func() {
		infralogger = createLogFile(time.Now().Format("2006-01-02") + "splunk-dashboard-creation.log")
	})
	return infralogger
}
