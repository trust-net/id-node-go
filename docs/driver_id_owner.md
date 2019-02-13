## test-idnode-owner
* [Introduction](#Introduction)
    * [Usage](#Usage)
    * [CLI commands](#CLI-commands)
* [Offline Transactions](#Offline-Transactions)
    * [Print PublicSECP256K1 registration](#Print-PublicSECP256K1-registration)
    * [Update transaction history](#Update-transaction-history)
* [Online Transactions](#Online-Transactions)
    * [Set idnode app base url](#Set-idnode-app-base-url)
    * [Ping idnode app](#Ping-idnode-app)
    * [Submit PublicSECP256K1 registration](#Submit-PublicSECP256K1-registration)

## Introduction
A test driver application is provided that has CLI commands for following:

* print transaction request on standard out, to use with offline tools like postman or curl
* update transaction history of test node submitter as result of offline submission
* submit transaction request via API of idnode application directly

### Usage
Test driver application is in the `github.com/trust-net/id-node-go/test/owner/` folder. It can be invoked as following:

```
go run $GOPATH/src/github.com/trust-net/id-node-go/test/owner/main.go
OWNER[01]: 
```

### CLI commands
Following CLI commands are provided:

```
OWNER[01]: help
Accepted commands...
"ping", "print_key", "submit_key", "update", "url"
```

## Offline Transactions
Following commands are available for offline identity attribute registration using the test driver...

### Print PublicSECP256K1 registration
Following CLI command is implemented to print transaction request for ID Node's API to register the PublicSECP256K1 attribute of a test identity:

```
OWNER[01]: help print_key
print transaction request for registering PublicSECP256K1 attribute with revision (default revision 1)
usage: print_key [<revision>]

OWNER[01]: print_key 1
{
  "payload": "eyJvcF9jb2RlIjoxLCJhcmdzIjoiZXlKdVlXMWxJam9pVUhWaWJHbGpVMFZEVURJMU5rc3hJaXdpZG1Gc2RXVWlPaUpDUXpGT01HZHNPRTQxVTNZMEwxUnJlbFpxWWxOUlVYQkRkbkp2TDFGM01sVTFkR2hYV25GR2NXWXhVa0kzT0c1NFVtazRNSFJaVWxsWFVtcHpUeTl0UkRGUloxUmpNRkJuUzJGNUwyMVRUalpaUm10dU5XczlJaXdpY21WMmFYTnBiMjRpT2pFc0luQnliMjltSWpvaWVqWlhjVlZXYTNJMk9YSkVZek5VUWs4M01YRm5iU3RHTVVZeWIxVTFjbXMxVWtSMFVrZFVkRzF0T1ZneWFUaEJaSFF5T0dnd1FuWnhlbk5DVlZKNWREQXhUV0ZzZWxKdU4yOVRXbFJCZWxkb2FsTkNUSGM5UFNKOSJ9",
  "shard_id": "74727573742d6e65742d6964656e746974792d706f63",
  "last_tx": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "submitter_id": "044d8757660c9d2a22e61e7e3996c330f7f79f137a16ca2c24e3800f2bf0c57d72a121b6d80549eb2df2547c517582796a36d414678ecce0b39dc1cdbea0db6e62",
  "submitter_seq": 1,
  "padding": 0,
  "signature": "2WbqRao5jzGHpDLZpdwZXMZ5KaxHpeIQZVvqM1ECk2xI4AA2Rr5lbenIKOKSet0TD1YBzTff7Fdv6wH+7ceKyA=="
}
```
> Output of the command printed as above can be used as the request body for transaction submission to an idnode application using an offline REST tool like postman or curl 

### Update transaction history
Following CLI command is implemented to update transaction history of the test submitter after a successful offline transaction submission:

```
OWNER[01]: help update
update transaction history of the test submitter using valid [64]byte hex encoded offline transaction submission
usage: update <tx_id>

OWNER[01]: update 8af5f8d5f01b553dbe205b326711d2733f55b4b08eca95f719a65bf0b63ea5346ea9f5d992b206d3612689fb5c16b29c4dbdf6eb89e8ed8f58113ac27397f313
OWNER[02]: 
```
> After any offline transaction submission, transaction history must be updated as above, to keep the test submitter compliant with [submitter sequencing rules](https://github.com/trust-net/dag-documentation#Submitter-Sequencing-Rules) of the DAG protocol.

## Online Transactions
Following commands are available to work with an idnode directly from the test driver's submitter and an implementation of idnode client...

> Below examples assume that idnode application is running on localhost, listening on port 1055

### Set idnode app base url
Following CLI command is implemented to point the driver client to an idnode app:

```
OWNER[01]: help url
point client to idnode application's base http/https url
usage: url <base url of idnode app>

OWNER[01]: url http://localhost:1055
OWNER[01]: 
```

### Ping idnode app
Following CLI command is implemented to ping (check connectivity) of registered base url, or an explicitly provided idnode url:

```
OWNER[01]: help ping
health check of registered url, or specified idnode base url
usage: ping [<base url of idnode app>]

OWNER[01]: ping http://localhost:1055
Connected!
OWNER[01]: 
```

### Submit PublicSECP256K1 registration
Following CLI command is implemented to submit PublicSECP256K1 registration transaction request via idnode application API to the registered url:

```
OWNER[01]: help submit_key
submit PublicSECP256K1 registration transaction request with revision (default revision 1) to idnode API
usage: submit_key [<revision>]

OWNER[01]: submit_key 1
OWNER[02]: 
```
> Online transaction submission automatically updates the test submitter's transaction history to keep up with [submitter sequencing rules](https://github.com/trust-net/dag-documentation#Submitter-Sequencing-Rules) of the DAG protocol.
