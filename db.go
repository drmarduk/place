package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func extractData(db string, start, end time.Time) ([]Pixel, error) {
	query := "select x, y, color from pixel where "
	if !start.IsZero() {
		query += " created < ? and created > ? "
	} else {
		query += " created < ? "
	}
	query += "order by created asc;"

	var result []Pixel

	c, err := sql.Open("sqlite3", db)
	if err != nil {
		return result, err
	}
	defer c.Close()

	rows, err := c.Query(query, end, start) // we add start even if we dont need it
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var x Pixel
		err = rows.Scan(&x.X, &x.Y, &x.Color)
		if err != nil {
			log.Printf("scan error: %v\n", err)
			continue
		}

		result = append(result, x)
	}
	return result, nil
}
