// utils.go
package utils

import (
	"os/exec"
	"strings"
)

// CheckDependency checks if a given command-line tool is available
func CheckDependency(toolName string) bool {
	var cmd *exec.Cmd
	if toolName == "kubectl" {
		cmd = exec.Command(toolName, "version", "--client")
	} else if toolName == "ksniff" {
		cmd = exec.Command("kubectl", "plugin", "list")
	} else {
		// For other tools, you can add more cases
	}

	if cmd == nil {
		return false
	}

	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// InstallInstructions returns a message with installation instructions
func InstallInstructions(toolName string) string {
	var instructions strings.Builder
	switch toolName {
	case "kubectl":
		instructions.WriteString("Download and install kubectl from: https://kubernetes.io/docs/tasks/tools/")
	case "ksniff":
		instructions.WriteString("Ensure kubectl is installed. Install ksniff via krew: kubectl krew install sniff")
		// Add more installation instructions for other tools if necessary
	}
	return instructions.String()
}
