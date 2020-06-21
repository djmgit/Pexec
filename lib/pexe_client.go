package lib

import (
	"os"
	"golang.org/x/crypto/ssh"
	"time"
)

const CUSTOM string = "CUSTOM"
const AWS string = "AWS"
const DEFAULT_USER = "root"

type PexecClient struct {

	TargetServers []Server

	ProviderOptions map[string]string

	Provider string

	Parallel bool

	Batch bool

	BatchSize int

	User string

	KeyPath string

	TimeOut time.Duration

	SSHConConfig *ssh.ClientConfig
}

func (client *PexecClient) getDefaults()  {

	if client.Provider == "" {
		client.Provider = CUSTOM
	}

	if client.User == "" {
		client.User = DEFAULT_USER
	}

	if client.KeyPath == "" {
		client.KeyPath = os.Getenv("HOME") + "/.ssh/id_rsa"
	}

	if client.TimeOut == 0 {
		client.TimeOut = 30
	}

	var sshconerror error

	client.SSHConConfig, sshconerror = PrepareSSHConConfig(client.User, client.KeyPath)

	if sshconerror != nil {
		panic(sshconerror.Error())
	}
}

func (client *PexecClient) populateTargetServers() {

	provider := GetProviderHandler(client.Provider)

	serverIps, err := provider.GetServers(client.ProviderOptions)

	if err != nil {
		panic(err.Error())
	}
}

func (client *PexecClient) Run(command string) ([]CommandResponseWithServer, error) {

	client.getDefaults()

	if client.Parallel {

		commandResponseWithServer, err := ParallelBatchExecute(command, client.SSHConConfig, client.TargetServers)

		if err != nil {
			return nil, err
		}

		return commandResponseWithServer, nil

	} else {

		commandResponseWithServer, err := SerialExecute(command, client.SSHConConfig, client.TargetServers)

		if err != nil {
			return nil, err
		} else {
			return commandResponseWithServer, nil
		}
	}

	return nil, nil
}
