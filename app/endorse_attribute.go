// Copyright 2019 The trust-net Authors
// ID Application identity attribute endorsement handler
package app

import (
	"fmt"
	"github.com/trust-net/id-node-go/dto"
)

// handle attribute endorsement operation
func endorseAttribute(argsBase64 string, state *idState) error {
	if attribute, err := dto.AttributeEndorsementFromBase64(argsBase64); err != nil {
		return err
	} else {
		// first check all mandatory attributes of the identity owner are registered
		// before any attributes can be endorsed
		if err := checkMandatoryAttributes(state); err != nil {
			return err
		}

		// validate revision against existing endorsement, if present
		if bytes, err := state.Get(attribute.Name); err == nil {
			if existing, err := dto.AttributeEndorsementFromBytes(bytes); err != nil {
				logger.Error("Failed to de-serialize world state resource: %s", err)
				return err
			} else if attribute.Revision != existing.Revision+1 {
				logger.Debug("Attempt to update with incorrect revision: %d", attribute.Revision)
				return fmt.Errorf("Incorrect update revision")
			}
		} else if attribute.Revision != 1 {
			logger.Debug("Initial registration with incorrect revision: %d", attribute.Revision)
			return fmt.Errorf("Incorrect initial revision")
		}

		// check if attribute is supported
		switch attribute.Name {
		case "PreferredEmail":
			return endorsePreferredEmail(attribute, state)
		default:
			return fmt.Errorf("unsupported standard attribute")
		}
	}
}

func endorsePreferredEmail(attribute *dto.AttributeEndorsement, state *idState) error {
	// mandatory validations are already done in caller, since common to all endorsements

	// success, update world state with attribute update
	if err := state.Put(attribute.Name, attribute.ToBytes()); err != nil {
		logger.Error("Failed to update world state: %s", err)
		return fmt.Errorf("world state update failed")
	}
	return nil
}
