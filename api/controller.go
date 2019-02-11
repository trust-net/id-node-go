// Copyright 2019 The trust-net Authors
// ID Application REST api controller
package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/trust-net/dag-lib-go/api"
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack"
	"net/http"
	"strconv"
)

type controller struct {
	logger log.Logger
	dlt    stack.DLT
}

func NewController(dlt stack.DLT) *controller {
	return &controller{
		logger: log.NewLogger("Controller"),
		dlt:    dlt,
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

func (c *controller) ping(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Recieved ping from: %s", r.RemoteAddr)
	setHeaders(w)
	json.NewEncoder(w).Encode("pong!")
}

func (c *controller) getAttributeByName(w http.ResponseWriter, r *http.Request) {
	// fetch request params
	params := mux.Vars(r)
	name := params["name"]
	id := params["id"]
	c.logger.Debug("Recieved get attribute '%s' for '%s' from: %s", name, id, r.RemoteAddr)

	// set headers
	setHeaders(w)

	// validate parameters
	if len(id) != 130 {
		c.logger.Debug("incorrect identity owner id length: %d", len(id))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("id must be hex encoded trust-net id")
		return
	}
	if len(name) < 1 {
		c.logger.Debug("incorrect attribute name length: %d", len(name))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing attribute name")
		return
	}

	// fetch resource
	// TBD
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode("get attribute not implemented")
}

func (c *controller) submitTransaction(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Recieved transaction submit request from: %s", r.RemoteAddr)
	// set headers
	setHeaders(w)
	// parse request body
	req, err := api.ParseSubmitRequest(r)
	if err != nil {
		c.logger.Debug("Failed to decode request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	c.logger.Debug("Transaction submitter: %x", req.DltRequest().SubmitterId)
	// submit transaction to app
	if tx, err := c.dlt.Submit(req.DltRequest()); err != nil {
		c.logger.Debug("failed to submit transaction: %s", err)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(err.Error())
	} else {
		c.logger.Debug("successfully submitted transaction: %x", tx.Id())
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(api.NewSubmitResponse(tx))
	}
}

func (c *controller) StartServer(listenPort int) error {
	// if not a valid port, do not start
	if listenPort < 1024 {
		return fmt.Errorf("Invalid port: %d", listenPort)
	}
	c.logger.Debug("Starting api controller on port: %d", listenPort)
	router := mux.NewRouter()
	router.HandleFunc("/ping", c.ping).Methods("GET")
	router.HandleFunc("/identity/{id}/attributes/{name}", c.getAttributeByName).Methods("GET")
	router.HandleFunc("/submit", c.submitTransaction).Methods("POST")
	err := http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	c.logger.Error("End of server: %s", err)
	return err
}
