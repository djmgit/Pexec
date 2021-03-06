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

// struct to pexec client properties
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

// function to set defaults in PexecClient struct
func (client *PexecClient) getDefaults()  {

	// create a default logger
	client.Logger = log.New(ioutil.Discard, "", 0)

	if client.Debug == true {

		// Set log output to Stderr if debug is true
		// Should allow writing to file in future
		client.Logger = log.New(os.Stderr, "", 0)
	}

	if client.Provider == "" {
		client.Provider = CUSTOM
	} else if client.Provider != CUSTOM {

		// Populate servesrs from upstream provider
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

}

func (client *PexecClient) populateTargetServers() {

	client.Logger.Printf("Using provider %s...\n", client.Provider)

	// Get the appropriate upstream provider handler
	provider := GetProviderHandler(client.Provider, client.Logger)

	// Get the actual live server IPs
	serverIps, err := provider.GetServers(client.ProviderOptions, client.Logger)

	if err != nil {
		panic(err.Error())
	}

	targetServers := make([]Server, 0, 0)

	// Iterate over the servers and create the appropriate server structs using
	// IP and Port
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

// Function to actually exeute the command remotely using the helper functions
func (client *PexecClient) Run(command string) ([]CommandResponseWithServer, error) {

	client.getDefaults()

	if client.Parallel {

		client.Logger.Printf("Preparing to execuite comand in parallel...")

		if client.BatchSize != 0 {

			client.Logger.Printf("Batch size is %d, command will be executed in parallel in all servers in batches of selected size...", client.BatchSize)
			commandResponseWithServer, err := BatchExecuter(command, client.SSHConConfig, client.TargetServers, client.BatchSize, client.Logger)

			if err != nil {
				return nil, err
			}

			return commandResponseWithServer, nil
		}

		client.Logger.Printf("Batch size is 0, command will be executed in all servers in parallel...")
		commandResponseWithServer, err := ParallelBatchExecute(command, client.SSHConConfig, client.TargetServers, client.Logger)

		if err != nil {
			return nil, err
		}

		return commandResponseWithServer, nil

	} else {

		client.Logger.Printf("Preparing to execute command serially on all servers one by one...")
		commandResponseWithServer, err := SerialExecute(command, client.SSHConConfig, client.TargetServers, client.Logger)

		if err != nil {
			return nil, err
		} else {
			return commandResponseWithServer, nil
		}
	}
}
