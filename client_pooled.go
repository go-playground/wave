package wave

import (
	"crypto/tls"
	"sync"
	"time"
)

// ClientPool allows for easy pooling of 'Client' connections
type ClientPool struct {
	pool *sync.Pool
}

// NewClientPooled returns a new 'Client' instance backed by a pool of connections
func NewClientPooled(serviceName, network, address string, timeout time.Duration, config *tls.Config) *ClientPool {

	return &ClientPool{
		pool: &sync.Pool{
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

// Get Retrieves and existing connection or creates a new one
// if none exist
func (p *ClientPool) Get() (*Client, error) {

	c := p.pool.Get()
	if err, ok := c.(error); ok {
		return nil, err
	}

	return c.(*Client), nil
}

// Put puts v back into the pool
func (p *ClientPool) Put(v interface{}) {
	p.pool.Put(v)
}
