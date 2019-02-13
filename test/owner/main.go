package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/trust-net/dag-lib-go/api"
	"github.com/trust-net/id-node-go/app"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	self    = "OWNER"
	owner   = app.TestSubmitter()
	baseUrl = ""
)

var commands = map[string][2]string{
	"print_key":  {"usage: print_key [<revision>]", "print transaction request for registering PublicSECP256K1 attribute with revision (default revision 1)"},
	"submit_key": {"usage: submit_key [<revision>]", "submit PublicSECP256K1 registration transaction request with revision (default revision 1) to idnode API"},
	"update":     {"usage: update <tx_id>", "update transaction history of the test submitter using valid [64]byte hex encoded offline transaction submission"},
	"url":        {"usage: url <base url of idnode app>", "point client to idnode application's base http/https url"},
	"ping":       {"usage: ping [<base url of idnode app>]", "health check of registered url, or specified idnode base url"},
}

func cmdPrompt() string {
	return fmt.Sprintf("%s[%02d]: ", self, owner.Seq())
}

type idnodeClient struct {
	baseUrl string
}

func (c *idnodeClient) Ping() bool {
	resp, err := http.Get(c.baseUrl + "/ping")
	return err == nil && resp.StatusCode == 200
}

func (c *idnodeClient) Submit(op *api.SubmitRequest) (*api.SubmitResponse, error) {
	opBytes, _ := json.Marshal(op)
	resp, err := http.Post(c.baseUrl+"/submit", "application/json", bytes.NewBuffer(opBytes))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			if resp.StatusCode == 201 {
				resp := &api.SubmitResponse{}
				if err := json.Unmarshal(body, resp); err != nil {
					return nil, err
				} else {
					return resp, nil
				}
			} else {
				return nil, fmt.Errorf("%s", body)
			}
		}
	}
}

func NewIdnodeClient(baseUrl string) (*idnodeClient, error) {
	var client *idnodeClient
	if url, err := url.ParseRequestURI(baseUrl); err == nil &&
		(url.Scheme == "http" || url.Scheme == "https") &&
		len(url.Host) > 0 {
		// strip any trailing '/' if present
		if baseUrl[len(baseUrl)-1] == '/' {
			baseUrl = baseUrl[:len(baseUrl)-1]
		}
		client = &idnodeClient{
			baseUrl: baseUrl,
		}

		// send ping health check
		if client.Ping() {
			return client, nil
		} else {
			return nil, fmt.Errorf("Failed to conect with url: %s", baseUrl)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Failed to parse url: %s", err)
	} else {
		return nil, fmt.Errorf("Bad url: %s", baseUrl)
	}
}

func main() {
	var (
		client *idnodeClient
	)
	for {
		fmt.Printf(cmdPrompt())
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
					case "submit_key":
						// get the revision
						rev := uint64(1)
						if wordScanner.Scan() {
							value, _ := strconv.Atoi(wordScanner.Text())
							rev = uint64(value)
						}
						if rev > 0 {
							// get a transaction for the key registration
							op := owner.PublicSECP256K1Op(rev)
							// submit transaction via API client
							if client != nil {
								if resp, err := client.Submit(op); err != nil {
									fmt.Printf("Failed to submit transaction: %s\n", err)
								} else {
									// get the transaction ID from response
									txId, _ := hex.DecodeString(resp.TxId)
									// update owner's transaction history
									owner.Update(txId)
								}
							} else {
								fmt.Printf("No base url registered, use \"url\" command to register first\n")
							}
						} else {
							fmt.Printf("%s\n", commands["submit_key"][1])
							fmt.Printf("%s\n", commands["submit_key"][0])
						}
					case "update":
						var tx_id []byte
						if wordScanner.Scan() {
							tx_id, _ = hex.DecodeString(wordScanner.Text())
						}

						if len(tx_id) == 64 {
							// update the submitter with successful transaction
							owner.Update(tx_id)
						} else {
							fmt.Printf("%s\n", commands["update"][1])
							fmt.Printf("%s\n", commands["update"][0])
						}
					case "url":
						if wordScanner.Scan() {
							newClient, err := NewIdnodeClient(wordScanner.Text())
							if err != nil {
								fmt.Printf("%s\n", err)
							} else {
								client = newClient
							}
						} else {
							fmt.Printf("%s\n", commands["url"][1])
							fmt.Printf("%s\n", commands["url"][0])
						}
					case "ping":
						pingClient := client
						var err error
						if wordScanner.Scan() {
							pingClient, err = NewIdnodeClient(wordScanner.Text())
							if err != nil {
								fmt.Printf("(%s)\n", err)
							}
						}
						if err != nil {
							fmt.Printf("%s\n", commands["ping"][1])
							fmt.Printf("%s\n", commands["ping"][0])
						} else if pingClient == nil {
							fmt.Printf("No base url registered, either specify url, or register using \"url\" command\n")
						} else {
							if pingClient.Ping() {
								fmt.Printf("Connected!\n")
							} else {
								fmt.Printf("Not connected!\n")
							}
						}
					case "help":
						if wordScanner.Scan() {
							cmd := wordScanner.Text()
							if _, found := commands[cmd]; found {
								fmt.Printf("%s\n", commands[cmd][1])
								fmt.Printf("%s\n", commands[cmd][0])
								break
							}
						}
						fmt.Printf("Accepted commands...\n")
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
					break
				}
				break
			}
			fmt.Printf("\n%s", cmdPrompt())
		}
	}
}
