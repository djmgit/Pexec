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
				Host: "52.4.32.118",
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
		fmt.Println("In main...")
		for _, commandResponse := range response {
			fmt.Println(commandResponse.Host + " : " + commandResponse.CommandResponse.StdOutput)
		}
	}


}
