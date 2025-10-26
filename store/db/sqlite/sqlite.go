package sqlite

import (
	"database/sql"
	"errors"

	"github.com/thetnaingtn/dirty-hand/internal/config"

	_ "modernc.org/sqlite"
)

var ErrNoDatabaseURL = errors.New("no database URL provided")

type DB struct {
	db     *sql.DB
	config *config.Config
}

func NewDB(cfg *config.Config) (*DB, error) {
	if cfg.Database.DSN == "" {
		return nil, ErrNoDatabaseURL
	}

	db, err := sql.Open("sqlite3", cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	schema := `CREATE TABLE IF NOT EXISTS products (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                description TEXT,
                price REAL NOT NULL,
                cover TEXT,
                created_at DATETIME NOT NULL,
                updated_at DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT NOT NULL UNIQUE,
                password_hash TEXT NOT NULL,
                role TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS sessions (
                user_id INTEGER NOT NULL,
                session_id CHAR(36) NOT NULL PRIMARY KEY,
                last_accessed_time DATETIME NOT NULL,
                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        );
		
		`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{
		db:     db,
		config: cfg,
	}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
