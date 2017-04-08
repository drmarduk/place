package main

import (
	"image/color"
	"log"
	"time"

	"encoding/json"

	"golang.org/x/net/websocket"
)

func main() {
	var image []Pixel
	ws, err := websocket.Dial("ws://pxls.space/ws", "", "http://pxls.space")
	if err != nil {
		panic(err)
	}

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

		image = append(image, p)
		log.Printf("Response: %s\n", string(msg[:n]))

	}
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
