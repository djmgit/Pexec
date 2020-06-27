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
}


