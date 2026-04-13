package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// logFile is a pointer (*os.File) which stores the memory address of the opened 
// file. This allows different functions (LogMsg, CloseLogger) to share the 
// same active file handle efficiently without copying the entire file object.
var logFile *os.File

var LOG_FILE_PATH string 

func InitLogger() error {
	// Capture the current working directory where the shell was launched
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	// Build the absolute path for the log file
	LOG_FILE_PATH = filepath.Join(root, "local", "shell.log")

	// Debug logs to terminal (since logFile isn't open yet)
	fmt.Printf("DEBUG: Target log path: %s\n", LOG_FILE_PATH)

	// Ensure the 'local' directory exists
	err = os.MkdirAll(filepath.Join(root, "local"), 0755)
	if err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Open the file
	f, err := os.OpenFile(LOG_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Assign the pointer to the global variable so LogMsg can start using it
	logFile = f

	// 5. Log internal messages to verify everything is working
	LogMsg("--- Logger Session Started ---")
	LogMsg(fmt.Sprintf("Current Working Directory (Root): %s", root))
	LogMsg(fmt.Sprintf("Log file successfully opened at: %s", LOG_FILE_PATH))

	return nil
}

func LogMsg(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, msg)
	
	// We check if the pointer is NOT nil to avoid a crash if LogMsg 
	// is called before InitLogger or after CloseLogger.
	if logFile != nil {
		logFile.WriteString(logLine)
	}
}

func CloseLogger() {
	if logFile != nil {
		LogMsg("--- Logger Session Closed ---")
		logFile.Close()
		logFile = nil // Reset pointer to nil for safety
	}
}