package client

import (
	"crypto/tls"
	"net"
	"net/rpc"
	"time"
)

// Client encapsulates the standard rpc.Client allowing for adding of aditional functions
// or hooking into existing rpc client functions.
type Client struct {
	*rpc.Client
}

// New creates and returns a new instance of 'Client'
// config is optional, if 'nil' it is ignored
func New(network, address string, timeout time.Duration, config *tls.Config) (*Client, error) {

	if config == nil {

		conn, err := net.DialTimeout(network, address, timeout)
		if err != nil {
			return nil, err
		}

		return &Client{rpc.NewClient(conn)}, nil
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	conn, err := tls.DialWithDialer(dialer, network, address, config)
	if err != nil {
		return nil, err
	}

	return &Client{rpc.NewClient(conn)}, nil
}
