package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	Proxy *httputil.ReverseProxy
}

type ServerPool struct {
	current  int
	backends []*Backend
}

// Add
func (p *ServerPool) Add(addr string) error {

	target, err := url.Parse(addr)

	if err != nil {
		return err
	}

	b := &Backend{
		URL:   target,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(target),
	}

	b.Proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, _ error) {

		// mark backend as dead
		b.Alive = false
		w.WriteHeader(502)
	}

	p.backends = append(p.backends, b)

	return nil
}

func (p *ServerPool) NextBackend() *Backend {

	next := (p.current + 1) % len(p.backends)

	// move full-cycle
	l := len(p.backends) + next
	for i := next; i < l; i++ {

		// normalize with moduli
		index := i % len(p.backends)

		backend := p.backends[index]

		if backend.Alive {
			p.current += 1
			return backend
		}
	}

	return nil
}

// ServeHTTP
func (p *ServerPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	backend := p.NextBackend()

	if backend == nil {
		w.WriteHeader(500)
		return
	}

	log.Printf("Backend [%s] | Path: %s\n", backend.URL, r.URL)

	backend.Proxy.ServeHTTP(w, r)
}
