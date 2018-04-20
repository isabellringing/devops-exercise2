package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"

	"github.com/digitallumens/devops-test/api"
)

func handler(rw http.ResponseWriter, r *http.Request) {
	// Must be POST
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get a copy of the client certificate from the Request object.
	if r.TLS == nil {
		writeError(rw, r, http.StatusBadRequest, "request has no tls.ConnectionState")
		return
	}
	if len(r.TLS.PeerCertificates) == 0 {
		writeError(rw, r, http.StatusBadRequest, "request did not present a client certificate")
		return
	}
	clientCert := r.TLS.PeerCertificates[0]
	if clientCert.SerialNumber.Cmp(big.NewInt(math.MaxUint32)) == 1 {
		writeError(rw, r, http.StatusBadRequest, "client certificate serial number exceeded MaxUint32")
		return
	}
	clientCertSerialNumber := uint32(clientCert.SerialNumber.Uint64())

	// Read body.
	if r.Body == nil {
		writeError(rw, r, http.StatusBadRequest, "empty request body")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(rw, r, http.StatusInternalServerError, "ioutil.ReadAll returned an error: %s\n", err)
		return
	}

	// Decode the JSON payload in body into an api.ClientRequest.
	var clientRequest api.ClientRequest
	err = json.Unmarshal(body, &clientRequest)
	if err != nil {
		writeError(rw, r, http.StatusBadRequest, "json.Unmarshal request body returned an error: %s\n", err)
		return
	}

	// Test serial numbers against each other.
	if clientRequest.SerialNumber != clientCertSerialNumber {
		writeError(rw, r, http.StatusUnauthorized, "client serial number mismatch")
		return
	}

	log.Printf("SUCCESS")

	reply := api.ClientReply{Msg: "hello", Secret: "seekrit"}
	msg, err := json.Marshal(reply)
	if err != nil {
		writeError(rw, r, http.StatusInternalServerError, "json.Marshal reply returned an error: %s", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(msg)
}

// writeError is a helper function that writes an error message to the provided ResponseWriter.
func writeError(rw http.ResponseWriter, r *http.Request, status int, detail string, detailArgs ...interface{}) {
	e := fmt.Sprintf(detail, detailArgs...)
	log.Printf("ERROR: %s\n", e)

	wrapper := make(map[string]interface{})
	wrapper["error"] = e

	msg, _ := json.Marshal(wrapper)
	rw.Header().Set("Content-Type", "application/json")
	http.Error(rw, string(msg), status)
}
