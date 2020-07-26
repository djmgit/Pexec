package lib

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"strconv"
	"log"
)

func GetSSHSession(config *ssh.ClientConfig, host string, port int) (*ssh.Session, error) {

	client, err := ssh.Dial("tcp", host + ":" + strconv.FormatInt(int64(port), 10), config)

	session, err := client.NewSession()
  	if err != nil {
    	return nil, err
	  }
	  
	  return session, nil
}

func ExecuteCommand(command string, session *ssh.Session, config *ssh.ClientConfig, host string, port int) (*CommandResponse, error) {

	if session == nil {
		var sessionErr error
		session, sessionErr = GetSSHSession(config, host, port)

		if sessionErr != nil {
			return nil, sessionErr
		}
	}

	var StdOutput bytes.Buffer
	var StdError bytes.Buffer

	session.Stdout = &StdOutput
	session.Stderr = &StdError

	cmdErr := session.Run(command)
	if cmdErr != nil {
		return nil, cmdErr
	}

	return &CommandResponse{
		StdOutput: StdOutput.String(),
		StdError: StdError.String(),
	}, nil
}

func SerialExecute(command string, sshClientConfig *ssh.ClientConfig,  targetServers []Server) ([]CommandResponseWithServer, error) {

	commandResponseWithServerList := make([]CommandResponseWithServer, 0, 0)

	for _, server := range(targetServers) {

		CommandResponse, err := ExecuteCommand(command, nil, sshClientConfig, server.Host, server.Port)
		var commandResponseWithServer CommandResponseWithServer

		if err == nil {

			commandResponseWithServer = CommandResponseWithServer{

				Host: server.Host,
				CommandResponse: *CommandResponse,
			}
		} else {
			commandResponseWithServer = CommandResponseWithServer{

				Host: server.Host,
				Err: err.Error(),
			}
		}

		commandResponseWithServerList = append(commandResponseWithServerList, commandResponseWithServer)
	}

	return commandResponseWithServerList, nil
}

func ParallelBatchExecute(command string, sshClientConfig *ssh.ClientConfig, targetServers []Server, logger *log.Logger) ([]CommandResponseWithServer, error) {

	commandResponseWithServerChan := make(chan CommandResponseWithServer)
	done := make(chan string)

	for _, server := range(targetServers) {

		go func(localServer Server) {

			CommandResponse, err := ExecuteCommand(command, nil, sshClientConfig, localServer.Host, localServer.Port)
			var commandResponseWithServer CommandResponseWithServer

			if err == nil {

				commandResponseWithServer = CommandResponseWithServer{

					Host: localServer.Host,
					CommandResponse: *CommandResponse,
				}
			} else {
				commandResponseWithServer = CommandResponseWithServer{

					Host: localServer.Host,
					Err: err.Error(),
				}
			}

			select {
				case commandResponseWithServerChan <- commandResponseWithServer:
				case <- done:
					return
			}
			//commandResponseWithServerChan <- commandResponseWithServer

		}(server)
	}

	commandResponseWithServer := make([]CommandResponseWithServer, 0, 0)

	for _, _ = range targetServers {
		individualServerResponse := <-commandResponseWithServerChan
		commandResponseWithServer = append(commandResponseWithServer, individualServerResponse)
	}

	//time.Sleep(client.TimeOut * time.Second)
	close(done)
	close(commandResponseWithServerChan)

	return commandResponseWithServer, nil
}

func BatchExecuter(command string, sshClientConfig *ssh.ClientConfig, targetServers []Server, batchSize int, logger *log.Logger) ([]CommandResponseWithServer, error) {

	commandResponseAllBatches := make([]CommandResponseWithServer, 0, 0)

	index := 0

	for index < len(targetServers) {
		serversBatch := targetServers[index: index + batchSize]

		commandResponseWithServer, err := ParallelBatchExecute(command, sshClientConfig, serversBatch, logger)

		if err != nil {
			return nil, err
		}

		commandResponseAllBatches = append(commandResponseAllBatches, commandResponseWithServer...)

		index += batchSize
	}

	return commandResponseAllBatches, nil
}
