# id-app-node
Trust-Net Identity App node implementation

## Contents
* [Introduction](#Introduction)
  * [Roles and Responsibilities](#Roles-and-Responsibilities)
  * [Standard Identity Attributes](#Standard-Identity-Attributes)
* [Identity Node APIs](#Identity-Node-APIs)
  * [Identity Transactions API](#Identity-Transactions-API)
  * [Identity Registration accessor API](#Identity-Registration-accessor-API)
  * [Identity Endorsement accessor API](#Identity-Endorsement-accessor-API)
* [Test Driver](#Test-Driver)

## Introduction
This is a proof of concept stand-alone "node" implementation that encapsulates Trust-Net's official [Identity App Component](./id-app-specs.md#id-app-specs), providing a simple self-hosted service on the trust-net network, that can be deployed to self-manage standard identity attributes supported by the official Trust-Net Identity Application.

### Roles and Responsibilities
As described in the official app specs ["Roles and Responsibilities"](./id-app-specs.md#Roles-and-Responsibilities), there are three types of entities involved in this implementation of self-managed identity service...
* **Identity Owner**: An identity owner is the end-user who owns an Identity. All identity attributes are managed by this entity and any access to these identity attributes is controlled by this entity.
* **Identity Partner**: An identity partner is an entity that can submit an endorsement for a specific identity attribute for an identity owner. Such endorsements, however, will need to be approved by the end-user (identity owner) before they are finalized and accepted by the network.
* **Identity Accessor**: An identity accessor is an entity interested in accessing specific attribute of some network identity. Such request will need to be approved by the identity owner before attribute can be accessed by identity accessor.

### Standard Identity Attributes
Refer to the official documentation for ["Standard Identity Attributes"](./id-app-specs.md#Standard-Identity-Attributes) supported by the Trust-Net Identity Application.

## Identity Node APIs
This stand-alone node implementation will provide wrapper API's to interact with the Identity Application component that this node encapsulates, so that this node can be used as a self-managed identity service over the trust-net network.

### Identity Transactions API

The Identity Application node will implement following API endpoint for submitting an identity application transaction:

```
POST /submit

{
    <Transaction DTO, encapsulating the application op-code in payload>
}
```
The specs for transaction submission request for Trust-Net's Identity application are defined at ["Identity application transaction"](./id-app-specs.md#Identity-App-Transactions) specs of official Identity Application. Following transaction op-codes (i.e. types of transactions) are supported by the underlying official Identity Application's transaction handler:
* [Op: Identity attribute registration](./id-app-specs.md#op-identity-attribute-registration)
* [Op: Identity attribute endorsement](./id-app-specs.md#op-identity-attribute-endorsement)

### Identity Registration accessor API
Once an identity attribute is registered, it can be accessed from the Identity Application's world state via an API as following:
```
GET /identity/<public-id>/registrations/<attribute-name>
```
Parameters in the above request are:
* `<public-id>`: [65]byte hex encoded trust-net public id of the identity owner
* `<attribute-name>`: url encoded plain text name of the attribute

**Response schema:** The response for registration of an existing identity attribute for specified trust-net public id would be the same ["Identity attribute registration"](./id-app-specs.md#Identity-attribute-registration) schema defined by the app specs to register the attribute.

> API handler will return error if the `attribute-name` is not one of the standard "self registered" attribute, e.g. "PreferredFirstName"

### Identity Endorsement accessor API
Once an identity attribute is endorsed, it can be accessed from the Identity Application's world state via an API as following:

```
GET /identity/<public-id>/endorsements/<attribute-name>
```
Parameters in the above request are:
* `<public-id>`: [65]byte hex encoded trust-net public id of the identity owner
* `<attribute-name>`: url encoded plain text name of the attribute

**Response schema:** The response for endorsement of an existing identity attribute for specified trust-net public id would be the same ["Identity attribute endorsement"](./id-app-specs.md#Identity-attribute-endorsement) schema defined by the app specs to submit endorsement for the attribute.

> API handler will return error if the `attribute-name` is not one of the standard "3rd party endorsed" attribute, e.g. "PreferredEmail"

## Test Driver
A test driver application is provided that has following CLI commands to test the idnode application functionality:

* print transaction request on standard out, to use with offline tools like postman or curl
* update transaction history of test node submitter
* submit transaction request via API of idnode application directly

Please refer to the [test identity owner driver](./driver_id_owner.md#test-idnode-owner) for details on how to use that application.
