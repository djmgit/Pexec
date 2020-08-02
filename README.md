# PExec

PExec is a lightweight CLI based tool and libray for executing commands on remote servers. It allows you to execute
commands on a list of remote hosts both sequentially - one after the other, parallely - execute the command on all the
hosts simultaneously all at once to speed up the process and lastly to group the list of servers into batches and execute
the command on the servers present in those individual batches in parallel.

Additionally, PExec also allows you to **discover**  instances from cloud (as of now only **AWS**) and execute commands on them.
Once PExec is provieded with the neccessary constraints and parameneters - for example : auto scaling group name in AWS or a tag key
and its corresponding value, PExec will use those provided information to discover the IPs of the desired instances, connect to
them over SSH on the desired port and execute the given command.

As of now PExec only provided integration with AWS for auto discovery of instances.

## QuickStart

### Using PExec as a library

** Executing command remotely on all the instances of an AWS Autoscaling group using automatic service discovery **

```
package main

import (
	pexec "github.com/djmgit/pexec/lib"
	"fmt"
)

func main() {

	providerOptions := map[string]string{
		"region": "aws_region_name",
		"addrType": "public_v4",
		"tagKey": "aws:autoscaling:groupName",
		"tagValue": "<aws_asg_name>",
		"accessKeyId": "<aws_access_key_id>",
		"secretAccessKey": "<aws_secret_access_key>",
	}

	pClient := pexec.PexecClient{
		Parallel: true,
		BatchSize: 0,
		User: "ubuntu",
		Port: 22,
		Provider: "AWS",
		ProviderOptions: providerOptions,
		KeyPath: "<path_to_ssh_key_trusted_by_the_target_insatnces>",
	}

	response, err := pClient.Run("echo 'Hellow World'")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("In main...")
		for _, commandResponse := range response {
			fmt.Println(commandResponse.Host + " : " + commandResponse.CommandResponse.StdOutput)
		}
	}
}

```
