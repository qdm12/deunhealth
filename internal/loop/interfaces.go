package loop

import (
	"context"

	"github.com/qdm12/deunhealth/internal/docker"
)

type Docker interface {
	GetLabeled(ctx context.Context, labels []string) (
		containers []docker.Container, err error)
	StreamLabeled(ctx context.Context, labels []string,
		containers chan<- docker.Container, crashed chan<- error)
	GetUnhealthy(ctx context.Context) (containers []docker.Container, err error)
	StreamUnhealthy(ctx context.Context,
		containers chan<- docker.Container, crashed chan<- error)
	RestartContainer(ctx context.Context, containerID string) (err error)
}
