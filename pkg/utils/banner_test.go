package utils

import (
	"testing"
)

func TestShowSyncshBanner(t *testing.T) {
	ShowSyncshBanner()
	t.Log("ShowSyncshBanner executed successfully")
}

func TestShowInitBanner(t *testing.T) {
	ShowInitBanner()
	t.Log("ShowInitBanner executed successfully")
}

func TestSyncshASCIIContent(t *testing.T) {
	if syncshASCII == "" {
		t.Error("syncshASCII constant should not be empty")
	}

	if len(syncshASCII) < 10 {
		t.Error("syncshASCII should contain meaningful content")
	}

	t.Logf("ASCII art length: %d characters", len(syncshASCII))
}

