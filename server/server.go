package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/digitallumens/devops-test/pem"
)

func main() {

	// Load environment variables
	portNumber, ok := os.LookupEnv("TEST_SERVER_PORT")
	if !ok {
		log.Fatal("TEST_SERVER_PORT environment variable not set")
	}

	CApem, err := pem.LoadCertFromEnv("TEST_SERVER_CA", "CERTIFICATE")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	CAcert, err := x509.ParseCertificate(CApem.Bytes)
	if err != nil {
		log.Fatalf("ParseCertificate returned an error: %s", err.Error())
	}

	clientCAPool := x509.NewCertPool()
	clientCAPool.AddCert(CAcert)

	certPEM, err := pem.LoadPEMFromEnv("TEST_SERVER_CERT")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	keyPEM, err := pem.LoadPEMFromEnv("TEST_SERVER_KEY")
	if err != nil {
		log.Fatalf("LoadFromEnv returned an error: %s", err.Error())
	}
	certificate, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("tls.X509KeyPair returned an error: %s", err.Error())
	}

	// Create a tls.Config using our client CA(s)
	tlsConfig := tls.Config{ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:                clientCAPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS10,
		Certificates:             []tls.Certificate{certificate}}

	// Set up a goroutine to cleanly exit on os.Interrupt signal (aka CTRL-C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// sig is a ^C, handle it
		log.Fatal("Exiting")
		os.Exit(0)
	}()

	log.Printf("Starting server on port %s", portNumber)
	s := &http.Server{
		Addr:      fmt.Sprintf(":%s", portNumber),
		TLSConfig: &tlsConfig,
		Handler:   http.HandlerFunc(handler),
	}
	tlsListener, err := tls.Listen("tcp", s.Addr, &tlsConfig)
	if err != nil {
		log.Fatalf("tls.Listen returned an error: %s\n", err.Error())
	}
	log.Fatal(s.Serve(tlsListener))
}
