package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	dbfile := flag.String("db", "", "the pixel database")
	// dateStart := flag.String("start", "", "ex: 2017-04-10 01:23:45")
	end := flag.String("end", time.Now().Format("2006-01-02 15:05:05"), "the time when to render the map ex: 2017-04-10 01:23:45")
	flag.Parse()

	if *dbfile == "" {
		log.Fatalln("no database found")
	}

	pixel, err := extractData(*dbfile, *end)
	if err != nil {
		log.Fatalf("error while extracting data: %v\n", err)
	}

	log.Printf("Found %d pixels to render\n", len(pixel))
	renderMap(pixel)
}

func renderMap(data []Pixel) {
	r := image.Rect(0, 0, 2000, 2000)
	img := image.NewPaletted(r, PixelColor)

	for _, p := range data {
		img.Set(p.X, p.Y, PixelColor[p.Color])
	}

	hwd, err := os.Create("output.png")
	if err != nil {
		log.Printf("error while creating output file: %v\n", err)
		return
	}
	// bmp.Encode(hwd, img)
	png.Encode(hwd, img)
	hwd.Close()
}
