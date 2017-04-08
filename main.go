package main

import (
	"compress/gzip"
	"errors"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"golang.org/x/image/bmp"

	"io"

	"golang.org/x/net/websocket"
)

var (
	width  = 1000
	height = 1000
)

func main() {

	pixelchan := make(chan Pixel)
	go PixelSetter(pixelchan)

	if err := loadInitialBoard(pixelchan, 5); err != nil {
		log.Printf("Could not load current board, start empty, %v\n", err)
	}
	SocketHandler(pixelchan)
}

func loadInitialBoard(c chan Pixel, maxtries int) error {

	log.Println("start initial load")
	var _tmp []byte
	var err error
	for i := 0; i < maxtries; i++ {
		_tmp, err = dlzip("http://pxls.space/boarddata")
		if err != nil {
			log.Printf("Try %d of %d failed, %v\n", i, maxtries, err)
			if i == maxtries-1 {
				// last try
				log.Println("we tried, return with empty pic")
				return errors.New("max attempts exceeded")
			}
			continue
		}
		break
	}
	data := string(_tmp)

	i := 0
	for y := 0; y < 1000; y++ {
		for x := 0; x < 1000; x++ {
			p := NewPixel()
			p.X = x
			p.Y = y
			p.Color = int(data[i])
			p.Type = "pixel"

			if p.Color > 15 {
				p.Color = 0
			}

			c <- p
			i++
		}
	}

	p := NewPixel()
	p.Type = "savedatshid"

	log.Println("initial load done.")
	return nil
}

// SocketHandler handles the server-client stuff
func SocketHandler(c chan Pixel) {
	ws, err := websocket.Dial("ws://pxls.space/ws", "", "http://pxls.space")
	if err != nil {
		panic(err)
	}
	log.Println("connected to websocket")
	defer ws.Close()

	for {

		msg := make([]byte, 512)
		n, err := ws.Read(msg)
		if err != nil {
			log.Printf("Error while reading websocket, %v\n", err)
			return
		}

		p := NewPixel()
		err = json.Unmarshal(msg[:n], &p)
		if err != nil {
			log.Printf("Error while parsing json string, %v\n", err)
			continue
		}

		c <- p
	}
}

// PixelSetter handels as goroutine
func PixelSetter(c chan Pixel) {
	r := image.Rect(0, 0, 999, 999)
	img := image.NewPaletted(r, PixelColor)

	for {
		p := <-c
		// receive pixel

		img.Set(p.X, p.Y, PixelColor[p.Color])
		// log.Printf("x: %d - y: %d - c: %d\n", p.X, p.Y, p.Color)

		if p.Type != "pixel" {
			hwd, _ := os.Create("test.bmp")
			bmp.Encode(hwd, img)
			hwd.Close()
		}
	}
}

func dl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("dl: error while connecting, %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Println(resp.Status)
	src, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("dl: error while reading body, %v\n", err)
		return nil, err
	}

	return src, nil
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
