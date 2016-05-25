package wrap

import (
	"io"
	"net/rpc"
	"sync"
	"time"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/go-playground/wave/client"
)

// RetryReconnectEndpoint is interface needed to comply with wrapping
// with Rety and Reconnect logic.
type RetryReconnectEndpoint interface {
	client.Endpoint

	// SetClient is called after a successful connection is established/reestablished
	// NOTE: should close the existing client in this method if not nil prior to
	// updating to the new one passed in.
	SetClient(*rpc.Client)

	// NewClient calls the RetryReconnectEndpoint's NewClient method
	// wrapping it with retry/reconnect logic
	NewClient() (*rpc.Client, error)
}

// RetryReconnect wraps the given RetryReconnectEndpoint endpoint and automatically
// handles logic to reconnect and retry
func RetryReconnect(endpoint RetryReconnectEndpoint, retryDuration time.Duration) (e client.Endpoint, err error) {

	if !log.HasHandlers() {
		cLog := console.New()
		cLog.SetBuffersAndWorkers(3, 3)
		log.RegisterHandler(cLog, log.AllLevels...)
	}

	rr := &retryReconnect{
		RetryReconnectEndpoint: endpoint,
		clientMutex:            new(sync.RWMutex),
		reconnectMutex:         new(sync.RWMutex),
		retryDuration:          retryDuration,
		isDisconnected:         true,
	}

	e = rr

	_, err = rr.NewClient()
	if err != nil {
		return
	}

	return
}

type retryReconnect struct {
	RetryReconnectEndpoint
	clientMutex    *sync.RWMutex
	reconnectMutex *sync.RWMutex
	retryDuration  time.Duration
	isDisconnected bool // defaults to true
}

var _ client.Endpoint = &retryReconnect{}
var _ RetryReconnectEndpoint = &retryReconnect{}

// NewClient calls the RetryReconnectEndpoint's NewClient method
// wrapping it with retry/reconnect logic
func (r *retryReconnect) NewClient() (c *rpc.Client, err error) {

	r.reconnectMutex.Lock()
	r.isDisconnected = true
	r.reconnectMutex.Unlock()

	// a bunch of calls could have gotten here if server was busy, lets double check we're still disconnected
	// right after lock is released
	r.clientMutex.Lock()
	defer r.clientMutex.Unlock()

	// check if still disconnected
	r.reconnectMutex.RLock()
	if !r.isDisconnected {
		r.reconnectMutex.RUnlock()
		return nil, err
	}
	r.reconnectMutex.RUnlock()

	for i := 1; i < 4; i++ {

		c, err = r.RetryReconnectEndpoint.NewClient()
		if err != nil {
			time.Sleep(r.retryDuration)
			log.WithFields(log.F("service", r.ServiceMethod()), log.F("attempt", i)).Warn("Attempting to established/re-establish RPC Connection")
			continue
		}

		r.SetClient(c)

		r.reconnectMutex.Lock()
		r.isDisconnected = false
		r.reconnectMutex.Unlock()

		return
	}

	log.WithFields(log.F("service", r.ServiceMethod()), log.F("err", err)).Alert("RPC Connection could not be established/re-established")

	go func() {

		r.clientMutex.Lock()
		defer r.clientMutex.Unlock()

		for {

			var client *rpc.Client
			var err2 error

			time.Sleep(r.retryDuration)
			client, err2 = r.RetryReconnectEndpoint.NewClient()
			if err2 != nil {
				time.Sleep(r.retryDuration)
				continue
			}

			r.SetClient(client)

			log.WithFields(log.F("service", r.ServiceMethod())).Notice("RPC connection re-established")
			r.reconnectMutex.Lock()
			r.isDisconnected = false
			r.reconnectMutex.Unlock()

			break
		}
	}()

	err = rpc.ErrShutdown
	return
}

// Call calls the RetryReconnectEndpoint's Call method
// wrapping it with retry/reconnect logic
func (r *retryReconnect) Call(args interface{}, reply interface{}) (err error) {

	// check if disconnected
	r.reconnectMutex.RLock()
	if r.isDisconnected {
		r.reconnectMutex.RUnlock()
		return rpc.ErrShutdown
	}
	r.reconnectMutex.RUnlock()

RETRY:
	// make rpc call
	r.clientMutex.RLock()
	err = r.RetryReconnectEndpoint.Call(args, reply)
	r.clientMutex.RUnlock()

	// if error indicates a disconnect of some sort, try and reconnect
	if err != nil && err == rpc.ErrShutdown || err == io.EOF || err == io.ErrUnexpectedEOF {

		_, err = r.NewClient()
		if err != nil {
			return
		}

		goto RETRY
	}

	return
}

// Go calls the RetryReconnectEndpoint's Call method
// wrapping it with retry/reconnect logic
func (r *retryReconnect) Go(args interface{}, reply interface{}, done chan *rpc.Call) (c *rpc.Call) {

	if done == nil {
		done = make(chan *rpc.Call, 10)
	}

	c = &rpc.Call{
		ServiceMethod: r.ServiceMethod(),
		Args:          args,
		Reply:         reply,
		Done:          done,
	}

	// check if disconnected
	r.reconnectMutex.RLock()
	if r.isDisconnected {
		r.reconnectMutex.RUnlock()

		c.Error = rpc.ErrShutdown

		go func() {
			c.Done <- c
		}()
		return
	}

	r.reconnectMutex.RUnlock()

	dc := make(chan *rpc.Call, cap(done))
	// make rpc call
	r.clientMutex.RLock()
	c2 := r.RetryReconnectEndpoint.Go(args, reply, dc)
	r.clientMutex.RUnlock()

	go func() {
		res := <-c2.Done

		if res.Error != nil && res.Error == rpc.ErrShutdown || res.Error == io.EOF || res.Error == io.ErrUnexpectedEOF {
			_, res.Error = r.NewClient()
			if res.Error != nil {
				c.Done <- res
				return
			}

			// RETRY
			r.clientMutex.RLock()
			c2 = r.RetryReconnectEndpoint.Go(args, reply, dc)
			r.clientMutex.RUnlock()

			res = <-c2.Done
		}

		c.Done <- res
	}()

	return
}
