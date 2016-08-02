package server

import "net/rpc"

// Server encapsulates the standard rpc.Server allowing for adding of aditional functions
// or hooking into existing rpc server functions.
type Server struct {
	*rpc.Server
}

// New creates and returns a new instance of 'Server'
func New() *Server {
	return new(Server)
}

// TODO: hook into `Register` and `RegisterName` functions to provide service discovery
// functionality
