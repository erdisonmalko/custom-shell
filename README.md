# custom-shell

A lightweight shell implementation written in Go that provides basic command execution, history tracking, and file path suggestions.

## Features

* **Real-time Key Listening**: Uses raw terminal mode to capture every keystroke, allowing for immediate feedback.
* **Interactive History Scroll**: Use the **Up** and **Down** arrow keys to traverse your command history.
* **Signal Awareness**: Intercepts `Ctrl+C` to clear the current buffer or interrupt running processes without crashing the shell.
* **Smart Path Suggestions**: Provides contextual file/directory hints when typing `ls`, `cat`, `cd`, `vim`, or `nano` without arguments.

* **Automated Logging**: Full audit trail of every command and internal shell state in `local/shell.log`.

- **Tilde Expansion**: Automatically converts home directory paths to `~` in the prompt

## Project Structure

```
custom-shell/
├── main.go              # Main entry point
├── go.mod              # Go module definition
├── helpers/            # Helper functions package
│   ├── log.go          # Logging functionality
│   ├── utils.go        # Utility functions (BuildPrefix, CompletePath)
│   ├── history_command.go # History management
│   └── execute_input.go  # Command execution logic
├── local/              # Local storage directory
│   ├── shell.log       # Shell activity logs
│   └── .simple_shell_history # Command history
└── README.md           # This file
```

## Building & Running

### Prerequisites

* Go 1.18 or higher
* A terminal supporting ANSI escape codes (Linux/macOS/WSL)

### Installation

```bash
# Clone the repository and sync dependencies
go mod tidy

# Build the binary
go build -o my-shell

# Run
./my-shell

```

##  How It Works

### 1. Terminal State Management

The shell uses **Raw Mode** during the input phase. This allows us to detect arrow keys and backspaces instantly. When a command is executed, the shell coordinates with the OS to manage standard input/output streams.

### 2. The Keyboard Listener

Inside `main.go`, a continuous listener evaluates `keys.Key` events:

* **Enter**: Finalizes the `inputBuffer` and sends it to the Execution Engine.
* **Up/Down Arrows**: Triggers `MoveBetweenHistoryLines`, updating the buffer with previous commands while clearing the current line using ANSI escape codes (`\x1b[K`).
* **Backspace**: Manually updates the buffer and handles the visual "erase" on the terminal using `\b \b`.
* **Space/Runes**: Appends printable characters to the active buffer.

### 3. Execution & Signal Handling

When a command runs, the shell creates a `signal.NotifyContext`. If you press `Ctrl+C` while a command is running, the shell sends a `SIGTERM` to the child process only, keeping your shell session alive.

### 4. **Logging**: Every step is logged with timestamps to `local/shell.log`

## Package Organization

**main.go**: 
- Initializes logger and shell loop
- Reads user input

**helpers/log.go**:
- `InitLogger()` - Setup logging
- `LogMsg()` - Write timestamped logs
- `CloseLogger()` - Cleanup

**helpers/utils.go**:
- `BuildPrefix()` - Build the shell prompt
- `CompletePath()` - Get file/directory suggestions

**helpers/history_command.go**:
- `StoreHistory()` - Save commands to history file
- `ShowHistory()` - Display command history
- `GetCurrentHistoryLine()` - Get last inserted command on history
- `MoveBetweenHistoryLines()` - Move between the lines of history

**helpers/execute_input.go**:
- `ExecuteInput()` - Main command execution logic
- Handles cd, exit, history, and external commands