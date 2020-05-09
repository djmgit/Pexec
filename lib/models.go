package lib

type Server struct {
	Host string

	Port int
}

type CommandResponse struct {
	StdOutput string

	StdError string
}
