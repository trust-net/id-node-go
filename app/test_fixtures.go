package app

import (
	"encoding/base64"
	"encoding/json"
)

func TestOperationPayload(opcode uint64, args string) []byte {
	op := Operation{
		OpCode: opcode,
		Args:   args,
	}
	payload, _ := json.Marshal(op)
	return payload
}

func TestOperationPayloadBase64(opcode uint64, args string) string {
	return base64.StdEncoding.EncodeToString(TestOperationPayload(opcode, args))
}
