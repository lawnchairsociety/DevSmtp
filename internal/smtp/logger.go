package smtp

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarning
	LogError
	LogDebug
)

type LogEntry struct {
	Time    time.Time
	Level   LogLevel
	Message string
}

func (l LogLevel) String() string {
	switch l {
	case LogInfo:
		return "INFO"
	case LogWarning:
		return "WARN"
	case LogError:
		return "ERROR"
	case LogDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

func (e LogEntry) String() string {
	return fmt.Sprintf("[%s] %s %s",
		e.Time.Format("15:04:05"),
		e.Level.String(),
		e.Message,
	)
}

type Logger struct {
	ch chan LogEntry
}

func NewLogger(bufferSize int) *Logger {
	return &Logger{
		ch: make(chan LogEntry, bufferSize),
	}
}

func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	entry := LogEntry{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprintf(format, args...),
	}

	select {
	case l.ch <- entry:
	default:
		// Channel full, drop oldest entry and try again
		select {
		case <-l.ch:
		default:
		}
		select {
		case l.ch <- entry:
		default:
		}
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(LogInfo, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(LogWarning, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(LogError, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(LogDebug, format, args...)
}

func (l *Logger) Channel() <-chan LogEntry {
	return l.ch
}
