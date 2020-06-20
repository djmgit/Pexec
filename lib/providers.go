package lib

import (
	discover "github.com/hashicorp/go-discover"
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
	
	discoverer := discover.Discover{
		Providers : map[string]discover.Provider{
			"aws": discover.Providers["aws"],
		},
	}

	cfg := "provider=aws region=%s access_key_id=%s secret_access_key=%s addr_type=%s tag_key=%s tag_value=%s", provider.Region, provider.AccessKeyId, provider.SecretAccessKey, provider.AddrType, provider.TagKey, provider.TagValue
}
