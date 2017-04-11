package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	width  = 2000
	height = 2000
)

func main() {
	templateFile := flag.String("template", "", "the original template file to diff agains current pixelmap")
	outfile := flag.String("difffile", "diff.png", "the name of the diff png-image/textfile")
	text := flag.Bool("noimage", false, "false (default) for image output, true for text output")
	trans := flag.Bool("trans", false, "orig image is in background")
	flag.Parse()

	if *templateFile == "" {
		log.Fatalf("no template found")
	}

	template, err := readTemplate(*templateFile)
	if err != nil {
		log.Fatalf("error while reading template: %v\n", err)
	}

	// load current boarddata
	pixel, err := loadCurrentBoard()
	if err != nil {
		log.Fatalf("error while getting current board: %v\n", err)
	}
	board := convertToImage(width, height, pixel)

	diff := compareImages(template, board, *text, *trans)
	filename := sanitizeOutfile(*text, *outfile)

	if !*text { // only write if we want an image
		if err = saveImage(filename, diff); err != nil {
			log.Fatalf("error while saving diff image: %v\n", err)
		}
	}
}

func sanitizeOutfile(text bool, outfile string) string {
	outfile = strings.Replace(outfile, ".png", "", 1)
	outfile = strings.Replace(outfile, ".txt", "", 1)

	if text {
		return outfile + ".txt"
	}
	return outfile + ".png"
}

func saveImage(filename string, src image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, src)
}

func compareImages(template, current image.Image, text, trans bool) image.Image {
	var re image.Rectangle
	var result *image.RGBA

	if !text {
		re = image.Rect(0, 0, width, height)
		result = image.NewRGBA(re) // , PixelColor)
	}

	f := func(c color.Color) string {
		r, g, b, _ := c.RGBA()
		return fmt.Sprintf("(%d, %d, %d)", uint8(r), uint8(g), uint8(b))
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			srcColor := template.At(x, y)
			boardColor := current.At(x, y)
			r, g, b, a := srcColor.RGBA()
			rr, gg, bb, _ := boardColor.RGBA()

			if trans {
				n := color.RGBA{R: uint8(rr), G: uint8(gg), B: uint8(bb), A: 40}
				result.SetRGBA(x, y, n)
			}
			if a != 0 { // only look at pixel with alphachannel <> 0

				if r != rr || g != gg || b != bb {
					// this is ugly, the func should accept a ReadWriter and write to this
					// but for now it is okayish
					// TODO: add ReadWriter
					if text {
						fmt.Printf("%d x %d is %s should be %s\n", x, y, f(srcColor), f(boardColor))
					} else {
						result.Set(x, y, srcColor)
					}
				}
			}
		}
	}
	// maybe write here to the writer
	x := image.Rect(950, 1500, 2000, 2000)
	return result.SubImage(x)
}

func readTemplate(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func convertToImage(w, h int, src []Pixel) image.Image {
	r := image.Rect(0, 0, w, h)
	img := image.NewPaletted(r, PixelColor)

	for _, p := range src {
		img.Set(p.X, p.Y, PixelColor[p.Color])
	}
	return img
}

// loadCurrentBoard downloads the current bitmap from the website
// and converts them in a Pixel array
func loadCurrentBoard() ([]Pixel, error) {
	var result []Pixel

	_tmp, err := dlzip("http://pxls.space/boarddata")
	if err != nil {
		return nil, err
	}
	data := string(_tmp)

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p := Pixel{
				X:     x,
				Y:     y,
				Color: int(data[i]),
			}
			if p.Color > 15 {
				p.Color = 0
			}
			i++
			result = append(result, p)
		}
	}

	return result, nil
}

func dlzip(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("dlzip: could not create request, %v\n", err)
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		log.Printf("dlzip: could not peform request, %v\n", err)
		return nil, err
	}

	defer response.Body.Close()

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			// log.Printf("dlzip: could not create gzipreader, %v\n", err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	src, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("dlzip: could not read from response.body, %v\n", err)
		return nil, err
	}
	return src, nil
}
