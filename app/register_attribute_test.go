// Copyright 2019 The trust-net Authors
package app

import (
	"encoding/base64"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack/dto"
	"testing"
)

// test that register attribute can handle error in args parsing
func TestRegisterAttribute_ErrorCase_InvalidArgs(t *testing.T) {
	log.SetLogLevel(log.NONE)

	if err := registerAttribute("not a base64 string", NewIdState(dto.TestSubmitter().Id, testWorldState())); err == nil {
		t.Errorf("Did not fail on non base64 string")
	}
}

// test that register attribute checks for correct revision
func TestRegisterAttribute_ErrorCase_InvalidRevision(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// try registering an attribute for revision 2, which does not have a prior revision 1
	if err := registerAttribute(TestAttributeRegistrationCustom("PublicSECP256K1", "value", 0x02, "proof").ToBase64(),
		NewIdState(dto.TestSubmitter().Id, testWorldState())); err == nil {
		t.Errorf("Did not fail on incorrect revision")
	}
}

// test register attribute happy path
func TestRegisterAttribute_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create a valid registration request
	sub := TestSubmitter()
	state := NewIdState(sub.Id(), testWorldState())
	if err := registerAttribute(TestAttributeRegistrationCustom("PublicSECP256K1",
		base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(sub.key.PublicKey.ExportECDSA())), 0x01,
		base64.StdEncoding.EncodeToString(sub.PublicSECP256K1Proof(0x01))).ToBase64(),
		state); err != nil {
		t.Errorf("Attribute registration failed: %s", err)
	}
	// validate that world state was updated
	if _, err := state.Get("PublicSECP256K1"); err != nil {
		t.Errorf("Attribute not in world state: %s", err)
	}
}
