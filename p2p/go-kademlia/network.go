package kademlia

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/anacrolix/utp"
	"github.com/sirupsen/logrus"
)

const (
	defaultConnDeadline = 3 * time.Second
)

// Network for distributed hash table
type Network struct {
	dht    *DHT          // the distributed hash table
	socket *utp.Socket   // the server socket for the network
	self   *Node         // local node itself
	done   chan struct{} // network is stopped
}

// NewNetwork returns a network service
func NewNetwork(dht *DHT, self *Node) (*Network, error) {
	s := &Network{
		dht:  dht,
		self: self,
		done: make(chan struct{}),
	}

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
			logrus.Debugf("close socket: %v", err)
		}
	}
}

func (s *Network) handleFindNode(ctx context.Context, conn net.Conn, message *Message) error {
	request, ok := message.Data.(*FindNodeRequest)
	if !ok {
		return errors.New("impossible: must be FindNodeRequest")
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
		return fmt.Errorf("encode response for find node: %v", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("conn write: %v", err)
	}

	return nil
}

func (s *Network) handleFindValue(ctx context.Context, conn net.Conn, message *Message) error {
	request, ok := message.Data.(*FindValueRequest)
	if !ok {
		return errors.New("impossible: must be FindValueRequest")
	}

	// add the sender to local hash table
	s.dht.addNode(message.Sender)

	data := &FindValueResponse{}
	// retrieve the value from local storage
	value := s.dht.store.Retrieve(ctx, request.Target)
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
		return fmt.Errorf("encode response for find value: %v", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("conn write: %v", err)
	}

	return nil
}

func (s *Network) handleStoreData(ctx context.Context, conn net.Conn, message *Message) error {
	request, ok := message.Data.(*StoreDataRequest)
	if !ok {
		return errors.New("impossible: must be StoreDataRequest")
	}
	logrus.Debugf("handle store data: %v", message.String())

	// add the sender to local hash table
	s.dht.addNode(message.Sender)

	// format the key
	key := s.dht.hashKey(request.Data)
	// expiration time for key
	expiration := s.dht.keyExpireTime(key)
	// replication time for key
	replication := time.Now().Add(s.dht.options.ReplicateTime)
	// store the data to local storage
	if err := s.dht.store.Store(ctx, key, request.Data, replication, expiration); err != nil {
		return fmt.Errorf("store the data: %v", err)
	}

	// new a response message
	response := s.dht.newMessage(StoreData, message.Sender, nil)
	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return fmt.Errorf("encode response for store data: %v", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("conn write: %v", err)
	}

	return nil
}

func (s *Network) handlePing(ctx context.Context, conn net.Conn, message *Message) error {
	// new a response message
	response := s.dht.newMessage(Ping, message.Sender, nil)
	// send the response to client
	encoded, err := encode(response)
	if err != nil {
		return fmt.Errorf("encode response for ping: %v", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("conn write: %v", err)
	}

	return nil
}

// handle the connection request
func (s *Network) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			logrus.Debugf("connection %s is done", conn.RemoteAddr().String())
			return
		default:
		}

		// read the request from connection
		request, err := decode(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			logrus.Errorf("read and decode: %v", err)
			return
		}

		// received message which is from myself
		if !CompareNodes(request.Receiver, s.self) {
			logrus.Errorf("impossible: receiver is myself")
			continue
		}

		switch request.MessageType {
		case FindNode:
			if err := s.handleFindNode(ctx, conn, request); err != nil {
				logrus.Errorf("handle find node request: %v", err)
			}
		case FindValue:
			// handle the request for finding value
			if err := s.handleFindValue(ctx, conn, request); err != nil {
				logrus.Errorf("handle find value request: %v", err)
			}
		case Ping:
			if err := s.handlePing(ctx, conn, request); err != nil {
				logrus.Errorf("handle ping request: %v", err)
			}
		case StoreData:
			// handle the request for storing data
			if err := s.handleStoreData(ctx, conn, request); err != nil {
				logrus.Errorf("handle store data request: %v", err)
			}
		default:
			logrus.Errorf("impossible: invalid message type: %v", request.MessageType)
		}
	}
}

// serve the incomming connection
func (s *Network) serve(ctx context.Context) {
	for {
		// accept the incomming connections
		conn, err := s.socket.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				return
			}

			logrus.Errorf("socket accept: %v", err)
			return
		}
		logrus.Debugf("%v: incomming connection: %v", s.self.String(), conn.RemoteAddr())

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
		return nil, fmt.Errorf("dial %q: %w", remoteAddr, err)
	}
	defer conn.Close()

	// set the deadline for read and write
	conn.SetDeadline(time.Now().Add(defaultConnDeadline))

	// encode and send the request message
	data, err := encode(request)
	if err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}
	if _, err := conn.Write(data); err != nil {
		return nil, fmt.Errorf("conn write: %v", err)
	}

	// receive and decode the response message
	response, err := decode(conn)
	if err != nil {
		return nil, fmt.Errorf("conn read: %v", err)
	}

	return response, nil
}
