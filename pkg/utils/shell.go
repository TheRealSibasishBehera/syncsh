package utils

import (
	"fmt"
	"os"
)

func GetShellKind() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable is not set")
	}

	switch {
	case shell == "/bin/bash":
		return "bash", nil
	case shell == "/bin/zsh":
		return "zsh", nil
	case shell == "/usr/bin/fish":
		return "fish", nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

func GetDefaultHistoryPath(shell string) (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable is not set")
	}

	switch shell {
	case "bash":
		return home + "/.bash_history", nil
	case "zsh":
		if histFile := os.Getenv("HISTFILE"); histFile != "" {
			return histFile, nil
		}
		return home + "/.zsh_history", nil
	case "fish":
		return home + "/.local/share/fish/fish_history", nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}
