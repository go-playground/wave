package wave

import (
	"crypto/tls"
	"sync"
	"time"
)

// ClientPool allows for easy pooling of 'Client' connections
type ClientPool struct {
	*sync.Pool
}

// NewClientPooled returns a new 'Client' instance backed by a pool of connections
func NewClientPooled(serviceName, network, address string, timeout time.Duration, config *tls.Config) *ClientPool {

	return &ClientPool{
		Pool: &sync.Pool{
			New: func() interface{} {

				c, err := NewClient(serviceName, network, address, timeout, config)
				if err != nil {
					return err
				}

				return c
			},
		},
	}
}
