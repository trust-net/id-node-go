package app

import (
	//	"fmt"
	"github.com/trust-net/dag-lib-go/stack/state"
	"github.com/trust-net/id-node-go/dto"
)

// handle attribute registration operation
func registerAttribute(argsBase64 string, submitterId []byte, state state.State) error {
	// parse the arguments for attribute registration request
	if _, err := dto.AttributeRegistrationFromBase64(argsBase64); err != nil {
		return err
	} else {
		// TBD
	}

	return nil
}
