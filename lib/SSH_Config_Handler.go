package lib

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

// Helper function for preparing SSH config for launching SSH session on
// target server
func PrepareSSHConConfig(user, keyPath string) (*ssh.ClientConfig, error) {

	signer, err := getPublicKeyAuth(keyPath)

	if err != nil {
		return nil, err
	}

	// Creating config for SSH
	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
		  signer,
		},
	}

	return config, nil
}

// Configure using existing SSH private key
func getPublicKeyAuth(file string) (ssh.AuthMethod, error) {
    buffer, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }

    key, err := ssh.ParsePrivateKey(buffer)
    if err != nil {
        return nil, err
    }
    return ssh.PublicKeys(key), nil
}
