package config

import (
	"github.com/qdm12/golibs/params"
)

type Health struct {
	Address string
}

func (h *Health) get(env params.Env) (warning string, err error) {
	h.Address, warning, err = h.getAddress(env)
	if err != nil {
		return warning, err
	}
	return warning, nil
}

func (h *Health) getAddress(env params.Env) (address, warning string, err error) {
	const key = "HEALTH_SERVER_ADDRESS"
	options := []params.OptionSetter{
		params.Default("127.0.0.1:9999"),
	}
	return env.ListeningAddress(key, options...)
}
