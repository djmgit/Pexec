package lib

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"strconv"
	"log"
)

// Helper function to establish a new SSH sessiont o target server
func GetSSHSession(config *ssh.ClientConfig, host string, port int) (*ssh.Session, error) {

	client, err := ssh.Dial("tcp", host + ":" + strconv.FormatInt(int64(port), 10), config)

	session, err := client.NewSession()
  	if err != nil {
    	return nil, err
	  }
	  
	  return session, nil
}

// Function responsible for actually executing the command reotely
// on the given server
func ExecuteCommand(command string, session *ssh.Session, config *ssh.ClientConfig, host string, port int, logger *log.Logger) (*CommandResponse, error) {

	logger.Printf("Execcuting command on %s...\n", host)

	// Check if there is already an existing session
	// TODO : persist session
	if session == nil {
		logger.Printf("No existing session for %s, openning new SSH session...\n", host)
		var sessionErr error
		session, sessionErr = GetSSHSession(config, host, port)

		if sessionErr != nil {
			logger.Printf("Failed to open session for %s, error : %s\n", host, sessionErr.Error())
			return nil, sessionErr
		}
	}

	var StdOutput bytes.Buffer
	var StdError bytes.Buffer

	session.Stdout = &StdOutput
	session.Stderr = &StdError

	cmdErr := session.Run(command)
	if cmdErr != nil {
		logger.Printf("Command execution failed with error %s on %s\n", cmdErr.Error(), host)
		return nil, cmdErr
	}

	logger.Printf("Command execution successfull on %s\n", host)

	return &CommandResponse{
		StdOutput: StdOutput.String(),
		StdError: StdError.String(),
	}, nil
}

// Function to execute command serially on the provided list of servers
func SerialExecute(command string, sshClientConfig *ssh.ClientConfig,  targetServers []Server, logger *log.Logger) ([]CommandResponseWithServer, error) {

	commandResponseWithServerList := make([]CommandResponseWithServer, 0, 0)

	// Iterate over the list if hosts and execute command remotely
	for _, server := range(targetServers) {

		CommandResponse, err := ExecuteCommand(command, nil, sshClientConfig, server.Host, server.Port, logger)
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

// Fucntion to execute commands parallely on list of given servers
func ParallelBatchExecute(command string, sshClientConfig *ssh.ClientConfig, targetServers []Server, logger *log.Logger) ([]CommandResponseWithServer, error) {

	commandResponseWithServerChan := make(chan CommandResponseWithServer)
	done := make(chan string)

	for _, server := range(targetServers) {

		// For each of the servers in the target Servers list, spawn a goroutine
		// for remotely executing the command on that server.
		go func(localServer Server) {

			CommandResponse, err := ExecuteCommand(command, nil, sshClientConfig, localServer.Host, localServer.Port, logger)
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

			// Send result to the above created channel
			// Done channel is used to terminate all the goroutines whenever requirted
			select {
				case commandResponseWithServerChan <- commandResponseWithServer:
				case <- done:
					return
			}
			//commandResponseWithServerChan <- commandResponseWithServer

		}(server)
	}

	commandResponseWithServer := make([]CommandResponseWithServer, 0, 0)

	// Iterate over the number of servers present in targetServers list
	for _, _ = range targetServers {

		// Pull out a result from the channel
		individualServerResponse := <-commandResponseWithServerChan
		commandResponseWithServer = append(commandResponseWithServer, individualServerResponse)
	}

	// Close the done channel
	close(done)

	// Close thse main response channel
	close(commandResponseWithServerChan)

	return commandResponseWithServer, nil
}

func BatchExecuter(command string, sshClientConfig *ssh.ClientConfig, targetServers []Server, batchSize int, logger *log.Logger) ([]CommandResponseWithServer, error) {

	commandResponseAllBatches := make([]CommandResponseWithServer, 0, 0)

	index := 0
	batchNumber := 1

	for index < len(targetServers) {
		logger.Printf("Executing Batch #%d...\n", batchNumber)
		serversBatch := targetServers[index: index + batchSize]

		commandResponseWithServer, err := ParallelBatchExecute(command, sshClientConfig, serversBatch, logger)

		if err != nil {
			return nil, err
		}

		commandResponseAllBatches = append(commandResponseAllBatches, commandResponseWithServer...)

		index += batchSize
		batchNumber += 1
	}

	return commandResponseAllBatches, nil
}
