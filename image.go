package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// PixelColor are the colors we use
var PixelColor color.Palette = []color.Color{
	color.RGBA{255, 255, 255, 255},
	color.RGBA{228, 228, 228, 255},
	color.RGBA{136, 136, 136, 255},
	color.RGBA{34, 34, 34, 255},
	color.RGBA{255, 167, 209, 255},
	color.RGBA{229, 0, 0, 255},
	color.RGBA{229, 149, 0, 255},
	color.RGBA{160, 106, 66, 255},
	color.RGBA{229, 217, 0, 255},
	color.RGBA{148, 224, 68, 255},
	color.RGBA{2, 190, 1, 255},
	color.RGBA{0, 211, 221, 255},
	color.RGBA{0, 131, 199, 255},
	color.RGBA{0, 0, 234, 255},
	color.RGBA{207, 110, 228, 255},
	color.RGBA{140, 0, 128, 255},
}

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

// NewPixel returns a pixel with a timestamp
func NewPixel() Pixel {
	return Pixel{Created: time.Now()}
}

// String overrides the formal string
func (p *Pixel) String() string {
	return fmt.Sprintf("X: %d, Y: %d, C: %d\n", p.X, p.Y, p.Color)
}

// Save saves the pixel in the database
func (p *Pixel) Save() error {
	initdb()

	_, err := db.Exec(
		"insert into pixel(id, x, y, color, created) values(NULL, ?, ?, ?, ?)",
		p.X, p.Y, p.Color, p.Created,
	)

	return err
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
		db, err = sql.Open("sqlite3", "pixel.db")
		return err
	}
	return nil

}

func closedb() error {
	return db.Close()
}

var db *sql.DB
