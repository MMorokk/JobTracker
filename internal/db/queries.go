package db

// GetJobs returns all job postings as rows of strings in column order:
// id, title, company, location, status.
// This matches the column layout expected by the table page.

import (
	"database/sql"
	"log"
	"strconv"
)

func (d *DB) GetJobs() ([][]string, error) {
	rows, err := d.conn.Query(`SELECT id, title, company, location, status FROM job_postings`)
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
