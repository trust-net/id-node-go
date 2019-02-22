// Copyright 2019 The trust-net Authors
// ID Application OpCode DTO for identity attribute endorsement
package dto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// arguments for attribute endorsement op-code
type AttributeEndorsement struct {
	// name of the attribute being endorsed
	Name string `json:"name"`
	// a base64 encoded [65]byte ECDSA public id/key of the endorsing identity partner
	EndorserId string `json:"endorser_id"`
	// an AES256 ([32]byte) secret key generated by the endorsing identity partner
	// and then encrypted using identity owner's PublicSECP256K1 key
	SecretKey string `json:"enc_secret"`
	// a base64 encoded attribute value as defined by each attribute,
	// encrypted as cipher text using the secret key above
	Value string `json:"enc_value"`
	// 64 bit revision number of the attribute
	Revision uint64 `json:"revision"`
	// a base64 encoded endorsement proof, as defined by each attribute
	Endorsement string `json:"endorsement"`
}

// decode attribute endorsement from a base64 encoded string
// ref: https://github.com/trust-net/id-node-go/issues/1
func AttributeEndorsementFromBase64(payload64 string) (*AttributeEndorsement, error) {
	// decode the base64 string
	if bytes, err := base64.StdEncoding.DecodeString(payload64); err != nil {
		logger.Debug("Failed to decode base64 payload: %s", err)
		return nil, fmt.Errorf("attribute endorsement base64 encoding incorrect")
	} else {
		// decode the json serialized structure
		return AttributeEndorsementFromBytes(bytes)
	}
}

func AttributeEndorsementFromBytes(bytes []byte) (*AttributeEndorsement, error) {
	args := &AttributeEndorsement{}
	// decode the json serialized structure
	if err := json.Unmarshal(bytes, args); err != nil {
		logger.Debug("Failed to json decode payload: %s", err)
		return nil, fmt.Errorf("attribute endorsement json encoding incorrect")
	}

	// validate all fields are present
	if len(args.Name) == 0 {
		logger.Debug("payload missing required field name")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect name")
	} else if len(args.EndorserId) == 0 {
		logger.Debug("payload missing required field endorser_id")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect endorser_id")
	} else if len(args.SecretKey) == 0 {
		logger.Debug("payload missing required field secret_key")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect secret_key")
	} else if len(args.Value) == 0 {
		logger.Debug("payload missing required field value")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect value")
	} else if args.Revision == 0 {
		logger.Debug("payload missing required field revision")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect revision")
	} else if len(args.Endorsement) == 0 {
		logger.Debug("payload missing required field endorsement")
		return nil, fmt.Errorf("attribute endorsement has missing or incorrect endorsement")
	}

	return args, nil
}

// encode attribute endorsement to bytes
func (a *AttributeEndorsement) ToBytes() []byte {
	if bytes, err := json.Marshal(a); err != nil {
		logger.Debug("Failed to json serialize: %s", err)
		return nil
	} else {
		return bytes
	}
}

// encode attribute endorsement to a base64 encoded string
func (a *AttributeEndorsement) ToBase64() string {
	return base64.StdEncoding.EncodeToString(a.ToBytes())
}
