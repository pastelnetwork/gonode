package kademlia

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/anacrolix/utp"
	"go.uber.org/ratelimit"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	defaultConnDeadline   = 5 * time.Second
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
}

// NewNetwork returns a network service
func NewNetwork(dht *DHT, self *Node) (*Network, error) {
	s := &Network{
		dht:  dht,
		self: self,
		done: make(chan struct{}),
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
	// close the socket
	if s.socket != nil {
		if err := s.socket.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("close socket failed")
		}
	}
}

func (s *Network) handleFindNode(_ context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*FindNodeRequest)
	if !ok {
		return nil, errors.New("impossible: must be FindNodeRequest")
	}

	// add the sender to local hash table
	s.dht.addNode(message.Sender)

	// the closest contacts
	closest := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})

	// new a response message
	response := s.dht.newMessage(
		FindNode,
		message.Sender,
		&FindNodeResponse{Closest: closest.Nodes},
	)

	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return nil, errors.Errorf("encode response for find node: %w", err)
	}

	return encoded, nil
}

func (s *Network) handleFindValue(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*FindValueRequest)
	if !ok {
		return nil, errors.New("impossible: must be FindValueRequest")
	}

	// add the sender to local hash table
	s.dht.addNode(message.Sender)

	data := &FindValueResponse{}
	// retrieve the value from local storage
	value, err := s.dht.store.Retrieve(ctx, request.Target)
	if err != nil {
		return nil, errors.Errorf("store retrieve: %w", err)
	}
	if value != nil {
		// return the value
		data.Value = value
	} else {
		// return the closest contacts
		closest := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})
		data.Closest = closest.Nodes
	}

	// new a response message
	response := s.dht.newMessage(FindValue, message.Sender, data)

	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return nil, errors.Errorf("encode response for find value: %w", err)
	}

	return encoded, nil
}

func (s *Network) handleStoreData(ctx context.Context, message *Message) ([]byte, error) {
	request, ok := message.Data.(*StoreDataRequest)
	if !ok {
		return nil, errors.New("impossible: must be StoreDataRequest")
	}
	log.WithContext(ctx).Debugf("handle store data: %v", message.String())

	// add the sender to local hash table
	s.dht.addNode(message.Sender)

	// format the key
	key := s.dht.hashKey(request.Data)
	// replication time for key
	replication := time.Now().Add(defaultReplicateTime)
	// store the data to local storage
	if err := s.dht.store.Store(ctx, key, request.Data, replication); err != nil {
		return nil, errors.Errorf("store the data: %w", err)
	}

	// new a response message
	response := s.dht.newMessage(StoreData, message.Sender, nil)
	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return nil, errors.Errorf("encode response for store data: %w", err)
	}

	return encoded, nil
}

func (s *Network) handlePing(_ context.Context, message *Message) ([]byte, error) {
	// new a response message
	response := s.dht.newMessage(Ping, message.Sender, nil)
	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return nil, errors.Errorf("encode response for ping: %w", err)
	}

	return encoded, nil
}

// handle the connection request
func (s *Network) handleConn(ctx context.Context, conn net.Conn) {
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
				continue
			}
			response = encoded
		case FindValue:
			// handle the request for finding value
			encoded, err := s.handleFindValue(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle find value request failed")
				continue
			}
			response = encoded
		case Ping:
			encoded, err := s.handlePing(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle ping request failed")
				continue
			}
			response = encoded
		case StoreData:
			// handle the request for storing data
			encoded, err := s.handleStoreData(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle store data request failed")
				continue
			}
			response = encoded
		default:
			log.WithContext(ctx).Errorf("impossible: invalid message type: %v", request.MessageType)
			return
		}

		// write the response
		if _, err := conn.Write(response); err != nil {
			log.WithContext(ctx).WithError(err).Error("conn write: failed")
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
		log.WithContext(ctx).Debugf("%v: incomming connection: %v", s.self.String(), conn.RemoteAddr())

		// handle the connection requests
		go s.handleConn(ctx, conn)
	}
}

// Call sends the request to target and receive the response
func (s *Network) Call(ctx context.Context, request *Message) (*Message, error) {
	remoteAddr := fmt.Sprintf("%s:%d", request.Sender.IP, request.Receiver.Port)

	// dial the remote address with udp network
	conn, err := utp.DialContext(ctx, remoteAddr)
	if err != nil {
		return nil, errors.Errorf("dial %q: %w", remoteAddr, err)
	}
	defer conn.Close()

	// set the deadline for read and write
	conn.SetDeadline(time.Now().Add(defaultConnDeadline))

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
