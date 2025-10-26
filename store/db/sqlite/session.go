package sqlite

import (
	"context"
	"time"

	"github.com/thetnaingtn/dirty-hand/store"
)

func (d *DB) CreateSession(ctx context.Context, session *store.Session) error {
	query := `INSERT INTO sessions (user_id, session_id, last_accessed_time) VALUES (?, ?, ?)`

	_, err := d.db.ExecContext(ctx, query, session.UserID, session.SessionID, session.LastAccessedTime)
	return err
}

func (d *DB) GetUserSessions(ctx context.Context, userId int64) ([]store.Session, error) {
	query := `SELECT session_id, user_id, last_accessed_time FROM sessions WHERE user_id = ?`

	rows, err := d.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []store.Session
	for rows.Next() {
		var session store.Session
		err := rows.Scan(&session.SessionID, &session.UserID, &session.LastAccessedTime)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (d *DB) UpdateLastAccessedTime(ctx context.Context, sessionId string, lastAccessTime time.Time) error {
	query := `UPDATE sessions SET last_accessed_time = ? WHERE session_id = ?`

	_, err := d.db.ExecContext(ctx, query, lastAccessTime, sessionId)
	return err
}
