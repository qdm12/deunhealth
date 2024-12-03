package settings

import (
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

func (s *Settings) Validate(validator govalid.Interface) error {
	if err := s.Docker.validate(); err != nil {
		return fmt.Errorf("validating Docker settings: %w", err)
	}

	if err := s.Health.validate(validator); err != nil {
		return fmt.Errorf("validating health settings: %w", err)
	}

	if err := s.Log.validate(); err != nil {
		return fmt.Errorf("validating log settings: %w", err)
	}

	return nil
}
