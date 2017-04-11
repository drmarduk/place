package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	dbfile := flag.String("db", "pxls.space.db", "the name of the database to store all pixels")
	websocket := flag.String("ws", "ws://pxls.space/ws", "the uri of the websocket to connect, eg. ws://url.tld/path")
	flag.Parse()

	if err := installTables(*dbfile); err != nil {
		log.Fatalf("error while createing database: %v\n", err)
	} else {
		log.Print("db already there")
	}

	savechan := make(chan Response)

	go PixelSaver(*dbfile, savechan)
	SocketHandler(*websocket, savechan)
}

// SocketHandler handles the server-client stuff
func SocketHandler(wsURL string, savechan chan Response) {
	var ws *websocket.Conn
	var err error

	for {
		if ws == nil {
			ws, err = connect(wsURL, "http://pxls.spcase")
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
	ws, err := websocket.Dial(url, "", loc)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// PixelSaver saves it tothe database
func PixelSaver(dbfile string, c chan Response) {
	var pixels []Pixel

	initdb(dbfile)
	for {
		r := <-c
		for _, p := range r.Pixels {
			p.Created = time.Now()
			pixels = append(pixels, p)
		}

		if len(pixels) > 100 {
			// save to db, and reset cache
			massInsert(pixels)
			pixels = pixels[:0]
		}
	}
}
