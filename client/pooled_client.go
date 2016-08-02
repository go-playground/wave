package client

import (
	"crypto/tls"
	"sync"
	"time"
)

// ClientPool allows for easy pooling of 'Client' connections
type ClientPool struct {
	pool *sync.Pool
}

// NewPooled returns a new 'Client' instance backed by a pool of connections
func NewPooled(network, address string, timeout time.Duration, config *tls.Config) *ClientPool {

	return &ClientPool{
		pool: &sync.Pool{
			New: func() interface{} {

				c, err := New(network, address, timeout, config)
				if err != nil {
					return err
				}

				return c
			},
		},
	}
}
