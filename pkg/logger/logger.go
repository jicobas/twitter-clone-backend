package logger

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger implements a custom logger
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an informative message
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.infoLogger.Println(msg)
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, err error, args ...interface{}) {
	errorMsg := fmt.Sprintf("%s: %v", msg, err)
	if len(args) > 0 {
		errorMsg = fmt.Sprintf("%s %v", errorMsg, args)
	}
	l.errorLogger.Println(errorMsg)
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.debugLogger.Println(msg)
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.warnLogger.Println(msg)
}
