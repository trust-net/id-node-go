// Copyright 2019 The trust-net Authors
// ID Application OpCode DTO for identity attribute registration
package dto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/trust-net/dag-lib-go/log"
)

var (
	logger = log.NewLogger("AttributeRegistration")
)

// A request to submit a transaction
type AttributeRegistration struct {
	// name of the attribute being registered
	Name string `json:"name"`
	// a base64 encoded value, as defined by each attribute
	Value string `json:"value"`
	// 64 bit revision number of the attribute
	Revision uint64 `json:"revision"`
	// a base64 encoded proof of ownership, as defined by each attribute
	Proof string `json:"proof"`
}

// decode attribute registration from a base64 encoded string
// ref: https://github.com/trust-net/id-node-go/issues/1
func AttributeRegistrationFromBase64(payload64 string) (*AttributeRegistration, error) {
	// decode the base64 string
	if bytes, err := base64.StdEncoding.DecodeString(payload64); err != nil {
		logger.Debug("Failed to decode base64 payload: %s", err)
		return nil, fmt.Errorf("attribute registration base64 encoding incorrect")
	} else {
		// decode the json serialized structure
		return AttributeRegistrationFromBytes(bytes)
	}
}

func AttributeRegistrationFromBytes(bytes []byte) (*AttributeRegistration, error) {
	reg := &AttributeRegistration{}
	// decode the json serialized structure
	if err := json.Unmarshal(bytes, reg); err != nil {
		logger.Debug("Failed to json decode payload: %s", err)
		return nil, fmt.Errorf("attribute registration json encoding incorrect")
	}

	// validate all fields are present
	if len(reg.Name) == 0 {
		logger.Debug("payload missing required field name")
		return nil, fmt.Errorf("attribute registration missing or incorrect name")
	} else if len(reg.Value) == 0 {
		logger.Debug("payload missing required field value")
		return nil, fmt.Errorf("attribute registration missing or incorrect value")
	} else if reg.Revision == 0 {
		logger.Debug("payload missing required field revision")
		return nil, fmt.Errorf("attribute registration missing or incorrect revision")
	}

	return reg, nil
}

// encode attribute registration to a bytes
func (a *AttributeRegistration) ToBytes() []byte {
	if bytes, err := json.Marshal(a); err != nil {
		logger.Debug("Failed to json serialize: %s", err)
		return nil
	} else {
		return bytes
	}
}

// encode attribute registration to a base64 encoded string
func (a *AttributeRegistration) ToBase64() string {
	return base64.StdEncoding.EncodeToString(a.ToBytes())
}
