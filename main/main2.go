package main

import (
	pexec "github.com/djmgit/pexec/lib"
	"fmt"
)

func main() {

	pClient := pexec.PexecClient{
		TargetServers : []pexec.Server{
			{
				Host: "52.72.75.225",
				Port: 22,
			},
			{
				Host: "100.25.139.36",
				Port: 22,
			},
		},

		Parallel: true,

		User: "ubuntu",

		KeyPath: "/home/deep/aws_keys/admin2.pem",
	}

	response, err := pClient.Run("ls")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, commandResponse := range response {
			fmt.Println(commandResponse.Host)
			fmt.Println(commandResponse.CommandResponse.StdOutput)
		}
	}


}
