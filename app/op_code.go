// Copyright 2019 The trust-net Authors
// ID Application operation codes
package app

import (
	"encoding/json"
)

// Op Codes for supported operations for identity app
const (
	OpCodeRegisterAttribute uint64 = iota + 0x01
)

// Transaction Operation
type Operation struct {
	// unsigned 64 bit integer specifying operation request type
	OpCode uint64 `json:"op_code"`
	// arguments for the requested operation, encoded as per specs for each specific op_code
	Args string `json:"args"`
}

// decode the json serialized operation from transaction's payload
func DecodeOperation(payload []byte) (*Operation, error) {
	op := &Operation{}
	if err := json.Unmarshal(payload, op); err != nil {
		return nil, err
	} else {
		return op, nil
	}
}
