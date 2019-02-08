// Copyright 2019 The trust-net Authors
// ID Application REST api controller
package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/trust-net/dag-lib-go/api"
	"github.com/trust-net/dag-lib-go/log"
	"net/http"
	"strconv"
)

var logger = log.NewLogger("Controller")

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

func ping(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Recieved ping from: %s", r.RemoteAddr)
	setHeaders(w)
	json.NewEncoder(w).Encode("pong!")
}

func submitTransaction(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Recieved transaction submit request from: %s", r.RemoteAddr)
	// set headers
	setHeaders(w)
	// parse request body
	req, err := api.ParseSubmitRequest(r)
	if err != nil {
		logger.Debug("Failed to decode request body: %s", err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	logger.Debug("Transaction submitter: %x", req.DltRequest().SubmitterId)
	// submit transaction to app
	// TBD

	logger.Debug("transaction handling not implemented")
	w.WriteHeader(http.StatusNotAcceptable)
	json.NewEncoder(w).Encode("transaction handling not implemented")
}

func StartServer(listenPort int) error {
	// if not a valid port, do not start
	if listenPort < 1024 {
		return fmt.Errorf("Invalid port: %d", listenPort)
	}
	logger.Debug("Starting api controller on port: %d", listenPort)
	router := mux.NewRouter()
	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/submit", submitTransaction).Methods("POST")
	err := http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	logger.Error("End of server: %s", err)
	return err
}
