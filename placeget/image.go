package main

import (
	"time"
)

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
