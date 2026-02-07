package main

import (
	"custom-shell/helpers"
	"fmt"
	"os"
	"strings"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func main() {
	if err := helpers.InitLogger(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize logger:", err)
		return
	}
	defer helpers.CloseLogger()

	helpers.LogMsg("Shell started")

	hostname, err := helpers.BuildPrefix(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
		return
	}

	var inputBuffer string
	// historyIndex represents our current position in history
	// It starts at the very end (total lines)
	historyIndex, _ := helpers.GetCurrentHistoryLine()

	fmt.Print(hostname + " ")

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			fmt.Println()
			trimmed := strings.TrimSpace(inputBuffer)
			if trimmed != "" {
				// Check if the user typed "exit"
				if trimmed == "exit" {
					fmt.Println("Goodbye!")
					return true, nil // This stops the listener and cleans up the terminal
				}

				err := helpers.ExecuteInput(trimmed)
				if err != nil {
					// If ExecuteInput returns our sentinel error, exit gracefully
					if err.Error() == "exit_requested" {
						return true, nil
					}
					fmt.Println("Error:", err)
				}
			}
			inputBuffer = ""
			historyIndex, _ = helpers.GetCurrentHistoryLine()
			fmt.Print(hostname + " ")
		case keys.Up:
			newIdx, cmd, err := helpers.MoveBetweenHistoryLines(historyIndex, "up")
			if err == nil {
				historyIndex = newIdx // Update state so next press knows where we are
				inputBuffer = cmd
				// \r = start of line, \x1b[K = clear to end of line
				fmt.Print("\r" + hostname + " \x1b[K" + inputBuffer)
			}

		case keys.Down:
			newIdx, cmd, err := helpers.MoveBetweenHistoryLines(historyIndex, "down")
			if err == nil {
				historyIndex = newIdx // Update state
				inputBuffer = cmd
				fmt.Print("\r" + hostname + " \x1b[K" + inputBuffer)
			}

		case keys.Backspace:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
				fmt.Print("\b \b") 
			}

		case keys.RuneKey:
			// If we were scrolling history and start typing, 
			// we are effectively creating a new command.
			inputBuffer += key.String()
			fmt.Print(key.String())

		case keys.Space:
			inputBuffer += " "
			fmt.Print(" ")
		case keys.CtrlC:
			// Instead of exiting the whole program, just clear the line
			inputBuffer = ""
			historyIndex, _ = helpers.GetCurrentHistoryLine()
			fmt.Print("\n" + hostname + " ")
			return false, nil
		default:
			// If it's not a control key we handled above, check if it's printable
			if key.Code == keys.Space || key.Code == keys.RuneKey {
				inputBuffer += key.String()
				fmt.Print(key.String())
			}

		}
		return false, nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}