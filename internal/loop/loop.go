package loop

import (
	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Runner
}

type Loop struct {
	docker docker.Dockerer
	logger logging.Logger
}

func New(docker docker.Dockerer, logger logging.Logger) *Loop {
	return &Loop{
		docker: docker,
		logger: logger,
	}
}
