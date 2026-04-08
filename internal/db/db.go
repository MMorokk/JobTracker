// Package db wraps the SQLite connection and exposes typed methods
// for all database operations used by JobTracker.
package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

// DB wraps a SQLite connection.
type DB struct {
	conn *sql.DB
}

// NewDB opens (or creates) a SQLite database at path and applies the schema.
func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	// recommended SQLite settings
	conn.Exec("PRAGMA journal_mode=WAL;")
	conn.Exec("PRAGMA foreign_keys=ON;")
	conn.Exec("PRAGMA busy_timeout=5000;")
	conn.SetMaxOpenConns(1)

	// create tables if they don't exist
	_, err = conn.Exec(`
        CREATE TABLE IF NOT EXISTS job_postings (
            id	INTEGER PRIMARY KEY AUTOINCREMENT,
            title	TEXT NOT NULL,
            company	TEXT NOT NULL,
            location	TEXT,
            type	TEXT,
            working_mode	TEXT,
            salary_min	REAL,
            salary_max	REAL,
            salary_type	TEXT,   -- 'monthly', 'hourly', 'annual', etc.
            salary_currency	TEXT,		-- ISO 4217 currency code
            description	TEXT,
            summary	TEXT,
            requirements	TEXT, -- json array of strings
            notes	TEXT,
            status TEXT, -- applied, interview, hired, rejected, offer, ghosted
            url	TEXT NOT NULL
        );
    `)

	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}
