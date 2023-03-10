package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type UnhealthyGetter interface {
	GetUnhealthy(ctx context.Context) (unhealthies []Container, err error)
}

func (d *Docker) GetUnhealthy(ctx context.Context) (unhealthies []Container, err error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("label", "deunhealth.restart.on.unhealthy=true")
	filtersArgs.Add("health", "unhealthy")

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filtersArgs,
	})

	if err != nil {
		return nil, err
	}

	unhealthies = make([]Container, len(containers))

	for i, container := range containers {
		unhealthies[i] = Container{
			ID:    container.ID,
			Name:  extractName(container),
			Image: container.Image,
		}
	}

	return unhealthies, nil
}

type UnhealthyStreamer interface {
	StreamUnhealthy(ctx context.Context, unhealthies chan<- Container, crashed chan<- error)
}

func (d *Docker) StreamUnhealthy(ctx context.Context, unhealthies chan<- Container, crashed chan<- error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("label", "deunhealth.restart.on.unhealthy=true")
	filtersArgs.Add("health", "unhealthy")

	// See https://github.com/moby/moby/blob/deda3d4933d3c0bd57f2cef672da5d28fc653706/client/events.go
	messages, errors := d.client.Events(ctx, types.EventsOptions{
		Filters: filtersArgs,
	})

	for {
		select {
		case <-ctx.Done():
			<-errors // wait for Events() to exit
			crashed <- ctx.Err()
			return

		case err := <-errors: // Events failed and has exit
			crashed <- err
			return

		case message := <-messages:
			if !isContainerMessage(message) || message.Action != "health_status: unhealthy" {
				break
			}

			unhealthy := Container{
				ID:    message.Actor.ID,
				Name:  extractNameFromActor(message.Actor),
				Image: message.Actor.Attributes["image"],
			}

			select {
			case unhealthies <- unhealthy:
			case <-ctx.Done(): // do not block
			}
		}
	}
}

type LinkedContainerGetter interface {
	GetLinkedContainer(ctx context.Context, unhealthy Container) (linkedContainers []Container, err error)
}

func (d *Docker) GetLinkedContainer(ctx context.Context, unhealthy Container) (linkedContainers []Container, err error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("label", "deunhealth.restart.with.unhealthy.container=" + unhealthy.Name)

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filtersArgs,
	})

	if err != nil {
		return nil, err
	}

	linkedContainers = make([]Container, len(containers))

	for i, container := range containers {
		linkedContainers[i] = Container{
			ID:    container.ID,
			Name:  extractName(container),
			Image: container.Image,
		}
	}

	return linkedContainers, nil
}
