package kademlia

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/anacrolix/utp"
	"go.uber.org/ratelimit"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"google.golang.org/grpc/credentials"
)

const (
	defaultConnDeadline   = 30 * time.Second
	defaultConnRate       = 500
	defaultMaxPayloadSize = 16 * 1024 * 1024 // 16MB
)

// Network for distributed hash table
type Network struct {
	dht     *DHT              // the distributed hash table
	socket  *utp.Socket       // the server socket for the network
	self    *Node             // local node itself
	limiter ratelimit.Limiter // the rate limit for accept socket
	done    chan struct{}     // network is stopped

	// For secure connection
	secureHelper credentials.TransportCredentials
	connPool     *ConnPool
	connPoolMtx  sync.Mutex

	// for authentication only
	authHelper *AuthHelper
}

// NewNetwork returns a network service
func NewNetwork(dht *DHT, self *Node, secureHelper credentials.TransportCredentials, authHelper *AuthHelper) (*Network, error) {
	s := &Network{
		dht:          dht,
		self:         self,
		done:         make(chan struct{}),
		secureHelper: secureHelper,
		connPool:     NewConnPool(),
		authHelper:   authHelper,
	}
	// init the rate limiter
	s.limiter = ratelimit.New(defaultConnRate)

	addr := fmt.Sprintf("%s:%d", self.IP, self.Port)
	// new the network socket
	socket, err := utp.NewSocket("udp", addr)
	if err != nil {
		return nil, err
	}
	s.socket = socket

	return s, nil
}

// Start the network
func (s *Network) Start(ctx context.Context) error {
	// serve the incoming connection
	go s.serve(ctx)

	return nil
}

// Stop the network
func (s *Network) Stop(ctx context.Context) {
	if s.secureHelper != nil {
		s.connPool.Release()
	}
	// close the socket
	if s.socket != nil {
		if err := s.socket.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("close socket failed")
		}
	}

}

func (s *Network) encodeMesage(mesage *Message) ([]byte, error) {
	// send the response to client
	encoded, err := encode(mesage)
	if err != nil {
		return nil, errors.Errorf("encode response: %w", err)
	}

	return encoded, nil
}

func (s *Network) handleFindNode(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*FindNodeRequest)
	if !ok {
		err := errors.New("invalid FindNodeRequest")
		response := &FindNodeResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		// new a response message
		resMsg := s.dht.newMessage(FindNode, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	// add the sender to local hash table
	s.dht.addNode(ctx, message.Sender)

	// the closest contacts
	closest := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})

	response := &FindNodeResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
		Closest: closest.Nodes,
	}

	// new a response message
	resMsg := s.dht.newMessage(FindNode, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleFindValue(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*FindValueRequest)
	if !ok {
		err := errors.New("invalid FindValueRequest")
		response := &FindValueResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		// new a response message
		resMsg := s.dht.newMessage(FindValue, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	// add the sender to local hash table
	s.dht.addNode(ctx, message.Sender)

	// retrieve the value from local storage
	value, err := s.dht.store.Retrieve(ctx, request.Target)
	if err != nil {
		err = errors.Errorf("store retrieve: %w", err)
		response := &FindValueResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		// new a response message
		resMsg := s.dht.newMessage(FindValue, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	response := &FindValueResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}

	if value != nil {
		// return the value
		response.Value = value
	} else {
		// return the closest contacts
		closest := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})
		response.Closest = closest.Nodes
	}

	// new a response message
	resMsg := s.dht.newMessage(FindValue, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleStoreData(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*StoreDataRequest)
	if !ok {
		err := errors.New("invalid StoreDataRequest")

		response := &StoreDataResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		resMsg := s.dht.newMessage(StoreData, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	log.WithContext(ctx).Debugf("handle store data: %v", message.String())

	// add the sender to local hash table
	s.dht.addNode(ctx, message.Sender)

	// format the key
	key := s.dht.hashKey(request.Data)

	// replication time for key
	replication := time.Now().Add(defaultReplicateTime)
	// store the data to local storage
	if err := s.dht.store.Store(ctx, key, request.Data, replication); err != nil {
		err = errors.Errorf("store the data: %w", err)
		response := &StoreDataResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		resMsg := s.dht.newMessage(StoreData, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	response := &StoreDataResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}

	// new a response message
	resMsg := s.dht.newMessage(StoreData, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handlePing(_ context.Context, message *Message) ([]byte, error) {
	// new a response message
	resMsg := s.dht.newMessage(Ping, message.Sender, nil)
	return s.encodeMesage(resMsg)
}

// handle the connection request
func (s *Network) handleConn(ctx context.Context, rawConn net.Conn) {
	var conn net.Conn
	var err error
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("conn:%s->%s", rawConn.LocalAddr(), rawConn.RemoteAddr()))

	// do secure handshaking
	if s.secureHelper != nil {
		conn, err = NewSecureServerConn(ctx, s.secureHelper, rawConn)
		if err != nil {
			rawConn.Close()
			log.WithContext(ctx).WithError(err).Error("server secure establish failed")
			return
		}
	} else {
		// if peer authentication is enabled
		if s.authHelper != nil {
			authHandshaker, _ := auth.NewServerHandshaker(ctx, s.authHelper, rawConn)
			conn, err = authHandshaker.ServerHandshake(ctx)
			if err != nil {
				rawConn.Close()
				log.WithContext(ctx).WithError(err).Error("server authentication failed")
				return
			}
		} else {
			conn = rawConn
		}
	}

	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// read the request from connection
		request, err := decode(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.WithContext(ctx).WithError(err).Error("read and decode failed")
			return
		}

		var response []byte
		switch request.MessageType {
		case FindNode:
			encoded, err := s.handleFindNode(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle find node request failed")
				return
			}
			response = encoded
		case FindValue:
			// handle the request for finding value
			encoded, err := s.handleFindValue(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle find value request failed")
				return
			}
			response = encoded
		case Ping:
			encoded, err := s.handlePing(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle ping request failed")
				return
			}
			response = encoded
		case StoreData:
			// handle the request for storing data
			encoded, err := s.handleStoreData(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle store data request failed")
				return
			}
			response = encoded
		default:
			log.WithContext(ctx).Errorf("invalid message type: %v", request.MessageType)
			return
		}

		// write the response
		if _, err := conn.Write(response); err != nil {
			log.WithContext(ctx).WithError(err).Error("write failed")
			return
		}
	}
}

// serve the incomming connection
func (s *Network) serve(ctx context.Context) {
	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		// rate limiter for the incomming connections
		s.limiter.Take()

		// accept the incomming connections
		conn, err := s.socket.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.WithContext(ctx).WithError(err).Errorf("socket accept failed, retrying in %v", tempDelay)

				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "closed") {
				return
			}

			log.WithContext(ctx).WithError(err).Error("socket accept failed")
			return
		}

		// handle the connection requests
		go s.handleConn(ctx, conn)
	}
}

// Call sends the request to target and receive the response
func (s *Network) Call(ctx context.Context, request *Message) (*Message, error) {
	var conn net.Conn
	var rawConn net.Conn
	var err error

	remoteAddr := fmt.Sprintf("%s:%d", request.Receiver.IP, request.Receiver.Port)

	// do secure handshaking
	if s.secureHelper != nil {
		s.connPoolMtx.Lock()
		conn, err = s.connPool.Get(remoteAddr)
		if err != nil {
			conn, err = NewSecureClientConn(ctx, s.secureHelper, remoteAddr)
			if err != nil {
				s.connPoolMtx.Unlock()
				return nil, errors.Errorf("client secure establish %q: %w", remoteAddr, err)
			}
			s.connPool.Add(remoteAddr, conn)
		}
		s.connPoolMtx.Unlock()
	} else {
		// dial the remote address with udp network
		rawConn, err = utp.DialContext(ctx, remoteAddr)
		if err != nil {
			return nil, errors.Errorf("dial %q: %w", remoteAddr, err)
		}
		defer rawConn.Close()

		// set the deadline for read and write
		rawConn.SetDeadline(time.Now().Add(defaultConnDeadline))

		// if peer authentication is enabled
		if s.authHelper != nil {
			authHandshaker, _ := auth.NewClientHandshaker(ctx, s.authHelper, rawConn)
			conn, err = authHandshaker.ClientHandshake(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			conn = rawConn
		}
	}

	defer func() {
		if err != nil && s.secureHelper != nil {
			s.connPoolMtx.Lock()
			defer s.connPoolMtx.Unlock()
			conn.Close()
			s.connPool.Del(remoteAddr)
		}
	}()

	// encode and send the request message
	data, err := encode(request)
	if err != nil {
		return nil, errors.Errorf("encode: %w", err)
	}
	if _, err := conn.Write(data); err != nil {
		return nil, errors.Errorf("conn write: %w", err)
	}

	// receive and decode the response message
	response, err := decode(conn)
	if err != nil {
		return nil, errors.Errorf("conn read: %w", err)
	}

	return response, nil
}
