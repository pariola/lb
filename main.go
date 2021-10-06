package main

import (
	"log"
	"os"
	"path/filepath"

	"lb/pkg/config"
)

var configFile = "lb.yml"

func main() {

	args := os.Args[1:]

	if args[0] != "" {
		configFile = args[0]
	}

	fPath, err := filepath.Abs(configFile)

	if err != nil {
		log.Fatal("no valid configuration file supplied")
	}

	f, err := os.Open(fPath)

	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(f)

	if err != nil {
		panic(err)
	}

	p := NewPool(cfg)

	p.Start()
}
