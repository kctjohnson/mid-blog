package db

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"

	"github.com/kctjohnson/mid-blog/assets"
)

type DB struct {
	*sqlx.DB
}

// New opens (creating if needed) the SQLite database file at path.
func New(path string, automigrate bool) (*DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	// Pragmas are applied to every pooled connection:
	//   foreign_keys - SQLite leaves FK enforcement (and ON DELETE CASCADE) off by default
	//   journal_mode - WAL lets readers proceed while a write is in progress
	//   busy_timeout - wait for locks instead of failing immediately on concurrent writes
	q := url.Values{}
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "journal_mode(WAL)")
	q.Add("_pragma", "busy_timeout(5000)")
	q.Add("_pragma", "synchronous(NORMAL)")
	dsn := "file:" + path + "?" + q.Encode()

	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", path, err)
	}

	if automigrate {
		driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
		if err != nil {
			return nil, err
		}

		source, err := iofs.New(assets.Migrations, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
		if err != nil {
			return nil, err
		}

		if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	return &DB{db}, nil
}
