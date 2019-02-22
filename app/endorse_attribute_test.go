// Copyright 2019 The trust-net Authors
package app

import (
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack/dto"
	"testing"
)

func TestXYZ(t *testing.T) {

}

// test that endorse attribute can handle error in args parsing
func TestEndorseAttribute_ErrorCase_InvalidArgs(t *testing.T) {
	log.SetLogLevel(log.NONE)

	if err := endorseAttribute("not a base64 string", NewIdState(dto.TestSubmitter().Id, testWorldState())); err == nil {
		t.Errorf("Did not fail on non base64 string")
	}
}

// test that endorse attribute checks for correct revision
func TestEndorseAttribute_ErrorCase_InvalidRevision(t *testing.T) {
	log.SetLogLevel(log.NONE)
	sub := TestSubmitter()
	state := NewIdState(sub.Id(), testWorldState())
	// register the mandatory attribute for identity owner first
	registerAttribute(sub.PublicSECP256K1Args(0x01).ToBase64(), state)
	// now submit an endorsement for the email attribute with revision 2, which does not have a prior revision 1
	endorser := TestSubmitter()
	if endorseAttribute(sub.PreferredEmailArgs(endorser, "test email", 0x02).ToBase64(), state) == nil {
		t.Errorf("Did not fail on incorrect revision")
	}
}

// test PreferredEmail endorsement checks for existing mandatory attribute PublicSECP256K1 for identity owner
func TestPreferredEmail_Error_MissingPublicSECP256K1_Owner(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// attempt to endorse PreferredEmail for submitter without having PublicSECP256K1 registered alrady
	endorser := TestSubmitter()
	if endorseAttribute(sub.PreferredEmailArgs(endorser, "test email", 0x01).ToBase64(), state) == nil {
		t.Errorf("Attribute endorsement did not check for owner's mandatory attributes")
	}
}

// test PreferredEmail endorsement checks for valid endorsement proof
func TestPreferredEmail_Error_Invalid_EncorsementProof(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// create test submitter
	sub := TestSubmitter()
	// initialize world state
	state := NewIdState(sub.Id(), testWorldState())
	// register the mandatory attribute for identity owner first
	registerAttribute(sub.PublicSECP256K1Args(0x01).ToBase64(), state)
	// attempt to endorse PreferredEmail for submitter with invalid/tempered endorsement proof
	endorser := TestSubmitter()
	args := sub.PreferredEmailArgs(endorser, "test email", 0x01)
	args.Endorsement = "tempered signature"
	if endorseAttribute(args.ToBase64(), state) == nil {
		t.Errorf("Attribute endorsement did not check for valid endorsement signature")
	}
}

// test endorse attribute happy path
func TestEndorseAttribute_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	sub := TestSubmitter()
	state := NewIdState(sub.Id(), testWorldState())
	// register the mandatory attribute for identity owner first
	if err := registerAttribute(sub.PublicSECP256K1Args(0x01).ToBase64(), state); err != nil {
		t.Errorf("PublicSECP256K1 registration failed: %s", err)
	}
	// now submit an endorsement for the email attribute
	endorser := TestSubmitter()
	if err := endorseAttribute(sub.PreferredEmailArgs(endorser, "test email", 0x01).ToBase64(), state); err != nil {
		t.Errorf("PreferredEmail endorsement failed: %s", err)
	}
	// validate that world state was updated
	if _, err := state.Get("PreferredEmail"); err != nil {
		t.Errorf("Attribute not in world state: %s", err)
	}
}
