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
