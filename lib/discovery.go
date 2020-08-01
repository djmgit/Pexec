package lib

import (
	"log"
)

// Provider interface to hold the desired interface
type Provider interface {
	GetServers(map[string]string, *log.Logger) ([]string, error)
}

// Fucntion to return the prober provider handler using the above interface
// As of now it only supports AWS
func GetProviderHandler(providerType string, logger *log.Logger) (Provider) {

	var provider Provider

	switch providerType {
	case "AWS":
		logger.Printf("Using %s provider... \n", "AWS")
		provider = AWSProvider{}
	}

	return provider
}
