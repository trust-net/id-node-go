// Copyright 2019 The trust-net Authors
// ID application full node
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/trust-net/dag-lib-go/stack/p2p"
	"github.com/trust-net/id-node-go/api"
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

	// instantiate DLT stack
	// TBD

	// start DLT stack
	// TBD

	// register app
	// TBD

	// start net server
	fmt.Printf("%s\n", api.StartServer(*apiPort))
}
