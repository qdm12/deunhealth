package env

import (
	"github.com/qdm12/deunhealth/internal/config/settings"
)

func (r *Reader) readLog() (s settings.Log) {
	s.Level = r.getEnv("LOG_LEVEL")
	return s
}
