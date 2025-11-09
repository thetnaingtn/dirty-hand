package store

import (
	"context"
	"time"
)

// Driver defines the data layer operations available to the store.
//
// Business logic should live in the store package and call these
// methods for persistence. Implementations of Driver should not
// contain application specific logic.
type Driver interface {
	Close() error

	CreateProduct(ctx context.Context, p *Product) (*Product, error)
	UpdateProduct(ctx context.Context, p *Product) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	CreateUser(ctx context.Context, user *User) error
	ListUsers(ctx context.Context, filter *FindUser) ([]User, error)
	GetUser(ctx context.Context, filter *FindUser) (*User, error)
	CreateSession(ctx context.Context, session *Session) error
	GetUserSessions(ctx context.Context, id int64) ([]Session, error)

	UpdateLastAccessedTime(ctx context.Context, sessionId string, lastAccessTime time.Time) error
}
