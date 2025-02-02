package logger

import (
	"errors"
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
		return NewTestLogger()
	default: // dev, development, or any other string
		return NewDevelopmentLogger()
	}
}

func NewProductionLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}

func NewTestLogger() *zerolog.Logger {
	rootPath, err := GetProjectRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error getting current working directory: %w", err).Error())
		os.Exit(1)
	}

	logDir := filepath.Join(rootPath, "log")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("error creating log director: %w %s", err, logDir).Error())
			os.Exit(1)
		}
	}

	logPath := filepath.Join(logDir, "test.log")

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

// GetProjectRoot attempts to locate the project root by searching upward
// from the current working directory for a marker file or directory.
// In this example, we use "go.mod" and ".git" as markers.
func GetProjectRoot() (string, error) {
	// Start from the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if "go.mod" exists in the current directory
		if exists, err := fileExists(filepath.Join(currentDir, "go.mod")); err == nil && exists {
			return currentDir, nil
		}

		// Alternatively, check if the ".git" folder exists in the current directory
		if exists, err := dirExists(filepath.Join(currentDir, ".git")); err == nil && exists {
			return currentDir, nil
		}

		// Get the parent directory
		parentDir := filepath.Dir(currentDir)
		// If we've reached the root (parent is same as current), stop searching
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", errors.New("project root not found: ensure that a marker file (e.g., go.mod or .git) exists in your project root")
}

// fileExists checks if the file at the given path exists and is not a directory.
func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// dirExists checks if the directory at the given path exists.
func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}
