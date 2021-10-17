package env

import (
	"github.com/qdm12/deunhealth/internal/config/settings"
)

func (r *Reader) readHealth() (s settings.Health) {
	s.Address = r.getEnv("HEALTH_SERVER_ADDRESS")
	return s
}
