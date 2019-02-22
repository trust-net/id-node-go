// Copyright 2019 The trust-net Authors
// ID Application identity attribute endorsement handler
package app

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/trust-net/id-node-go/dto"
	"math/big"
)

// create the standardized endorsement proof verification bytes
// (we explicitly exclude revision, so that same endorsement can be re-used
// for updating new revisions, e.g., when secret key is re-encrypted by
// identity owner using a new version of its PublicSECP256K1 key)
func EndorsementBytes(ownerId []byte, name string, value []byte) []byte {
	bytes := make([]byte, 0, len(ownerId)+len(name)+len(value))
	bytes = append(bytes, ownerId...)
	bytes = append(bytes, []byte(name)...)
	bytes = append(bytes, value...)
	return bytes
}

// handle attribute endorsement operation
func endorseAttribute(argsBase64 string, state *idState) error {
	if attribute, err := dto.AttributeEndorsementFromBase64(argsBase64); err != nil {
		return err
	} else {
		// first check all mandatory attributes of the identity owner are registered
		// before any attributes can be endorsed
		if err := checkMandatoryAttributes(state); err != nil {
			return err
		}

		// validate revision against existing endorsement, if present
		if bytes, err := state.Get(attribute.Name); err == nil {
			if existing, err := dto.AttributeEndorsementFromBytes(bytes); err != nil {
				logger.Error("Failed to de-serialize world state resource: %s", err)
				return err
			} else if attribute.Revision != existing.Revision+1 {
				logger.Debug("Attempt to update with incorrect revision: %d", attribute.Revision)
				return fmt.Errorf("Incorrect update revision")
			}
		} else if attribute.Revision != 1 {
			logger.Debug("Initial registration with incorrect revision: %d", attribute.Revision)
			return fmt.Errorf("Incorrect initial revision")
		}

		// check if attribute is supported
		switch attribute.Name {
		case "PreferredEmail":
			return endorsePreferredEmail(attribute, state)
		default:
			return fmt.Errorf("unsupported standard attribute")
		}
	}
}

func validateEndorsement(attribute *dto.AttributeEndorsement, state *idState) error {
	// decode the base64 encoded endorser's public key from attribute value
	bytes, err := base64.StdEncoding.DecodeString(attribute.EndorserId)
	if err != nil {
		logger.Debug("Failed to decode base64 endoser id: %s", err)
		return fmt.Errorf("malformed base64 endorser id")
	}
	pubKey := crypto.ToECDSAPub(bytes)
	if pubKey == nil {
		logger.Debug("Failed to decode endorser's ECDSA public key")
		return fmt.Errorf("Invalid endroser EC public key")
	}

	// decode the base64 encoded endorsement proof
	bytes, err = base64.StdEncoding.DecodeString(attribute.Endorsement)
	if err != nil {
		logger.Debug("Failed to decode base64 endorsement: %s", err)
		return fmt.Errorf("malformed base64 endorsement")
	}
	// adjust for extra byte, if present
	if len(bytes) == 65 {
		bytes = bytes[1:]
	}
	if len(bytes) != 64 {
		logger.Debug("Incorrect endorsement length: %d", len(bytes))
		return fmt.Errorf("Incorrect endorsement length")
	}

	// regenerate signature parameters
	s := signature{
		R: &big.Int{},
		S: &big.Int{},
	}
	s.R.SetBytes(bytes[0:32])
	s.S.SetBytes(bytes[32:64])

	// decode the base64 encoded cipher text in value
	bytes, err = base64.StdEncoding.DecodeString(attribute.Value)
	if err != nil {
		logger.Debug("Failed to decode base64 cipher text in value: %s", err)
		return fmt.Errorf("malformed base64 value cipher text")
	}

	// create the message to verify signature
	message := EndorsementBytes(state.SubmitterId, attribute.Name, bytes)

	// we want to validate the hash of the message
	hash := sha256.Sum256(message)

	// verify the signature
	if !ecdsa.Verify(pubKey, hash[:], s.R, s.S) {
		return fmt.Errorf("proof validation failed")
	}

	// success
	return nil
}

func endorsePreferredEmail(attribute *dto.AttributeEndorsement, state *idState) error {
	// mandatory atribute validation is already done in caller since common to all endorsements

	// validate endorsement signature
	if err := validateEndorsement(attribute, state); err != nil {
		logger.Error("Endorsement validation failed: %s", err)
		return fmt.Errorf("invalid or incorrect endorsement")
	}

	// success, update world state with attribute update
	if err := state.Put(attribute.Name, attribute.ToBytes()); err != nil {
		logger.Error("Failed to update world state: %s", err)
		return fmt.Errorf("world state update failed")
	}
	return nil
}
