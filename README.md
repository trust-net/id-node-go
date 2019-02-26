# id-node-go
Trust-Net Identity Application [component specs](./doc/id-app-specs.md#id-app-specs) and [stand-alone node](./docs/id-app-node.md#id-app-node) implementation.

## Contents
* [Introduction](#Introduction)
* [Build Instructions](#Build-Instructions)
    * [Clone Repo](#Clone-Repo)
    * [Install Dependencies](#Install-Dependencies)
    * [Build Application Node](#Build-Application-Node)
* [Run Instructions](#Run-Instructions)
    * [Stage Application Node](#Stage-Application-Node)
    * [Create Network Config](#Create-Network-Config)
    * [Run Application Node](#Run-Application-Node)
* [Usage Instructions](#Usage-Instructions)
    * [Identity Node APIs](#Identity-Node-APIs)
    * [CLI Client](#CLI-Client)
* [Application Architecture](#Application-Architecture)
    * [Application Node](#Application-Node)
    * [Client API Controller](#Client-API-Controller)
    * [App Spec Tx Handler](#App-Spec-Tx-Handler)
    * [DLT Stack](#DLT-Stack)

## Introduction
This is an implementation of Trust-Net's official Identity Application. This implementation consists of two parts:
* an official [application spec](./doc/id-app-specs.md#id-app-specs) (a.k.a. app component), and
* a sample stand-alone [application node](./docs/id-app-node.md#id-app-node)

The application node implementation can be used "as is", to self-host a stand alone Trust-Net node for global Trust-Net network identities. Alternatively, the application spec can be used as a "component" by other application node implementations, to build more complex enterprise applications that work/rely-on the global identities of Trust-Net network.

## Build Instructions
Below instructions assume you have:
* platform specific distribution of [golang](https://golang.org/) installed
* env variable `GOPATH` set to `golang` workspace directory (e.g. `$HOME/go`)
* platform specific `gcc` or `CC` compiler is installed

### Clone Repo
```
mkdir -p $GOPATH/src/github.com/trust-net
cd $GOPATH/src/github.com/trust-net
git clone git@github.com:trust-net/id-node-go.git
```

### Install Dependencies
Project uses Ethereum's `go-ethereum` for low level p2p and crypto libraries from `release/1.7	` branch. Install these dependencies as following:

```
mkdir -p $GOPATH/src/github.com/ethereum
cd $GOPATH/src/github.com/ethereum
git clone --single-branch --branch release/1.7  https://github.com/ethereum/go-ethereum.git 
```

After above step, install remaining dependencies using `go get` as following:

```
cd $GOPATH/src/github.com/trust-net/id-node-go/idnode
go get
```

Above will install remaining dependencies into your `golang` workspace.

> Note: Ethereum dependency requires gcc/CC installed for compiling and building crypto library. Hence, `go get` may fail if gcc/CC is not found. Install the platform appropriate compiler and then re-run `go get`.

### Build Application Node

```
cd $GOPATH/src/github.com/trust-net/id-node-go/
(cd idnode; go build)
```

## Run Instructions

> Following instructions assume that you've built the application node as instructed [above](#Build-Instructions)

### Stage Application Node
Create a staging directory for your application:
```
mkdir -p $USER/trust-net/idnode
```

Copy the build binaries into staging area:
```
cp $GOPATH/src/github.com/trust-net/id-node-go/idnode/idnode $USER/trust-net/idnode
```

### Create Network Config
Copy the following example config files into staging directory:
```
cd $USER/trust-net/idnode

cat << EOF > config.json
{
	"key_file": "node.key",
	"key_type": "ECDSA_S256",
	"max_peers": 10,
	"node_name": "<name your node here>",
	"listen_port": "50808",
	"boot_nodes": [
     "enode://c3da24ed70538b731b9734e4e0b8206e441089ab4fcd1d0faadb1031e736491b70de0b70e1d581958b28eb43444491b3b9091bd8a81d1767bf7d4ebc3e7bd108@<other.node.IP.address>:<other-node-port>"
   ]
}
EOF
```
Please make sure:
* `node_name` has appropriate string name to identify your node
* `<other.node.IP.address>:<other-node-port>` in the `boot_nodes` array use appropriate IP and Port of another instance of a trust-node application node (either this app, or some other app), if you want this node to join/discover the network
* if your other instance of a trust-node application node is running on same host, then `listen_port` and `key_file` values are different between the two instances

### Run Application Node
Start the application node instance as following, to listen on HTTP port 1055 (or a different port of your choice) for client API and use [above](#Create-Network-Config) created network configuration:
```
./idnode -h
Usage of ../idnode:
  -apiPort int
    	port for client API
  -config string
    	config file name

./idnode -apiPort 1055 -config config.json 
```

## Usage Instructions
The application node can be used with a remote client, or with the test driver client provided with application as following:

### Identity Node APIs
Refer to the stand alone application node [documentation](./docs/id-app-node.md#id-app-node) for Identity Node APIs that can be used with a remote submitter client to submit transactions for identity attribute management, as well as to access the identity attributes for a network identity.

### CLI Client
A test CLI client is provided to submit transactions as a network identity owner. This client can be used as following:
```
cd $GOPATH/src/github.com/trust-net/id-node-go/

go run test/owner/main.go 
```
Refer to the test identity owner client CLI [documentation](./docs/driver_id_owner.md#test-idnode-owner) for details of various CLI commands to use this client.

## Application Architecture
In a quick summary, the architecture for trust-net applications consists of following layers:
```
+--------------------------------------------------+
|                 Application Node                 |
|  +--------------------+  +--------------------+  |
|  |                    |  |                    |  |
|  |    Client API      |  |      App Spec      |  |
|  |    Controller      |  |     Tx Handler     |  |
|  |                    |  |                    |  |
|  +----------||--------+  +----------||--------+  |
|  +----------||----------------------||--------+  |
|  |                                            |  |
|  |                DLT Stack                   |  |
|  |                                            |  |
|  +--------------------||----------------------+  |
|                       ||                         |
+-----------------------||-------------------------+
                 Trust-Net Network
```

### Application Node
The whole application is encapsulated within an application node. This is a stand alone node in the trust-net network, and it implements the client APIs for interfacing with the application. Developers may choose to implement and/or host their own instances of the application nodes. However, only the Client API controllers (and API access control) are the only differentiators among these implementations. Undreneath, the "App Spec" and "DLT Stack" would typically be re-used from official implementations.

> This `README` documents how to build, deploy and use one such application node for Trust-Net's official Identity Application.

### Client API Controller
Each application node implementation needs to provide the client APIs for submitting application specific transactions to DLT stack, and to access application's world state from the DLT stack.

> Documentation for client API controller of the Trust-Net Identity Application node is [here](./docs/id-app-node.md#Identity-Node-APIs).

### App Spec Tx Handler
This is the core application business logic, which implements methods to decode transaction payload as per application specifications, process transaction instruction(s) as per application's specifications, and update application's world state based on the result of the transaction processing.

> Documentation for the application specs and transaction handler of the Trust-Net Identity Application component is [here](./docs/id-app-specs.md#Identity-App-Transactions).

### DLT Stack
Application node implementation's instantiate the official Trust-Net DAG protocol stack, to connect with the Trust-Net network across the globe.

> Documentation for the DLT Stack and DAG protocol is [here](https://github.com/trust-net/dag-documentation#dag-documentation).
