package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"custom-shell/helpers"
)

func main() {
	if err := helpers.InitLogger(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize logger:", err)
		return
	}
	defer helpers.CloseLogger()
	
	helpers.LogMsg("Shell started")
	reader := bufio.NewReader(os.Stdin)

	hostname, err := helpers.BuildPrefix(reader)
	if err != nil {
		helpers.LogMsg(fmt.Sprintf("Error building prefix: %v", err))
		fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
		return
	}
	helpers.LogMsg(fmt.Sprintf("Prompt built: %s", strings.TrimSpace(hostname)))

	for {
		fmt.Print(hostname + " ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
			continue
		}

		if err := ExecuteInput(input); err != nil {
			helpers.LogMsg(fmt.Sprintf("Error executing input '%s': %v", strings.TrimSpace(input), err))
			fmt.Fprintln(os.Stderr, "SHELL ➜ ", err)
		}
	}
}

