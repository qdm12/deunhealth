package loop

import (
	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Runner
}

type Loop struct {
	unhealthy Runner
}

func New(docker docker.Dockerer, logger logging.Logger) *Loop {
	return &Loop{
		unhealthy: newUnhealthyLoop(docker, logger),
	}
}
