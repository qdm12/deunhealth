package docker

import (
	"github.com/moby/moby/client"
)

type Docker struct {
	client client.CommonAPIClient
}

func New(dockerHost string) (d *Docker, err error) {
	client, err := client.NewClientWithOpts(client.WithHost(dockerHost))
	if err != nil {
		return nil, err
	}

	return &Docker{
		client: client,
	}, nil
}
