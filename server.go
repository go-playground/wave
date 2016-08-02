package wave

import "net/rpc"

// Server encapsulates the standard rpc.Server allowing for adding of aditional functions
// or hooking into existing rpc server functions.
type Server struct {
	*rpc.Server
}

// NewServer creates and returns a new instance of 'Server'
func NewServer() *Server {
	return new(Server)
}

// TODO: hook into `Register` and `RegisterName` functions to provide service discovery
// functionality
