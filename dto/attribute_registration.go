// Copyright 2019 The trust-net Authors
// ID Application OpCode DTO for identity attribute registration
package dto

import (
	"fmt"
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
func AttributeRegistrationFromBase64(payload string) (*AttributeRegistration, error) {
	// TBD
	return nil, fmt.Errorf("not implemented")
}

// encode attribute registration to a base64 encoded string
func (a *AttributeRegistration) ToBase64() string {
	// TBD
	return ""
}