package lib

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

func GetSSHSession(config *ssh.ClientConfig, host string, port int) (*ssh.Session, error) {

	client, err := ssh.Dial("tcp", host + ":" + string(port), config)

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

func SerialExecute(command string, sshClientConfig *ssh.ClientConfig,  targetServers []Server) (*[]CommandResponseWithServer, error) {

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

	return &commandResponseWithServerList, nil
}

func ParallelBatchExecute(command string, sshClientConfig *ssh.ClientConfig, targetServers []Server, done <-chan string) (<-chan CommandResponseWithServer, error) {

	commandResponseWithServerChan := make(chan CommandResponseWithServer)
	defer close(commandResponseWithServerChan)

	for _, server := range(targetServers) {

		go func() {

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

			for {
				select {
				case commandResponseWithServerChan <- commandResponseWithServer:
				case <- done:
					return
				}
			}

		}()
	}

	return commandResponseWithServerChan, nil
}
