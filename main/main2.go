package main

import (
	pexec "github.com/djmgit/pexec/lib"
	"fmt"
)

func main() {

	providerOptions := map[string]string{
		"region": "us-east-1",
		"addrType": "public_v4",
		"tagKey": "aws:autoscaling:groupName",
		"tagValue": "asg_1",
		"accessKeyId": "",
		"secretAccessKey": "",
	}

	pClient := pexec.PexecClient{
		TargetServers : []pexec.Server{
			{
				Host: "52.72.75.225",
				Port: 22,
			},
			{
				Host: "52.4.32.118",
				Port: 22,
			},
		},

		Parallel: true,

		User: "ubuntu",

		Provider: "AWS",

		ProviderOptions: providerOptions,

		KeyPath: "/home/deep/aws_keys/admin2.pem",
	}

	response, err := pClient.Run("ls /usr/bin")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("In main...")
		for _, commandResponse := range response {
			fmt.Println(commandResponse.Host + " : " + commandResponse.CommandResponse.StdOutput)
		}
	}


}
