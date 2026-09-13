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

func TestSessionQueries(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	store, err := Open()
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer store.Close()

	startedAt := time.Now().Truncate(time.Second)

	sessionID, err := store.CreateSession("query-test", startedAt)
	if err != nil {
		t.Fatalf("CreateSession() returned error: %v", err)
	}

	event := CommandEvent{
		SessionID: &sessionID,
		Command:   "git status",
		Cwd:       "/tmp/query-test",
		ExitCode:  0,
		StartedAt: startedAt.Add(2 * time.Second),
		EndedAt:   startedAt.Add(3 * time.Second),
	}

	if err := store.InsertCommandEvent(event); err != nil {
		t.Fatalf("InsertCommandEvent() returned error: %v", err)
	}

	session, err := store.SessionByName("query-test")
	if err != nil {
		t.Fatalf("SessionByName() returned error: %v", err)
	}

	if session.ID != sessionID {
		t.Errorf("session ID = %d, want %d", session.ID, sessionID)
	}

	if session.Name != "query-test" {
		t.Errorf("session name = %q, want %q", session.Name, "query-test")
	}

	events, err := store.CommandEventsForSession(sessionID)
	if err != nil {
		t.Fatalf("CommandEventsForSession() returned error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}

	if events[0].Command != "git status" {
		t.Errorf(
			"command = %q, want %q",
			events[0].Command,
			"git status",
		)
	}

	if events[0].ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", events[0].ExitCode)
	}
}
