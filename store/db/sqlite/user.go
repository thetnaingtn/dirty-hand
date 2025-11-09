package sqlite

import (
	"context"
	"strings"

	"github.com/thetnaingtn/dirty-hand/store"
)

func (d *DB) CreateUser(ctx context.Context, user *store.User) error {
	fields := []string{"username", "password_hash", "role"}
	placeholders := []string{"?", "?", "?"}
	values := []any{user.Username, user.PasswordHash, user.Role}

	stmt := "INSERT INTO users (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ") RETURNING id"

	if err := d.db.QueryRowContext(ctx, stmt, values...).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (d *DB) ListUsers(ctx context.Context, filter *store.FindUser) ([]store.User, error) {
	var users []store.User

	stmt := "SELECT id, username, password_hash, role FROM users WHERE 1=1"
	var args []any

	if filter != nil {
		if filter.Role != nil {
			stmt += " AND role = ?"
			args = append(args, *filter.Role)
		}
	}

	rows, err := d.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user store.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (d *DB) GetUser(ctx context.Context, filter *store.FindUser) (*store.User, error) {
	var user store.User
	where, args := []string{"1 = 1"}, []any{}

	if filter != nil {
		if v := filter.Role; v != nil {
			where = append(where, "role = ?")
			args = append(args, *v)
		}

		if v := filter.Username; v != nil {
			where = append(where, "username = ?")
			args = append(args, *v)
		}

		if v := filter.ID; v != nil {
			where = append(where, "id = ?")
			args = append(args, *v)
		}
	}

	stmt := "SELECT id, username, password_hash, role FROM users WHERE " + strings.Join(where, " AND ")

	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
