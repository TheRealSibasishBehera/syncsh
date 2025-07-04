package watcher

import (
	"github.com/fsnotify/fsnotify"
)

// Watcher is a struct that holds the fsnotify watcher instance
type Watcher struct {
	Watcher *fsnotify.Watcher
	Paths   []string
}

// NewWatcher creates a new watcher instance
func NewWatcher(paths []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			return nil, err
		}
	}

	return &Watcher{
		Watcher: watcher,
		Paths:   paths,
	}, nil
}
