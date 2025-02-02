// pkg/logger/goose_writer.go

package logger

import (
	"strings"

	"github.com/rs/zerolog"
)

// GooseWriter implements io.Writer to forward log messages to Zerolog with level mapping.
type GooseWriter struct {
	logger *zerolog.Logger
}

// NewGooseWriter creates a new GooseWriter instance.
func NewGooseWriter(logger *zerolog.Logger) *GooseWriter {
	return &GooseWriter{logger: logger}
}

// Write captures Goose's log output and forwards it to Zerolog with appropriate log levels.
func (gw GooseWriter) Write(logData []byte) (int, error) {
	msg := strings.TrimSpace(string(logData))
	if len(msg) == 0 {
		return len(logData), nil
	}

	msg = strings.TrimSpace(msg[19:]) // Remove timestamp from the message

	lowerMsg := strings.ToLower(msg)
	switch {
	case strings.Contains(lowerMsg, "error"):
		gw.logger.Error().Msg(msg)
	case strings.Contains(lowerMsg, "warn"):
		gw.logger.Warn().Msg(msg)
	default:
		gw.logger.Info().Msg(msg)
	}

	return len(logData), nil
}
