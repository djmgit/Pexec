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
