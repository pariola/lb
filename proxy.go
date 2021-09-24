package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	Proxy *httputil.ReverseProxy

	weight        uint8
	currentWeight uint8
}

type ServerPool struct {
	current  int
	backends []*Backend

	// scheduling
	totalWeight uint8
}

// Add
func (p *ServerPool) Add(addr string, weight uint8) error {

	target, err := url.Parse(addr)

	if err != nil {
		return err
	}

	b := &Backend{
		URL:   target,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(target),

		weight:        weight,
		currentWeight: weight,
	}

	b.Proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, _ error) {

		// mark backend as dead
		b.Alive = false
		w.WriteHeader(502)
	}

	p.totalWeight += weight
	p.backends = append(p.backends, b)

	return nil
}

// NextBackend returns the next available backend based on Weighted Round Robin selection
// Reference: https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
func (p *ServerPool) NextBackend() *Backend {

	var big *Backend

	for _, backend := range p.backends {

		backend.currentWeight += backend.weight

		if big == nil || backend.currentWeight > big.currentWeight {
			big = backend
		}
	}

	if big != nil {
		big.currentWeight -= p.totalWeight
	}

	return big
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

// health
func (p *ServerPool) health() {

	for _, backend := range p.backends {

		conn, err := net.DialTimeout("tcp", backend.URL.Host, 2*time.Second)

		if err != nil {

			if backend.Alive {
				log.Printf("Backend [%s] no longer alive.\n", backend.URL)
			}

			backend.Alive = false
			continue
		}

		_ = conn.Close()

		if !backend.Alive {
			log.Printf("Backend [%s] now alive.\n", backend.URL)
		}

		backend.Alive = true
	}
}

// HealthCheck
func (p *ServerPool) HealthCheck() {

	t := time.NewTicker(30 * time.Second)

	for range t.C {
		log.Println("...starting health check...")
		p.health()
		log.Println("...health check done...")
	}
}
