package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initdb(dbfile string) (err error) {
	if db == nil {
		db, err = sql.Open("sqlite3", dbfile)
		return err
	}
	return nil
}

func closedb() error {
	err := db.Close()
	db = nil
	return err
}

func installTables(dbfile string) error {
	query := `CREATE TABLE IF NOT EXISTS pixel (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	x INTEGER NOT NULL,
	 y	INTEGER NOT NULL,
	color	INTEGER NOT NULL,
	created	DATETIME NOT NULL
);`

	queryCreateIndex := `CREATE INDEX IF NOT EXISTS created ON pixel (created ASC)`
	initdb(dbfile)
	defer closedb()

	_, err := db.Exec(query)
	if err != nil {
		log.Println("could not create table")
		return err
	}
	_, err = db.Exec(queryCreateIndex)
	if err != nil {
		log.Println("could note create index")
		return err
	}
	return nil
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
