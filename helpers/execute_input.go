package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"context"
    "os/signal"
    "syscall"
)

func ExecuteInput(input string) error {
	input = strings.TrimSuffix(input, "\n")
	LogMsg(fmt.Sprintf("Executing input: %s", input))
	args := strings.Fields(input)
	if len(args) == 0 {
		LogMsg("Empty input, returning")
		return nil
	}

	// store command in history
	if err := StoreHistory(input); err != nil {
		LogMsg(fmt.Sprintf("Error storing history: %v", err))
		return err
	}

	cmdName := args[0]
	LogMsg(fmt.Sprintf("Command: %s with %d args", cmdName, len(args)-1))
	
	// Show suggestions only if command has no arguments
	if len(args) == 1 {
		fileCmds := map[string]bool{"ls": true, "cat": true, "cd": true, "vim": true, "nano": true}
		if fileCmds[cmdName] {
			LogMsg(fmt.Sprintf("File command '%s' with no args, showing suggestions", cmdName))
			
			matches, err := CompletePath("")
			if err == nil && len(matches) > 0 {
				LogMsg(fmt.Sprintf("Found %d completion suggestions", len(matches)))
				fmt.Println("Suggestions:")
				for _, m := range matches {
					fmt.Println("  " + m)
				}
			} else if err != nil {
				LogMsg(fmt.Sprintf("Error getting completions: %v", err))
			}
			return nil
		}
	}

	switch cmdName {
		case "", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "{", "}", "|", ":", "\"", "<", ">", "?", "~", "\\":
			LogMsg("Skipping special character input")
			return nil
		case "history":
			return ShowHistory()
		case "cd":
			// Handle cd specially - use os.Chdir instead of exec.Command
			if len(args) < 2 {
				LogMsg("cd: no argument provided")
				return fmt.Errorf("cd: no argument provided")
			}
			path := args[1]
			LogMsg(fmt.Sprintf("Changing directory to: %s", path))
			err := os.Chdir(path)
			if err != nil {
				LogMsg(fmt.Sprintf("cd failed: %v", err))
			} else {
				LogMsg(fmt.Sprintf("Successfully changed directory to: %s", path))
			}
			return err
		case "exit":
		LogMsg("Exit command received")
		return fmt.Errorf("exit_requested") // Return a custom sentinel error
	}

	// Create a context that is naturally cancelled if we receive an Interrupt (Ctrl+C)
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

	LogMsg(fmt.Sprintf("Running external command: %s", args[0]))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin // Crucial for interactive commands like 'top' or 'vim'
	
	err := cmd.Run()
	
    // Check if the command was killed by a signal
    if ctx.Err() != nil {
        fmt.Println("^C")
        return nil 
    }

	if err != nil {
		LogMsg(fmt.Sprintf("Command '%s' failed: %v", args[0], err))
	} else {
		LogMsg(fmt.Sprintf("Command '%s' completed successfully", args[0]))
	}
	
    return err
}
