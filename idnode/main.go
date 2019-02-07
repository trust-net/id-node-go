// Copyright 2019 The trust-net Authors
// ID application full node
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/trust-net/dag-lib-go/stack/p2p"
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

	//	// start net server
	//	if err := StartServer(*apiPort); err != nil {
	//		fmt.Printf("Did not start client API: %s\n", err)
	//	}

	//	// instantiate two DLT stacks
	//	if localDlt, err := stack.NewDltStack(config, db.NewInMemDbProvider()); err != nil {
	//		fmt.Printf("Failed to create 1st DLT stack: %s", err)
	//	} else if remoteDlt, err := stack.NewDltStack(config2, db.NewInMemDbProvider()); err != nil {
	//		fmt.Printf("Failed to create 2nd DLT stack: %s", err)
	//	} else if err = cli(localDlt, remoteDlt); err != nil {
	//		fmt.Printf("Error in CLI: %s", err)
	//	} else {
	//		fmt.Printf("Shutdown cleanly")
	//	}
	//	fmt.Printf("\n")
}
