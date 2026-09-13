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

func TestCreateAndEndSession(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	store, err := Open()
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer store.Close()

	startedAt := time.Now()

	sessionID, err := store.CreateSession("homelab-debug", startedAt)
	if err != nil {
		t.Fatalf("CreateSession() returned error: %v", err)
	}

	if sessionID <= 0 {
		t.Fatalf("CreateSession() returned invalid ID: %d", sessionID)
	}

	endedAt := startedAt.Add(5 * time.Minute)

	if err := store.EndSession(sessionID, endedAt); err != nil {
		t.Fatalf("EndSession() returned error: %v", err)
	}

	var (
		name      string
		storedEnd string
	)

	err = store.db.QueryRow(
		`SELECT name, ended_at FROM sessions WHERE id = ?`,
		sessionID,
	).Scan(&name, &storedEnd)

	if err != nil {
		t.Fatalf("query session: %v", err)
	}

	if name != "homelab-debug" {
		t.Errorf("name = %q, want %q", name, "homelab-debug")
	}

	if storedEnd == "" {
		t.Error("ended_at is empty")
	}
}
