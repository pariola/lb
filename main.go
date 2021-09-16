package main

import (
	"net/http"
)

func main() {

	p := &ServerPool{
		backends: []*Backend{},
	}

	_ = p.Add("https://www.ask.com")
	_ = p.Add("https://www.bing.com")
	_ = p.Add("https://www.google.com")

	_ = http.ListenAndServe(":8080", p)
}
