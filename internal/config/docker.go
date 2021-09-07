package config

import (
	"fmt"

	"github.com/moby/moby/client"
	"github.com/qdm12/golibs/params"
)

type Docker struct {
	Host string
}

func (h *Docker) get(env params.Env) (err error) {
	h.Host, err = env.Get("DOCKER_HOST", params.Default(client.DefaultDockerHost))
	if err != nil {
		return fmt.Errorf("environment variable DOCKER_HOST: %w", err)
	}

	return nil
}
