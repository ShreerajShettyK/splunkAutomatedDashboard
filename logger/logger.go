package logger

import (
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

// // check if log folder exists
func logFolderExists() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	FolderPath = filepath.Join(wd, "logger", "logs")
	// fmt.Println("FolderPath:", FolderPath)

	if err := os.MkdirAll(FolderPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating log folder: %v", err)
	}
}

func createLogFile(filename string) *infraLogger {
	file, err := os.OpenFile(filepath.Join(FolderPath, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}

	return &infraLogger{
		fileName: filename,
		Logger:   log.New(file, time.Now().Format("2006-01-02 15:04:05 ")+": Log: ", log.Lshortfile),
	}
}

func CreateLogger() *infraLogger {
	logFolderExists()
	once.Do(func() {
		infralogger = createLogFile(time.Now().Format("2006-01-02") + "splunk-dashboard-creation.log")
	})
	return infralogger
}
