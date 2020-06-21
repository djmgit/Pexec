package lib

type Provider interface {
	GetServers() ([]string, error)
}

func GetProviderHandler(providerType string) (Provider) {

	var provider Provider

	switch providerType {
	case "AWS":
		provider = AWSProvider{}
	}

	return provider
}
