// Package config takes care of reading and checking the program configuration
// from environment variables.
package config

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
)

type Config struct {
	Log    Log
	Health Health
}

var (
	ErrLogConfig    = errors.New("cannot obtain log config")
	ErrHealthConfig = errors.New("cannot obtain health config")
)

func (c *Config) get(env params.Env) (warnings []string, err error) {
	err = c.Log.get(env)
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrLogConfig, err)
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
