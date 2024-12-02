package docker

import (
	"context"
)

func (d *Docker) RestartContainer(ctx context.Context, name string) (err error) {
	return d.client.ContainerRestart(ctx, name, nil)
}
