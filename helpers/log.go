package helpers

import (
 "fmt"
 "os"
 "time"
)

var logFile *os.File

const LOG_FILE_PATH = "./local/shell.log"

func InitLogger() error {
	var err error
	logFile, err = os.OpenFile(LOG_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LogMsg(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, msg)
	if logFile != nil {
		logFile.WriteString(logLine)
	}
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}