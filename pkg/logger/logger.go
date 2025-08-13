package logger

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger implementa un logger personalizado
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewLogger crea una nueva instancia del logger
func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info registra un mensaje informativo
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.infoLogger.Println(msg)
}

// Error registra un mensaje de error
func (l *SimpleLogger) Error(msg string, err error, args ...interface{}) {
	errorMsg := fmt.Sprintf("%s: %v", msg, err)
	if len(args) > 0 {
		errorMsg = fmt.Sprintf("%s %v", errorMsg, args)
	}
	l.errorLogger.Println(errorMsg)
}

// Debug registra un mensaje de debug
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.debugLogger.Println(msg)
}

// Warn registra un mensaje de advertencia
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf("%s %v", msg, args)
	}
	l.warnLogger.Println(msg)
}
