package store

import (
	"context"
	"time"
)

type Session struct {
	SessionID        string
	UserID           int64
	LastAccessedTime time.Time
}

func (s *Store) CreateSession(ctx context.Context, session *Session) error {
	return s.driver.CreateSession(ctx, session)
}

func (s *Store) UpdateLastAccessedTime(ctx context.Context, sessionId string, lastAccessTime time.Time) error {
	return s.driver.UpdateLastAccessedTime(ctx, sessionId, lastAccessTime)
}
