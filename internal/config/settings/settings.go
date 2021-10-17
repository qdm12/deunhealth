package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/govalid"
)

type Settings struct {
	Docker Docker
	Health Health
	Log    Log
}

func (s *Settings) SetDefaults() {
	s.Docker.setDefaults()
	s.Health.setDefaults()
	s.Log.setDefaults()
}

func (s *Settings) MergeWith(other Settings) {
	s.Docker.mergeWith(other.Docker)
	s.Health.mergeWith(other.Health)
	s.Log.mergeWith(other.Log)
}

var (
	ErrValidatingDocker = errors.New("error validating docker settings")
	ErrValidatingHealth = errors.New("error validating health settings")
	ErrValidatingLog    = errors.New("error validating log settings")
)

func (s *Settings) Validate(validator govalid.Interface) error {
	if err := s.Docker.validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrValidatingDocker, err)
	}

	if err := s.Health.validate(validator); err != nil {
		return fmt.Errorf("%w: %s", ErrValidatingHealth, err)
	}

	if err := s.Log.validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrValidatingLog, err)
	}

	return nil
}
