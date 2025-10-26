package store

import "context"

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
	Role *Role
}

func (s *Store) CreateUser(ctx context.Context, user *User) error {
	return s.driver.CreateUser(ctx, user)
}

func (s *Store) ListUsers(ctx context.Context, filter *FindUser) ([]User, error) {
	return s.driver.ListUsers(ctx, filter)
}

func (s *Store) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.driver.GetUser(ctx, id)
}

func (s *Store) GetUserSessions(ctx context.Context, userId int64) ([]Session, error) {
	return s.driver.GetUserSessions(ctx, userId)
}
