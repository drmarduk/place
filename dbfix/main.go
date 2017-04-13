package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbfile := flag.String("db", "", "the database to scan")
	flag.Parse()

	times, err := extractDates(*dbfile)
	if err != nil {
		log.Fatalf("error while extracting times: %v\n", err)
	}

	filtered := filterTimes(times)

	for k, v := range filtered {
		fmt.Printf("%s\t%d\n", k, v)
	}

}

func filterTimes(times []time.Time) (result map[string]int) {
	result = make(map[string]int)
	for _, t := range times {
		result[t.Format("02.01.2006 15:04:05")]++
	}
	return result
}
func extractDates(file string) ([]time.Time, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var times []time.Time
	query := "select distinct created from pixel order by created asc;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t time.Time

		err := rows.Scan(&t)
		if err != nil {
			log.Println(err)
			continue
		}
		times = append(times, t)
	}
	return times, nil
}
