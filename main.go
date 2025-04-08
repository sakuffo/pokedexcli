package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sakuffo/pokedexcli/internal/app"
	"github.com/sakuffo/pokedexcli/internal/config"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/repl"
)

func main() {
	// Parse command-line flags for log level
	logLevelStr := flag.String("loglevel", "NONE", "Set log level (DEBUG, INFO, ERROR, FATAL, NONE)")
	flag.Parse()

	// Convert string log level to logger.LogLevel
	level := parseLogLevel(*logLevelStr)

	// Initialize the application
	cfg, err := app.Initialize(level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling
	setupSignalHandling(cfg)

	// Start REPL
	repl.StartRepl(cfg)

	// Final save on exit
	cfg.Logger.Info("REPL exited normally. Performing final save...")
	if err := app.SaveData(cfg); err != nil {
		cfg.Logger.Error("Failed to perform final save on exit: %v", err)
	} else {
		cfg.Logger.Info("Final save completed.")
	}

	cfg.Logger.Info("Exiting Pokedex CLI.")
}

// parseLogLevel converts a string log level to logger.LogLevel
func parseLogLevel(levelStr string) logger.LogLevel {
	switch levelStr {
	case "NONE":
		return logger.NONE
	case "DEBUG":
		return logger.DEBUG
	case "INFO":
		return logger.INFO
	case "ERROR":
		return logger.ERROR
	case "FATAL":
		return logger.FATAL
	default:
		fmt.Fprintf(os.Stderr, "Invalid log level: %s. Using INFO.\n", levelStr)
		return logger.INFO
	}
}

// setupSignalHandling sets up signal handling for graceful shutdown
func setupSignalHandling(cfg *config.Config) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		cfg.Logger.Info("Received signal: %v. Saving data before exiting...", sig)

		if err := app.SaveData(cfg); err != nil {
			cfg.Logger.Error("Failed to save data during signal handling: %v", err)
			os.Exit(1)
		}

		cfg.Logger.Info("Data saved successfully. Exiting now.")
		os.Exit(0)
	}()
}
