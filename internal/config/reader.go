package config

import (
	"fmt"

	"github.com/qdm12/deunhealth/internal/config/env"
	"github.com/qdm12/deunhealth/internal/config/settings"
	"github.com/qdm12/govalid"
)

type Reader struct {
	env       *env.Reader
	validator govalid.Interface
}

func New() *Reader {
	return &Reader{
		env:       env.New(),
		validator: govalid.New(),
	}
}

func (r *Reader) Read() (s settings.Settings, err error) {
	s, err = r.env.Read()
	if err != nil {
		return s, fmt.Errorf("reading environment settings: %w", err)
	}

	s.SetDefaults()

	err = s.Validate(r.validator)
	if err != nil {
		return s, fmt.Errorf("validating settings: %w", err)
	}

	return s, nil
}
