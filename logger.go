package main

import(
	"os"
	"time"
	"github.com/rs/zerolog"
)

// Initialize logging
func initLogger() {
	if isDev {  // Development mode: console output with colors
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.DebugLevel).With().Timestamp().Caller().Logger()
	} else {  // Production mode: JSON format
		logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}
}