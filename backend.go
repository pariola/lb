package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	m sync.RWMutex

	alive bool

	target *url.URL

	weight, currentWeight int32

	proxy *httputil.ReverseProxy
}

func NewBackend(target *url.URL, weight, maxRetries int32) *Backend {

	b := &Backend{
		alive:         true,
		target:        target,
		weight:        weight,
		currentWeight: weight,
	}

	b.proxy = httputil.NewSingleHostReverseProxy(b.target)

	b.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, _ error) {

		retries, _ := r.Context().Value("x-retries").(int32)

		if retries < maxRetries {
			time.Sleep(10 * time.Millisecond)

			// increment retries
			ctx := context.WithValue(r.Context(), "x-retries", retries+1)

			b.proxy.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// mark backend as dead
		b.SetAlive(false)
		w.WriteHeader(502)
	}

	return b
}

// SetAlive for this backend
func (b *Backend) SetAlive(v bool) {
	b.m.Lock()
	defer b.m.Unlock()
	b.alive = v
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() bool {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.alive
}
