package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	logLevelNames = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}

	currentLogLevel = INFO
	logger          *log.Logger
)

// InitLogger initializes the logger with the specified configuration
func InitLogger(level string, output io.Writer) {
	if output == nil {
		output = os.Stdout
	}

	logger = log.New(output, "", 0)

	// Set log level from string
	switch strings.ToUpper(level) {
	case "DEBUG":
		currentLogLevel = DEBUG
	case "INFO":
		currentLogLevel = INFO
	case "WARN", "WARNING":
		currentLogLevel = WARN
	case "ERROR":
		currentLogLevel = ERROR
	default:
		currentLogLevel = INFO
	}
}

// logMessage logs a message at the specified level
func logMessage(level LogLevel, format string, args ...interface{}) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	if level < currentLogLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := logLevelNames[level]
	message := fmt.Sprintf(format, args...)

	logger.Printf("[%s] %s - %s", timestamp, levelName, message)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	logMessage(DEBUG, format, args...)
}

// LogInfo logs an informational message
func LogInfo(format string, args ...interface{}) {
	logMessage(INFO, format, args...)
}

// LogWarn logs a warning message
func LogWarn(format string, args ...interface{}) {
	logMessage(WARN, format, args...)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	logMessage(ERROR, format, args...)
}

// LogRequest logs an HTTP request
func LogRequest(method, path, remoteAddr string, statusCode int, duration time.Duration) {
	LogInfo("HTTP %s %s from %s - %d (%v)", method, path, remoteAddr, statusCode, duration)
}
