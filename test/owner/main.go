package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/trust-net/id-node-go/app"
	"os"
	"strconv"
	"strings"
)

var (
	cmdPrompt = "OWNER: "
	owner     = app.TestSubmitter()
)

var commands = map[string][2]string{
	"print_key": {"usage: print_key [<revision>]", "print transaction request for registering PublicSECP256K1 attribute with revision (default revision 1)"},
}

func main() {
	for {
		fmt.Printf(cmdPrompt)
		lineScanner := bufio.NewScanner(os.Stdin)
		for lineScanner.Scan() {
			line := lineScanner.Text()
			if len(line) != 0 {
				wordScanner := bufio.NewScanner(strings.NewReader(line))
				wordScanner.Split(bufio.ScanWords)
				for wordScanner.Scan() {
					cmd := wordScanner.Text()
					switch cmd {
					case "quit":
						fallthrough
					case "q":
						return
					case "print_key":
						// get the revision
						rev := uint64(1)
						if wordScanner.Scan() {
							value, _ := strconv.Atoi(wordScanner.Text())
							rev = uint64(value)
						}
						if rev > 0 {
							// get a transaction for the key registration
							op := owner.PublicSECP256K1Op(rev)
							text, _ := json.MarshalIndent(op, "", "  ")
							fmt.Printf("%s\n", text)
						} else {
							fmt.Printf("%s\n", commands["print_key"][1])
							fmt.Printf("%s\n", commands["print_key"][0])
						}
					default:
						fmt.Printf("Unknown Command: %s", cmd)
						for wordScanner.Scan() {
							fmt.Printf(" %s", wordScanner.Text())
						}
						fmt.Printf("\n\nAccepted commands...\n")
						isFirst := true
						for k, _ := range commands {
							if !isFirst {
								fmt.Printf(", ")
							} else {
								isFirst = false
							}
							fmt.Printf("\"%s\"", k)
						}
						fmt.Printf("\n")
						break
					}
				}
			}
			fmt.Printf("\n%s", cmdPrompt)
		}
	}
}
