package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

// In production mode, log to stdout in JSON format
// in development mode, log to stdout in human readable format
// in test mode, log to "log/test.log" in human readable format.
func NewLogger(mode string) *zerolog.Logger {
	switch strings.ToLower(mode) {
	case "prod", "production":
		return NewProductionLogger()
	case "test":
		return NewTestLogger("log", "test.log")
	default: // dev, development, or any other string
		return NewDevelopmentLogger()
	}
}

func NewProductionLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}

// NewTestLogger creates a logger that writes to a file in the log directory.
// first argument is the log directory, second argument is the log file name.
func NewTestLogger(paths ...string) *zerolog.Logger {
	filename := ""
	logDir := ""

	if len(paths) > 0 {
		logDir = paths[0]
	}

	if len(paths) > 1 {
		filename = paths[1]
	}

	if filename == "" {
		filename = "test.log"
	}

	if logDir == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("error getting working directory: %w", err).Error())
			os.Exit(1)
		}

		logDir = filepath.Join(workingDir, "log")
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("error creating log director: %w %s", err, logDir).Error())
			os.Exit(1)
		}
	}

	logPath := filepath.Join(logDir, filename)

	fmt.Println("Logging to file:", logPath)

	const logFilePerm = 0o644
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, logFilePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error opening log file: %w %s", err, logPath).Error())
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	writer := zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: "02-01-2006 15:04:05 Z07:00",
		NoColor:    true,
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

func NewDevelopmentLogger() *zerolog.Logger {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 Z07:00",
		NoColor:    !isTerminal,
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}
