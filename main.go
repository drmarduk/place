package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	outfile := flag.String("outfile", "output.png", "name of the output png file")
	dbfile := flag.String("db", "", "the pixel database")
	start := flag.String("start", "", "the start date 31.12.2017 23:45:12")
	end := flag.String("end", time.Now().Format("02.01.2006 15:04:05"), "the time when to render the map ex: 31.12.2017 23:45:12")

	iterate := flag.Bool("iterate", false, "wether to output many images from start to end in step seconds")
	step := flag.Int("step", 10, "seconds to render together in one image, in seconds")
	flag.Parse()

	if *dbfile == "" {
		log.Fatalln("no database found")
	}

	_start := parseTime(*start, time.Time{})
	_end := parseTime(*end, time.Now())

	if !*iterate {
		// output a single image
		pixel, err := extractData(*dbfile, _start, _end)
		if err != nil {
			log.Fatalf("error while extracting data: %v\n", err)
		}

		log.Printf("Found %d pixels to render\n", len(pixel))
		renderMap(pixel, *outfile)
	} else {
		_outfile := strings.TrimSuffix(*outfile, ".png")
		i := 0
		_step := time.Duration(*step) * time.Second

		img = image.NewPaletted(image.Rect(0, 0, 2000, 2000), PixelColor)

		for current := _start; _end.Sub(current).Seconds() > 0; current = current.Add(_step) {
			a := current
			b := current.Add(_step)
			pixel, err := extractData(*dbfile, a, b)
			if err != nil {
				log.Printf("error while stepping: current: %s - %v\n", current.String(), err)
				continue
			}
			_tmp := fmt.Sprintf("%s-%05d.png", _outfile, i)
			log.Printf("Found %d between %s and %s - %s\n", len(pixel), a.Format("02.01.2006 15:04:05"), b.Format("02.01.2006 15:04:05"), _tmp)
			i++
			renderMap(pixel, _tmp)
		}
	}
}

var img *image.Paletted

func renderMap(data []Pixel, out string) {
	for _, p := range data {
		img.Set(p.X, p.Y, PixelColor[p.Color])
	}

	hwd, err := os.Create(out)
	if err != nil {
		log.Printf("error while creating output file: %v\n", err)
		return
	}
	// bmp.Encode(hwd, img)
	png.Encode(hwd, img)
	hwd.Close()
}

func parseTime(s string, _default time.Time) time.Time {
	t, err := time.Parse("02.01.2006 15:04:05", s)
	if err != nil {
		log.Printf("could not parse time %s: %s\n", s, err)
		return _default
	}
	return t
}
