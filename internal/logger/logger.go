package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	currentLevel Level
	jsonFormat   bool
)

// Init initializes the logger
func Init(level, format string) {
	switch level {
	case "debug":
		currentLevel = DebugLevel
	case "info":
		currentLevel = InfoLevel
	case "warn":
		currentLevel = WarnLevel
	case "error":
		currentLevel = ErrorLevel
	case "fatal":
		currentLevel = FatalLevel
	default:
		currentLevel = InfoLevel
	}

	jsonFormat = format == "json"
	log.SetFlags(0)
}

// Debug logs debug message
func Debug(message string, fields map[string]interface{}) {
	if currentLevel <= DebugLevel {
		logMessage("DEBUG", message, fields)
	}
}

// Info logs info message
func Info(message string, fields map[string]interface{}) {
	if currentLevel <= InfoLevel {
		logMessage("INFO", message, fields)
	}
}

// Warn logs warning message
func Warn(message string, fields map[string]interface{}) {
	if currentLevel <= WarnLevel {
		logMessage("WARN", message, fields)
	}
}

// Error logs error message
func Error(message string, fields map[string]interface{}) {
	if currentLevel <= ErrorLevel {
		logMessage("ERROR", message, fields)
	}
}

// Fatal logs fatal message and exits
func Fatal(message string, fields map[string]interface{}) {
	logMessage("FATAL", message, fields)
	os.Exit(1)
}

func logMessage(level, message string, fields map[string]interface{}) {
	if jsonFormat {
		logJSON(level, message, fields)
	} else {
		logText(level, message, fields)
	}
}

func logJSON(level, message string, fields map[string]interface{}) {
	entry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     level,
		"message":   message,
	}

	if fields != nil {
		for k, v := range fields {
			entry[k] = v
		}
	}

	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func logText(level, message string, fields map[string]interface{}) {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	output := fmt.Sprintf("[%s] %s - %s", timestamp, level, message)

	if fields != nil && len(fields) > 0 {
		output += " |"
		for k, v := range fields {
			output += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	fmt.Println(output)
}
