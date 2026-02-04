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

	// Get the prompt string
	hostname, err := helpers.BuildPrefix(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
		return
	}

	var inputBuffer string
	// Get initial history line count
	historyIndex, _ := helpers.GetCurrentHistoryLine()

	fmt.Print(hostname + " ")

	// Start the keyboard listener
	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			fmt.Println()
			trimmed := strings.TrimSpace(inputBuffer)
			if trimmed != "" {
				helpers.ExecuteInput(trimmed)
			}
			inputBuffer = ""
			// Refresh history index after execution
			historyIndex, _ = helpers.GetCurrentHistoryLine()
			fmt.Print(hostname + " ")

		case keys.Up:
			newIdx, cmd, err := helpers.MoveBetweenHistoryLines(historyIndex, "up")
			if err == nil {
				historyIndex = newIdx
				// \r = start of line, \x1b[K = clear to end of line
				fmt.Print("\r" + hostname + " \x1b[K")
				inputBuffer = cmd
				fmt.Print(inputBuffer)
			}

		case keys.Down:
			newIdx, cmd, err := helpers.MoveBetweenHistoryLines(historyIndex, "down")
			if err == nil {
				historyIndex = newIdx
				fmt.Print("\r" + hostname + " \x1b[K")
				inputBuffer = cmd
				fmt.Print(inputBuffer)
			}

		case keys.Backspace:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
				fmt.Print("\b \b") // Visual erase
			}

		case keys.RuneKey:
			inputBuffer += key.String()
			fmt.Print(key.String())

		case keys.CtrlC:
			fmt.Println("\nExiting shell...")
			return true, nil // Stop listener
		}
		return false, nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
