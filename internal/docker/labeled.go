package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type LabeledGetter interface {
	GetLabeled(ctx context.Context, labels []string) (
		containerNames []string, err error)
	StreamLabeled(ctx context.Context, labels []string,
		containerNames chan<- string, crashed chan<- error)
}

var ErrListContainers = errors.New("cannot list containers")

func (d *Docker) GetLabeled(ctx context.Context, labels []string) (
	containerNames []string, err error) {
	// See https://docs.docker.com/engine/reference/commandline/ps/#filtering
	filtersArgs := filters.NewArgs()
	for _, label := range labels {
		filtersArgs.Add("label", label)
	}

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filtersArgs,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrListContainers, err)
	}

	containerNames = make([]string, len(containers))

	for i, container := range containers {
		containerNames[i] = extractName(container)
	}

	return containerNames, nil
}

func (d *Docker) StreamLabeled(ctx context.Context, labels []string,
	containerNames chan<- string, crashed chan<- error) {
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
			if message.Type != "container" || message.Action != "starting" { // TODO starting
				break
			}

			containerName := extractNameFromActor(message.Actor)

			select {
			case containerNames <- containerName:
			case <-ctx.Done(): // do not block
			}
		}
	}
}
