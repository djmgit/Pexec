package lib

type Provider interface {
	GetServers() ([]string, error)
}

func GetProviderHandler(providerType string) (Provider) {

	var provider Provider

	switch providerType {
	case "aws":
		provider = AWSProvider{}
	}

	return provider
}
