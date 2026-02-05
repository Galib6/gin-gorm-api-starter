package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

// LogLevel represents log severity
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l LogLevel) Color() string {
	switch l {
	case DEBUG:
		return Gray
	case INFO:
		return Green
	case WARN:
		return Yellow
	case ERROR:
		return Red
	case FATAL:
		return BgRed + White
	default:
		return White
	}
}

// Logger is a structured logger with colored output
type Logger struct {
	minLevel LogLevel
	output   io.Writer
	prefix   string
}

var defaultLogger = NewLogger(INFO, os.Stdout, "")

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, output io.Writer, prefix string) *Logger {
	return &Logger{
		minLevel: level,
		output:   output,
		prefix:   prefix,
	}
}

// SetLevel sets the minimum log level
func SetLevel(level LogLevel) {
	defaultLogger.minLevel = level
}

// SetOutput sets the output writer
func SetOutput(w io.Writer) {
	defaultLogger.output = w
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := fmt.Sprintf("%-5s", level.String())
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	// Get caller info
	_, file, line, ok := runtime.Caller(2)
	caller := ""
	if ok {
		// Extract just the filename
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			file = parts[len(parts)-1]
		}
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	prefix := ""
	if l.prefix != "" {
		prefix = fmt.Sprintf("[%s] ", l.prefix)
	}

	output := fmt.Sprintf("%s%s | %s%s%s | %s%-20s%s | %s%s\n",
		Gray, timestamp,
		level.Color(), levelStr, Reset,
		Cyan, caller, Reset,
		prefix, formattedMsg,
	)

	fmt.Fprint(l.output, output)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

// Package-level functions using default logger

func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}

// WithPrefix creates a new logger with a prefix
func WithPrefix(prefix string) *Logger {
	return NewLogger(defaultLogger.minLevel, defaultLogger.output, prefix)
}

// LogHTTPRequest logs an HTTP request with colored status
func LogHTTPRequest(method, path string, statusCode int, latency time.Duration, clientIP string) {
	statusColor := Green
	if statusCode >= 400 && statusCode < 500 {
		statusColor = Yellow
	} else if statusCode >= 500 {
		statusColor = Red
	}

	methodColor := Blue
	switch method {
	case "GET":
		methodColor = Green
	case "POST":
		methodColor = Cyan
	case "PUT", "PATCH":
		methodColor = Yellow
	case "DELETE":
		methodColor = Red
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("%s%s%s | %s%3d%s | %13v | %15s | %s%-7s%s %s\n",
		Gray, timestamp, Reset,
		statusColor, statusCode, Reset,
		latency,
		clientIP,
		methodColor, method, Reset,
		path,
	)
}
