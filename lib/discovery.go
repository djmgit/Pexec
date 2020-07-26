package lib

import (
	"log"
)

type Provider interface {
	GetServers(map[string]string, *log.Logger) ([]string, error)
}

func GetProviderHandler(providerType string, logger *log.Logger) (Provider) {

	var provider Provider

	switch providerType {
	case "AWS":
		logger.Printf("Using %s provider... \n", "AWS")
		provider = AWSProvider{}
	}

	return provider
}
