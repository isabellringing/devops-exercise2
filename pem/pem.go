package pem

import (
	"encoding/pem"
	"fmt"
	"os"
)

func LoadCertFromEnv(envVar, expectedType string) (*pem.Block, error) {
	varString, ok := os.LookupEnv(envVar)
	if !ok {
		return nil, fmt.Errorf("%s environment variable not set", envVar)
	}

	rawcert, _ := pem.Decode([]byte(varString))
	if rawcert == nil {
		return nil, fmt.Errorf("pem.Decode failed %s", varString)
	}
	if rawcert.Type != expectedType {
		return nil, fmt.Errorf("Improperly formed %s: %s", expectedType, varString)
	}
	return rawcert, nil
}

// returns the raw PEM without decoding it
func LoadPEMFromEnv(envVar string) ([]byte, error) {
	varString, ok := os.LookupEnv(envVar)
	if !ok {
		return []byte{}, fmt.Errorf("%s environment variable not set", envVar)
	}

	return []byte(varString), nil
}
