package main

import (
	"log"
	"time"

	"encoding/json"

	"golang.org/x/net/websocket"
)

var (
	width  = 4000
	height = 4000
)

func main() {
	savechan := make(chan Response)

	go PixelSaver(savechan)
	SocketHandler(savechan)
}

// SocketHandler handles the server-client stuff
func SocketHandler(savechan chan Response) {
	var ws *websocket.Conn
	var err error

	for {
		if ws == nil {
			ws, err = connect("ws://plxs.space/ws", "http://pxls.spcase")
			if err != nil {
				log.Printf("connection not established, %v\n", err)
				time.Sleep(time.Second * 5)
				continue
			}
			log.Println("connected to websocket.")
		}

		msg := make([]byte, 4096)
		n, err := ws.Read(msg)
		if err != nil {
			log.Printf("Error while reading websocket, %v (%s)\n", err, string(msg[:n]))
			ws = nil
			continue
		}

		r := NewResponse()
		err = json.Unmarshal(msg[:n], &r)
		if err != nil {
			log.Printf("Error while parsing json string, %v (%s)\n", err, string(msg[:n]))
			continue
		}
		if r.Type != "pixel" {
			continue
		}
		savechan <- r
	}
}

func connect(url, loc string) (*websocket.Conn, error) {
	ws, err := websocket.Dial("ws://pxls.space/ws", "", "http://pxls.space")
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// PixelSaver saves it tothe database
func PixelSaver(c chan Response) {
	var pixels []Pixel

	initdb()
	for {
		r := <-c
		for _, p := range r.Pixels {
			p.Created = time.Now()
			pixels = append(pixels, p)
		}

		if len(pixels) > 100 {
			// save to db
			massInsert(pixels)
			pixels = pixels[:0]
		}
	}
}
