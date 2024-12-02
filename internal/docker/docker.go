package docker

import (
	"errors"
	"fmt"

	"github.com/moby/moby/client"
)

type Docker struct {
	client client.CommonAPIClient
}

var (
	ErrCreateDockerClient = errors.New("cannot create Docker client")
)

func New(dockerHost string) (d *Docker, err error) {
	client, err := client.NewClientWithOpts(client.WithHost(dockerHost))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateDockerClient, err)
	}

	return &Docker{
		client: client,
	}, nil
}
