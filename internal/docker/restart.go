package docker

import (
	"context"
)

type ContainerRestarter interface {
	RestartContainer(ctx context.Context, name string) error
}

func (d *Docker) RestartContainer(ctx context.Context, name string) (err error) {
	return d.client.ContainerRestart(ctx, name, nil)
}
