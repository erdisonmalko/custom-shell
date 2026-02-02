package helpers

import (
 "bufio"
 "fmt"
 "os"
 "os/exec"
 "strings"
 "path/filepath"
)

func BuildPrefix(reader *bufio.Reader) (string, error) {
	LogMsg("Building prompt prefix")
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		LogMsg(fmt.Sprintf("Error running whoami: %v", err))
		return "", err
	}
	LogMsg(fmt.Sprintf("Got username: %s", strings.TrimSpace(string(output))))

	cmd1 := exec.Command("hostname")
	output1, err1 := cmd1.Output()
	if err1 != nil {
		LogMsg(fmt.Sprintf("Error running hostname: %v", err1))
		return "", err1
	}
	LogMsg(fmt.Sprintf("Got hostname: %s", strings.TrimSpace(string(output1))))

    hostname := strings.TrimSuffix(string(output), "\n") + "@" + strings.TrimSuffix(string(output1), "\n")
    
	cmd2 := exec.Command("pwd")
	output, err2 := cmd2.Output()
	if err2 != nil {
		LogMsg(fmt.Sprintf("Error running pwd: %v", err2))
		return "", err2
	}
	pwd_string := strings.TrimSuffix(string(output), "\n")
	pwd_string = strings.ReplaceAll(pwd_string, os.Getenv("HOME"), "~")
	pwd_string += "$"

	hostname = fmt.Sprintf("➜ %s: %s ", hostname, pwd_string)
	LogMsg(fmt.Sprintf("Prompt prefix built: %s", hostname))
	
	return hostname, nil
}

// CompletePath returns matching entries for a given prefix.
func CompletePath(prefix string) ([]string, error) {
    // expand ~ to HOME (simple handling)
    if strings.HasPrefix(prefix, "~") {
        prefix = strings.Replace(prefix, "~", os.Getenv("HOME"), 1)
    }

    var dir, base string
    if prefix == "" || prefix == "." {
        dir = "."
        base = ""
    } else if strings.HasSuffix(prefix, string(os.PathSeparator)) {
        dir = prefix
        base = ""
    } else if strings.Contains(prefix, string(os.PathSeparator)) {
        dir = filepath.Dir(prefix)
        base = filepath.Base(prefix)
    } else {
        dir = "."
        base = prefix
    }

    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    var matches []string
    for _, e := range entries {
        name := e.Name()
        if strings.HasPrefix(name, base) {
            candidate := name
            if dir != "." {
                candidate = filepath.Join(dir, name)
            }
            if e.IsDir() {
                candidate += string(os.PathSeparator)
            }
            matches = append(matches, candidate)
        }
    }
    return matches, nil
}