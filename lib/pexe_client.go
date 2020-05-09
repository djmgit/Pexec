package lib

import (
	"os"
	"golang.org/x/crypto/ssh"
)

const CUSTOM string = "CUSTOM"
const AWS string = "AWS"
const DEFAULT_USER = "root"

type PexecClient struct {

	TargetServers []Server

	Provider string

	Parallel bool

	Batch bool

	BatchSize int

	User string

	KeyPath string

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

	var sshconerror error

	client.SSHConConfig = PrepareSSHConConfig()
}
