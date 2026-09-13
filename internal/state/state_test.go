package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestState(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", tempDir)

	return tempDir
}

func TestSaveAndLoadSession(t *testing.T) {
	setupTestState(t)

	startedAt := time.Now().Truncate(time.Second)

	expected := Session{
		Name:      "test-session",
		StartedAt: startedAt,
	}

	if err := SaveSession(expected); err != nil {
		t.Fatalf("SaveSession() returned error: %v", err)
	}

	actual, err := LoadSession()
	if err != nil {
		t.Fatalf("LoadSession() returned error: %v", err)
	}

	if actual.Name != expected.Name {
		t.Errorf("Name = %q, want %q", actual.Name, expected.Name)
	}

	if !actual.StartedAt.Equal(expected.StartedAt) {
		t.Errorf(
			"StartedAt = %v, want %v",
			actual.StartedAt,
			expected.StartedAt,
		)
	}
}

func TestHasActiveSession(t *testing.T) {
	setupTestState(t)

	active, err := HasActiveSession()
	if err != nil {
		t.Fatalf("HasActiveSession() returned error: %v", err)
	}

	if active {
		t.Fatal("HasActiveSession() = true, want false")
	}

	session := Session{
		Name:      "active-test",
		StartedAt: time.Now(),
	}

	if err := SaveSession(session); err != nil {
		t.Fatalf("SaveSession() returned error: %v", err)
	}

	active, err = HasActiveSession()
	if err != nil {
		t.Fatalf("HasActiveSession() returned error: %v", err)
	}

	if !active {
		t.Fatal("HasActiveSession() = false, want true")
	}
}

func TestClearSession(t *testing.T) {
	tempDir := setupTestState(t)

	session := Session{
		Name:      "clear-test",
		StartedAt: time.Now(),
	}

	if err := SaveSession(session); err != nil {
		t.Fatalf("SaveSession() returned error: %v", err)
	}

	if err := ClearSession(); err != nil {
		t.Fatalf("ClearSession() returned error: %v", err)
	}

	active, err := HasActiveSession()
	if err != nil {
		t.Fatalf("HasActiveSession() returned error: %v", err)
	}

	if active {
		t.Fatal("session still active after ClearSession()")
	}

	path := filepath.Join(tempDir, "waketrail", "active-session.json")

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("session file still exists at %s", path)
	}
}

func TestClearSessionWhenNoneExists(t *testing.T) {
	setupTestState(t)

	if err := ClearSession(); err != nil {
		t.Fatalf("ClearSession() returned error: %v", err)
	}
}
