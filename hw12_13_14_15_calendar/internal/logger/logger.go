package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

const (
	DEBUG int = iota
	INFO
	WARNING
	ERROR
)

const logFlags = log.LstdFlags

func toIntLevel(level string) int {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

type Printer interface {
	Println(v ...interface{})
	Printf(fmt string, v ...interface{})
}

type Logger struct {
	level   int
	printer Printer
	closer  io.Closer
}

func New(level string, filePath string) (*Logger, error) {
	if len(filePath) == 0 {
		return &Logger{
			level:   toIntLevel(level),
			printer: log.New(os.Stderr, "", logFlags),
		}, nil
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		level:   toIntLevel(level),
		printer: log.New(file, "", logFlags),
		closer:  file,
	}, nil
}

func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}

	return nil
}

func (l *Logger) log(level int, v ...interface{}) {
	if level >= l.level {
		l.printer.Println(v...)
	}
}

func (l *Logger) logf(level int, fmt string, v ...interface{}) {
	if level >= l.level {
		l.printer.Printf(fmt, v...)
	}
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, msg)
}

func (l *Logger) Debugf(fmt string, v ...interface{}) {
	l.logf(DEBUG, fmt, v...)
}

func (l Logger) Info(msg string) {
	l.log(INFO, msg)
}

func (l Logger) Infof(fmt string, v ...interface{}) {
	l.logf(INFO, fmt, v...)
}

func (l Logger) Warning(msg string) {
	l.log(WARNING, msg)
}

func (l Logger) Warningf(fmt string, v ...interface{}) {
	l.logf(WARNING, fmt, v...)
}

func (l Logger) Error(msg string) {
	l.log(ERROR, msg)
}

func (l Logger) Errorf(fmt string, v ...interface{}) {
	l.logf(ERROR, fmt, v...)
}
