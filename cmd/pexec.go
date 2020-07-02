package main

import (
	"flag"
	lib "github.com/djmgit/pexec/lib"
	"fmt"
	"err"
	"strings"
)

func main() {

	cmdParams := lib.CmdParams{}

	cmdParams.AccessKeyId = *flag.String("access_key_id", "", "AWS Access key id")
	cmdParams.SecretAccessKey = *flag.String("secret_access_key", "", "AWS secret access key")
	cmdParams.AsgName = *flag.String("asg", "", "AWS Auto scaling group name")
	cmdParams.TagKey = *flag.String("tag_key", "", "tag key name")
	cmdParams.TagValue = *flag.String("tag_value", "", "tag value")
	cmdParams.Servers = *flag.String("servers", "", "Server ip and port in format <IP>:<PORT>, multuple values can be separated by ','. If port is not provided then 22 will be used as default SSH port")
	cmdParams.Port = *flag.Int("port", 0, "Port to override for all")
	cmdParams.Provider = *flag.String("provider", "", "Provider of servers - can be one of CUSTOM | AWS")
	cmdParams.Parallel = *flag.Bool("parallel", true, "If true then commands will be exected in parallel on the discovered or provided servers")
	cmdParams.BatchSize = *flag.Int("batch_size", 0, "If more than one, then batches of that many servers will be executed in parallel")
	cmdParams.User = *flag.String("user", "", "User which will be used to login to the server")
	cmdParams.KeyPath = *flag.String("user", "", "If not provided then default key path for rsa key will be used - /home/<user>/.ssh/id_rsa")
}

func getDefaults(cmdParams *lib.CmdParams) error {

	if cmdParams.Provider == "CUSTOM" {
		if cmdParams.Servers == "" {
			return error.Error("Please provide target servers")
		}

		servers, err := getServers(cmdParams.Servers)
	}

}

func getServers(servers string) ([]lib.Server) {

	serverStrings := strings.Split(servers, ",")

	servers := make(lib.Server. 0, 0)

	for _, serverString := range serverStrings {
		serverIP := strings.Split(serverString, ":")[0]
		serverPort := strings.Split(serverString, ":")[1]

		servers = append(servers, lib.Server {
			Host: serverIP,
			Port: serverPort
		})
	}

	return servers
}


