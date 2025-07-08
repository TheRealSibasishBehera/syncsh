package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateShellKind(t *testing.T) {
	tests := []struct {
		name        string
		shellKind   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid bash",
			shellKind:   "bash",
			expectError: false,
		},
		{
			name:        "valid zsh",
			shellKind:   "zsh",
			expectError: false,
		},
		{
			name:        "valid fish",
			shellKind:   "fish",
			expectError: false,
		},
		{
			name:        "invalid shell",
			shellKind:   "invalid",
			expectError: true,
			errorMsg:    "unsupported shell type",
		},
		{
			name:        "empty shell",
			shellKind:   "",
			expectError: true,
			errorMsg:    "shell type cannot be empty",
		},
		{
			name:        "case sensitive bash",
			shellKind:   "BASH",
			expectError: true,
			errorMsg:    "unsupported shell type",
		},
		{
			name:        "shell with spaces",
			shellKind:   " zsh ",
			expectError: true,
			errorMsg:    "shell type cannot contain spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShellKind(tt.shellKind)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateHistoryPath(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create test files
	existingFile := filepath.Join(tmpDir, "existing_history")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Clean up permissions
	
	tests := []struct {
		name         string
		historyPath  string
		shellKind    ShellKind
		expectError  bool
		errorMsg     string
		shouldCreate bool
	}{
		{
			name:        "empty path with valid shell",
			historyPath: "",
			shellKind:   ShellZsh,
			expectError: false,
		},
		{
			name:        "existing file",
			historyPath: existingFile,
			shellKind:   ShellBash,
			expectError: false,
		},
		{
			name:         "non-existing file in existing directory",
			historyPath:  filepath.Join(tmpDir, "new_history"),
			shellKind:    ShellFish,
			expectError:  false,
			shouldCreate: true,
		},
		{
			name:        "invalid path - directory instead of file",
			historyPath: tmpDir,
			shellKind:   ShellBash,
			expectError: true,
			errorMsg:    "path is a directory",
		},
		{
			name:        "non-existing directory",
			historyPath: "/non/existing/directory/history",
			shellKind:   ShellZsh,
			expectError: true,
			errorMsg:    "parent directory does not exist",
		},
		{
			name:        "permission denied directory",
			historyPath: filepath.Join(readOnlyDir, "history"),
			shellKind:   ShellFish,
			expectError: true,
			errorMsg:    "cannot write to directory",
		},
		{
			name:        "absolute path validation",
			historyPath: "relative/path/history",
			shellKind:   ShellBash,
			expectError: true,
			errorMsg:    "history path must be absolute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHistoryPath(tt.historyPath, tt.shellKind)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				
				// If we expect the file to be created, check it exists
				if tt.shouldCreate && tt.historyPath != "" {
					if _, err := os.Stat(tt.historyPath); os.IsNotExist(err) {
						t.Errorf("expected file to be created at %s", tt.historyPath)
					}
				}
			}
		})
	}
}

func TestValidateShellKindAndHistoryPath(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name        string
		shellKind   string
		historyPath string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid shell with empty history path",
			shellKind:   "zsh",
			historyPath: "",
			expectError: false,
		},
		{
			name:        "valid shell with valid history path",
			shellKind:   "bash",
			historyPath: filepath.Join(tmpDir, "bash_history"),
			expectError: false,
		},
		{
			name:        "invalid shell with valid history path",
			shellKind:   "invalid",
			historyPath: filepath.Join(tmpDir, "history"),
			expectError: true,
			errorMsg:    "shell type",
		},
		{
			name:        "valid shell with invalid history path",
			shellKind:   "fish",
			historyPath: "/non/existing/path",
			expectError: true,
			errorMsg:    "history path",
		},
		{
			name:        "both invalid",
			shellKind:   "invalid",
			historyPath: "/non/existing/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShellKindAndHistoryPath(tt.shellKind, tt.historyPath)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestShellKindString(t *testing.T) {
	tests := []struct {
		shell    ShellKind
		expected string
	}{
		{ShellBash, "bash"},
		{ShellZsh, "zsh"},
		{ShellFish, "fish"},
	}

	for _, tt := range tests {
		t.Run(string(tt.shell), func(t *testing.T) {
			if string(tt.shell) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.shell))
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      hasSubstring(s, substr))))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}