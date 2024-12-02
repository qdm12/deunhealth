package loop

import (
	"github.com/qdm12/deunhealth/internal/loop/info"
	"github.com/qdm12/log"
)

type Loop struct {
	runners []runner
}

func New(docker Docker, logger log.LeveledLogger) *Loop {
	return &Loop{
		runners: []runner{
			info.NewUnhealthyLoop(docker, logger),
			newUnhealthyLoop(docker, logger),
		},
	}
}
