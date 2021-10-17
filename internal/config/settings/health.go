package settings

import (
	"os"

	"github.com/qdm12/govalid"
	"github.com/qdm12/govalid/address"
)

type Health struct {
	Address string
}

func (h *Health) setDefaults() {
	if h.Address == "" {
		h.Address = "127.0.0.1:9999"
	}
}

func (h *Health) mergeWith(other Health) {
	if h.Address == "" {
		h.Address = other.Address
	}
}

func (h *Health) validate(validator govalid.Interface) (err error) {
	uid := os.Getuid()
	h.Address, err = validator.ValidateAddress(h.Address, address.OptionListening(uid))
	return err
}
