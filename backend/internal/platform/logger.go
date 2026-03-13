package platform

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func NewLogger(level string) zerolog.Logger {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return logger.Level(parsedLevel)
}
