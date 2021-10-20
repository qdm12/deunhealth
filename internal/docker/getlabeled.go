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
		containerNames[i] = container.ID
		if len(container.Names) > 0 && container.Names[0] != "" {
			containerNames[i] = container.Names[0]
		}
	}

	return containerNames, nil
}
