package lib

import (
	"os"
	"golang.org/x/crypto/ssh"
	"time"
	"log"
	"io/ioutil"
)

const CUSTOM string = "CUSTOM"
const AWS string = "AWS"
const DEFAULT_USER = "root"

type PexecClient struct {

	TargetServers []Server

	Port int

	ProviderOptions map[string]string

	Provider string

	Parallel bool

	Batch bool

	BatchSize int

	User string

	KeyPath string

	TimeOut time.Duration

	SSHConConfig *ssh.ClientConfig

	Debug bool

	Logger *log.Logger
}

func (client *PexecClient) getDefaults()  {

	if client.Provider == "" {
		client.Provider = CUSTOM
	} else {
		client.populateTargetServers()
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

	client.Logger = log.New(ioutil.Discard, "", 0)

	if client.Debug == true {
		client.Logger = log.New(os.Stderr, "", 0)
	}
}

func (client *PexecClient) populateTargetServers() {

	provider := GetProviderHandler(client.Provider)

	serverIps, err := provider.GetServers(client.ProviderOptions, client.Logger)

	if err != nil {
		panic(err.Error())
	}

	targetServers := make([]Server, 0, 0)

	for _, serverIp := range serverIps {

		server := Server {
			Host: serverIp,
			Port: 22,
		}

		if client.Port != 0 {
			server.Port = client.Port
		}

		targetServers = append(targetServers, server)
	}

	client.TargetServers = targetServers

}

func (client *PexecClient) Run(command string) ([]CommandResponseWithServer, error) {

	client.getDefaults()

	if client.Parallel {

		if client.BatchSize != 0 {
			commandResponseWithServer, err := BatchExecuter(command, client.SSHConConfig, client.TargetServers, client.BatchSize)

			if err != nil {
				return nil, err
			}

			return commandResponseWithServer, nil
		}

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
}
