package settings

import "github.com/moby/moby/client"

type Docker struct {
	Host string
}

func (d *Docker) setDefaults() {
	if d.Host == "" {
		d.Host = client.DefaultDockerHost
	}
}

func (d *Docker) mergeWith(other Docker) {
	if d.Host == "" {
		d.Host = other.Host
	}
}

func (d *Docker) validate() error {
	return nil
}
