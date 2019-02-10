package app

import (
	"github.com/trust-net/dag-lib-go/log"
	"testing"
)

// test that register attribute can handle error in args parsing
func TestRegisterAttribute_ErrorCase_InvalidArgs(t *testing.T) {
	log.SetLogLevel(log.NONE)

	if err := registerAttribute("not a base64 string", nil, testWorldState()); err == nil {
		t.Errorf("Did not fail on non base64 string")
	}
}

// test that register attribute checks for correct revision
func TestRegisterAttribute_ErrorCase_InvalidRevision(t *testing.T) {
	log.SetLogLevel(log.NONE)
	// try registering an attribute for revision 2, which does not have a prior revision 1
	if err := registerAttribute(TestAttributeRegistrationCustom("name", "value", 0x02, "proof").ToBase64(), nil, testWorldState()); err == nil {
		t.Errorf("Did not fail on incorrect revision")
	}
}
