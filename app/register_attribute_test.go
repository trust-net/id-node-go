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

// test PreferredFirstName registration checks for existing mandatory attribute PublicSECP256K1
func TestPreferredFirstName_Error_MissingPublicSECP256K1(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// attempt to register PreferredFirstName for submitter without having PublicSECP256K1 registered alrady
	if registerAttribute(TestAttributeRegistrationCustom("PreferredFirstName", "first_name", 0x01, "").ToBase64(), state) == nil {
		t.Errorf("Attribute registration did not check for existing mandatory params")
	}
}

func TestPreferredFirstName_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// register the PublicSECP256K1 attribute for submitter
	registerAttribute(TestAttributeRegistrationCustom("PublicSECP256K1",
		base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(sub.key.PublicKey.ExportECDSA())), 0x01,
		base64.StdEncoding.EncodeToString(sub.PublicSECP256K1Proof(0x01))).ToBase64(),
		state)
	// attempt to register PreferredFirstName for submitter with PublicSECP256K1 registered alrady
	if err := registerAttribute(TestAttributeRegistrationCustom("PreferredFirstName", "first_name", 0x01, "").ToBase64(), state); err != nil {
		t.Errorf("Attribute registration failed: %s", err)
	}
	// validate that world state was updated
	if _, err := state.Get("PreferredFirstName"); err != nil {
		t.Errorf("Attribute not in world state: %s", err)
	}
}

// test PreferredLastName registration checks for existing mandatory attribute PublicSECP256K1
func TestPreferredLastName_Error_MissingPublicSECP256K1(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// attempt to register PreferredLastName for submitter without having PublicSECP256K1 registered alrady
	if registerAttribute(TestAttributeRegistrationCustom("PreferredLastName", "last_name", 0x01, "").ToBase64(), state) == nil {
		t.Errorf("Attribute registration did not check for existing mandatory params")
	}
}

func TestPreferredLastName_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// register the PublicSECP256K1 attribute for submitter
	registerAttribute(TestAttributeRegistrationCustom("PublicSECP256K1",
		base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(sub.key.PublicKey.ExportECDSA())), 0x01,
		base64.StdEncoding.EncodeToString(sub.PublicSECP256K1Proof(0x01))).ToBase64(),
		state)
	// attempt to register PreferredLastName for submitter with PublicSECP256K1 registered alrady
	if err := registerAttribute(TestAttributeRegistrationCustom("PreferredLastName", "last_name", 0x01, "").ToBase64(), state); err != nil {
		t.Errorf("Attribute registration failed: %s", err)
	}
	// validate that world state was updated
	if _, err := state.Get("PreferredLastName"); err != nil {
		t.Errorf("Attribute not in world state: %s", err)
	}
}
