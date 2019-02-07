// Copyright 2019 The trust-net Authors
// ID Application transaction handler
package app

import (
	"fmt"
	stackdto "github.com/trust-net/dag-lib-go/stack/dto"
	"github.com/trust-net/dag-lib-go/stack/state"
	appdto "github.com/trust-net/id-node-go/dto"
)

func txHandler(tx stackdto.Transaction, state state.State) error {
	// deserialize the requested operation in the transaction
	var op *Operation
	var err error
	if op, err = appdto.DecodeOperation(tx.Request().Payload); err != nil {
		fmt.Printf("Invalid TX from %x\n%s", tx.Anchor().NodeId, err)
		return err
	}

	// handle the opcode specific operation
	// TBD

	return nil
}
