package info

import (
	"context"

	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/loop/helpers"
)

func NewUnhealthyLoop(docker docker.Dockerer, infoer Infoer) *UnhealthyLoop {
	return &UnhealthyLoop{
		docker:       docker,
		infoer:       infoer,
		monitoredIDs: make(map[string]struct{}),
	}
}

type Infoer interface {
	Info(s string)
}

type UnhealthyLoop struct {
	infoer       Infoer
	docker       docker.Dockerer
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

	l.infoer.Info("Monitoring containers " + helpers.BuildEnum(containerNames) + " to restart when becoming unhealthy")

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
			l.infoer.Info("Monitoring new container " + container.Name + " to restart when becoming unhealthy")
		}
	}
}
