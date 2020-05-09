package lib

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func PrepareSSHConConfig(user, keyPath string) (*ssh.ClientConfig, error) {

	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
		  ssh.PublicKeys(signer),
		},
	}
}

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
