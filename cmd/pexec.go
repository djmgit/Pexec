package main

import (
	"flag"
	"github.com/djmgit/pexec/lib"
	"fmt"
)

func main() {

	accessKeyId := flag.String("access_key_id", "", "AWS Access key id")
	secretAccessKey := flag.String("secret_access_key", "", "AWS secret access key")
	asgName := flag.String("asg", "", "AWS Auto scaling group name")
	tagKey := flag.String("tag_key", "", "tag key name")
	tagValue := flag.String("tag_value", "", "tag value")
	servers := flag.String("servers", "", "Server ip and port in format <IP>:<PORT>, multuple values can be separated by ','. If port is not provided then 22 will be used as default SSH port")
	port := flag.Int("port", 0, "Port to override for all")
	provider := flag.String("provider", "", "Provider of servers - can be one of CUSTOM | AWS")
	parallel := flag.Bool("parallel", true, "If true then commands will be exected in parallel on the discovered or provided servers")
	batchSize := flag.Int("batch_size", 0, "If more than one, then batches of that many servers will be executed in parallel")
	user := flag.String("user", "", "User which will be used to login to the server")
	keyPath := flag.String("user", "", "If not provided then default key path for rsa key will be used - /home/<user>/.ssh/id_rsa")
}


