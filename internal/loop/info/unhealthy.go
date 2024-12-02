package info

import (
	"context"

	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/loop/helpers"
)

func NewUnhealthyLoop(docker Docker, infoer Logger) *UnhealthyLoop {
	return &UnhealthyLoop{
		docker:       docker,
		logger:       infoer,
		monitoredIDs: make(map[string]struct{}),
	}
}

type UnhealthyLoop struct {
	logger       Logger
	docker       Docker
	monitoredIDs map[string]struct{}
}

func (l *UnhealthyLoop) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	healthMonitorLabels := []string{"deunhealth.restart.on.unhealthy=true"}
	onUnhealthyContainers, err := l.docker.GetLabeled(ctx, healthMonitorLabels)
	if err != nil {
		return err
	}

	containerNames := make([]string, len(onUnhealthyContainers))
	for i, container := range onUnhealthyContainers {
		l.monitoredIDs[container.ID] = struct{}{}
		containerNames[i] = container.Name
	}

	switch len(containerNames) {
	case 0:
		l.logger.Infof("No container found to restart when becoming unhealthy")
	case 1:
		l.logger.Infof("Monitoring container %s to restart when becoming unhealthy",
			containerNames[0])
	default:
		l.logger.Infof("Monitoring containers %s to restart when becoming unhealthy",
			helpers.BuildEnum(containerNames))
	}

	healthMonitored := make(chan docker.Container)
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

		case container := <-healthMonitored:
			_, alreadyMonitored := l.monitoredIDs[container.ID]
			if alreadyMonitored {
				break
			}
			l.monitoredIDs[container.ID] = struct{}{}
			l.logger.Infof("Monitoring new container %s to restart when becoming unhealthy",
				container.Name)
		}
	}
}
