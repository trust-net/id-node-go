# id-app-specs
Trust-Net Identity Application Specifications

## Contents
* [Introduction](#Introduction)
* [Roles and Responsibilities](#Roles-and-Responsibilities)
    * [Identity Owner](#Identity-Owner)
    * [Identity Partner](#Identity-Partner)
    * [Identity Accessor](#Identity-Accessor)
* [Standard Identity Attributes](#Standard-Identity-Attributes)
    * [Proof of ownership](#Proof-of-ownership)
    * [3rd party certification](#3rd-party-certification)
    * [Zero Knowledge Endorsement Proof](#Zero-Knowledge-Endorsement-Proof)
* [Identity App Transactions](#Identity-App-Transactions)
    * [Payload schema](#Payload-schema)
    * [Op: Identity attribute registration](#Op-Identity-attribute-registration)
    * [Op: Identity attribute endorsement](#Op-Identity-attribute-endorsement)
* [Standard Attributes Payload Schema](#Standard-Attributes-Payload-Schema)
    * [PublicSECP256K1 registration payload](#PublicSECP256K1-registration-payload)
    * [PublicSECP256K1 registration proof](#PublicSECP256K1-registration-proof)
    * [PreferredFirstName registration payload](#PreferredFirstName-registration-payload)
    * [PreferredFirstName proof of ownership](#PreferredFirstName-proof-of-ownership)
    * [PreferredLastName registration payload](#PreferredLastName-registration-payload)
    * [PreferredLastName proof of ownership](#PreferredLastName-proof-of-ownership)
    * [PreferredEmail endorsement payload](#PreferredEmail-endorsement-payload)
    * [PreferredEmail proof of ownership](#PreferredEmail-proof-of-ownership)

* [Test Driver](#Test-Driver)
    * [Test identity owner](#Test-identity-owner)

## Introduction
This document defines specifications for the official trust-net Identity Application component that can be used either as a [stand-alone Identity App Node](./id-app-node.md#id-app-node), or can be combined with other trust-net app components to compose complex enterprise applications that rely on global identities managed with this app component.

## Roles and Responsibilities
There are three types of entities involved in the official trust-net Identity Application...

### Identity Owner
An identity owner is the end-user who owns an Identity. All identity attributes are managed by this entity and any access to these identity attributes is controlled by this entity.

### Identity Partner
An identity partner is an entity that can submit an endorsement for a specific identity attribute for an identity owner. Such endorsements, however, will need to be approved by the end-user (identity owner) before they are finalized and accepted by the network.

### Identity Accessor
An identity accessor is an entity interested in accessing specific attribute of some network identity. Such request will need to be approved by the identity owner before attribute can be accessed by identity accessor.

## Standard Identity Attributes
Identity App will define certain "standard" attributes that all Trust-Net applications can work with. Each of these standard attributes will have well defined specification and convention. Following are the currently supported standard attributes by application:

| Attribute | Purpose | Attribute Scope | Certification | Comments |
|----------:|---------|----------------|---------------|----------|
| **`PublicSECP256K1`** | registers an ECIES public key over `secp256k1` curve | Global | self registered | This is a mandatory attribute, before any other attribute can be registered, this one should be registered|
| **`PreferredFirstName`** | registers a self declared first name for the submitter | Personal | self registered | (optional) |
| **`PreferredLastName`** | registers a self declared last name for the submitter | Personal | self registered | (optional) |
| **`PreferredEmail`** | registers a self declared email address for the submitter | Global | 3rd-party endorsed | (optional) |

### Proof of ownership
When registering an attribute, if the attribute is meant to be a globally unique attribute (e.g. a public key) then the registration process should ensure the following:
* registration transaction should have proof that submitter really owns the registered attribute value
* it should not be possible to "replay" a registration request proof from another submitter

### 3rd party certification
While its possible to include self certified proof for things like public key, it may not be possible to include self certified proof for things like email address. In such cases, a 3rd party endorsement for that attribute would be necessary.

### Zero Knowledge Endorsement Proof
Not having an endorsement proof that can be independently verified by each Identity Application node's transaction handler is very weak!!! A strong endorsement scheme is possible **_if the endorser encrypts the attribute and then signs over the encrypted cipher text of the attribute value_** that it is endorsing. This will ensure that an identity owner can only submit an endorsement transaction for the attribute value that was endorsed by the identity partner, without requiring disclosure of the attribute value to public network. Hence, for 3rd party endorsements, following protocol will be used:
* Identity owner submits the attribute value for endorsement to identity partner in plain text
* Identity partner completes "off-the-chain" verification of attribute value ownership
* Identity partner gets the `PublicSECP256K1` for identity owner from trust-net's Identity Service
* Identity partner generates an AES256 (32 byte) secret key and encrypts the plain text attribute value
* Identity partner encrypts the secret key using identity owner's `PublicSECP256K1` key
* Identity partner signs the cipher text of encrypted attribute value as endorsement
* Identity partner returns the encrypted secret key, encrypted attribute value and signed endorsement back
* Identity owner creates a new endorsement transaction using returned values above

## Identity App Transactions

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
> Each standard attribute with certification type "self registered" will define the attribute specific rules/semantics for contents of value and proof fields, as applicable.

### Op: Identity attribute endorsement
The transaction to submit endorsement for an identity attribute will have following payload:

**op-code**:
`0x02`

**args**:
op-code will have its `args` field content set to a base64 string of json serialized structure of below schema type:
<a id="Identity-attribute-endorsement"></a>
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Identity attribute endorsement",
  "description": "A Trust-Net identity attribute endorsement request",
  "type": "object",
  "properties": {
    "name": {
      "description": "name of the attribute being endorsed",
      "type": "string"
    },
    "endorser_id": {
      "description": "a base64 encoded [65]byte ECDSA public id/key of the endorsing identity partner",
      "type": "string"
    },
    "enc_secret": {
      "description": "a base64 encoded AES256 secret key, encrypted using identity owner's PublicSECP256K1 key",
      "type": "string"
    },
    "enc_value": {
      "description": "a base64 encoded attribute value as defined by each attribute encrypted as cipher text using the secret key above",
      "type": "string"
    },
    "revision": {
      "description": "unsigned 64 bit revision number of the attribute",
      "type": "integer"
    },
    "endorsement": {
      "description": "a base64 encoded endorsement proof, as defined by each attribute",
      "type": "string"
    }
  },
  "required": [ "name", "endorser_id", "enc_secret", "enc_value", "revision", "endorsement" ]
}
```
> Each standard attribute with certification type "3rd party endorsed" will define the attribute specific rules/semantics for contents of enc_value and endorsement fields, as applicable.

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

### PreferredFirstName registration payload
The payload for identity's `PreferredFirstName` attribute registration would consist of following:

  | Contents | Encoding | Semantic
--:| -- | -- | --
**`name`** | "PreferredFirstName" | plain text | a personal attribute registering the self declared first name of the identity owner
**`value`** | string | plain text| first name in plain text string
**`revision`** | 64 bit revision number | plain number | revision number for the attribute update
**`proof`** | `null` | n/a | no additional proof required beyond the transaction ownership of the request submitter

### PreferredFirstName proof of ownership
This is a self-declared attribute with personal scope and hence does not require any proof. A transaction submitter's signature on the request is sufficient to validate that submitter identity has requested the attribute registration to their desired value.

> Note: this is not official name -- that would require some kind of 3rd party certification, from an identity partner in the Trust-Net system.

### PreferredLastName registration payload
The payload for identity's `PreferredLastName` attribute registration would consist of following:

  | Contents | Encoding | Semantic
--:| -- | -- | --
**`name`** | "PreferredLastName" | plain text | a personal attribute registering the self declared last name of the identity owner
**`value`** | string | plain text| last name in plain text string
**`revision`** | 64 bit revision number | plain number | revision number for the attribute update
**`proof`** | `null` | n/a | no additional proof required beyond the transaction ownership of the request submitter

### PreferredLastName proof of ownership
This is a self-declared attribute with personal scope and hence does not require any proof. A transaction submitter's signature on the request is sufficient to validate that submitter identity has requested the attribute registration to their desired value.

> Note: this is not official name -- that would require some kind of 3rd party certification, from an identity partner in the Trust-Net system.

### PreferredEmail endorsement payload
The payload for identity's `PreferredEmail` attribute endorsement would consist of following:

| |Contents | Encoding | Semantic
--:|-- | -- | --
**name** | "PreferredEmail" | plain text | a name indicating this is an endorsement for the preferred email of the identity owner
**endorser_id** | [65]byte | base64 | ECDSA public id/key of the endorsing identity partner
**enc_secret** | []byte | base64 | an AES256 ([32]byte) secret key **generated by the endorsing identity partner** and then encrypted using identity owner's PublicSECP256K1 key
**enc_value** | []byte | base64 | email address for the identity owner as cipher text, **AES256 encrypted by the endorsing identity partner** using the secret key above
**revision** | 64 bit revision number | plain number | revision number for the attribute update
**endorsement** | [64]byte | base64 | ECDSA secpk256 signature using endorsing identity partner's private key over SHA256 hash of ([65]byte identity owner's public id + "PreferredEmail" + ~[8]byte revision number~ + cipher text email address value)

### PreferredEmail proof of ownership
The endorsement proof above guarantees that identity owner can only submit an endorsement that was actually endorsed and signed by the identity partner. Identity owner will make sure to decrypt the secret and cipher text to make sure attribute values are as it expected, before submitting the transaction to the network. We are not including revision number in the endorsement proof so that same endorsement can be used to update the revision, e.g., when secret_key is re-encrypted after a new PublicSECP256K1 key revision.
