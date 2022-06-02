package main

import (
	"log"
	"net/url"
)

var addr = "127.0.0.1:2001"

func main() {
	rs1 := "http://127.0.0.1:2002"
	url1, err := url.Parse(rs1)
	if err != nil {
		log.Println(err)
	}
	urls := []*url.URL{url1}
	print(urls)
}
