package main

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {

	go main()

	time.Sleep(2 * time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = http.Get("http://localhost:8080")
		}()
	}

	wg.Wait()
}
