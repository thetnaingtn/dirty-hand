package store

import (
	"context"
	"log/slog"
	"strconv"
	"time"
)

type Session struct {
	SessionID        string
	UserID           int64
	LastAccessedTime time.Time
}

func (s *Store) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	res, err := s.driver.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	sessions, err := s.GetUserSessions(ctx, session.UserID)
	if err != nil {
		slog.Error("can't get all user sessions")
	}

	key := strconv.FormatInt(session.UserID, 10)
	s.sessionCache.Set(key, sessions)

	return res, err
}

func (s *Store) UpdateLastAccessedTime(ctx context.Context, sessionId string, lastAccessTime time.Time) error {
	return s.driver.UpdateLastAccessedTime(ctx, sessionId, lastAccessTime)
}
