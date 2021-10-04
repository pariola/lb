package main

import (
	"fmt"
	"net/http"
	"os"

	"lb/pkg/config"
)

func main() {

	f, err := os.Open("lb.yml")

	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(f)

	if err != nil {
		panic(err)
	}

	p := NewPool(cfg)

	go p.HealthCheck()

	_ = http.ListenAndServe(fmt.Sprintf(":%d", p.cfg.Port), p)
}
