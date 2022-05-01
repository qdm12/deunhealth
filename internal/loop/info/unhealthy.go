package info

import (
	"context"

	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/loop/helpers"
	"github.com/qdm12/golibs/logging"
)

func NewUnhealthyLoop(docker docker.Dockerer, logger logging.Logger) *UnhealthyLoop {
	return &UnhealthyLoop{
		docker: docker,
		logger: logger,
	}
}

type UnhealthyLoop struct {
	logger logging.Logger
	docker docker.Dockerer
}

func (l *UnhealthyLoop) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	healthMonitorLabels := []string{"deunhealth.restart.on.unhealthy=true"}
	onUnhealthyNames, err := l.docker.GetLabeled(ctx, healthMonitorLabels)
	if err != nil {
		return err
	}
	l.logger.Info("Monitoring containers " + helpers.BuildEnum(onUnhealthyNames) + " to restart when becoming unhealthy")

	healthMonitored := make(chan string)
	healthStreamCrashed := make(chan error)

	go l.docker.StreamLabeled(ctx, healthMonitorLabels, healthMonitored, healthStreamCrashed)

	for {
		select {
		case <-ctx.Done():
			<-healthStreamCrashed
			close(healthStreamCrashed)
			close(healthMonitored)

			return ctx.Err()

		case err := <-healthStreamCrashed:
			close(healthStreamCrashed)
			close(healthMonitored)

			return err

		case healthMonitorName := <-healthMonitored:
			l.logger.Info("Monitoring new container " + healthMonitorName + " to restart when becoming unhealthy")
		}
	}
}
