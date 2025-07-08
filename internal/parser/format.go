package parser

type HistoryEntry struct {
	ID        int64  // Auto-increment primary key
	Timestamp int64  // Unix timestamp when command was executed
	MachineID string // Identifier for the machine that executed it
	Command   string // The actual shell command
	Duration  int    // Command execution duration in seconds
	ExitCode  int    // Command exit code (0 = success)
	Hash      string // SHA256 hash for deduplication
}

type ShellParser interface {
	ParseLine(line string) (command string, timestamp int64, skip bool)
	GetHistoryPath() []string
}

// so there would be a history file
// by default we will only look at a default shell
// optionally the user can set in the config the kind of shell they use
// shell = [zsh , bash] // valid for now is zsh , bash

// parsing techinque
// 1st get the shell kind (default or config)
// get the history file
// read them in lines (after time stamp x , x = last read time stored int the sqlite db)
// clean the white space
// dont read starting with whitespace (private commands)
// find new entries
//what is a new entry ?
