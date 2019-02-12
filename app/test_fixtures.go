package app

import (
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

func (s *idSubmitter) PublicSECP256K1Tx(rev uint64) dltdto.Transaction {
	return s.sub.NewTransaction(dltdto.TestAnchor(), string(TestOperationPayload(OpCodeRegisterAttribute,
		TestAttributeRegistrationCustom("PublicSECP256K1",
			base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(s.key.PublicKey.ExportECDSA())), rev,
			base64.StdEncoding.EncodeToString(s.PublicSECP256K1Proof(rev))))))
}

func (s *idSubmitter) PublicSECP256K1Op(rev uint64) *api.SubmitRequest {
	txReq := s.sub.NewRequest(string(TestOperationPayload(OpCodeRegisterAttribute,
		TestAttributeRegistrationCustom("PublicSECP256K1",
			base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(s.key.PublicKey.ExportECDSA())), rev,
			base64.StdEncoding.EncodeToString(s.PublicSECP256K1Proof(rev))))))
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
