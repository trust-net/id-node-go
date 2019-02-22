package dto

import (
	"encoding/base64"
	"encoding/json"
	"github.com/trust-net/dag-lib-go/log"
	"testing"
)

func testEndorsement() *AttributeEndorsement {
	return &AttributeEndorsement{
		Name:        "test attribute name",
		EndorserId:  "a test endorser",
		SecretKey:   "a test secret key",
		Value:       "test attribute value",
		Revision:    0x21,
		Endorsement: "test endorsement",
	}
}

func Test_AttributeEndorsement_ToBase64(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	jsonArgs, _ := json.Marshal(args)
	base64Args := base64.StdEncoding.EncodeToString(jsonArgs)
	if args.ToBase64() != base64Args {
		t.Errorf("Incorrect base64 encoding\nActual: %s\nexpected: %s", args.ToBase64(), base64Args)
	}
}

func Test_AttributeEndorsement_FromBase64_Error_NotBase64(t *testing.T) {
	log.SetLogLevel(log.NONE)
	if _, err := AttributeEndorsementFromBase64("not a bas64 string"); err == nil {
		t.Errorf("Failed to detect a non base64 string")
	}
}

func Test_AttributeEndorsement_FromBase64_Error_NotJsonSerialized(t *testing.T) {
	log.SetLogLevel(log.NONE)
	if _, err := AttributeEndorsementFromBase64(base64.StdEncoding.EncodeToString([]byte("not a json encoded structure"))); err == nil {
		t.Errorf("Failed to detect a non json encoded structure")
	}
}

func Test_AttributeEndorsement_FromBase64_Error_MissingName(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	args.Name = ""
	if _, err := AttributeEndorsementFromBase64(args.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing name in decoded structure")
	}
}

func Test_AttributeEndorsement_FromBase64_Error_MissingValue(t *testing.T) {
	args := testEndorsement()
	args.Value = ""
	if _, err := AttributeEndorsementFromBase64(args.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing value in decoded structure")
	}
}

func Test_AttributeEndorsement_FromBase64_Error_MissingRevision(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	args.Revision = 0
	if _, err := AttributeEndorsementFromBase64(args.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing revision in decoded structure")
	}
}

// make sure that we do fail for missing secret key,
// (secret key is a required field for attribute endorsement)
func Test_AttributeEndorsement_FromBase64_Success_MissingSecretKey(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	args.SecretKey = ""
	if _, err := AttributeEndorsementFromBase64(args.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing secret key in decoded structure")
	}
}

// make sure that we do fail for missing endorsement,
// (endorsement is a required field for attribute endorsement)
func Test_AttributeEndorsement_FromBase64_Error_MissingEndorsement(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	args.Endorsement = ""
	if _, err := AttributeEndorsementFromBase64(args.ToBase64()); err == nil {
		t.Errorf("Failed to detect missing endorsement in decoded structure")
	}
}

func Test_AttributeEndorsement_FromBase64_Success_HappyPath(t *testing.T) {
	log.SetLogLevel(log.NONE)
	args := testEndorsement()
	jsonArgs, _ := json.Marshal(args)
	base64Args := base64.StdEncoding.EncodeToString(jsonArgs)
	if decodedArgs, err := AttributeEndorsementFromBase64(base64Args); err != nil {
		t.Errorf("Failed to decode payload: %s", err)
	} else {
		if decodedArgs.Name != args.Name {
			t.Errorf("decoded attribute name do not match")
		}
		if decodedArgs.EndorserId != args.EndorserId {
			t.Errorf("decoded attribute endorser_id do not match")
		}
		if decodedArgs.SecretKey != args.SecretKey {
			t.Errorf("decoded attribute secret_key do not match")
		}
		if decodedArgs.Value != args.Value {
			t.Errorf("decoded attribute value do not match")
		}
		if decodedArgs.Revision != args.Revision {
			t.Errorf("decoded attribute revision do not match")
		}
		if decodedArgs.Endorsement != args.Endorsement {
			t.Errorf("decoded attribute endorsements do not match")
		}
	}
}
