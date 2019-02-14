// Copyright 2019 The trust-net Authors
// ID Application identity attribute registration handler
package app

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/trust-net/dag-lib-go/common"
	"github.com/trust-net/id-node-go/dto"
	"math/big"
)

// handle attribute registration operation
func registerAttribute(argsBase64 string, state *idState) error {
	// parse the arguments for attribute registration request
	if attribute, err := dto.AttributeRegistrationFromBase64(argsBase64); err != nil {
		return err
	} else {
		// validate revision against existing attribute, if present
		if bytes, err := state.Get(attribute.Name); err == nil {
			if existing, err := dto.AttributeRegistrationFromBytes(bytes); err != nil {
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
		case "PublicSECP256K1":
			return registerPublicSECP256K1(attribute, state)
		case "PreferredFirstName":
			return registerPreferredFirstName(attribute, state)
		default:
			return fmt.Errorf("unsupported standard attribute")
		}
	}
}

type signature struct {
	R *big.Int
	S *big.Int
}

func registerPublicSECP256K1(opCode *dto.AttributeRegistration, state *idState) error {
	// decode the base64 encoded encryption public key from attribute value
	bytes, err := base64.StdEncoding.DecodeString(opCode.Value)
	if err != nil {
		logger.Debug("Failed to decode base64 public key: %s", err)
		return fmt.Errorf("malformed base64 public key")
	}
	pubKey := crypto.ToECDSAPub(bytes)
	if pubKey == nil {
		logger.Debug("Failed to decode ECDSA public key")
		return fmt.Errorf("Invalid EC public key")
	}

	// decode the base64 encoded signature
	bytes, err = base64.StdEncoding.DecodeString(opCode.Proof)
	if err != nil {
		logger.Debug("Failed to decode base64 proof: %s", err)
		return fmt.Errorf("malformed base64 proof")
	}
	// adjust for extra byte, if present
	if len(bytes) == 65 {
		bytes = bytes[1:]
	}
	if len(bytes) != 64 {
		logger.Debug("Incorrect signature length: %d", len(bytes))
		return fmt.Errorf("Incorrect proof signature")
	}

	// regenerate signature parameters
	s := signature{
		R: &big.Int{},
		S: &big.Int{},
	}
	s.R.SetBytes(bytes[0:32])
	s.S.SetBytes(bytes[32:64])

	// create the message to verify signature
	message := append(state.SubmitterId, common.Uint64ToBytes(opCode.Revision)...)

	// we want to validate the hash of the message
	hash := sha256.Sum256(message)

	// verify the signature
	if !ecdsa.Verify(pubKey, hash[:], s.R, s.S) {
		return fmt.Errorf("proof validation failed")
	}

	// success, update world state with identity attribute
	if err := state.Put(opCode.Name, opCode.ToBytes()); err != nil {
		logger.Error("Failed to update world state: %s", err)
		return fmt.Errorf("world state update failed")
	}
	return nil
}

func checkMandatoryAttributes(state *idState) error {
	// make sure mandatory attribute PublicSECP256K1 is already registered
	if _, err := state.Get("PublicSECP256K1"); err != nil {
		logger.Error("Mandatory parameter PublicSECP256K1 not in world state: %s", err)
		return fmt.Errorf("PublicSECP256K1 not registered")
	}
	// we are all good
	return nil
}

func registerPreferredFirstName(opCode *dto.AttributeRegistration, state *idState) error {
	if err := checkMandatoryAttributes(state); err != nil {
		return err
	}

	// update world state with identity attribute
	if err := state.Put(opCode.Name, opCode.ToBytes()); err != nil {
		logger.Error("Failed to update world state: %s", err)
		return fmt.Errorf("world state update failed")
	}
	return nil
}
