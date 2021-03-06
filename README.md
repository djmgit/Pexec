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

In the above example, the tagKey field's value can be changed to some other valid key and the corresponding tagvalue can be provided
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

**Executing command remotely on custom provided servers**

PExec also allows you to provide host IPs manually for remote command execution.

```
package main

import (
	pexec "github.com/djmgit/pexec/lib"
	"fmt"
)

func main() {

	pClient := pexec.PexecClient{
		TargetServers : []pexec.Server{
			{
				Host: "52.87.231.249",
				Port: 22,
			},
			{
				Host: "35.174.213.9",
				Port: 22,
			},
		},
		Parallel: true,
		BatchSize: 0,
		User: "ubuntu",
		Provider: "CUSTOM",
		KeyPath: "<path_to_ssh_key_trusted_by_the_target_insatnces>",
	}

	response, err := pClient.Run("echo 'Hello World'")

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

In order to execute command on custom list of hosts, the list of hosts has to be provided as the ```TargetServers``` field. The list consists of
structs of type ```pexec.Server``` which has only two fields : ```Host``` and ```Port```. Remaining all the other field have the same meaning.

In both the examples presented above, the field ```Provider``` as of now can be either of ```AWS``` or ```CUSTOM```.

Also it is to be noted that PExec executes commands on the remote servers using ```SSH``` and it allows SSH only via key based authentication and
not password. So you need to provide you ```private key``` path or path to your ```pem``` file in ```keyPath``` field of ```PexecClient``` struct.
Also the key you are providing should be trusted by the target hosts. If you omit the ```keyPath``` field then the path to the current users
SSH RSA key will be used.

### Using PExec as a CLI based tool

**Executing command remotely on all the instances of an AWS Autoscaling group using automatic service discovery**

```
./pexec -provider=AWS -port=22 \
> -asg=[aws_asg_name] -region=[aws_region] \
> --access_key_id=[aws_access_key_id] -secret_access_key=[aws_secret_access_key] \
> -key=[path_to_ssh_key] \
>  -user=ubuntu -cmd="echo 'Hello World'"
```

If you want to provide custom tag key and tag value for filtering instances, you can do so by using ```tag_key``` and ```tag_value``` parameters.

**Executing command remotely on custom provided servers**

```
./pexec -provider=CUSTOM \
> -key=[path_to_ssh_key] \
> -user=ubuntu \
> -cmd="echo 'Hello World'" \
> -servers=52.87.231.249:22,35.174.213.9:22
```

## Installing PExec for using it as a CLI tool

### Using prebuilt binary

You can use the prebuilt binary from github :

- Download the compressed prebuilt binary from <a href="https://github.com/djmgit/Pexec/releases/download/v0.0.1/pexec-v0.0.1.tar.xz">here</a>
- Extract the binary using ```tar -xvf pexec.tar.xz```
- Copy the binary to your path for example ```/usr/local/bin```
- Open terminal and execute ```pexec -h```, it should print all the relevant options.

### Building PExec from source

you can build PExec from source:

- Clone this repository or download the code from <a href="https://github.com/djmgit/Pexec/archive/v0.0.1.tar.gz">here</a>
- Extract the source if requited and enter into the source root.
- You will require the **go toolchain** for this, preferably version>=1.11
- From the project root execute ```go build -o ./cmd/pexec cmd/pexec.go```
- You will find the binary at cmd/
- Run ```./pexec -h```, it should show you all the relevant options.
- Optionally you can copy it to your path.


## PExec library reference

Given below is the ```PexecClient``` struct you will have to initialise in order to use PExec as a library.
Not all fields are required to be used as they will be assigned a default value.

```
type PexecClient struct {

	TargetServers []Server

	Port int

	ProviderOptions map[string]string

	Provider string

	Parallel bool

	BatchSize int

	User string

	KeyPath string

	Debug bool

	Logger *log.Logger
}
```

- TargetServers is a list of structs of type Server. You will have to initialise this in case you want to execute commands
  remote custom provided hosts. Server is struct is as shown below:
  
  ```
  type Server struct {
	Host string

	Port int
  }
  ```
- Port : Port to use for ssh, this will be used when you are using a provider other than **CUSTOM**. Default will be **22**. For custom provided
  list of servers, you will have to provided port for the individual servers.
  
- ProviderOptions: Options to be used when using a provider (AWS as of now). Options required as of now:
	- region
	- addrType
	- tagKey
	- tagValue
	- accessKeyId
	- secretAccessKey

- Provider : This can be either of **CUSTOM** or **AWS** as of now. If nothing is provided then it will be set to CUSTOM.

- Parallel : If set to **true** then command will be executed in parallel in all the provided/discovered servers or else it will be executed
             sequentially
	 
- BatchSize : If set to non-zero value and if **Parallel** is set to true, then teh list of provided/discovered hosts will be distributed in batches
              of the provided number of hosts and command will be executed parallely on the members of the individual groups. It is to be noted that the
	      groups themselves will be iterated upon sequentially. Please see the example in the **QuickStart** section above for further explaination.
	      
- User : User to be used by PExec for openning SSH connection. Default will be **root**

- KeyPath : The path to SSH key to be used for SSH. It can be a pem file containing the key or the key itself trusted by the target servers.

- Debug : If set to true, PExexc will generate debug statements.

- Logger : This can be custom provided logger. However if Debug is set to true, then a logger will be created and assigned to this field.

Apart from the above fields the struct also contains some other fileds which are used internally by PExec.

## PExec CLI tool refernce

The follow options are provided by the CLI tool:

```
Usage of ./pexec:
  -access_key_id string
    	AWS Access key id
  -addr_type string
    	Command to execute on the servers (default "public_v4")
  -asg string
    	AWS Auto scaling group name
  -batch_size int
    	If more than one, then batches of that many servers will be executed in parallel
  -cmd string
    	Command to execute on the servers
  -debug
    	Enable debugging, by default false
  -key string
    	If not provided then default key path for rsa key will be used - /home/<user>/.ssh/id_rsa
  -parallel
    	If true then commands will be exected in parallel on the discovered or provided servers (default true)
  -port int
    	Port to override for all
  -provider string
    	Provider of servers - can be one of CUSTOM | AWS (default "CUSTOM")
  -region string
    	AWS region if provider is aws, default is us-east-1 (default "us-east-1")
  -secret_access_key string
    	AWS secret access key
  -servers string
    	Server ip and port in format <IP>:<PORT>, multuple values can be separated by ','. If port is not provided then 22 will be used as default SSH port
  -tag_key string
    	tag key name
  -tag_value string
    	tag value
  -user string
    	User which will be used to login to the server (default "root")

```

Most of the options directly refers to a field in the PexecClient struct as already described in the previous section.

## Contributing to PExec
	     
Contribution to PExec are most welcome. If you have found a bug or if you have any feature in mind or if you want to
help in integrating other cloud platforms or other providers/prchestrators like Consul, DCOS, etc, please create an issue
send a PR. I would be extremely happy to review it.
