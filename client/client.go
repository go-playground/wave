package client

import "net/rpc"

// Endpoint is interface for use with Client methods
type Endpoint interface {
	ServiceMethod() string
	Call(args interface{}, reply interface{}) error
	Go(args interface{}, reply interface{}, done chan *rpc.Call) (c *rpc.Call)
}
