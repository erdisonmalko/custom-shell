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
    // maybe having big history file is not optimal, so let's limit it to last 1000 commands
	if counter >= 1000 {
		LogMsg("History file exceeded 1000 entries, truncating")
		lines := strings.Split(string(data), "\n")
		if len(lines) > 1000 {
			lines = lines[len(lines)-999:] // keep last 999 lines
		}
		err = os.WriteFile(SHELL_HISTORY_FILE, []byte(strings.Join(lines, "\n")), 0644)
		if err != nil {
			LogMsg(fmt.Sprintf("Error truncating history file: %v", err))
			return err
		}
		counter = len(lines)
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

// MoveBetweenHistoryLines allows navigating through the command history based on the current line index 
// and the direction (up/down). 
// It returns the new line index and the corresponding command from the history.
func MoveBetweenHistoryLines(currentLine int, direction string) (int, string, error) {
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
        // If we are at the very end (totalLines), first 'up' goes to totalLines - 1
        if newIdx > 0 {
            newIdx--
        }
    } else if direction == "down" {
        if newIdx < totalLines {
            newIdx++
        }
    }

    // If 'down' takes us back to the bottom (where the user was typing something new)
    if newIdx >= totalLines {
        return totalLines, "", nil
    }

    // Extract command (removing the "1 " prefix from your file format)
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
    data, err := os.ReadFile(SHELL_HISTORY_FILE)
    if err != nil {
        return 0, nil // Return 0 if file doesn't exist yet
    }
    content := strings.TrimSpace(string(data))
    if content == "" {
        return 0, nil
    }
    lines := strings.Split(content, "\n")
    return len(lines), nil 
}