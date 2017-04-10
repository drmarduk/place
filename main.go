package main

import (
	"database/sql"
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/image/bmp"
)

func main() {
	dbfile := flag.String("db", "", "the pixel database")
	// dateStart := flag.String("start", "", "ex: 2017-04-10 01:23:45")
	end := flag.String("end", "", "the time when to render the map ex: 2017-04-10 01:23:45")
	flag.Parse()

	if *dbfile == "" {
		return
	}

	pixel, err := extractData(*dbfile, *end)
	if err != nil {
		log.Printf("Could not extract data from databas, %v\n", err)
		return
	}

	log.Printf("Found %d pixels to render\n", len(pixel))
	renderMap(pixel)
}

// Pixel represents the current Pixel in the image
type Pixel struct {
	Created time.Time
	Type    string  `json:"type"` // pixel, cooldown, users
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Color   int     `json:"color"`
	Wait    float64 `json:"wait"`
	Count   int     `json:"count"`
}

// PixelColor is a lookuptable for the used colors
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

func renderMap(data []Pixel) {
	r := image.Rect(0, 0, 2000, 2000)
	img := image.NewPaletted(r, PixelColor)

	for _, p := range data {
		img.Set(p.X, p.Y, PixelColor[p.Color])
	}

	hwd, _ := os.Create("output.bmp")
	bmp.Encode(hwd, img)
	hwd.Close()
}

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
