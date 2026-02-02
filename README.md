# custom-shell

A lightweight shell implementation written in Go that provides basic command execution, history tracking, and file path suggestions.

## Features

- **Command Execution**: Execute system commands with arguments
- **Built-in Commands**:
  - `cd <path>` - Change directory using `os.Chdir()`
  - `exit` - Exit the shell gracefully
  - `history` - Display all executed commands with timestamps
  
- **File Path Suggestions**: Show file/directory suggestions when running:
  - `ls` (without arguments)
  - `cd` (without arguments)
  - `cat`, `vim`, `nano` (without arguments)

- **Command History**: All commands are logged to `.simple_shell_history` with line numbers

- **Logging**: Every shell action is logged to `./local/shell.log` with timestamps for debugging

- **Tilde Expansion**: Automatically converts home directory paths to `~` in the prompt

## Project Structure

```
custom-shell/
├── main.go              # Main entry point
├── go.mod              # Go module definition
├── helpers/            # Helper functions package
│   ├── log.go          # Logging functionality
│   ├── utils.go        # Utility functions (BuildPrefix, CompletePath)
│   ├── historyCommand.go # History management
│   └── executeInput.go  # Command execution logic
├── local/              # Local storage directory
│   ├── shell.log       # Shell activity logs
│   └── .simple_shell_history # Command history
└── README.md           # This file
```

## Building & Running

### Prerequisites
- Go 1.18 or higher
- Linux/macOS/Windows (any OS with Go support)

### Build
```bash
cd custom-shell/
go build
```

### Run
```bash
./custom-shell
```

## How It Works

1. **Initialization**: Starts the logger, gets current user/hostname/directory
2. **Main Loop**: 
   - Displays prompt with username@hostname:pwd$ format
   - Reads user input
   - Stores input in history
   - Executes command or shows suggestions if no arguments

3. **Command Handling**:
   - Special characters are ignored
   - Built-in commands are handled specially (cd, exit, history)
   - External commands are executed via `exec.Command`
   - File suggestions are shown for file-related commands with no arguments

4. **Logging**: Every step is logged with timestamps to `local/shell.log`

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

**helpers/historyCommand.go**:
- `StoreHistory()` - Save commands to history file
- `ShowHistory()` - Display command history

**helpers/executeInput.go**:
- `ExecuteInput()` - Main command execution logic
- Handles cd, exit, history, and external commands

## Future Enhancements

- Arrow key history navigation
- Pipe (`|`) and redirection (`>`, `>>`) support