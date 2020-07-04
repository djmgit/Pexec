package main

import (
	"flag"
	lib "github.com/djmgit/pexec/lib"
	"fmt"
	"strconv"
	"strings"
	"errors"
	"os"
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
	cmdParams.Provider = *flag.String("provider", "CUSTOM", "Provider of servers - can be one of CUSTOM | AWS")
	cmdParams.Parallel = *flag.Bool("parallel", true, "If true then commands will be exected in parallel on the discovered or provided servers")
	cmdParams.BatchSize = *flag.Int("batch_size", 0, "If more than one, then batches of that many servers will be executed in parallel")
	cmdParams.User = *flag.String("user", "root", "User which will be used to login to the server")
	cmdParams.KeyPath = *flag.String("key", "", "If not provided then default key path for rsa key will be used - /home/<user>/.ssh/id_rsa")
	cmdParams.Command = *flag.String("cmd", "", "Command to execute on the servers")
	cmdParams.Region = *flag.String("region", "us-east-1", "AWS region if provider is aws, default is us-east-1")
	cmdParams.AddrType = *flag.String("addr_type", "public_v4", "Command to execute on the servers")


	err := getDefaults(&cmdParams)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	pexecClient := lib.PexecClient{
		Provider: cmdParams.Provider,
		Port: cmdParams.Port,
		Parallel: cmdParams.Parallel,
		BatchSize: cmdParams.BatchSize,
		User: cmdParams.User,
		KeyPath: cmdParams.KeyPath,
	}

	if cmdParams.Provider == "CUSTOM" {
		pexecClient.TargetServers = cmdParams.TargetServers
	}

	if cmdParams.Provider == "AWS" {
		pexecClient.ProviderOptions = map[string]string{
			"region": cmdParams.Region,
			"addrType": cmdParams.AddrType,
			"tagKey": cmdParams.TagKey,
			"tagVal": cmdParams.TagValue,
			"accessKeyId": cmdParams.AccessKeyId,
			"secretAccessKey": cmdParams.SecretAccessKey,
		}
	}

	response, err := pexecClient.Run(cmdParams.Command)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	for _, commandResponse := range response {
		fmt.Println(commandResponse.Host + " : " + commandResponse.CommandResponse.StdOutput)
	}
}

func getDefaults(cmdParams *lib.CmdParams) error {

	if cmdParams.Provider == "CUSTOM" {
		if cmdParams.Servers == "" {
			return errors.New("Please provide target servers")
		}

		servers := getServers(cmdParams.Servers)
		cmdParams.TargetServers = servers

		if cmdParams.Command == "" {
			return errors.New("Please provide a command to execute on remote server")
		}
	} else if (cmdParams.Provider == "AWS") {
		if cmdParams.AccessKeyId == "" || cmdParams.SecretAccessKey == "" {
			return errors.New("Please provide aws access creds")
		}

		if cmdParams.TagKey != "" {
			if cmdParams.TagValue == "" {
				return errors.New("Please provide Tag value")
			}
		} else if cmdParams.AsgName != "" {
			cmdParams.TagKey = "aws:autoscaling:groupName"
			cmdParams.TagValue = cmdParams.AsgName
		} else {

			return errors.New("Please provide either ASG name or Tag key,value pair to discover aws instances")
		}
	}
	return nil
}

func getServers(serversList string) ([]lib.Server) {

	serverStrings := strings.Split(serversList, ",")

	servers := make([]lib.Server, 0, 0)

	for _, serverString := range serverStrings {
		serverIP := strings.Split(serverString, ":")[0]
		serverPort, _ := strconv.ParseInt(strings.Split(serverString, ":")[1], 10, 64)

		servers = append(servers, lib.Server {
			Host: serverIP,
			Port: int(serverPort),
		})
	}

	return servers
}


