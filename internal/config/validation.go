package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateShellKind validates the shell type string
func ValidateShellKind(shellKind string) error {
	if shellKind == "" {
		return fmt.Errorf("shell type cannot be empty")
	}
	
	// Check for spaces or other whitespace
	if strings.TrimSpace(shellKind) != shellKind {
		return fmt.Errorf("shell type cannot contain spaces or whitespace")
	}
	
	// Validate against supported shells
	shell := ShellKind(shellKind)
	switch shell {
	case ShellBash, ShellZsh, ShellFish:
		return nil
	default:
		return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh, fish)", shellKind)
	}
}

// ValidateHistoryPath validates the history file path for a given shell
func ValidateHistoryPath(historyPath string, shellKind ShellKind) error {
	// Empty path is valid - will use shell default
	if historyPath == "" {
		return nil
	}
	
	// Path must be absolute
	if !filepath.IsAbs(historyPath) {
		return fmt.Errorf("history path must be absolute, got: %s", historyPath)
	}
	
	// Check if path exists
	info, err := os.Stat(historyPath)
	if err == nil {
		// Path exists - check if it's a file
		if info.IsDir() {
			return fmt.Errorf("history path is a directory, expected file: %s", historyPath)
		}
		// Existing file is valid
		return nil
	}
	
	if !os.IsNotExist(err) {
		// Some other error (permission, etc.)
		return fmt.Errorf("cannot access history path %s: %w", historyPath, err)
	}
	
	// File doesn't exist - check if we can create it
	parentDir := filepath.Dir(historyPath)
	
	// Check if parent directory exists
	parentInfo, err := os.Stat(parentDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("parent directory does not exist: %s", parentDir)
	}
	if err != nil {
		return fmt.Errorf("cannot access parent directory %s: %w", parentDir, err)
	}
	
	// Check if parent is actually a directory
	if !parentInfo.IsDir() {
		return fmt.Errorf("parent path is not a directory: %s", parentDir)
	}
	
	// Check if we can write to the parent directory
	if err := checkWritePermission(parentDir); err != nil {
		return fmt.Errorf("cannot write to directory %s: %w", parentDir, err)
	}
	
	// Try to create the file to validate the path
	if err := createHistoryFileIfNotExists(historyPath); err != nil {
		return fmt.Errorf("cannot create history file %s: %w", historyPath, err)
	}
	
	return nil
}

// ValidateShellKindAndHistoryPath validates both shell kind and history path together
func ValidateShellKindAndHistoryPath(shellKind, historyPath string) error {
	// First validate shell kind
	if err := ValidateShellKind(shellKind); err != nil {
		return fmt.Errorf("invalid shell type: %w", err)
	}
	
	// Convert to ShellKind for history path validation
	shell := ShellKind(shellKind)
	
	// Then validate history path
	if err := ValidateHistoryPath(historyPath, shell); err != nil {
		return fmt.Errorf("invalid history path: %w", err)
	}
	
	return nil
}

// checkWritePermission checks if we can write to a directory
func checkWritePermission(dir string) error {
	// Try to create a temporary file
	tempFile := filepath.Join(dir, ".syncsh_test_write")
	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	file.Close()
	
	// Clean up
	os.Remove(tempFile)
	return nil
}

// createHistoryFileIfNotExists creates an empty history file if it doesn't exist
func createHistoryFileIfNotExists(path string) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil // File exists, nothing to do
	}
	
	// Create empty file with appropriate permissions
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return nil
}

// GetSupportedShells returns a list of supported shell types
func GetSupportedShells() []string {
	return []string{
		string(ShellBash),
		string(ShellZsh),
		string(ShellFish),
	}
}

// IsValidShellKind checks if a shell kind is valid without returning an error
func IsValidShellKind(shellKind string) bool {
	return ValidateShellKind(shellKind) == nil
}