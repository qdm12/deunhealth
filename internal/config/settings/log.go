package settings

import (
	"errors"
	"fmt"
)

type Log struct {
	Level string
}

func (l *Log) setDefaults() {
	if l.Level == "" {
		l.Level = "info"
	}
}

func (l *Log) mergeWith(other Log) {
	if l.Level == "" {
		l.Level = other.Level
	}
}

var (
	ErrLogLevel = errors.New("invalid log level")
)

func (l *Log) validate() (err error) {
	switch l.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("%w: %s", ErrLogLevel, l.Level)
	}
	return nil
}
