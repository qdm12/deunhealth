package loop

import (
	"context"

	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/loop/info"
	"github.com/qdm12/log"
)

func newUnhealthyLoop(docker docker.Dockerer, logger log.LeveledLogger) *unhealthyLoop {
	return &unhealthyLoop{
		docker: docker,
		logger: logger,
		info:   info.NewUnhealthyLoop(docker, logger),
	}
}

type unhealthyLoop struct {
	logger log.LeveledLogger
	docker docker.Dockerer
	info   Runner
}

func (l *unhealthyLoop) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	infoStreamCrashed := make(chan error)
	go func() {
		infoStreamCrashed <- l.info.Run(ctx)
	}()

	existingUnhealthies, err := l.docker.GetUnhealthy(ctx)
	if err != nil {
		return err
	}
	for _, unhealthy := range existingUnhealthies {
		l.restartUnhealthy(ctx, unhealthy)
	}

	unhealthies := make(chan docker.Container)
	unhealthyStreamCrashed := make(chan error)

	go l.docker.StreamUnhealthy(ctx, unhealthies, unhealthyStreamCrashed)

	for {
		select {
		case <-ctx.Done():
			<-unhealthyStreamCrashed
			close(unhealthyStreamCrashed)
			<-infoStreamCrashed
			close(infoStreamCrashed)
			close(unhealthies)

			return ctx.Err()

		case err := <-unhealthyStreamCrashed:
			close(unhealthyStreamCrashed)
			cancel()
			<-infoStreamCrashed
			close(infoStreamCrashed)
			close(unhealthies)

			return err

		case err := <-infoStreamCrashed:
			close(infoStreamCrashed)
			cancel()
			<-unhealthyStreamCrashed
			close(unhealthyStreamCrashed)
			close(unhealthies)

			return err

		case unhealthy := <-unhealthies:
			l.restartUnhealthy(ctx, unhealthy)
		}
	}
}

func (l *unhealthyLoop) restartUnhealthy(ctx context.Context, unhealthy docker.Container) {
	l.logger.Info("container " + unhealthy.Name +
		" (image " + unhealthy.Image + ") is unhealthy, restarting it...")
	err := l.docker.RestartContainer(ctx, unhealthy.Name)
	if err != nil {
		l.logger.Error("failed restarting container: " + err.Error())
	} else {
		l.logger.Info("container " + unhealthy.Name + " restarted successfully")
	}
	l.restartLinked(ctx, unhealthy)
}

func (l *unhealthyLoop) restartLinked(ctx context.Context, unhealthy docker.Container) {
	linkedContainers, _ := l.docker.GetLinkedContainer(ctx, unhealthy)
	for _, linkedContainer := range linkedContainers {
		l.logger.Info("container " + linkedContainer.Name +
			" (image " + linkedContainer.Image + ") is linked to unhealthy container " + 
			unhealthy.Name + ", restarting it...")
		err := l.docker.RestartContainer(ctx, linkedContainer.Name)
		if err != nil {
			l.logger.Error("failed restarting linked container: " + err.Error())
		} else {
			l.logger.Info("linked container " + linkedContainer.Name + " restarted successfully")
		}	
	}
}
