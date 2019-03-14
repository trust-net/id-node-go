// Copyright 2019 The trust-net Authors
// ID application full node
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/trust-net/dag-lib-go/dbp"
	"github.com/trust-net/dag-lib-go/stack"
	"github.com/trust-net/dag-lib-go/stack/p2p"
	"github.com/trust-net/id-node-go/api"
	"github.com/trust-net/id-node-go/app"
	"os"
)

func main() {
	fileName := flag.String("config", "", "config file name")
	apiPort := flag.Int("apiPort", 0, "port for client API")
	flag.Parse()
	if len(*fileName) == 0 {
		fmt.Printf("Missing required parameter \"config\"\n")
		return
	}
	// open the config file
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Failed to open config file: %s\n", err)
		return
	}
	data := make([]byte, 2048)
	// read config data from file
	config := p2p.Config{}
	if count, err := file.Read(data); err == nil {
		data = data[:count]
		// parse json data into structure
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("Failed to parse config data: %s\n", err)
			return
		}
	} else {
		fmt.Printf("Failed to read config file: %s\n", err)
		return
	}

	// instantiate DLT stack and register app
	if ldbp, err := dbp.NewDbp("database"); err != nil {
		fmt.Printf("Failed to create database: %s", err)
	} else if dlt, err := stack.NewDltStack(config, ldbp); err != nil {
		fmt.Printf("Failed to create dlt stack: %s", err)
	} else if err := dlt.Start(); err != nil {
		fmt.Printf("Failed to start dlt stack: %s", err)
	} else if err := dlt.Register(app.AppShard, app.AppName, app.TxHandler); err != nil {
		fmt.Printf("Failed to register app: %s", err)
	} else {
		fmt.Printf("Register app: %s\nShard ID: %x\n", app.AppName, app.AppShard)
		// start net server and wait
		fmt.Printf("%s\n", api.NewController(dlt).StartServer(*apiPort))
	}
}
