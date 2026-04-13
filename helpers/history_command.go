package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SHELL_HISTORY_FILE stores the absolute path to the history file.
// We keep this global so it's calculated only once.
var SHELL_HISTORY_FILE string

func getHistoryPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".simple_shell_history")
}

// ensureHistoryPath ensures our global path is set before any file operations.
// each function that interacts with history should call this first to guarantee 
// the path is initialized. we dont need to wory about each function calling it multiple times 
// as the file is short lived and the check is lightweight.

func ensureHistoryPath() {
	if SHELL_HISTORY_FILE == "" {
		SHELL_HISTORY_FILE = getHistoryPath()
	}
}

func StoreHistory(input string) error {
	ensureHistoryPath()
	LogMsg(fmt.Sprintf("Storing command in history: %s", input))

	var lines []string
	data, err := os.ReadFile(SHELL_HISTORY_FILE)
	if err == nil {
		content := strings.TrimSpace(string(data))
		if content != "" {
			lines = strings.Split(content, "\n")
		}
	}

	// Handle Truncation
	if len(lines) >= 1000 {
		LogMsg("History file exceeded 1000 entries, truncating")
		lines = lines[len(lines)-999:]
		newContent := strings.Join(lines, "\n") + "\n"
		os.WriteFile(SHELL_HISTORY_FILE, []byte(newContent), 0644)
	}

	// Use = instead of := to avoid shadowing if you choose to keep the global var,
	// but here we just use a local 'f' for clarity.
	f, err := os.OpenFile(SHELL_HISTORY_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogMsg(fmt.Sprintf("Error opening history file: %v", err))
		return err
	}
	defer f.Close()

	newLine := fmt.Sprintf("%d %s\n", len(lines)+1, input)
	if _, err := f.WriteString(newLine); err != nil {
		LogMsg(fmt.Sprintf("Error writing to history file: %v", err))
		return err
	}

	return nil
}

// ShowHistory reads the history file and prints its contents to the terminal.
func ShowHistory() error {
	ensureHistoryPath()
	LogMsg("Showing command history")
	
	data, err := os.ReadFile(SHELL_HISTORY_FILE)
	if err != nil {
		return fmt.Errorf("history is empty or file missing")
	}
	
	fmt.Print(string(data))
	return nil
}

// MoveBetweenHistoryLines allows navigating through the command history based on the current line index 
// and the direction (up/down). 
// It returns the new line index and the corresponding command from the history.
func MoveBetweenHistoryLines(currentLine int, direction string) (int, string, error) {
	ensureHistoryPath()
	data, err := os.ReadFile(SHELL_HISTORY_FILE)
	if err != nil {
		return currentLine, "", err
	}
	
	content := strings.TrimSpace(string(data))
	if content == "" {
		return 0, "", nil
	}
	
	lines := strings.Split(content, "\n")
	totalLines := len(lines)
	newIdx := currentLine

	if direction == "up" {
		if newIdx > 0 {
			newIdx--
		}
	} else if direction == "down" {
		if newIdx < totalLines {
			newIdx++
		}
	}

	if newIdx >= totalLines {
		return totalLines, "", nil
	}

	targetLine := lines[newIdx]
	parts := strings.SplitN(targetLine, " ", 2)
	cmd := targetLine
	if len(parts) > 1 {
		cmd = parts[1]
	}

	return newIdx, cmd, nil
}

// GetCurrentHistoryLine returns the current line index in the history file, 
// which is essentially the count of commands stored. 
// This is used to navigate through history with up/down keys.
func GetCurrentHistoryLine() (int, error) {
	ensureHistoryPath()
	data, err := os.ReadFile(SHELL_HISTORY_FILE)
	if err != nil {
		return 0, nil 
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return 0, nil
	}
	lines := strings.Split(content, "\n")
	return len(lines), nil 
}