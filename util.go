package main

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

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
