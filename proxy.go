package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type ServerPool struct {

	// scheduling
	weight int32

	requests uint32

	backends []*Backend
}

// Add
func (p *ServerPool) Add(addr string, weight int32) error {

	target, err := url.Parse(addr)

	if err != nil {
		return err
	}

	b := NewBackend(target, weight)

	p.weight += b.weight
	p.backends = append(p.backends, b)

	return nil
}

// NextBackend returns the next available backend based on Weighted Round Robin selection
// Reference: https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
func (p *ServerPool) NextBackend() *Backend {

	var big *Backend

	for _, backend := range p.backends {

		if !backend.IsAlive() {
			continue
		}

		atomic.AddInt32(&backend.currentWeight, backend.weight)

		if big == nil ||
			atomic.LoadInt32(&backend.currentWeight) > atomic.LoadInt32(&big.currentWeight) {
			big = backend
		}
	}

	if big != nil {
		atomic.AddInt32(&big.currentWeight, -p.weight)
	}

	return big
}

// ServeHTTP
func (p *ServerPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	atomic.AddUint32(&p.requests, 1)

	backend := p.NextBackend()

	if backend == nil {
		w.WriteHeader(500)
		return
	}

	log.Printf("Backend [%s] | Path: %s\n", backend.target, r.URL)

	backend.proxy.ServeHTTP(w, r)
}

// health
func (p *ServerPool) health() {

	for _, backend := range p.backends {

		conn, err := net.DialTimeout("tcp", backend.target.Host, 2*time.Second)

		if err != nil {

			// backend was alive
			if backend.IsAlive() {
				log.Printf("Backend [%s] no longer alive.\n", backend.target)
			}

			backend.SetAlive(false)
			continue
		}

		_ = conn.Close()

		// backend was dead
		if !backend.IsAlive() {
			log.Printf("Backend [%s] now alive.\n", backend.target)
		}

		backend.SetAlive(true)
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
