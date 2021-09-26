package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	alive  bool
	target *url.URL

	proxy *httputil.ReverseProxy

	weight        uint8
	currentWeight uint8
}

func NewBackend(target *url.URL, weight uint8) *Backend {

	b := &Backend{
		alive:         true,
		target:        target,
		weight:        weight,
		currentWeight: weight,
	}

	b.proxy = httputil.NewSingleHostReverseProxy(b.target)

	b.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, _ error) {

		// mark backend as dead
		b.alive = false
		w.WriteHeader(502)
	}

	return b
}
