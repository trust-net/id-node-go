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
func Test_ErrorCase_InvalidPayload(t *testing.T) {
	log.SetLogLevel(log.NONE)

	// create a transaction with invalid payload
	if err := TxHandler(dto.TestSignedTransaction("invalid operation"), testWorldState()); err == nil {
		t.Errorf("did not fail in invalid transaction")
	}
}

// test that tx_handler error case when op-code in transaction request is invalid
func Test_ErrorCase_InvalidOpCode(t *testing.T) {
	defer log.SetLogLevel(log.NONE)
	// create a transaction with invalid payload
	if err := TxHandler(dto.TestSignedTransaction(string(TestOperationPayload(0xffffff, "test args"))), testWorldState()); err == nil {
		t.Errorf("did not fail in invalid op-code")
	}
}

// test that tx_handler handles happy path for a valid operation in payload
func Test_SuccessCase_HappyPath(t *testing.T) {
	log.SetLogLevel(log.DEBUG)
	defer log.SetLogLevel(log.NONE)
	// create a transaction with invalid payload
	if err := TxHandler(dto.TestSignedTransaction(string(TestOperationPayload(OpCodeRegisterAttribute, "test args"))), testWorldState()); err != nil {
		t.Errorf("transaction handler failed: %s", err)
	}
}
