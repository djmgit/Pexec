package lib

type Provider interface {
	GetServers() ([]string, error)
}
