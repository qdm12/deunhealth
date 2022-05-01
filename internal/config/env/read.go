package env

import (
	"github.com/qdm12/deunhealth/internal/config/settings"
)

func (r *Reader) Read() (s settings.Settings, err error) {
	s.Docker = r.readDocker()
	s.Health = r.readHealth()

	s.Log, err = r.readLog()
	if err != nil {
		return s, err
	}

	return s, nil
}
