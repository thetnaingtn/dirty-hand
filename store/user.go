package store

import (
	"context"
	"strconv"
)

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleProductView Role = "product:view"
	RoleProductEdit Role = "product:edit"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         Role
}

type FindUser struct {
	Role     *Role
	ID       *int64
	Username *string
}

func (s *Store) CreateUser(ctx context.Context, user *User) (*User, error) {
	user, err := s.driver.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	key := strconv.FormatInt(user.ID, 10)
	s.userCache.Set(key, user)

	return user, nil
}

func (s *Store) ListUsers(ctx context.Context, filter *FindUser) ([]User, error) {
	return s.driver.ListUsers(ctx, filter)
}

func (s *Store) GetUser(ctx context.Context, filter *FindUser) (*User, error) {
	if filter.ID != nil {
		key := strconv.FormatInt(*filter.ID, 10)
		item, exist := s.userCache.Get(key)
		if exist {
			if user, ok := item.(*User); ok {
				return user, nil
			}
		}
	}

	return s.driver.GetUser(ctx, filter)
}

func (s *Store) GetUserSessions(ctx context.Context, userId int64) ([]Session, error) {
	return s.driver.GetUserSessions(ctx, userId)
}
