package lib

type Server struct {
	Host string

	Port int
}

// struct to hold response of executed command
type CommandResponse struct {
	StdOutput string

	StdError string
}

type CommandResponseWithServer struct {
	Host string

	Err string

	CommandResponse CommandResponse
}

type CmdParams struct {
	AccessKeyId string

	SecretAccessKey string

	AsgName string

	TagKey string

	TagValue string

	Servers string

	Command string

	Port int

	Provider string

	Parallel bool

	BatchSize int

	User string

	KeyPath string

	TargetServers []Server

	Region string

	AddrType string

	Debug bool

}
