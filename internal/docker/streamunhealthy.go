package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type UnhealthyStreamer interface {
	StreamUnhealthy(ctx context.Context, unhealthies chan<- UnhealthyContainer, crashed chan<- error)
}

type UnhealthyContainer struct {
	Name  string
	Image string
}

func (d *Docker) StreamUnhealthy(ctx context.Context, unhealthies chan<- UnhealthyContainer, crashed chan<- error) {
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
			if message.Type != "container" || message.Action != "health_status: unhealthy" {
				break
			}

			unhealthy := UnhealthyContainer{
				Name:  message.Actor.Attributes["name"],
				Image: message.Actor.Attributes["image"],
			}

			if unhealthy.Name == "" {
				unhealthy.Name = message.Actor.ID
			}

			select {
			case unhealthies <- unhealthy:
			case <-ctx.Done(): // do not block
			}
		}
	}
}
