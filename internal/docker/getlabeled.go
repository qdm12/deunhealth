package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type LabeledGetter interface {
	GetLabeled(ctx context.Context, labels []string) (
		containerNames []string, err error)
}

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
		return nil, err
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
