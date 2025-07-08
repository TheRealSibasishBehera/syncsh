package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/TheRealSibasishBehera/syncsh/internal/parser"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *Store {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	
	store := New(db)
	
	// Initialize schema
	ctx := context.Background()
	if err := store.InitSchema(ctx); err != nil {
		t.Fatal(err)
	}
	
	return store
}

func TestSchemaParserAlignment(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()
	
	// Test that HistoryEntry struct maps correctly to database schema
	entry := parser.HistoryEntry{
		Timestamp: time.Now().Unix(),
		MachineID: "test-machine-1",
		Command:   "ls -la",
		Duration:  150,
		ExitCode:  0,
		Hash:      "", // Will be auto-generated
	}
	
	// Create entry
	err := store.CreateEntry(ctx, &entry)
	if err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}
	
	// Verify the hash was generated
	if entry.Hash == "" {
		t.Error("Hash should have been generated")
	}
	
	// Retrieve by hash
	retrieved, err := store.GetEntryByHash(ctx, entry.Hash)
	if err != nil {
		t.Fatalf("Failed to retrieve entry: %v", err)
	}
	
	// Verify all fields match
	if retrieved.Timestamp != entry.Timestamp {
		t.Errorf("Timestamp mismatch: got %d, want %d", retrieved.Timestamp, entry.Timestamp)
	}
	if retrieved.MachineID != entry.MachineID {
		t.Errorf("MachineID mismatch: got %s, want %s", retrieved.MachineID, entry.MachineID)
	}
	if retrieved.Command != entry.Command {
		t.Errorf("Command mismatch: got %s, want %s", retrieved.Command, entry.Command)
	}
	if retrieved.Duration != entry.Duration {
		t.Errorf("Duration mismatch: got %d, want %d", retrieved.Duration, entry.Duration)
	}
	if retrieved.ExitCode != entry.ExitCode {
		t.Errorf("ExitCode mismatch: got %d, want %d", retrieved.ExitCode, entry.ExitCode)
	}
	if retrieved.Hash != entry.Hash {
		t.Errorf("Hash mismatch: got %s, want %s", retrieved.Hash, entry.Hash)
	}
	if retrieved.ID == 0 {
		t.Error("ID should be auto-generated and non-zero")
	}
}

func TestCreateEntry(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()
	
	entry := parser.HistoryEntry{
		Timestamp: time.Now().Unix(),
		MachineID: "machine-1",
		Command:   "echo hello",
		Duration:  50,
		ExitCode:  0,
	}
	
	err := store.CreateEntry(ctx, &entry)
	if err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}
	
	// Test duplicate hash protection
	entry2 := entry // Same values should generate same hash
	err = store.CreateEntry(ctx, &entry2)
	if err != ErrDuplicateHash {
		t.Errorf("Expected ErrDuplicateHash, got: %v", err)
	}
}

func TestListEntries(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()
	
	// Create test entries
	entries := []parser.HistoryEntry{
		{
			Timestamp: time.Now().Unix() - 100,
			MachineID: "machine-1",
			Command:   "ls",
			Duration:  10,
			ExitCode:  0,
		},
		{
			Timestamp: time.Now().Unix() - 50,
			MachineID: "machine-2",
			Command:   "pwd",
			Duration:  5,
			ExitCode:  0,
		},
		{
			Timestamp: time.Now().Unix(),
			MachineID: "machine-1",
			Command:   "echo test",
			Duration:  20,
			ExitCode:  0,
		},
	}
	
	for i := range entries {
		err := store.CreateEntry(ctx, &entries[i])
		if err != nil {
			t.Fatalf("Failed to create entry: %v", err)
		}
	}
	
	// Test listing all entries
	allEntries, err := store.ListEntries(ctx, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to list entries: %v", err)
	}
	
	if len(allEntries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(allEntries))
	}
	
	// Test filtering by machine
	machine1Entries, err := store.ListEntries(ctx, "machine-1", 0, 0)
	if err != nil {
		t.Fatalf("Failed to list machine-1 entries: %v", err)
	}
	
	if len(machine1Entries) != 2 {
		t.Errorf("Expected 2 entries for machine-1, got %d", len(machine1Entries))
	}
	
	// Test filtering by timestamp
	since := time.Now().Unix() - 60
	recentEntries, err := store.ListEntries(ctx, "", since, 0)
	if err != nil {
		t.Fatalf("Failed to list recent entries: %v", err)
	}
	
	if len(recentEntries) != 2 {
		t.Errorf("Expected 2 recent entries, got %d", len(recentEntries))
	}
	
	// Test limit
	limitedEntries, err := store.ListEntries(ctx, "", 0, 1)
	if err != nil {
		t.Fatalf("Failed to list limited entries: %v", err)
	}
	
	if len(limitedEntries) != 1 {
		t.Errorf("Expected 1 limited entry, got %d", len(limitedEntries))
	}
}

func TestSyncState(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()
	
	machineID := "test-machine"
	
	// Test getting non-existent sync state
	timestamp, err := store.GetLastSyncTimestamp(ctx, machineID)
	if err != nil {
		t.Fatalf("Failed to get non-existent sync timestamp: %v", err)
	}
	if timestamp != 0 {
		t.Errorf("Expected 0 for non-existent sync timestamp, got %d", timestamp)
	}
	
	// Test updating sync state
	newTimestamp := time.Now().Unix()
	err = store.UpdateLastSyncTimestamp(ctx, machineID, newTimestamp)
	if err != nil {
		t.Fatalf("Failed to update sync timestamp: %v", err)
	}
	
	// Test retrieving updated sync state
	retrieved, err := store.GetLastSyncTimestamp(ctx, machineID)
	if err != nil {
		t.Fatalf("Failed to get updated sync timestamp: %v", err)
	}
	if retrieved != newTimestamp {
		t.Errorf("Expected timestamp %d, got %d", newTimestamp, retrieved)
	}
	
	// Test updating again (should replace, not duplicate)
	newerTimestamp := time.Now().Unix()
	err = store.UpdateLastSyncTimestamp(ctx, machineID, newerTimestamp)
	if err != nil {
		t.Fatalf("Failed to update sync timestamp again: %v", err)
	}
	
	retrieved, err = store.GetLastSyncTimestamp(ctx, machineID)
	if err != nil {
		t.Fatalf("Failed to get newest sync timestamp: %v", err)
	}
	if retrieved != newerTimestamp {
		t.Errorf("Expected newest timestamp %d, got %d", newerTimestamp, retrieved)
	}
}

func TestDeleteEntry(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()
	
	entry := parser.HistoryEntry{
		Timestamp: time.Now().Unix(),
		MachineID: "machine-1",
		Command:   "rm test",
		Duration:  10,
		ExitCode:  0,
	}
	
	err := store.CreateEntry(ctx, &entry)
	if err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}
	
	// Delete the entry
	err = store.DeleteEntry(ctx, entry.Hash)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}
	
	// Verify it's gone
	_, err = store.GetEntryByHash(ctx, entry.Hash)
	if err != ErrEntryNotFound {
		t.Errorf("Expected ErrEntryNotFound, got: %v", err)
	}
	
	// Test deleting non-existent entry
	err = store.DeleteEntry(ctx, "non-existent-hash")
	if err != ErrEntryNotFound {
		t.Errorf("Expected ErrEntryNotFound for non-existent entry, got: %v", err)
	}
}

func TestHashGeneration(t *testing.T) {
	// Test that same inputs generate same hash
	hash1 := generateHash(1234567890, "machine-1", "ls -la")
	hash2 := generateHash(1234567890, "machine-1", "ls -la")
	
	if hash1 != hash2 {
		t.Error("Same inputs should generate same hash")
	}
	
	// Test that different inputs generate different hashes
	hash3 := generateHash(1234567890, "machine-2", "ls -la")
	if hash1 == hash3 {
		t.Error("Different machine IDs should generate different hashes")
	}
	
	hash4 := generateHash(1234567890, "machine-1", "pwd")
	if hash1 == hash4 {
		t.Error("Different commands should generate different hashes")
	}
	
	hash5 := generateHash(1234567891, "machine-1", "ls -la")
	if hash1 == hash5 {
		t.Error("Different timestamps should generate different hashes")
	}
}