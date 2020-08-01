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
