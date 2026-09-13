package storage

import (
	"testing"
	"time"
)

func TestInsertCommandEvent(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	store, err := Open()
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer store.Close()

	event := CommandEvent{
		Command:   "git status",
		Cwd:       "/tmp/project",
		ExitCode:  0,
		StartedAt: time.Now(),
		EndedAt:   time.Now(),
	}

	if err := store.InsertCommandEvent(event); err != nil {
		t.Fatalf("InsertCommandEvent() returned error: %v", err)
	}
}
