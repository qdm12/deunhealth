package loop

import (
	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/log"
)

type Looper interface {
	Runner
}

type Loop struct {
	unhealthy Runner
}

func New(docker docker.Dockerer, logger log.LeveledLogger) *Loop {
	return &Loop{
		unhealthy: newUnhealthyLoop(docker, logger),
	}
}
