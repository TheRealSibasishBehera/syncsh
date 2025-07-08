CREATE TABLE history_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    machine_id TEXT NOT NULL,
    command TEXT NOT NULL,
    duration INTEGER DEFAULT 0,
    exit_code INTEGER DEFAULT 0,
    hash TEXT NOT NULL UNIQUE
);

-- Sync state tracking
CREATE TABLE sync_state (
    machine_id TEXT PRIMARY KEY,
    last_sync_timestamp INTEGER NOT NULL
);

-- Indexes for performance
CREATE INDEX idx_history_timestamp ON history_entries(timestamp);
CREATE INDEX idx_history_machine ON history_entries(machine_id);
CREATE INDEX idx_history_hash ON history_entries(hash);
CREATE INDEX idx_history_command ON history_entries(command);
