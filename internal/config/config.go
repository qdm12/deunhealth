// Package config takes care of reading and checking the program configuration
// from environment variables.
package config

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
)

type Config struct {
	Docker Docker
	Log    Log
	Health Health
}

var (
	ErrDockerConfig = errors.New("cannot obtain docker config")
	ErrLogConfig    = errors.New("cannot obtain log config")
	ErrHealthConfig = errors.New("cannot obtain health config")
)

func (c *Config) get(env params.Interface) (warnings []string, err error) {
	err = c.Docker.get(env)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDockerConfig, err)
	}

	err = c.Log.get(env)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrLogConfig, err)
	}

	warning, err := c.Health.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrHealthConfig, err)
	}

	return warnings, nil
}
