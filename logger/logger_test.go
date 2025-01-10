package logger

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var filename = "execution_log.log"

func TestCreateLogger(t *testing.T) {
	t.Run("test to validate log folder", func(t *testing.T) {
		FolderPath = "./logs/"
		logFolderExists()
		_, err := os.Stat(FolderPath)
		assert.False(t, os.IsNotExist(err))
	})
	t.Run("test for log file creation", func(t *testing.T) {
		infralogger = createLogFile(filename)
		assert.NotNil(t, infralogger)
	})
	t.Run("test for logging in log file", func(t *testing.T) {
		logMsg := "This is a log message"
		infralogger.Print(logMsg)
		time.Sleep(1 * time.Second)

		file, err := os.Open(FolderPath + filename)
		assert.Nil(t, err)

		fileInfo, err := file.Stat()
		assert.Nil(t, err)
		buffer := make([]byte, fileInfo.Size())
		_, err = file.Read(buffer)
		assert.Nil(t, err)
		assert.Contains(t, string(buffer), logMsg)
	})
}
