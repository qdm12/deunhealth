package env

import (
	"fmt"

	"github.com/qdm12/deunhealth/internal/config/settings"
	"github.com/qdm12/log"
)

func (r *Reader) readLog() (s settings.Log, err error) {
	levelString := r.getEnv("LOG_LEVEL")
	if levelString != "" {
		level, err := log.ParseLevel(levelString)
		if err != nil {
			return s, fmt.Errorf("environment variable LOG_LEVEL: %w", err)
		}
		s.Level = &level
	}

	return s, nil
}
