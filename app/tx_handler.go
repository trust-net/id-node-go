// Copyright 2019 The trust-net Authors
// ID Application transaction handler
package app

import (
	"fmt"
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack/dto"
	"github.com/trust-net/dag-lib-go/stack/state"
)

var (
	AppName  = "trust-net-identity-poc"
	AppShard = []byte(AppName)
	logger   = log.NewLogger("TxHandler")
)

// ID Application's transaction handler callback
func TxHandler(tx dto.Transaction, state state.State) error {
	// deserialize the requested operation in the transaction
	var op *Operation
	var err error
	if op, err = DecodeOperation(tx.Request().Payload); err != nil {
		logger.Debug("Operation decode failed: %s", err)
		return err
	}

	// handle the opcode specific operation
	switch op.OpCode {
	case OpCodeRegisterAttribute:
		err = registerAttribute(op.Args, NewIdState(tx.Request().SubmitterId, state))
	case OpCodeEndorseAttribute:
		err = endorseAttribute(op.Args, NewIdState(tx.Request().SubmitterId, state))
	default:
		err = fmt.Errorf("unsupported op-code: %d", op.OpCode)
	}

	return err
}
