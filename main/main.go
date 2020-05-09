
package main

import (
  "bytes"
  "golang.org/x/crypto/ssh"
  "fmt"
  "io/ioutil"
  "os"
)

func main() {
  pk, _ := ioutil.ReadFile(os.Getenv("HOME") + "/aws_keys/admin2.pem")
  signer, err := ssh.ParsePrivateKey(pk)

  if err != nil {
    panic(err)
  }

  config := &ssh.ClientConfig{
	User: "ubuntu",
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    Auth: []ssh.AuthMethod{
      ssh.PublicKeys(signer),
    },
  }

  client, err := ssh.Dial("tcp", "54.175.61.187:22", config)
  
  if err != nil {
    panic("Failed to dial: " + err.Error())
  }

  // Each ClientConn can support multiple interactive sessions,
  // represented by a Session.
  session, err := client.NewSession()
  if err != nil {
    panic("Failed to create session: " + err.Error())
  }
  defer session.Close()

  // Once a Session is created, you can execute a single command on
  // the remote side using the Run method.
  var b bytes.Buffer
  session.Stdout = &b
  if err := session.Run("ls"); err != nil {
    panic("Failed to run: " + err.Error())
  }
  fmt.Println(b.String())
}
