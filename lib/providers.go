package lib

import (
	"github.com/hashicorp/go-discover"
)

type AWSProvider struct {
	Region string
	TagKey string
	TagValue string
	AddrType string 
	AccessKeyId string
	SecretAccessKey string
}

func (provider *AWSProvider) GetServers() {

	
}
