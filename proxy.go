package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"lb/pkg/config"
)

type ServerPool struct {

	// scheduling
	weight int32

	requests uint32

	cfg *config.Config

	backends []*Backend
}

// NewPool creates a new ServerPool with the specified configuration
func NewPool(cfg *config.Config) *ServerPool {

	if cfg == nil {
		return nil
	}

	p := &ServerPool{
		cfg: cfg,
	}

	// TODO

	return p
}

// Add creates a new Backend instance
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

// ServeHTTP handles incoming requests and forwards to a backend
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

// healthCheck
func (p *ServerPool) healthCheck() {

	var dead, alive int32

	for _, backend := range p.backends {

		conn, err := net.DialTimeout("tcp", backend.target.Host, 2*time.Second)

		if err != nil {

			// backend was alive
			if backend.IsAlive() {
				log.Printf("health: Backend [%s] no longer alive.\n", backend.target)
			}

			dead++
			backend.SetAlive(false)
			continue
		}

		_ = conn.Close()

		// backend was dead
		if !backend.IsAlive() {
			log.Printf("health: Backend [%s] now alive.\n", backend.target)
		}

		alive++
		backend.SetAlive(true)
	}

	log.Printf("health: total %d | alive %d | dead %d\n", alive+dead, alive, dead)
}

// healthChecker triggers the health check function on a specified interval
func (p *ServerPool) healthChecker() {

	t := time.NewTicker(p.cfg.Health.Interval)

	for range t.C {
		p.healthCheck()
	}
}

// Start
func (p *ServerPool) Start() {

	go p.healthChecker()

	if http.ListenAndServe(":"+p.cfg.Port, p) != nil {
		log.Fatal("failed to start proxy!")
	}
}
