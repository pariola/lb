package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pariola/lb/pkg/config"
)

var configFile = "lb.yml"

func main() {

	args := os.Args[1:]

	if len(args) > 0 && args[0] != "" {
		configFile = args[0]
	}

	fPath, err := filepath.Abs(configFile)

	if err != nil {
		log.Fatal("no valid configuration file supplied")
	}

	f, err := os.Open(fPath)

	if err != nil {
		log.Fatal("failed to open configuration file")
	}

	cfg, err := config.Load(f)

	if err != nil {
		log.Fatal("failed to parse configuration file")
	}

	p := NewPool(cfg)

	p.Start()
}
