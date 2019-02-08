// Copyright 2019 The trust-net Authors
// ID Application transaction handler
package app

import (
	"fmt"
	"github.com/trust-net/dag-lib-go/stack/dto"
	"github.com/trust-net/dag-lib-go/stack/state"
)

var (
	AppName  = "trust-net-identity-poc"
	AppShard = []byte(AppName)
)

// ID Application's transaction handler callback
func TxHandler(tx dto.Transaction, state state.State) error {
	// deserialize the requested operation in the transaction
	var op *Operation
	var err error
	if op, err = DecodeOperation(tx.Request().Payload); err != nil {
		fmt.Printf("Invalid TX from %x\n%s", tx.Anchor().NodeId, err)
		return err
	}

	// handle the opcode specific operation
	switch op.OpCode {
	case OpCodeRegisterAttribute:
		// TBD
		err = fmt.Errorf("operation not implemented")
	default:
		err = fmt.Errorf("unsupported op-code: %d", op.OpCode)
	}

	return err
}
