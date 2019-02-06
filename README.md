# id-node-go
Identity node implementation

## Contents
* [Introduction](#Introduction)
* [Roles and Responsibilities](#Roles-and-Responsibilities)
    * [Identity Owner](#Identity-Owner)
    * [Identity Partner](#Identity-Partner)
    * [Identity Accessor](#Identity-Accessor)

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
