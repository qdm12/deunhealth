package env

import (
	"github.com/qdm12/deunhealth/internal/config/settings"
)

func (r *Reader) readDocker() (s settings.Docker) {
	s.Host = r.getEnv("DOCKER_HOST")
	return s
}
