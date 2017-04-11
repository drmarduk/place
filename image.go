package main

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Response comes from the websocket server as json
type Response struct {
	Type   string  `json:"type"` // pixel
	Pixels []Pixel `json:"pixels"`
}

// NewResponse creates a new struct
func NewResponse() Response {
	return Response{}
}

// Pixel is a json string from the websocket
type Pixel struct {
	Created time.Time
	Type    string  `json:"type"` // pixel, cooldown, users
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Color   int     `json:"color"`
	Wait    float64 `json:"wait"`
	Count   int     `json:"count"`
}

func massInsert(pixels []Pixel) error {
	//initdb()
	//defer closedb()

	query := "insert into pixel(id, x, y, color, created) values"
	vals := []interface{}{}

	for _, p := range pixels {
		query += "(NULL, ?, ?, ?, ?),"
		vals = append(vals, p.X, p.Y, p.Color, p.Created)
	}
	query = strings.TrimSuffix(query, ",")

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("massInsert: statement prepare failed, %v\n", err)
		return err
	}
	_, err = stmt.Exec(vals...)
	log.Printf("Saved %d pixels\n", len(pixels))
	return err
}

func initdb() (err error) {
	if db == nil {
		db, err = sql.Open("sqlite3", "reset.db")
		return err
	}
	return nil

}

func closedb() error {
	return db.Close()
}

var db *sql.DB
