package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func extractData(db, end string) ([]Pixel, error) {
	query := "select x, y, color from pixel where created < ? order by created asc;"
	var result []Pixel

	c, err := sql.Open("sqlite3", db)
	if err != nil {
		return result, err
	}
	defer c.Close()

	rows, err := c.Query(query, end)
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
