package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type CommandEvent struct {
	Command   string
	Cwd       string
	ExitCode  int
	StartedAt time.Time
	EndedAt   time.Time
}

type Store struct {
	db *sql.DB
}

func stateHome() (string, error) {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return filepath.Join(xdg, "waketrail"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".local", "state", "waketrail"), nil
}

func Open() (*Store, error) {
	dir, err := stateHome()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, "waketrail.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}

	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS command_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	command TEXT NOT NULL,
	cwd TEXT NOT NULL,
	exit_code INTEGER NOT NULL,
	started_at TEXT NOT NULL,
	ended_at TEXT NOT NULL
);
`

	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) InsertCommandEvent(event CommandEvent) error {
	const query = `
INSERT INTO command_events (
	command,
	cwd,
	exit_code,
	started_at,
	ended_at
)
VALUES (?, ?, ?, ?, ?);
`

	_, err := s.db.Exec(
		query,
		event.Command,
		event.Cwd,
		event.ExitCode,
		event.StartedAt.Format(time.RFC3339Nano),
		event.EndedAt.Format(time.RFC3339Nano),
	)

	return err
}
