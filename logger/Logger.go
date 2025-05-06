// Package logger provides a simple and customizable logging system with support for context-based logging
// and structured logging in both text and JSON formats.
//
// This package allows you to log messages with various log levels (DEBUG, INFO, WARN, ERROR), including additional
// contextual information such as file name, function name, line number, and request ID.
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity level of the log.
type LogLevel int

const (
	// DebugLevel Log levels for the logger.
	DebugLevel LogLevel = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger defines the interface for a logger with methods for different log levels.
type Logger interface {
	Debug(msg string, args ...any)                    // Logs a message with DEBUG level
	Info(msg string, args ...any)                     // Logs a message with INFO level
	Warn(msg string, args ...any)                     // Logs a message with WARN level
	Error(msg string, args ...any)                    // Logs a message with ERROR level
	InfoWithFields(msg string, fields map[string]any) // Logs a message with INFO level and additional fields
	WithContext(ctx context.Context) Logger           // Returns a logger with added context
}

// StdLogger is the default implementation of the Logger interface.
type StdLogger struct {
	currentLevel LogLevel    // Minimum log level to log messages
	format       string      // Log format: "text" or "json"
	logger       *log.Logger // Underlying logger instance
}

type logContextKey string

const (
	// RequestIDKey is the context key for storing the request ID in logs.
	RequestIDKey logContextKey = "requestID"
)

// DefaultLogger creates a new logger instance with the specified configuration.
// Example:
//
//	logger := DefaultLogger()
func DefaultLogger() Logger {
	return NewLogger(DefaultConfig())
}

// NewLogger creates a new logger instance with the specified configuration.
// Example:
//
//	logger := NewLogger(&Config{
//	    LogLevel: "info",
//	    LogFormat: "json",
//	    LogType: "file",
//	    LogDir: "./logs",
//	})
func NewLogger(cfg *Config) Logger {
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// If log type is set to file, create the log directory and file
	if strings.ToLower(cfg.LogType) == "file" {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create log dir: %v\n", err)
		} else {
			f, err := os.OpenFile(cfg.LogDir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open log file: %v, fallback to stdout only\n", err)
			} else {
				writers = append(writers, f)
			}
		}
	}

	multiWriter := io.MultiWriter(writers...)
	return &StdLogger{
		logger:       log.New(multiWriter, "", log.LstdFlags|log.Lshortfile),
		currentLevel: parseLogLevel(cfg.LogLevel),
		format:       strings.ToLower(cfg.LogFormat),
	}
}

// parseLogLevel converts a string log level to LogLevel type.
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// shouldLog checks if the message level is greater than or equal to the configured log level.
func (l *StdLogger) shouldLog(level LogLevel) bool {
	return level >= l.currentLevel
}

// log logs a message at the specified level, format, and context.
func (l *StdLogger) log(level LogLevel, levelStr string, msg string, ctx context.Context, args ...any) {
	if !l.shouldLog(level) {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	// Get file name, function name, and line number for debugging
	_, file, line, _ := runtime.Caller(3)
	funcName := "unknown"
	if pc, _, _, ok := runtime.Caller(2); ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	// Get requestID from context
	requestID := ctx.Value(RequestIDKey)

	// Log in JSON or text format
	if l.format == "json" {
		entry := map[string]any{
			"level":     levelStr,
			"time":      time.Now().Format(time.RFC3339),
			"message":   formatted,
			"file":      file,
			"function":  funcName,
			"line":      line,
			"requestID": requestID,
		}
		jsonData, _ := json.Marshal(entry)
		l.logger.Println(string(jsonData))
	} else {
		l.logger.Printf("[%s] %s | file=%s, function=%s, line=%d, requestID=%v | %s",
			strings.ToUpper(levelStr), time.Now().Format(time.RFC3339), file, funcName, line, requestID, formatted)
	}
}

// Debug logs a message with DEBUG level.
func (l *StdLogger) Debug(msg string, args ...any) {
	l.log(DebugLevel, "debug", msg, nil, args...)
}

// Info logs a message with INFO level.
func (l *StdLogger) Info(msg string, args ...any) {
	l.log(InfoLevel, "info", msg, nil, args...)
}

// Warn logs a message with WARN level.
func (l *StdLogger) Warn(msg string, args ...any) {
	l.log(WarnLevel, "warn", msg, nil, args...)
}

// Error logs a message with ERROR level.
func (l *StdLogger) Error(msg string, args ...any) {
	l.log(ErrorLevel, "error", msg, nil, args...)
}

// InfoWithFields logs a message with INFO level along with additional fields.
func (l *StdLogger) InfoWithFields(msg string, fields map[string]any) {
	if !l.shouldLog(InfoLevel) {
		return
	}
	if l.format == "json" {
		entry := map[string]any{
			"level":   "info",
			"time":    time.Now().Format(time.RFC3339),
			"message": msg,
		}
		for k, v := range fields {
			entry[k] = v
		}
		jsonData, _ := json.Marshal(entry)
		l.logger.Println(string(jsonData))
	} else {
		var parts []string
		for k, v := range fields {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		l.logger.Printf("[INFO] %s | %s", msg, strings.Join(parts, ", "))
	}
}

// WithContext adds context to the logger, allowing you to include information like requestID in your logs.
// Example:
//
//	ctx := context.WithValue(context.Background(), logger.RequestIDKey, "my-request-id")
//	logger := logger.WithContext(ctx)
//	logger.Info("Request received")
func (l *StdLogger) WithContext(ctx context.Context) Logger {
	return &StdLoggerWithContext{
		StdLogger: l,
		ctx:       ctx,
	}
}

type StdLoggerWithContext struct {
	*StdLogger
	ctx context.Context
}

// log logs a message using the context provided, overriding the default context.
func (l *StdLoggerWithContext) log(level LogLevel, levelStr string, msg string, ctx context.Context, args ...any) {
	if ctx == nil {
		ctx = l.ctx
	}
	l.StdLogger.log(level, levelStr, msg, ctx, args...)
}

// Debug logs a message with DEBUG level, using the context provided.
func (l *StdLoggerWithContext) Debug(msg string, args ...any) {
	l.log(DebugLevel, "debug", msg, l.ctx, args...)
}

// Info logs a message with INFO level, using the context provided.
func (l *StdLoggerWithContext) Info(msg string, args ...any) {
	l.log(InfoLevel, "info", msg, l.ctx, args...)
}

// Warn logs a message with WARN level, using the context provided.
func (l *StdLoggerWithContext) Warn(msg string, args ...any) {
	l.log(WarnLevel, "warn", msg, l.ctx, args...)
}

// Error logs a message with ERROR level, using the context provided.
func (l *StdLoggerWithContext) Error(msg string, args ...any) {
	l.log(ErrorLevel, "error", msg, l.ctx, args...)
}

// InfoWithFields logs a message with INFO level and additional fields, using the context provided.
func (l *StdLoggerWithContext) InfoWithFields(msg string, fields map[string]any) {
	l.StdLogger.InfoWithFields(msg, fields)
}
