package electrum

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// ClientVersion identifies the client version/name to the remote server
	ClientVersion = "go-electrum1.0"

	// ProtocolVersion identifies the support protocol version to the remote server
	ProtocolVersion = "1.4"

	nl = byte('\n')
)

var (
	// DebugMode provides debug output on communications with the remote server if enabled.
	DebugMode bool

	// ErrServerConnected throws an error if remote server is already connected.
	ErrServerConnected = errors.New("server is already connected")

	// ErrServerShutdown throws an error if remote server has shutdown.
	ErrServerShutdown = errors.New("server has shutdown")

	// ErrTimeout throws an error if request has timed out
	ErrTimeout = errors.New("request timeout")

	// ErrNotImplemented throws an error if this RPC call has not been implemented yet.
	ErrNotImplemented = errors.New("RPC call is not implemented")

	// ErrDeprecated throws an error if this RPC call is deprecated.
	ErrDeprecated = errors.New("RPC call has been deprecated")
)

// Transport provides interface to server transport.
type Transport interface {
	SendMessage([]byte) error
	Responses() <-chan []byte
	Errors() <-chan error
	Close() error
}

type container struct {
	content []byte
	err     error
}

type ClientOptions struct {
	ConnTimeout time.Duration
	ReqTimeout  time.Duration
}

// Client stores information about the remote server.
type Client struct {
	transport Transport
	opts      *ClientOptions

	handlers     map[uint64]chan *container
	handlersLock sync.RWMutex

	pushHandlers     map[string][]chan *container
	pushHandlersLock sync.RWMutex

	Error chan error
	quit  chan struct{}

	nextID uint64
}

// NewServer initialize a new remote server.
func NewClient(opts *ClientOptions) *Client {
	s := &Client{
		handlers:     make(map[uint64]chan *container),
		pushHandlers: make(map[string][]chan *container),

		Error: make(chan error),
		quit:  make(chan struct{}),

		opts: opts,
	}

	return s
}

// ConnectTCP connects to the remote server using TCP.
// 		TODO: bring timeout with ctx
func (s *Client) ConnectTCP(ctx context.Context, addr string) error {
	if s.transport != nil {
		return ErrServerConnected
	}

	transport, err := NewTCPTransport(ctx, addr)
	if err != nil {
		return err
	}

	s.transport = transport
	go s.listen()

	return nil
}

// ConnectSSL connects to the remote server using SSL.
func (s *Client) ConnectSSL(addr string, config *tls.Config) error {
	if s.transport != nil {
		return ErrServerConnected
	}

	transport, err := NewSSLTransport(addr, config, s.opts.ConnTimeout)
	if err != nil {
		return err
	}

	s.transport = transport
	go s.listen()

	return nil
}

type apiErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *apiErr) Error() string {
	return fmt.Sprintf("errNo: %d, errMsg: %s", e.Code, e.Message)
}

type response struct {
	ID     uint64  `json:"id"`
	Method string  `json:"method"`
	Error  *apiErr `json:"error"`
}

func (s *Client) listen() {
	for {
		if s.IsShutdown() {
			break
		}
		if s.transport == nil {
			break
		}
		select {
		case <-s.quit:
			break
		case err := <-s.transport.Errors():
			s.Error <- err
			s.Shutdown()
		case bytes := <-s.transport.Responses():
			result := &container{
				content: bytes,
			}

			msg := &response{}
			err := json.Unmarshal(bytes, msg)
			if err != nil {
				if DebugMode {
					log.Printf("Unmarshal received message failed: %v", err)
				}
				result.err = fmt.Errorf("Unmarshal received message failed: %v", err)
			} else if msg.Error != nil {
				result.err = msg.Error
			}

			if len(msg.Method) > 0 {
				s.pushHandlersLock.RLock()
				handlers := s.pushHandlers[msg.Method]
				s.pushHandlersLock.RUnlock()

				for _, handler := range handlers {
					select {
					case handler <- result:
					default:
					}
				}
			}

			s.handlersLock.RLock()
			c, ok := s.handlers[msg.ID]
			s.handlersLock.RUnlock()

			if ok {
				c <- result
			}
		}
	}
}

func (s *Client) listenPush(method string) <-chan *container {
	c := make(chan *container, 1)
	s.pushHandlersLock.Lock()
	s.pushHandlers[method] = append(s.pushHandlers[method], c)
	s.pushHandlersLock.Unlock()

	return c
}

type request struct {
	ID     uint64        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func (s *Client) request(method string, params []interface{}, v interface{}) error {
	select {
	case <-s.quit:
		return ErrServerShutdown
	default:
	}

	msg := request{
		ID:     atomic.AddUint64(&s.nextID, 1),
		Method: method,
		Params: params,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	bytes = append(bytes, nl)

	err = s.transport.SendMessage(bytes)
	if err != nil {
		s.Shutdown()
		return err
	}

	c := make(chan *container, 1)

	s.handlersLock.Lock()
	s.handlers[msg.ID] = c
	s.handlersLock.Unlock()

	var resp *container
	select {
	case resp = <-c:
	case <-time.After(s.opts.ReqTimeout):
		return ErrTimeout
	}

	if resp.err != nil {
		return resp.err
	}

	s.handlersLock.Lock()
	delete(s.handlers, msg.ID)
	s.handlersLock.Unlock()

	if v != nil {
		err = json.Unmarshal(resp.content, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Client) Shutdown() {
	if !s.IsShutdown() {
		close(s.quit)
	}
	if s.transport != nil {
		_ = s.transport.Close()
	}
	s.transport = nil
	s.handlers = nil
	s.pushHandlers = nil
}

func (s *Client) IsShutdown() bool {
	select {
	case <-s.quit:
		return true
	default:
	}
	return false
}
