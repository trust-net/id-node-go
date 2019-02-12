# id-node-go
Identity node implementation

## Contents
* [Introduction](#Introduction)
* [Roles and Responsibilities](#Roles-and-Responsibilities)
    * [Identity Owner](#Identity-Owner)
    * [Identity Partner](#Identity-Partner)
    * [Identity Accessor](#Identity-Accessor)
* [Standard Identity Attributes](#Standard-Identity-Attributes)
    * [Proof of ownership](#Proof-of-ownership)
    * [3rd party certification](#3rd-party-certification)
* [Identity Transactions API](#Identity-Transactions-API)
    * [API endpoint](#API-endpoint)
    * [Payload schema](#Payload-schema)
    * [Op: Identity attribute registration](#Op-Identity-attribute-registration)
* [Identity Access API](#Identity-Access-API)
    * [API endpoint](#API-endpoint-2)
    * [Response schema](#Response-schema)
* [Standard Attributes Payload Schema](#Standard-Attributes-Payload-Schema)
    * [PublicSECP256K1 registration payload](#PublicSECP256K1-registration-payload)
    * [PublicSECP256K1 registration proof](#PublicSECP256K1-registration-proof)

## Introduction
A proof of concept for self-managed identity service over Trust-Net.

## Roles and Responsibilities
There are three types of entities involved in this implementation of self-managed identity service...

### Identity Owner
An identity owner is the end-user who owns an Identity. All identity attributes are managed by this entity and any access to these identity attributes is controlled by this entity.

### Identity Partner
An identity partner is an entity that can submit an endorsement for a specific identity attribute for an identity owner. Such endorsements, however, will need to be approved by the end-user (identity owner) before they are finalized and accepted by the network.

### Identity Accessor
An identity accessor is an entity interested in accessing specific attribute of some network identity. Such request will need to be approved by the identity owner before attribute can be accessed by identity accessor.

## Standard Identity Attributes
Identity service will define certain "standard" attributes that all Trust-Net applications can work with. Each of these standard attributes will have well defined specification and convention. Following are the currently supported standard attributes by application:

| Attribute | Purpose | Attribute Scope | Certification | Comments |
|----------:|---------|----------------|---------------|----------|
| **`PublicSECP256K1`** | registers an ECIES public key over `secp256k1` curve | Global | self | This is a mandatory attribute, before any other attribute can be registered, this one should be registered|
| **`PreferredFirstName`** | registers a self declared first name for the submitter | Personal | self | (optional) |
| **`PreferredLastName`** | registers a self declared last name for the submitter | Personal | self | (optional) |
| **`PreferredEmail`** | registers a self declared email address for the submitter | Global | 3rd-party | (optional) |

### Proof of ownership
When registering an attribute, if the attribute is meant to be a globally unique attribute (e.g. a public key) then the registration process should ensure the following:
* registration transaction should have proof that submitter really owns the registered attribute value
* it should not be possible to "replay" a registration request proof from another submitter

### 3rd party certification
While its possible to include self certified proof for things like public key, it may not be possible to include self certified proof for things like email address. In such cases, a 3rd party certification for that attribute endorsement would be necessary.

## Identity Transactions API

### API endpoint
The ID Application will implement following API endpoint for submitting an identity application transaction:

```
POST /submit

{
    <Transaction DTO, encapsulating the application op-code in payload>
}
```
A transaction submission request for Trust-Net's Identity application would follow the same transaction submission spec as defined in [Op: Transaction Submission](https://github.com/trust-net/dag-lib-go/blob/master/docs/SpendrApp.md#op-submit-transaction). The payload for each transaction would consist of different requested Identity application specific operations.

Above request will be validated for each specific transaction op-code as defined in later sections. For success cases, the transaction ID will be returned back with `HTTP 201` response.

### Payload schema
The payload field for transaction submission request would be a base64 string of json serialized structure of below schema type:
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Identity Operation",
  "description": "A Trust-Net identity application operation request",
  "type": "object",
  "properties": {
    "op_code": {
      "description": "unsigned 64 bit integer specifying operation request type",
      "type": "integer"
    },
    "args": {
      "description": "arguments for the requested operation, encoded as per specs for each specific op_code",
      "type": "string"
    }
  },
  "required": [ "op_code", "args" ]
}
```
Above paylaod structure will provide following:
* `op_code`: this field will indicate specific Identity service operation requested
* `args`: this field will provide arguments required for the requested operation

A list of supported operations for Trust-Net's Identity application is provided in below sections.

### Op: Identity attribute registration
The transaction to register an identity attribute will have following payload:

**op-code**:
`0x01`

**args**:
op-code will have its `args` field content set to a base64 string of json serialized structure of below schema type:
<a id="Identity-attribute-registration"></a>
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Identity attribute registration",
  "description": "A Trust-Net identity attribute registration request",
  "type": "object",
  "properties": {
    "name": {
      "description": "name of the attribute being registered",
      "type": "string"
    },
    "value": {
      "description": "a base64 encoded value, as defined by each attribute",
      "type": "string"
    },
    "revision": {
      "description": "unsigned 64 bit revision number of the attribute",
      "type": "integer"
    },
    "proof": {
      "description": "a base64 encoded proof of ownership, as defined by each attribute",
      "type": "string"
    }
  },
  "required": [ "name", "value", "revision", "proof" ]
}
```
> Each standard attribute will define the attribute specific rules/semantics for contents of value and proof fields, as applicable.

## Identity Access API
Once an identity attribute is registered, it can be accessed from the Identity Application's world state via an API as following:

### API endpoint
```
GET /identity/<public-id>/attributes/<attribute-name>
```
Parameters in the above request are:
* `<public-id>`: [65]byte hex encoded trust-net public id of the identity owner
* `<attribute-name>`: url encoded plain text name of the attribute

### Response schema
The response for an existing identity attribute for specified trust-net public id would be the same **Identity attribute registration** schema [defined above](#Identity-attribute-registration) and used to register the attribute.

> Note: the attribute value might be encrypted, and hence requestor may have to request/fetch an additional access grant from the identity owner which will provide decryption key. That flow would be covered in a separate ticket.

## Standard Attributes Payload Schema
Following are the op-code payload schema and semantics for supported standard attributes...

### PublicSECP256K1 registration payload

The payload for identity's `PublicSECP256K1` attribute registration would consist of following:

|  | Contents | Encoding | Semantic |
|------:|:------|:----------|-----------|
|**name**|  "PublicSECP256K1" | plain text | a global attribute registering the public key to send encrypted content to identity owner |
|**value**| `[65]byte` | base64 | ECIES public key over secp256k1 curve |
|**revision**| 64 bit revision number | plain number | revision number for the attribute update|
|**proof**|`[64]byte`| base64 | ~ECIES encrypted cipher text~ ECDSA secpk256 signature using corresponding private key over SHA256 digest of [65]byte public ID of submitter + [8]byte revision number|

### PublicSECP256K1 registration proof
Application node's transaction handler will validated the "PublicSECP256K1" attribute registration request  by verifying the ECDSA secpk256 signature in `proof` with the registered public key over the SHA256 digest of 65 bytes of transaction submitter's public ID and the 8 bytes of revision number as following:

```
// decode the base64 encoded encryption public key from attribute value
bytes, _ := base64.StdEncoding.DecodeString(opCode.Value)
pubKey := crypto.ToECDSAPub(bytes)

// decode the base64 encoded signature
proof, _ := base64.StdEncoding.DecodeString(opCode.Proof)
s := signature{
  R: &big.Int{},
  S: &big.Int{},
}
s.R.SetBytes(proof[0:32])
s.S.SetBytes(proof[32:64])

// create the message to verify signature
message := append(tx.Payload().SubmitterId, common.Uint64ToBytes(opCode.Revision)...)

// we want to validate the hash of the message
hash := sha256.Sum256(message)

// verify the signature
if !ecdsa.Verify(pubKey, hash[:], s.R, s.S) {
    return fmt.Errorf("proof validation failed")
}
```
