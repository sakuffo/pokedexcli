package main

import (
	"flag"
	"time"

	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
	"golang.org/x/exp/rand"
)

func main() {
	// Define flags
	logLevel := flag.String("log", "none", "Set log level (none, debug, info, error, fatal)")
	flag.Parse()

	var level logger.LogLevel
	switch *logLevel {
	case "debug":
		level = logger.DEBUG
	case "info":
		level = logger.INFO
	case "error":
		level = logger.ERROR
	case "fatal":
		level = logger.FATAL
	default:
		level = logger.NONE
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	// Initialize configuration with persistence
	cfg := pokedata.New(level)

	startRepl(cfg)
}
