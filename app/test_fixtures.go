package app

import (
	"encoding/base64"
	"encoding/json"
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
