// Package logger exposes logging capabilities, that can be configured through
// the application config
package logger

import (
	"fmt"

	"github.com/adrienlucbert/gofeur/config"
)

type logLevel uint

const (
	debugLevel logLevel = iota
	infoLevel
	warnLevel
	errorLevel
	noneLevel
)

func logLevelFromString(level string) logLevel {
	return map[string]logLevel{
		"Debug": debugLevel,
		"Info":  infoLevel,
		"Warn":  warnLevel,
		"Error": errorLevel,
		"None":  noneLevel,
	}[level]
}

func shouldLog(level logLevel) bool {
	return logLevelFromString(config.GetOr("logLevel", "Debug").(string)) <= level
}

// SetLogLevel sets the debug level in the application config
func SetLogLevel(level string) {
	config.Set("logLevel", level)
	// TODO: return error if level string is invalid
}

func log(level logLevel, format string, a ...any) {
	if !shouldLog(level) {
		return
	}
	fmt.Printf(format, a...)
}

// Debug logs a value if debug logs are enabled
func Debug(format string, a ...any) {
	log(debugLevel, format, a...)
}

// Info logs a value if info logs are enabled
func Info(format string, a ...any) {
	log(infoLevel, format, a...)
}

// Warn logs a value if warn logs are enabled
func Warn(format string, a ...any) {
	log(warnLevel, format, a...)
}

// Error logs a value if error logs are enabled
func Error(format string, a ...any) {
	log(errorLevel, format, a...)
}
