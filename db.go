package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "modernc.org/sqlite" // pure Go, no gcc needed
)

type Database struct {
	db *sql.DB
}

func tableViewQuery(db *sql.DB) ([][]string, error) {
	rows, err := db.Query(`SELECT id, title, company, location, status FROM job_postings`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var result [][]string
	for rows.Next() {
		var id int
		var title, company string
		var location, status sql.NullString
		err := rows.Scan(&id, &title, &company, &location, &status)
		if err != nil {
			return nil, err
		}
		result = append(result, []string{
			strconv.Itoa(id),
			title,
			company,
			location.String,
			status.String,
		})
	}
	return result, rows.Err()
}

// NewDatabase initializes a new SQLite database connection with the specified path and ensures required tables exist.
func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// recommended SQLite settings
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA foreign_keys=ON;")
	db.Exec("PRAGMA busy_timeout=5000;")
	db.SetMaxOpenConns(1)

	// create tables if they don't exist
	_, err = db.Exec(`
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

	return &Database{db: db}, nil
}
