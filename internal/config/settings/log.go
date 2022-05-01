package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/log"
)

type Log struct {
	Level *log.Level
}

func (l *Log) setDefaults() {
	if l.Level == nil {
		defaultLevel := log.LevelInfo
		l.Level = &defaultLevel
	}
}

func (l *Log) mergeWith(other Log) {
	if l.Level == nil {
		l.Level = other.Level
	}
}

var (
	ErrLogLevel = errors.New("invalid log level")
)

func (l *Log) validate() (err error) {
	switch *l.Level {
	case log.LevelError, log.LevelWarn, log.LevelInfo, log.LevelDebug:
	default:
		return fmt.Errorf("%w: %s", ErrLogLevel, l.Level)
	}
	return nil
}
