// Copyright 2019 The trust-net Authors
package app

import (
	"github.com/trust-net/dag-lib-go/db"
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack/dto"
	"github.com/trust-net/dag-lib-go/stack/state"
	"testing"
)

func testWorldState() state.State {
	s, _ := state.NewWorldState(db.NewInMemDbProvider(), []byte("test shard"))
	return s
}

// test that tx_handler handles error case when payload in the transaction request is not a valid ID App operation
func TestTxHandler_ErrorCase_InvalidPayload(t *testing.T) {
	log.SetLogLevel(log.NONE)

	// create a transaction with invalid payload
	if err := TxHandler(dto.TestSignedTransaction("invalid operation"), testWorldState()); err == nil {
		t.Errorf("did not fail in invalid transaction")
	}
}

// test that tx_handler error case when op-code in transaction request is invalid
func TestTxHandler_ErrorCase_InvalidOpCode(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create a transaction with invalid payload
	if err := TxHandler(dto.TestSignedTransaction(string(TestOperationPayload(0xffffff,
		TestAttributeRegistration("name", "value")))), testWorldState()); err == nil {
		t.Errorf("did not fail in invalid op-code")
	}
}

// test that tx_handler handles happy path for a valid OpCodeRegisterAttribute in payload
func TestTxHandler_SuccessCase_OpCodeRegisterAttribute(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create a transaction with valid payload
	sub := TestSubmitter()
	if err := TxHandler(sub.PublicSECP256K1Tx(1), testWorldState()); err != nil {
		t.Errorf("OpCodeRegisterAttribute transaction handler failed: %s", err)
	}
}

// test that tx_handler handles happy path for a valid OpCodeEndorseAttribute in payload
func TestTxHandler_SuccessCase_OpCodeEndorseAttribute(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create a transaction with valid payload
	sub := TestSubmitter()
	worldState := testWorldState()
	// register mandatory attribute for identity owner
	if err := TxHandler(sub.PublicSECP256K1Tx(1), worldState); err != nil {
		t.Errorf("OpCodeRegisterAttribute transaction handler failed: %s", err)
	}
	// now submit transaction for endorsement
	if err := TxHandler(sub.PreferredEmailTx(TestSubmitter(), "test email", 1), worldState); err != nil {
		t.Errorf("OpCodeEndorseAttribute transaction handler failed: %s", err)
	}
}
