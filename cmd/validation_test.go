package cmd

import (
	"testing"

	"github.com/TheRealSibasishBehera/syncsh/internal/config"
)

func TestValidationIntegration(t *testing.T) {
	tests := []struct {
		name        string
		shellKind   string
		historyPath string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid zsh",
			shellKind:   "zsh",
			historyPath: "",
			expectError: false,
		},
		{
			name:        "invalid shell",
			shellKind:   "invalid",
			historyPath: "",
			expectError: true,
			errorMsg:    "unsupported shell type",
		},
		{
			name:        "empty shell",
			shellKind:   "",
			historyPath: "",
			expectError: true,
			errorMsg:    "shell type cannot be empty",
		},
		{
			name:        "shell with spaces",
			shellKind:   " bash ",
			historyPath: "",
			expectError: true,
			errorMsg:    "cannot contain spaces",
		},
		{
			name:        "relative history path",
			shellKind:   "fish",
			historyPath: "relative/path",
			expectError: true,
			errorMsg:    "must be absolute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateShellKindAndHistoryPath(tt.shellKind, tt.historyPath)
			
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

// containsString checks if string contains substring (from validation_test.go)
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