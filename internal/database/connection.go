package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(home, ".passvault")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "vault.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) InitSchema() error {
	query := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS secrets (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at DATETIME,
		updated_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS fields (
		id TEXT,
		secret_id TEXT,
		key TEXT,
		value BLOB,
		is_sensitive BOOLEAN,
		FOREIGN KEY(secret_id) REFERENCES secrets(id) ON DELETE CASCADE
	);
	`
	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}
