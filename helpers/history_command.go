package helpers

import (
 "fmt"
 "os"
 "strings"
)

const SHELL_HISTORY_FILE = "./local/.simple_shell_history"

func StoreHistory(input string) error {
	LogMsg(fmt.Sprintf("Storing command in history: %s", input))
	// Open the history file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(SHELL_HISTORY_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogMsg(fmt.Sprintf("Error opening history file: %v", err))
		return err
	}
	defer file.Close()
    
    counter := 0
    // Count existing lines in the history file
    data, err := os.ReadFile(SHELL_HISTORY_FILE)
    if err == nil {
        counter = strings.Count(string(data), "\n")
    }
    
    // Write the input to the file with a newline and counter
    if _, err := file.WriteString(fmt.Sprintf("%d %s\n", counter+1, input)); err != nil {
        return err
    }

    return nil
}

func ShowHistory() error {
	LogMsg("Showing command history")
	data, err := os.ReadFile(SHELL_HISTORY_FILE)
	if err != nil {
		LogMsg(fmt.Sprintf("Error reading history file: %v", err))
		return err
	}
	fmt.Print(string(data))
	LogMsg(fmt.Sprintf("Displayed %d bytes of history", len(data)))
	return nil
}
