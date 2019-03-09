// Copyright 2019 The trust-net Authors
// ID Application REST api controller
package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/trust-net/dag-lib-go/api"
	"github.com/trust-net/dag-lib-go/log"
	"github.com/trust-net/dag-lib-go/stack"
	"github.com/trust-net/id-node-go/app"
	"github.com/trust-net/id-node-go/dto"
	"net/http"
	"strconv"
	"strings"
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

func (c *controller) validateAttributeParams(r *http.Request) (idBytes []byte, name string, err error) {
	// fetch request params
	params := mux.Vars(r)
	name = params["name"]
	id := params["id"]
	c.logger.Debug("Recieved get attribute '%s' for '%s' from: %s", name, id, r.RemoteAddr)

	// validate parameters
	if len(id) != 130 {
		c.logger.Debug("incorrect identity owner id length: %d", len(id))
		err = fmt.Errorf("incorrect public id length")
		return
	}
	idBytes, err = hex.DecodeString(id)
	if err != nil {
		c.logger.Debug("failed to decode hex id: %s", err)
		err = fmt.Errorf("id must be hex encoded trust-net id")
		return
	}
	if len(name) < 1 {
		c.logger.Debug("incorrect attribute name length: %d", len(name))
		err = fmt.Errorf("missing attribute name")
		return
	}
	return
}

func (c *controller) GetEndorsementByName(w http.ResponseWriter, r *http.Request) {
	c.getEndorsementByName(w, r)
}

func (c *controller) getEndorsementByName(w http.ResponseWriter, r *http.Request) {
	// set headers
	setHeaders(w)

	// fetch attribute parameters
	idBytes, name, err := c.validateAttributeParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
	}

	// attribute registration is only applicable to specific attributes
	isValid := false
	switch strings.ToLower(name) {
	case "preferredemail":
		isValid = true
	}
	if !isValid {
		c.logger.Debug("attribute type '%s' not support for endorsement", name)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("attribute does not support endorsement")
		return
	}

	// fetch resource
	if res, err := c.dlt.GetState(app.NewIdState(idBytes, nil).Prefixed(name)); err != nil {
		c.logger.Debug("failed to get resource from DLT: %s", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	} else {
		if attribute, err := dto.AttributeEndorsementFromBytes(res.Value); err != nil {
			c.logger.Error("Failed to de-serialize world state resource: %s", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		} else {
			c.logger.Error("responding back with attribute revision: %d", attribute.Revision)
			json.NewEncoder(w).Encode(attribute)
			return
		}
	}
}

func (c *controller) GetRegistrationByName(w http.ResponseWriter, r *http.Request) {
	c.getRegistrationByName(w, r)
}

func (c *controller) getRegistrationByName(w http.ResponseWriter, r *http.Request) {
	// set headers
	setHeaders(w)

	// fetch attribute parameters
	idBytes, name, err := c.validateAttributeParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
	}

	// attribute registration is only applicable to specific attributes
	isValid := false
	switch strings.ToLower(name) {
	case "publicsecp256k1":
		fallthrough
	case "preferredfirstname":
		fallthrough
	case "preferredlastname":
		isValid = true
	}
	if !isValid {
		c.logger.Debug("attribute type '%s' not support for registration", name)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("attribute does not support registration")
		return
	}

	// fetch resource
	if res, err := c.dlt.GetState(app.NewIdState(idBytes, nil).Prefixed(name)); err != nil {
		c.logger.Debug("failed to get resource from DLT: %s", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	} else {
		if attribute, err := dto.AttributeRegistrationFromBytes(res.Value); err != nil {
			c.logger.Error("Failed to de-serialize world state resource: %s", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		} else {
			c.logger.Error("responding back with attribute revision: %d", attribute.Revision)
			json.NewEncoder(w).Encode(attribute)
			return
		}
	}
}

func (c *controller) SubmitTransaction(w http.ResponseWriter, r *http.Request) {
	c.submitTransaction(w, r)
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
	router.HandleFunc("/identity/{id}/registrations/{name}", c.getRegistrationByName).Methods("GET")
	router.HandleFunc("/identity/{id}/endorsements/{name}", c.getEndorsementByName).Methods("GET")
	router.HandleFunc("/submit", c.submitTransaction).Methods("POST")
	err := http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	c.logger.Error("End of server: %s", err)
	return err
}
