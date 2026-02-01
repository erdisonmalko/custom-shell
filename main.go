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



func executeInput(input string) error {
    input = strings.TrimSuffix(input, "\n")
    args := strings.Split(input, " ")
    switch args[0] {
        //in case of empty input, just return
        case "":
        return nil
        case "history":
        fmt.Println("history command is not implemented yet.")
        return nil
        //will work on tab completion later
        case "tab":
        fmt.Println("tab completion is not implemented yet.")
        return nil
        // will work on arrow keys later, up and down
        case "arrow":
        fmt.Println("arrow key functionality is not implemented yet.")
        return nil
        // SHELL ➜  exec: "\t": executable file not found in $PATH
        //this errors happen as the input is not sanitized properly
        //exit command to exit the shell
        case "exit":
        os.Exit(0)
    }
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    return cmd.Run()

}