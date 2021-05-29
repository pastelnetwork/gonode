package kademlia

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"

	"github.com/anacrolix/utp"
	"github.com/ccding/go-stun/stun"
)

type networking interface {
	sendMessage(context.Context, *message, bool, int64) (*expectedResponse, error)
	getMessage() chan (*message)
	messagesFin()
	timersFin()
	getDisconnect() chan (int)
	init(self *NetworkNode)
	createSocket(host string, port int, useStun bool, stunAddr string) (publicHost string, publicPort int, err error)
	listen(ctx context.Context) error
	disconnect() error
	cancelResponse(*expectedResponse)
	isInitialized() bool
	getNetworkAddr() string
}

type realNetworking struct {
	socket        *utp.Socket
	sendChan      chan (*message)
	recvChan      chan (*message)
	dcStartChan   chan (int)
	dcEndChan     chan (int)
	dcTimersChan  chan (int)
	dcMessageChan chan (int)
	mutex         *sync.Mutex
	connected     bool
	initialized   bool
	responseMap   map[int64]*expectedResponse
	aliveConns    *sync.WaitGroup
	self          *NetworkNode
	msgCounter    int64
	remoteAddress string
}

type expectedResponse struct {
	ch    chan (*message)
	query *message
	node  *NetworkNode
	id    int64
}

func (rn *realNetworking) init(self *NetworkNode) {
	rn.self = self
	rn.mutex = &sync.Mutex{}
	rn.sendChan = make(chan (*message))
	rn.recvChan = make(chan (*message))
	rn.dcStartChan = make(chan (int), 10)
	rn.dcEndChan = make(chan (int))
	rn.dcTimersChan = make(chan (int))
	rn.dcMessageChan = make(chan (int))
	rn.responseMap = make(map[int64]*expectedResponse)
	rn.aliveConns = &sync.WaitGroup{}
	rn.connected = false
	rn.initialized = true
}

func (rn *realNetworking) isInitialized() bool {
	return rn.initialized
}

func (rn *realNetworking) getMessage() chan (*message) {
	return rn.recvChan
}

func (rn *realNetworking) getNetworkAddr() string {
	return rn.remoteAddress
}

func (rn *realNetworking) messagesFin() {
	rn.dcMessageChan <- 1
}

func (rn *realNetworking) getDisconnect() chan (int) {
	return rn.dcStartChan
}

func (rn *realNetworking) timersFin() {
	rn.dcTimersChan <- 1
}

func (rn *realNetworking) createSocket(host string, port int, useStun bool, stunAddr string) (publicHost string, publicPort int, err error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	if rn.connected {
		return "", 0, errors.New("already connected")
	}

	remoteAddress := net.JoinHostPort(host, strconv.Itoa(port))
	socket, err := utp.NewSocket("udp", remoteAddress)
	if err != nil {
		return "", 0, errors.Errorf("failed to create a new socket %q: %w", remoteAddress, err)
	}

	if useStun {
		c := stun.NewClientWithConnection(socket)

		if stunAddr != "" {
			c.SetServerAddr(stunAddr)
		}

		_, h, err := c.Discover()
		if err != nil {
			return "", 0, errors.Errorf("failed to contact the STUN server: %w", err)
		}

		_, err = c.Keepalive()
		if err != nil {
			return "", 0, errors.Errorf("failed to enable keepalive: %w", err)
		}

		host := h.IP()
		port := int(h.Port())
		remoteAddress = net.JoinHostPort(host, strconv.Itoa(port))
	}

	rn.remoteAddress = remoteAddress
	rn.connected = true
	rn.socket = socket

	return host, port, nil
}

func (rn *realNetworking) sendMessage(ctx context.Context, msg *message, expectResponse bool, id int64) (*expectedResponse, error) {
	rn.mutex.Lock()
	if id == -1 {
		id = rn.msgCounter
		rn.msgCounter++
	}
	msg.ID = id
	rn.mutex.Unlock()

	host := msg.Receiver.IP.String()
	port := strconv.Itoa(msg.Receiver.Port)
	remoteAddress := net.JoinHostPort(host, port)

	conn, err := rn.socket.DialContext(ctx, "", remoteAddress)
	if err != nil {
		return nil, errors.Errorf("failed to dial %q: %w", remoteAddress, err)
	}

	data, err := serializeMessage(msg)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		return nil, errors.Errorf("failed to write: %w", err)
	}

	if expectResponse {
		rn.mutex.Lock()
		defer rn.mutex.Unlock()
		expectedResponse := &expectedResponse{
			ch:    make(chan (*message)),
			node:  msg.Receiver,
			query: msg,
			id:    id,
		}
		// TODO we need a way to automatically clean these up as there are
		// cases where they won't be removed manually
		rn.responseMap[id] = expectedResponse
		return expectedResponse, nil
	}

	return nil, nil
}

func (rn *realNetworking) cancelResponse(res *expectedResponse) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	close(rn.responseMap[res.query.ID].ch)
	delete(rn.responseMap, res.query.ID)
}

func (rn *realNetworking) disconnect() error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	if !rn.connected {
		return errors.New("not connected")
	}
	rn.dcStartChan <- 1
	rn.dcStartChan <- 1
	<-rn.dcTimersChan
	<-rn.dcMessageChan
	close(rn.sendChan)
	close(rn.recvChan)
	close(rn.dcTimersChan)
	close(rn.dcMessageChan)
	err := rn.socket.CloseNow()
	rn.connected = false
	rn.initialized = false
	close(rn.dcEndChan)
	if err != nil {
		return errors.Errorf("failed to close socket: %w", err)
	}
	return nil
}

func (rn *realNetworking) listen(ctx context.Context) error {
	for {
		conn, err := rn.socket.Accept()
		if err != nil {
			rn.disconnect()
			<-rn.dcEndChan

			if err.Error() == "closed" {
				return nil
			}
			return errors.Errorf("failed to accept new connection: %s", err)
		}

		go func(conn net.Conn) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Wait for messages
					msg, err := deserializeMessage(conn)
					if err != nil {
						// if err.Error() == "EOF" {
						// 	// Node went bye bye
						// }
						// TODO should we penalize this node somehow ? Ban it ?
						return
					}

					isPing := msg.Type == messageTypePing

					if !areNodesEqual(msg.Receiver, rn.self, isPing) {
						// TODO should we penalize this node somehow ? Ban it ?
						continue
					}

					if msg.ID < 0 {
						// TODO should we penalize this node somehow ? Ban it ?
						continue
					}

					rn.mutex.Lock()
					if rn.connected {
						if msg.IsResponse {
							if rn.responseMap[msg.ID] == nil {
								// We were not expecting this response
								rn.mutex.Unlock()
								continue
							}

							if !areNodesEqual(rn.responseMap[msg.ID].node, msg.Sender, isPing) {
								// TODO should we penalize this node somehow ? Ban it ?
								rn.mutex.Unlock()
								continue
							}

							if msg.Type != rn.responseMap[msg.ID].query.Type {
								close(rn.responseMap[msg.ID].ch)
								delete(rn.responseMap, msg.ID)
								rn.mutex.Unlock()
								continue
							}

							if !msg.IsResponse {
								close(rn.responseMap[msg.ID].ch)
								delete(rn.responseMap, msg.ID)
								rn.mutex.Unlock()
								continue
							}

							resChan := rn.responseMap[msg.ID].ch
							rn.mutex.Unlock()
							resChan <- msg
							rn.mutex.Lock()
							close(rn.responseMap[msg.ID].ch)
							delete(rn.responseMap, msg.ID)
							rn.mutex.Unlock()
						} else {
							assertion := false
							switch msg.Type {
							case messageTypeFindNode:
								_, assertion = msg.Data.(*queryDataFindNode)
							case messageTypeFindValue:
								_, assertion = msg.Data.(*queryDataFindValue)
							case messageTypeStore:
								_, assertion = msg.Data.(*queryDataStore)
							default:
								assertion = true
							}

							if !assertion {
								fmt.Printf("Received bad message %v from %+v", msg.Type, msg.Sender)
								close(rn.responseMap[msg.ID].ch)
								delete(rn.responseMap, msg.ID)
								rn.mutex.Unlock()
								continue
							}

							rn.recvChan <- msg
							rn.mutex.Unlock()
						}
					} else {
						rn.mutex.Unlock()
					}
				}
			}
		}(conn)
	}
}
