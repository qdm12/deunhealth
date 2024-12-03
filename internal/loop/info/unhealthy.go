package info

import (
	"context"
	"fmt"

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
	cancel       context.CancelFunc
	done         <-chan struct{}
}

func (l *UnhealthyLoop) String() string {
	return "unhealthy info loop"
}

func (l *UnhealthyLoop) Start(ctx context.Context) (runError <-chan error, err error) {
	ready := make(chan struct{})
	done := make(chan struct{})
	l.done = done
	runErrorCh := make(chan error)
	runCtx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	go l.run(runCtx, ready, done, runErrorCh) //nolint:contextcheck
	select {
	case <-ctx.Done():
		l.cancel()
		<-done
		return nil, ctx.Err()
	case <-ready:
		return runErrorCh, nil
	}
}

func (l *UnhealthyLoop) Stop() (err error) {
	l.cancel()
	<-l.done
	return nil
}

func (l *UnhealthyLoop) run(ctx context.Context, ready chan<- struct{},
	done chan<- struct{}, runError chan<- error) {
	defer close(done)
	close(ready)

	healthMonitorLabels := []string{"deunhealth.restart.on.unhealthy=true"}
	onUnhealthyContainers, err := l.docker.GetLabeled(ctx, healthMonitorLabels)
	if err != nil {
		runError <- fmt.Errorf("getting health monitored containers: %w", err)
		return
	}

	containerNames := make([]string, len(onUnhealthyContainers))
	for i, container := range onUnhealthyContainers {
		l.monitoredIDs[container.ID] = struct{}{}
		containerNames[i] = container.Name
	}
	logContainerNames(l.logger, containerNames)

	healthMonitored := make(chan docker.Container)
	healthStreamCrashed := make(chan error)

	go l.docker.StreamLabeled(ctx, healthMonitorLabels, healthMonitored, healthStreamCrashed)

	for {
		select {
		case <-ctx.Done(): // stop requested
			<-healthStreamCrashed
			return

		case err := <-healthStreamCrashed:
			runError <- fmt.Errorf("streaming unhealthy containers: %w", err)
			return

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

func logContainerNames(logger Logger, containerNames []string) {
	switch len(containerNames) {
	case 0:
		logger.Infof("No container found to restart when becoming unhealthy")
	case 1:
		logger.Infof("Monitoring container %s to restart when becoming unhealthy",
			containerNames[0])
	default:
		logger.Infof("Monitoring containers %s to restart when becoming unhealthy",
			helpers.BuildEnum(containerNames))
	}
}
