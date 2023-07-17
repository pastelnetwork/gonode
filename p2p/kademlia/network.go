package kademlia

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"go.uber.org/ratelimit"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"google.golang.org/grpc/credentials"
)

const (
	defaultConnDeadline   = 30 * time.Second
	defaultConnRate       = 1000
	defaultMaxPayloadSize = 32 * 1024 * 1024 // 32MB
)

// Network for distributed hash table
type Network struct {
	dht      *DHT              // the distributed hash table
	listener net.Listener      // the server socket for the network
	self     *Node             // local node itself
	limiter  ratelimit.Limiter // the rate limit for accept socket
	done     chan struct{}     // network is stopped

	// For secure connection
	secureHelper credentials.TransportCredentials
	connPool     *ConnPool
	connPoolMtx  sync.Mutex

	// for authentication only
	authHelper *AuthHelper
}

// NewNetwork returns a network service
func NewNetwork(ctx context.Context, dht *DHT, self *Node, secureHelper credentials.TransportCredentials, authHelper *AuthHelper) (*Network, error) {
	s := &Network{
		dht:          dht,
		self:         self,
		done:         make(chan struct{}),
		secureHelper: secureHelper,
		connPool:     NewConnPool(ctx),
		authHelper:   authHelper,
	}
	// init the rate limiter
	s.limiter = ratelimit.New(defaultConnRate)

	addr := fmt.Sprintf("%s:%d", self.IP, self.Port)
	// new tcp listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.P2P().Debugf("Error trying to get tcp socket: %s", err)
		return nil, err
	}
	s.listener = listener
	log.P2P().Debugf("Listening on: %s", addr)

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
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Errorf("close socket failed")
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

	if len(value) > 0 {
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

	log.P2P().WithContext(ctx).Debugf("handle store data: %v", message.String())

	// add the sender to local hash table
	s.dht.addNode(ctx, message.Sender)

	// format the key
	key := s.dht.hashKey(request.Data)

	value, err := s.dht.store.Retrieve(ctx, key)
	if err == nil && len(value) > 0 {
		log.WithContext(ctx).WithField("key", hex.EncodeToString(key)).Info("data already exists")
	} else {
		// store the data to local storage
		if err := s.dht.store.Store(ctx, key, request.Data, request.Type, false); err != nil {
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

func (s *Network) handleReplicate(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*ReplicateDataRequest)
	if !ok {
		err := errors.New("invalid ReplicateDataRequest")

		response := &ReplicateDataResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		resMsg := s.dht.newMessage(Replicate, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	log.P2P().WithContext(ctx).Debugf("handle replicate data: %v", message.String())

	if err := s.handleReplicateRequest(ctx, request, message.Sender.ID, message.Sender.IP, message.Sender.Port); err != nil {
		response := &ReplicateDataResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}
		resMsg := s.dht.newMessage(Replicate, message.Sender, response)
		return s.encodeMesage(resMsg)
	}

	response := &ReplicateDataResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}

	// new a response message
	resMsg := s.dht.newMessage(Replicate, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleReplicateRequest(ctx context.Context, req *ReplicateDataRequest, id []byte, ip string, port int) error {
	keys, err := decompressKeys(req.Keys)
	if err != nil {
		return fmt.Errorf("unable to decode keys: %w", err)
	}

	var keysToStore [][]byte
	for i := 0; i < len(keys); i++ {
		value, err := s.dht.store.Retrieve(ctx, keys[i])
		if err == nil && len(value) > 0 {
			log.WithContext(ctx).WithField("key", hex.EncodeToString(keys[i])).Info("data already exists")
		} else {
			keysToStore = append(keysToStore, keys[i])
		}
	}

	if len(keysToStore) > 0 {
		if err := s.dht.store.StoreBatchRepKeys(keysToStore, string(id), ip, port); err != nil {
			return fmt.Errorf("unable to store batch replication keys: %w", err)
		}
	}
	log.WithContext(ctx).WithField("keys", len(keysToStore)).Info("store batch replication keys count")

	return nil
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
			log.P2P().WithContext(ctx).WithError(err).Error("server secure establish failed")
			return
		}
	} else {
		// if peer authentication is enabled
		if s.authHelper != nil {
			authHandshaker, _ := auth.NewServerHandshaker(ctx, s.authHelper, rawConn)
			conn, err = authHandshaker.ServerHandshake(ctx)
			if err != nil {
				rawConn.Close()
				log.P2P().WithContext(ctx).WithError(err).Error("server authentication failed")
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
			log.P2P().WithContext(ctx).WithError(err).Error("read and decode failed")
			return
		}

		var response []byte
		switch request.MessageType {
		case FindNode:
			encoded, err := s.handleFindNode(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle find node request failed")
				return
			}
			response = encoded
		case FindValue:
			// handle the request for finding value
			encoded, err := s.handleFindValue(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle find value request failed")
				return
			}
			response = encoded
		case Ping:
			encoded, err := s.handlePing(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle ping request failed")
				return
			}
			response = encoded
		case StoreData:
			// handle the request for storing data
			encoded, err := s.handleStoreData(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle store data request failed")
				return
			}
			response = encoded
		case Replicate:
			// handle the request for replicate request
			encoded, err := s.handleReplicate(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle replicate request failed")
				return
			}
			response = encoded
		default:
			log.P2P().WithContext(ctx).Errorf("invalid message type: %v", request.MessageType)
			return
		}

		// write the response
		if _, err := conn.Write(response); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("write failed")
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
		conn, err := s.listener.Accept()
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
				log.P2P().WithContext(ctx).WithError(err).Errorf("socket accept failed, retrying in %v", tempDelay)

				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "closed") {
				return
			}

			log.P2P().WithContext(ctx).WithError(err).Error("socket accept failed")
			return
		}

		// handle the connection requests
		go s.handleConn(ctx, conn)
	}
}

// Call sends the request to target and receive the response
func (s *Network) Call(ctx context.Context, request *Message) (*Message, error) {
	if request.Receiver != nil && request.Receiver.Port == 50052 {
		log.P2P().WithContext(ctx).Error("invalid port")
		return nil, errors.New("invalid port")
	}
	if request.Sender != nil && request.Sender.Port == 50052 {
		log.P2P().WithContext(ctx).Error("invalid port")
		return nil, errors.New("invalid port")
	}

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
		// dial the remote address with tcp network
		var d net.Dialer
		rawConn, err = d.DialContext(ctx, "tcp", remoteAddr)
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
