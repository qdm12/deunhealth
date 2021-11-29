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
		name := container.ID
		if len(container.Names) > 0 && container.Names[0] != "" {
			name = container.Names[0]
		}

		unhealthies[i] = Container{
			ID:    container.ID,
			Name:  name,
			Image: container.Image,
		}
	}

	return unhealthies, nil
}
