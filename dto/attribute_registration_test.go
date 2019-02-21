package dto

import (
	"encoding/base64"
	"encoding/json"
	"github.com/trust-net/dag-lib-go/log"
	"testing"
)

func Test_AttributeRegistration_ToBase64(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Name:     "test attribute name",
		Value:    "test attribute value",
		Revision: 0x21,
		Proof:    "test proof",
	}
	jsonReg, _ := json.Marshal(reg)
	base64Reg := base64.StdEncoding.EncodeToString(jsonReg)
	if reg.ToBase64() != base64Reg {
		t.Errorf("Incorrect base64 encoding\nActual: %s\nexpected: %s", reg.ToBase64(), base64Reg)
	}
}

func Test_AttributeRegistration_FromBase64_Error_NotBase64(t *testing.T) {
	log.SetLogLevel(log.NONE)
	if _, err := AttributeRegistrationFromBase64("not a bas64 string"); err == nil {
		t.Errorf("Failed to detect a non base64 string")
	}
}

func Test_AttributeRegistration_FromBase64_Error_NotJsonSerialized(t *testing.T) {
	log.SetLogLevel(log.NONE)
	if _, err := AttributeRegistrationFromBase64(base64.StdEncoding.EncodeToString([]byte("not a json encoded structure"))); err == nil {
		t.Errorf("Failed to detect a non json encoded structure")
	}
}

func Test_AttributeRegistration_FromBase64_Error_MissingName(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Value:    "test attribute value",
		Revision: 0x21,
		Proof:    "test proof",
	}
	if _, err := AttributeRegistrationFromBase64(reg.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing name in decoded structure")
	}
}

func Test_AttributeRegistration_FromBase64_Error_MissingValue(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Name:     "test attribute name",
		Revision: 0x21,
		Proof:    "test proof",
	}
	if _, err := AttributeRegistrationFromBase64(reg.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing value in decoded structure")
	}
}

func Test_AttributeRegistration_FromBase64_Error_MissingRevision(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Name:  "test attribute name",
		Value: "test attribute value",
		Proof: "test proof",
	}
	if _, err := AttributeRegistrationFromBase64(reg.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing revision in decoded structure")
	}
}

// make sure that we do not fail for missing proof,
// let the application's transaction handler decide if proof is needed or not
func Test_AttributeRegistration_FromBase64_Success_MissingProof(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Name:     "test attribute name",
		Value:    "test attribute value",
		Revision: 0x21,
	}
	if _, err := AttributeRegistrationFromBase64(reg.ToBase64()); err != nil {
		t.Errorf("Failed upon missing proof in decoded structure")
	}
}

func Test_AttributeRegistration_FromBase64_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	reg := &AttributeRegistration{
		Name:     "test attribute name",
		Value:    "test attribute value",
		Revision: 0x21,
		Proof:    "test proof",
	}
	jsonReg, _ := json.Marshal(reg)
	base64Reg := base64.StdEncoding.EncodeToString(jsonReg)
	if decodedReg, err := AttributeRegistrationFromBase64(base64Reg); err != nil {
		t.Errorf("Failed to decode payload: %s", err)
	} else {
		if decodedReg.Name != reg.Name {
			t.Errorf("decoded attribute name do not match")
		}
		if decodedReg.Value != reg.Value {
			t.Errorf("decoded attribute value do not match")
		}
		if decodedReg.Revision != reg.Revision {
			t.Errorf("decoded attribute revision do not match")
		}
		if decodedReg.Proof != reg.Proof {
			t.Errorf("decoded attribute proof do not match")
		}
	}
}
