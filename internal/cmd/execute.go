package cmd

import (
    "os/exec"
)

// ExecuteCommand executes a shell command and returns the output as a string
func ExecuteCommand(command string, args ...string) (string, error) {
    cmd := exec.Command(command, args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}