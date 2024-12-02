package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

var ErrListContainers = errors.New("cannot list containers")

func (d *Docker) GetLabeled(ctx context.Context, labels []string) (
	containers []Container, err error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	for _, label := range labels {
		filtersArgs.Add("label", label)
	}

	list, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filtersArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrListContainers, err)
	}

	containers = make([]Container, len(list))
	for i, container := range list {
		containers[i] = Container{
			ID:    container.ID,
			Image: container.Image,
		}
		if len(container.Names) > 0 {
			containers[i].Name = container.Names[0]
		}
	}

	return containers, nil
}

func (d *Docker) StreamLabeled(ctx context.Context, labels []string,
	containers chan<- Container, crashed chan<- error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	for _, label := range labels {
		filtersArgs.Add("label", label)
	}

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
			if !isContainerMessage(message) {
				break
			}

			container := Container{
				ID:    message.ID,
				Name:  extractNameFromActor(message.Actor),
				Image: extractImageFromActor(message.Actor),
			}

			select {
			case containers <- container:
			case <-ctx.Done(): // do not block
			}
		}
	}
}
