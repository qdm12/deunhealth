package config

import (
	"errors"
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

var (
	ErrReadingEnv = errors.New("error reading environment variables")
	ErrValidation = errors.New("error validating settings")
)

func (r *Reader) Read() (s settings.Settings, err error) {
	s, err = r.env.Read()
	if err != nil {
		return s, fmt.Errorf("%w: %s", ErrReadingEnv, err)
	}

	s.SetDefaults()

	err = s.Validate(r.validator)
	if err != nil {
		return s, fmt.Errorf("%w: %s", ErrValidation, err)
	}

	return s, nil
}
