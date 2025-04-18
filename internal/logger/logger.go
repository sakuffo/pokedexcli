package logger

// TODO: Does the logger also keep normal fmt.Println() output?

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

const (
	NONE LogLevel = iota
	DEBUG
	INFO
	ERROR
	FATAL
)

type Logger struct {
	level  LogLevel
	writer io.Writer
}

func New(level LogLevel) *Logger {
	return &Logger{level: level, writer: os.Stdout}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= DEBUG {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.writer, "[%s] [DEBUG] "+format+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= INFO {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.writer, "[%s] [INFO] "+format+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= ERROR {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.writer, "[%s] [ERROR] "+format+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.level >= FATAL {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.writer, "[%s] [FATAL] "+format+"\n", append([]interface{}{timestamp}, args...)...)
		os.Exit(1)
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetWriter(writer io.Writer) {
	l.writer = writer
}

// ParseLogLevel converts a string representation of a log level to LogLevel type.
// Returns an error if the provided string is not a valid log level.
func ParseLogLevel(level string) (LogLevel, error) {
	switch level {
	case "NONE":
		return NONE, nil
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return INFO, fmt.Errorf("invalid log level: %s", level)
	}
}
