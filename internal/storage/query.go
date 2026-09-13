package storage

import (
	"database/sql"
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

func (s *Store) LatestSession() (SessionRecord, error) {
	const query = `
SELECT
	id,
	name,
	started_at,
	ended_at
FROM sessions
ORDER BY id DESC
LIMIT 1;
`

	return s.scanSession(s.db.QueryRow(query))
}

func (s *Store) SessionByName(name string) (SessionRecord, error) {
	const query = `
SELECT
	id,
	name,
	started_at,
	ended_at
FROM sessions
WHERE name = ?
ORDER BY id DESC
LIMIT 1;
`

	return s.scanSession(s.db.QueryRow(query, name))
}

func (s *Store) scanSession(row *sql.Row) (SessionRecord, error) {
	var (
		session       SessionRecord
		startedAtText string
		endedAtText   sql.NullString
	)

	err := row.Scan(
		&session.ID,
		&session.Name,
		&startedAtText,
		&endedAtText,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return SessionRecord{}, ErrSessionNotFound
	}

	if err != nil {
		return SessionRecord{}, err
	}

	startedAt, err := time.Parse(time.RFC3339Nano, startedAtText)
	if err != nil {
		return SessionRecord{}, err
	}

	session.StartedAt = startedAt

	if endedAtText.Valid {
		endedAt, err := time.Parse(time.RFC3339Nano, endedAtText.String)
		if err != nil {
			return SessionRecord{}, err
		}

		session.EndedAt = &endedAt
	}

	return session, nil
}

func (s *Store) CommandEventsForSession(sessionID int64) ([]CommandEvent, error) {
	const query = `
SELECT
	command,
	cwd,
	exit_code,
	started_at,
	ended_at
FROM command_events
WHERE session_id = ?
ORDER BY started_at ASC;
`

	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []CommandEvent

	for rows.Next() {
		var (
			event         CommandEvent
			startedAtText string
			endedAtText   string
		)

		err := rows.Scan(
			&event.Command,
			&event.Cwd,
			&event.ExitCode,
			&startedAtText,
			&endedAtText,
		)
		if err != nil {
			return nil, err
		}

		startedAt, err := time.Parse(time.RFC3339Nano, startedAtText)
		if err != nil {
			return nil, err
		}

		endedAt, err := time.Parse(time.RFC3339Nano, endedAtText)
		if err != nil {
			return nil, err
		}

		event.SessionID = &sessionID
		event.StartedAt = startedAt
		event.EndedAt = endedAt

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
