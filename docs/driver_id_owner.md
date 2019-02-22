## test-idnode-owner
* [Introduction](#Introduction)
    * [Usage](#Usage)
    * [CLI commands](#CLI-commands)
* [Offline Transactions](#Offline-Transactions)
    * [Update transaction history](#Update-transaction-history)
    * [Print PublicSECP256K1 registration](#Print-PublicSECP256K1-registration)
    * [Print PreferredFirstName registration](#Print-PreferredFirstName-registration)
    * [Print PreferredLastName registration](#Print-PreferredLastName-registration)
    * [Print PreferredEmail endorsement](#Print-PreferredEmail-endorsement)
* [Online Transactions](#Online-Transactions)
    * [Set idnode app base url](#Set-idnode-app-base-url)
    * [Ping idnode app](#Ping-idnode-app)
    * [Submit PublicSECP256K1 registration](#Submit-PublicSECP256K1-registration)
    * [Submit PreferredFirstName registration](#Submit-PreferredFirstName-registration)
    * [Submit PreferredLastName registration](#Submit-PreferredLastName-registration)
    * [Submit PreferredEmail endorsement](#Submit-PreferredEmail-endorsement)

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
"submit_first_name", "print_last_name", "submit_last_name", "url", "ping", "print_first_name", "print_key", "submit_key", "update"
```

## Offline Transactions
Following commands are available for offline identity attribute registration using the test driver...

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

### Print PreferredFirstName registration
Following CLI command is implemented to print transaction request for ID Node's API to register the PreferredFirstName attribute of a test identity:

```
OWNER[01]: help print_first_name
print transaction request for registering PreferredFirstName attribute with revision (default revision 1)
usage: print_first_name <first name> [<revision>]

OWNER[01]: print_first_name amit 1
{
  "payload": "eyJvcF9jb2RlIjoxLCJhcmdzIjoiZXlKdVlXMWxJam9pVUhKbFptVnljbVZrUm1seWMzUk9ZVzFsSWl3aWRtRnNkV1VpT2lKaGJXbDBJaXdpY21WMmFYTnBiMjRpT2pFc0luQnliMjltSWpvaUluMD0ifQ==",
  "shard_id": "74727573742d6e65742d6964656e746974792d706f63",
  "last_tx": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "submitter_id": "049d639f4e965f96ff459844f320d5b04df67beee73525df22bea6c27978593465baeb192e3fabab5ad50a9c64b5b2da558b648b55ba6621a2ac260f3134c2926d",
  "submitter_seq": 1,
  "padding": 0,
  "signature": "lhrMsLFDLzDWExJ4MVxFZnn5UExvurkSL/s21cdOGhlgYkWno/ANUbgV+zo3jpwSefjgGpaSWASR1303J37Dvg=="
}
```
> Output of the command printed as above can be used as the request body for transaction submission to an idnode application using an offline REST tool like postman or curl 

### Print PreferredLastName registration
Following CLI command is implemented to print transaction request for ID Node's API to register the PreferredLastName attribute of a test identity:

```
OWNER[01]: help print_last_name
print transaction request for registering PreferredLastName attribute with revision (default revision 1)
usage: print_last_name <last name> [<revision>]

OWNER[01]: print_last_name bhadoria 1
{
  "payload": "eyJvcF9jb2RlIjoxLCJhcmdzIjoiZXlKdVlXMWxJam9pVUhKbFptVnljbVZrVEdGemRFNWhiV1VpTENKMllXeDFaU0k2SW1Kb1lXUnZjbWxoSWl3aWNtVjJhWE5wYjI0aU9qRXNJbkJ5YjI5bUlqb2lJbjA9In0=",
  "shard_id": "74727573742d6e65742d6964656e746974792d706f63",
  "last_tx": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "submitter_id": "049d639f4e965f96ff459844f320d5b04df67beee73525df22bea6c27978593465baeb192e3fabab5ad50a9c64b5b2da558b648b55ba6621a2ac260f3134c2926d",
  "submitter_seq": 1,
  "padding": 0,
  "signature": "85Blms04Ev2PM+GPSogReqPvnAsPuDo067yGsAtiUB9CYn3Q85UeM1kiWzD7VNBEOtdbSivYOFYJ3Yd5N5M8rQ=="
}
```
> Output of the command printed as above can be used as the request body for transaction submission to an idnode application using an offline REST tool like postman or curl 

### Print PreferredEmail endorsement
Following CLI command is implemented to print transaction request for ID Node's API to submit and endorsement for the PreferredEmail attribute of a test identity:

```
OWNER[01]: help print_email
print transaction request for PreferredEmail attribute endorsement with revision (default revision 1)
usage: print_email <email address> [<revision>]

OWNER[01]: print_email test@example.com 1
{
  "payload": "eyJvcF9jb2RlIjoyLCJhcmdzIjoiZXlKdVlXMWxJam9pVUhKbFptVnljbVZrUlcxaGFXd2lMQ0psYm1SdmNuTmxjbDlwWkNJNklrSklNemwxZVVKYU9HcDBZbWh6VURaV2RtcFVXbU5KWWs5U05qQjNZVTlDU0dGT1RVY3diRTE1SzNNMmNUSnZVVGwyVVVKWlozVlBiR3BKUzJ0blVGUmlOak15Y0ZJMlNrSlhSSE5qVlhOSGJrZzVTRUpyV1QwaUxDSmxibU5mYzJWamNtVjBJam9pUWtGRkwxcEVUWFZ2VUhoUFNEVm1VWE5aWWtwR1NuUm1Zazg0ZW1WWmMxcHphRTAyTDFCbFlYbGllSEJSTHk4M1NHeFZkMFk0UlRaWFNrbGhNRlIzVEVkak1UazFUVTR6VVU5RFRqSnVha3R1WVhGMVMzSm5SM1Z3Ymt4SmFWSmxTWFpNTjBsdlVtazBRMHh5U1hwTVlUaDJWbHBYZURGbU9FRmxZbkpMTjJsNGVYaGtXR0Z6YkRST1RrVnRlakoxZGl0TkswMWlTWGRzWldaNldtcDRUazB6T1ZremJFNVVjQ3RJWW01dFYxQlZkWGxDVFdabUwySmpSRTlZWVhCMlRHYzlQU0lzSW1WdVkxOTJZV3gxWlNJNklqVjZlVmRLY0RCSk9WUlpVV1JoUkVkRWNUbHhaeTgxTjFWbVNtWTNaM0ZUWkUxck5rd3pTRzl5ZGpROUlpd2ljbVYyYVhOcGIyNGlPakVzSW1WdVpHOXljMlZ0Wlc1MElqb2lUbTVxUm1Gek1FbE9hbUptTjFKWmIzSm9NVTFGVWtkbFQzQTJUSEpyVDNoVllWTktOMmREVVhCTEt5dHlhekpxV1cwMGJqSnpia2h2Y0ZsSFQwOXFSazlhYUhkRGIzTkhZVmxqWVN0Mk1uaHJkRUZhWldjOVBTSjkifQ==",
  "shard_id": "74727573742d6e65742d6964656e746974792d706f63",
  "last_tx": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "submitter_id": "04da61c728e587a73957a24497fea050b6d3d4318df20b54ba3e462a529775bbfea1a06fad4cb4e39fcc2647ef3c72455ef60a8e1a9829453d83df1136bd2d221e",
  "submitter_seq": 1,
  "padding": 0,
  "signature": "Yfo6X5/37LXHvNetRQI/xHhIBwaChRaAQpQtumpZZBQlNDm7rWHh+e1vW1kotlqkTt8XR6cNYn7qFppfw6w2Aw=="
}
```
> Output of the command printed as above can be used as the request body for transaction submission to an idnode application using an offline REST tool like postman or curl 

## Online Transactions
Following commands are available to work with an idnode directly from the test driver's submitter and an implementation of idnode client...

> Below examples assume that idnode application is running on localhost, listening on port 1055. All the attribute registrations done using below online commands can be verified by calling the [Identity Access API](../README.md#Identity-Access-API) of the idnode application from postman or curl, e.g.:
```
$ curl -v http://localhost:1055/identity/049d639f4e965f96ff459844f320d5b04df67beee73525df22bea6c27978593465baeb192e3fabab5ad50a9c64b5b2da558b648b55ba6621a2ac260f3134c2926d/attributes/PublicSECP256K1
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 1055 (#0)
> GET /identity/049d639f4e965f96ff459844f320d5b04df67beee73525df22bea6c27978593465baeb192e3fabab5ad50a9c64b5b2da558b648b55ba6621a2ac260f3134c2926d/attributes/PublicSECP256K1 HTTP/1.1
> Host: localhost:1055
> User-Agent: curl/7.55.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 14 Feb 2019 06:23:27 GMT
< Content-Length: 238
< 
{"name":"PublicSECP256K1","value":"BCqMj5q4Eg8egWONgQxYRF4c3tWxSTegmLE4IWZk2u78PhUXAE+ie+QrAG91mMHS7vMfUx+s7KK2TxVdAVACw+g=","revision":1,"proof":"418l9X87fJy7vnMkoWIeIGaHu1rd4jh5V2BFsvL89msm/5htwuP6+wxwxj5UK66f3+sviXn/oOrA0DwWXcRWng=="}
* Connection #0 to host localhost left intact
```

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

### Submit PreferredFirstName registration
Following CLI command is implemented to submit PreferredFirstName registration transaction request via idnode application API to the registered url:

```
OWNER[01]: help submit_first_name
submit PreferredFirstName registration transaction request with revision (default revision 1) to idnode API
usage: submit_first_name <first name> [<revision>]

OWNER[01]: submit_first_name amit 1
Failed to submit transaction: "PublicSECP256K1 not registered"

OWNER[01]: submit_key
OWNER[02]: submit_first_name amit 1
OWNER[03]:  
```
> Before registering optional attribute `PreferredFirstName`, identity owner needs to register the mandatory attribute `PublicSECP256K1` for the identity, as done with the `submit_key` command in above example.

### Submit PreferredLastName registration
Following CLI command is implemented to submit PreferredLastName registration transaction request via idnode application API to the registered url:

```
OWNER[03]: help submit_last_name
submit PreferredLastName registration transaction request with revision (default revision 1) to idnode API
usage: submit_last_name <last name> [<revision>]

OWNER[03]: submit_last_name bhadoria 1
OWNER[04]: 
```
> Before registering optional attribute `PreferredLastName`, identity owner needs to register the mandatory attribute `PublicSECP256K1` for the identity, which was already done in the [PreferredFirstName example](#Submit-PreferredFirstName-registration) above.

### Submit PreferredEmail registration
Following CLI command is implemented to submit PreferredEmail endorsement transaction request via idnode application API to the registered url:

```
OWNER[04]: help submit_email
submit PreferredEmail endorsement transaction request with revision (default revision 1) to idnode API
usage: submit_email <email address> [<revision>]

OWNER[04]: submit_email test@example.com 1
OWNER[05]: 
```
> Before submitting endorsement for attribute `PreferredEmail`, identity owner needs to register the mandatory attribute `PublicSECP256K1` for the identity, which will be used to encrypt the shared secret in the endorsement request payload.
