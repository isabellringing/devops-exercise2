package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/digitallumens/devops-test/api"
	"github.com/digitallumens/devops-test/pem"
)

func main() {
	// Load environment variables
	serverAddr, ok := os.LookupEnv("TEST_SERVER_ADDR")
	if !ok {
		log.Fatal("TEST_SERVER_ADDR environment variable not set")
	}

	clientCA, err := pem.LoadPEMFromEnv("TEST_CLIENT_CA")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	clientCAPool := x509.NewCertPool()
	ok = clientCAPool.AppendCertsFromPEM(clientCA)
	if !ok {
		log.Fatal("Unable to append cert to pool")
	}

	certPEM, err := pem.LoadPEMFromEnv("TEST_CLIENT_CERT")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	keyPEM, err := pem.LoadPEMFromEnv("TEST_CLIENT_KEY")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	certificate, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		log.Fatalf("tls.X509KeyPair returned an error: %s", err.Error())
	}

	tlsConfig := &tls.Config{
		RootCAs:      clientCAPool,
		Certificates: []tls.Certificate{certificate}}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	clientRequest := api.ClientRequest{SerialNumber: 1}
	payloadBytes, err := json.Marshal(clientRequest)
	request, err := http.NewRequest("POST", serverAddr, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Fatalf("http.NewRequest returned an error: %s", err)
	}

	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("client.Do returned an error: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll returned an error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		var error api.ClientError
		err = json.Unmarshal(body, &error)
		if err != nil {
			log.Fatalf("json.Unmarshal returned an error: %s", err)
			return
		}
		log.Printf("Server returned an error: %d %s\n", resp.StatusCode, error.Msg)
		return
	}
	var reply api.ClientReply
	err = json.Unmarshal(body, &reply)
	if err != nil {
		log.Fatalf("json.Unmarshal returned an error: %s", err)
		return
	}
	log.Println("SUCCESS!")
	log.Printf(" msg: %s\n", reply.Msg)
	log.Printf(" secret: %s\n", reply.Secret)
}
