package wave

import (
	"crypto/tls"
	"net"
	"net/rpc"
	"time"
)

// Client encapsulates the standard rpc.Client allowing for adding of aditional functions
// or hooking into existing rpc client functions.
type Client struct {
	serviceName string
	*rpc.Client
}

// NewClient creates and returns a new instance of 'Client'
// config is optional, if 'nil' it is ignored
func NewClient(serviceName, network, address string, timeout time.Duration, config *tls.Config) (*Client, error) {

	if config == nil {

		conn, err := net.DialTimeout(network, address, timeout)
		if err != nil {
			return nil, err
		}

		return &Client{
			serviceName: serviceName,
			Client:      rpc.NewClient(conn),
		}, nil
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	conn, err := tls.DialWithDialer(dialer, network, address, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		serviceName: serviceName,
		Client:      rpc.NewClient(conn),
	}, nil
}

// ServiceName returns the Client's service name
// this allow for that to be hidden away during implementation
// if desired.
func (c *Client) ServiceName() string {
	return c.serviceName
}
