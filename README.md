# PExec

PExec is a lightweight CLI based tool and libray for executing commands on remote servers. It allows you to execute
commands on a list of remote hosts both sequentially - one after the other, parallely - execute the command on all the
hosts simultaneously all at once to speed up the process and lastly to group the list of servers into batches and execute
the command on the servers present in those individual batches in parallel.

Additionally, PExec also allows you to **discover**  instances from cloud (as of now only **AWS**) and execute commands on them.
Once PExec is provieded with the neccessary constraints and parameneters - for example : auto scaling group name in AWS or a tag key
and its corresponding value, PExec will use those provided information to discover the IPs of the desired instances, connect to
them over SSH on the desired port and execute the given command.

PExec can be used for automating usual tasks, or as a part of a bigger tool or CI/CD where you want to execute commands remotely
and you want to make it fast.

As of now PExec only provided integration with AWS for auto discovery of instances.

## QuickStart

### Using PExec as a library

**Executing command remotely on all the instances of an AWS Autoscaling group using automatic service discovery**

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

In the above example, the tagKey parameter's value can be changed to some other valid key and the corresponding tagvalue can be provided
to filter desired instances.
In the above Pexec client configuration, parallel is set to ```true``` and batchSize is set to ```0``` which means, command will be executed
simultaneously on all the discovered servers. Setting the batchSize to non-zero, for example ```2``` would distribute the discovered servers
into groups of ```2``` servers and the command would be executed in parallel on all the individual servers of a group but each group will be processed
sequentially.

For example : if there were 4 servers and batchSize would be set to 2. then PExec would distribute the 4 discovered servers into 2 grousp of two
servers each, it would iterate sequentially over each group and execute the command on both the servers of the current group in parallel.
This can be usefull when you dont want to execute your command on all the servers at once but still want to speed up the process. For example you
want to restart apache or redis slaves, but you dont want to take down all of them at once but in batches.

Right now there is no way to configure a delay between the processing of two groups, but that will be added soon.

In case you do not want PExec to execute the command in parallel on all the servers, you can set ```Parallel``` to ```false```. Doing that
will cause PExec to iterate sequentially over the list of discovered servers and execute the provided command.

