package loop

import (
	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/loop/info"
	"github.com/qdm12/log"
)

type Looper interface {
	Runner
}

type Loop struct {
	runners []Runner
}

func New(docker docker.Dockerer, logger log.LeveledLogger) *Loop {
	return &Loop{
		runners: []Runner{
			info.NewUnhealthyLoop(docker, logger),
			newUnhealthyLoop(docker, logger),
		},
	}
}
