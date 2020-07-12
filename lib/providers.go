package lib

import (
	discover "github.com/hashicorp/go-discover"
	"fmt"
	"log"
	"io/ioutil"
)

type AWSProvider struct {}

func (provider AWSProvider) GetServers(providerOptions map[string]string) ([]string, error) {
	
	discoverer := discover.Discover{
		Providers : map[string]discover.Provider{
			"aws": discover.Providers["aws"],
		},
	}

	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	cfg := fmt.Sprintf("provider=aws region=%s access_key_id=%s secret_access_key=%s addr_type=%s tag_key=%s tag_value=%s", providerOptions["region"], providerOptions["accessKeyId"], providerOptions["secretAccessKey"], providerOptions["addrType"], providerOptions["tagKey"], providerOptions["tagValue"])
	serverIps, err := discoverer.Addrs(cfg, logger)

	if err != nil {
		return nil, err
	}

	return serverIps, nil
}
