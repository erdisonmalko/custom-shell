package main

import (
 "bufio"
 "fmt"
 "os"
 "os/exec"
 "strings"
)

func main() {
    //global reader for stdin
    reader := bufio.NewReader(os.Stdin)

    hostname, err := build_prefix(reader)
    if err != nil {
        fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
        return
    }

    for {
        fmt.Print(hostname + " ")
        input, err := reader.ReadString('\n')

        if err != nil {
            fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
        }

        if err := executeInput(input); err != nil {
            fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
        }
    }

}

func build_prefix(reader *bufio.Reader) (string, error) {
    cmd := exec.Command("whoami")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    cmd1 := exec.Command("hostname")
    output1, err1 := cmd1.Output()
    if err1 != nil {
        return "", err1
    }

    hostname := strings.TrimSuffix(string(output), "\n") + "@" + strings.TrimSuffix(string(output1), "\n")
    
    cmd2 := exec.Command("pwd")
    output, err2 := cmd2.Output()
    if err2 != nil {
        return "", err2
    }
    pwd_string := strings.TrimSuffix(string(output), "\n")
    pwd_string = strings.ReplaceAll(pwd_string, os.Getenv("HOME"), "~")
    pwd_string += "$"

    hostname = fmt.Sprintf("➜ %s: %s ", hostname, pwd_string)
    
    return hostname, nil
}

func store_history(input string) error {
    // Open the history file in append mode, create it if it doesn't exist
    file, err := os.OpenFile(".simple_shell_history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    
    counter := 0
    // Count existing lines in the history file
    data, err := os.ReadFile(".simple_shell_history")
    if err == nil {
        counter = strings.Count(string(data), "\n")
    }
    
    // Write the input to the file with a newline and counter
    if _, err := file.WriteString(fmt.Sprintf("%d %s\n", counter+1, input)); err != nil {
        return err
    }

    return nil
}

func show_history() error {
    data, err := os.ReadFile(".simple_shell_history")
    if err != nil {
        return err
    }
    fmt.Print(string(data))
    return nil
}


func get_most_possible_command(input string) (string, error) {
    
    data, err := os.ReadFile(".simple_shell_history")
    if err != nil {
        return "", err
    }
    lines := strings.Split(string(data), "\n")
    var possible_commands []string

    for _, line := range lines {
        if strings.HasPrefix(line, input) {
            // Extract the command part after the line number
            parts := strings.SplitN(line, " ", 2)
            if len(parts) == 2 {
                possible_commands = append(possible_commands, parts[1])
            }
        }
    }


    if len(possible_commands) == 0 {
        return "", nil
    }

    // Return the first matched command as the most possible command
    return possible_commands[0], nil
}

func executeInput(input string) error {
    input = strings.TrimSuffix(input, "\n")
    args := strings.Split(input, " ")

    //store command in history
    if err := store_history(input); err != nil {
        return err
    }

    switch args[0] {
        //in case of empty input, just return
        // !, @, #, $, %, ^, &, *, (, ), _, +, {, }, |, :, ", <, >, ?, ~, \
        case "", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "{", "}", "|", ":", "\"", "<", ">", "?", "~", "\\":
        return nil
        case strings.Contains(input, "\t"):
        fmt.Println("tab completion functionality is not implemented yet.")
        return nil
        case "history":
        if err := show_history(); err != nil {
            return err
        }
        return nil
        //will work on tab completion later 
        case "9 0011 0x09":
        if possible_cmd, err := get_most_possible_command(input); err != nil {
            return err
        } else {
            fmt.Println(possible_cmd)
        }
        return nil
        // will work on arrow keys later, up and down
        case "\x1b[A", "\x1b[B":
        fmt.Println("arrow key functionality is not implemented yet.")
        return nil
        //exit command to exit the shell
        case "exit":
        os.Exit(0)
    }
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    return cmd.Run()

}