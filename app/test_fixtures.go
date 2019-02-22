// Copyright 2019 The trust-net Authors
// ID Application test identity owner client
package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/trust-net/dag-lib-go/api"
	"github.com/trust-net/dag-lib-go/common"
	dltdto "github.com/trust-net/dag-lib-go/stack/dto"
	"github.com/trust-net/id-node-go/dto"
	"io"
	"math/big"
)

func TestOperationPayload(opcode uint64, args interface{}) []byte {
	jsonArgs, _ := json.Marshal(args)
	op := Operation{
		OpCode: opcode,
		Args:   base64.StdEncoding.EncodeToString(jsonArgs),
	}
	payload, _ := json.Marshal(op)
	return payload
}

func TestOperationPayloadBase64(opcode uint64, args string) string {
	return base64.StdEncoding.EncodeToString(TestOperationPayload(opcode, args))
}

func TestAttributeRegistration(name, value string) *dto.AttributeRegistration {
	return TestAttributeRegistrationCustom(name, value, 0x01, "test proof")
}

func TestAttributeRegistrationCustom(name, value string, rev uint64, proof string) *dto.AttributeRegistration {
	return &dto.AttributeRegistration{
		Name:     name,
		Value:    value,
		Revision: rev,
		Proof:    proof,
	}
}

func TestAttributeEndorsementCustom(name, endorserId, secretKey, value string, rev uint64, endorsement string) *dto.AttributeEndorsement {
	return &dto.AttributeEndorsement{
		Name:        name,
		EndorserId:  endorserId,
		SecretKey:   secretKey,
		Value:       value,
		Revision:    rev,
		Endorsement: endorsement,
	}
}

type idSubmitter struct {
	sub *dltdto.Submitter
	key *ecies.PrivateKey
}

func TestSubmitter() *idSubmitter {
	key, _ := ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
	sub := dltdto.TestSubmitter()
	sub.ShardId = AppShard
	return &idSubmitter{
		sub: sub,
		key: key,
	}
}

func (s *idSubmitter) Seq() uint64 {
	return s.sub.Seq
}

func (s *idSubmitter) Update(newTx []byte) uint64 {
	copy(s.sub.LastTx[:], newTx)
	s.sub.Seq += 1
	return s.sub.Seq
}

func (s *idSubmitter) Id() []byte {
	return s.sub.Id
}

func (s *idSubmitter) PublicSECP256K1Proof(rev uint64) []byte {
	// create the message to sign
	message := append(s.sub.Id, common.Uint64ToBytes(rev)...)
	// we want to sign the hash of the message
	hash := sha256.Sum256(message)

	// sign using the ECIS private key
	sig := signature{}
	sig.R, sig.S, _ = ecdsa.Sign(rand.Reader, s.key.ExportECDSA(), hash[:])
	return append(sig.R.Bytes(), sig.S.Bytes()...)
}

func (s *idSubmitter) PublicSECP256K1Args(rev uint64) *dto.AttributeRegistration {
	return TestAttributeRegistrationCustom("PublicSECP256K1",
		base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(s.key.PublicKey.ExportECDSA())), rev,
		base64.StdEncoding.EncodeToString(s.PublicSECP256K1Proof(rev)))
}

func (s *idSubmitter) PublicSECP256K1Payload(rev uint64) string {
	return string(TestOperationPayload(OpCodeRegisterAttribute, s.PublicSECP256K1Args(rev)))
}

func (s *idSubmitter) PublicSECP256K1Tx(rev uint64) dltdto.Transaction {
	return s.sub.NewTransaction(dltdto.TestAnchor(), s.PublicSECP256K1Payload(rev))
}

func (s *idSubmitter) SubmitRequest(txReq *dltdto.TxRequest) *api.SubmitRequest {
	return &api.SubmitRequest{
		// payload for transaction's operations
		Payload: base64.StdEncoding.EncodeToString(txReq.Payload),
		// shard id for the transaction
		ShardId: hex.EncodeToString(txReq.ShardId),
		// submitter's last transaction
		LastTx: hex.EncodeToString(txReq.LastTx[:]),
		// Submitter's public ID
		SubmitterId: hex.EncodeToString(txReq.SubmitterId),
		// submitter's transaction sequence
		SubmitterSeq: txReq.SubmitterSeq,
		// a padding to meet challenge for network's DoS protection
		Padding: txReq.Padding,
		// signature of the transaction request's contents using submitter's private key
		Signature: base64.StdEncoding.EncodeToString(txReq.Signature),
	}
}

func (s *idSubmitter) PublicSECP256K1Op(rev uint64) *api.SubmitRequest {
	txReq := s.sub.NewRequest(s.PublicSECP256K1Payload(rev))
	return s.SubmitRequest(txReq)
}

func (s *idSubmitter) PreferredFirstNamePayload(name string, rev uint64) string {
	return string(TestOperationPayload(OpCodeRegisterAttribute,
		TestAttributeRegistrationCustom("PreferredFirstName", name, rev, "")))
}

func (s *idSubmitter) PreferredFirstNameTx(name string, rev uint64) dltdto.Transaction {
	return s.sub.NewTransaction(dltdto.TestAnchor(), s.PreferredFirstNamePayload(name, rev))
}

func (s *idSubmitter) PreferredFirstNameOp(name string, rev uint64) *api.SubmitRequest {
	txReq := s.sub.NewRequest(s.PreferredFirstNamePayload(name, rev))
	return s.SubmitRequest(txReq)
}

func (s *idSubmitter) PreferredLastNamePayload(name string, rev uint64) string {
	return string(TestOperationPayload(OpCodeRegisterAttribute,
		TestAttributeRegistrationCustom("PreferredLastName", name, rev, "")))
}

func (s *idSubmitter) PreferredLastNameOp(name string, rev uint64) *api.SubmitRequest {
	txReq := s.sub.NewRequest(s.PreferredLastNamePayload(name, rev))
	return s.SubmitRequest(txReq)
}

func (s *idSubmitter) SignSha256(bytes []byte) []byte {
	// sign the request using SHA256 digest and ECDSA private key
	type signature struct {
		R *big.Int
		S *big.Int
	}
	sig := signature{}
	// sign the request
	hash := sha256.Sum256(bytes)
	sig.R, sig.S, _ = ecdsa.Sign(rand.Reader, s.sub.Key, hash[:])
	return append(sig.R.Bytes(), sig.S.Bytes()...)
}

// symmetric encryption using an AES-256 32 byte secret key
func EncryptAES256(plaintext []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

func (s *idSubmitter) Endorse(ownerId []byte, ownerKey *ecies.PublicKey, name string, value string) (cipherText, encSecret, signature []byte, err error) {
	// encrypt the value using a secret key
	secret, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	cipherText, err = EncryptAES256([]byte(value), secret)
	if err != nil {
		return
	}

	// encrypt the secret key using owner's public key
	encSecret,_ = ecies.Encrypt(rand.Reader, ownerKey, secret, nil, nil)

	// return the encrypted secret key, cipher text of value, and the endorsement signature
	signature = s.SignSha256(EndorsementBytes(ownerId, name, cipherText))
	return
}

func (s *idSubmitter) PreferredEmailArgs(endorser *idSubmitter, email string, rev uint64) *dto.AttributeEndorsement {
	cipherText, encSecret, signature, _ := endorser.Endorse(s.Id(), &s.key.PublicKey, "PreferredEmail", email)
	return TestAttributeEndorsementCustom("PreferredEmail",
		base64.StdEncoding.EncodeToString(endorser.Id()),
		base64.StdEncoding.EncodeToString(encSecret),
		base64.StdEncoding.EncodeToString(cipherText), rev,
		base64.StdEncoding.EncodeToString(signature))
}

func (s *idSubmitter) PreferredEmailPayload(endorser *idSubmitter, email string, rev uint64) string {
	return string(TestOperationPayload(OpCodeEndorseAttribute, s.PreferredEmailArgs(endorser, email, rev)))
}

func (s *idSubmitter) PreferredEmailOp(endorser *idSubmitter, email string, rev uint64) *api.SubmitRequest {
	txReq := s.sub.NewRequest(s.PreferredEmailPayload(endorser, email, rev))
	return s.SubmitRequest(txReq)
}

func (s *idSubmitter) PreferredEmailTx(endorser *idSubmitter, email string, rev uint64) dltdto.Transaction {
	return s.sub.NewTransaction(dltdto.TestAnchor(), s.PreferredEmailPayload(endorser, email, rev))
}
