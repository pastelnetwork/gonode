package kademlia

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/google/uuid"
	"go.uber.org/ratelimit"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc/credentials"
)

const (
	defaultConnDeadline                = 10 * time.Minute
	defaultConnRate                    = 1000
	defaultMaxPayloadSize              = 200 // MB
	errorBusy                          = "Busy"
	maxConcurrentFindBatchValsRequests = 25
)

// Network for distributed hash table
type Network struct {
	dht      *DHT              // the distributed hash table
	listener net.Listener      // the server socket for the network
	self     *Node             // queries node itself
	limiter  ratelimit.Limiter // the rate limit for accept socket
	done     chan struct{}     // network is stopped

	// For secure connection
	secureHelper credentials.TransportCredentials
	connPool     *ConnPool
	connPoolMtx  sync.Mutex

	// for authentication only
	authHelper *AuthHelper
	sem        *semaphore.Weighted
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
		sem:          semaphore.NewWeighted(maxConcurrentFindBatchValsRequests),
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

func (s *Network) handleFindNode(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, FindNode); response != nil || err != nil {
			res = response
		}
	}()

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

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)

	// the closest contacts
	hashedTargetID, _ := utils.Sha3256hash(request.Target)
	closest, _ := s.dht.ht.closestContacts(K, hashedTargetID, []*Node{message.Sender})

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

func (s *Network) handleFindValue(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, FindValue); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*FindValueRequest)
	if !ok {
		err := errors.New("invalid FindValueRequest")
		return s.generateResponseMessage(FindValue, message.Sender, ResultFailed, err.Error())
	}

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)

	// retrieve the value from queries storage
	value, err := s.dht.store.Retrieve(ctx, request.Target)
	if err != nil {
		err = errors.Errorf("store retrieve: %w", err)
		response := &FindValueResponse{
			Status: ResponseStatus{
				Result: ResultFailed,
				ErrMsg: err.Error(),
			},
		}

		closest, _ := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})
		response.Closest = closest.Nodes

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
		closest, _ := s.dht.ht.closestContacts(K, request.Target, []*Node{message.Sender})
		response.Closest = closest.Nodes
	}

	// new a response message
	resMsg := s.dht.newMessage(FindValue, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleStoreData(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, StoreData); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*StoreDataRequest)
	if !ok {
		err := errors.New("invalid StoreDataRequest")
		return s.generateResponseMessage(StoreData, message.Sender, ResultFailed, err.Error())
	}

	log.P2P().WithContext(ctx).Debugf("handle store data: %v", message.String())

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)

	// format the key
	key, _ := utils.Sha3256hash(request.Data)

	value, err := s.dht.store.Retrieve(ctx, key)
	if err != nil || len(value) == 0 {
		// store the data to queries storage
		if err := s.dht.store.Store(ctx, key, request.Data, request.Type, false); err != nil {
			err = errors.Errorf("store the data: %w", err)
			return s.generateResponseMessage(StoreData, message.Sender, ResultFailed, err.Error())
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

func (s *Network) handleReplicate(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, Replicate); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*ReplicateDataRequest)
	if !ok {
		err := errors.New("invalid ReplicateDataRequest")
		return s.generateResponseMessage(Replicate, message.Sender, ResultFailed, err.Error())
	}

	log.P2P().WithContext(ctx).Debugf("handle replicate data: %v", message.String())

	if err := s.handleReplicateRequest(ctx, request, message.Sender.ID, message.Sender.IP, message.Sender.Port); err != nil {
		return s.generateResponseMessage(Replicate, message.Sender, ResultFailed, err.Error())
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
	keysToStore, err := s.dht.store.RetrieveBatchNotExist(ctx, req.Keys, 5000)
	if err != nil {
		log.WithContext(ctx).WithField("keys", len(req.Keys)).WithField("from-ip", ip).Errorf("unable to retrieve batch replication keys: %v", err)
		return fmt.Errorf("unable to retrieve batch replication keys: %w", err)
	}

	log.WithContext(ctx).WithField("to-store-keys", len(keysToStore)).WithField("rcvd-keys", len(req.Keys)).WithField("from-ip", ip).Debug("store batch replication keys to be stored")

	if len(keysToStore) > 0 {
		if err := s.dht.store.StoreBatchRepKeys(keysToStore, string(id), ip, port); err != nil {
			return fmt.Errorf("unable to store batch replication keys: %w", err)
		}

		log.WithContext(ctx).WithField("to-store-keys", len(keysToStore)).WithField("rcvd-keys", len(req.Keys)).WithField("from-ip", ip).Info("store batch replication keys stored")
	}

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
		reqID := uuid.New().String()

		var response []byte
		switch request.MessageType {
		case FindNode:
			encoded, err := s.handleFindNode(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle find node request failed")
				return
			}
			response = encoded
		case BatchFindNode:
			encoded, err := s.handleBatchFindNode(ctx, request)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("handle batch find node request failed")
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
		case BatchFindValues:
			// handle the request for finding value
			encoded, err := s.handleBatchFindValues(ctx, request, reqID)
			if err != nil {
				log.P2P().WithContext(ctx).WithField("p2p-req-id", reqID).WithError(err).Error("handle batch find values request failed")
				return
			}
			response = encoded
		case BatchStoreData:
			// handle the request for storing data
			encoded, err := s.handleBatchStoreData(ctx, request)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("handle batch store data request failed")
				return
			}
			response = encoded
		case BatchGetValues:
			// handle the request for finding value
			encoded, err := s.handleGetValuesRequest(ctx, request, reqID)
			if err != nil {
				log.WithContext(ctx).WithField("p2p-req-id", reqID).WithError(err).Error("handle batch get values request failed")
				return
			}
			response = encoded
		default:
			log.P2P().WithContext(ctx).Errorf("invalid message type: %v", request.MessageType)
			return
		}

		// write the response
		if _, err := conn.Write(response); err != nil {
			log.P2P().WithField("p2p-req-id", reqID).WithContext(ctx).WithError(err).Errorf("write failed %d", request.MessageType)
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
func (s *Network) Call(ctx context.Context, request *Message, isLong bool) (*Message, error) {
	timeout := 10 * time.Second

	if request.MessageType == BatchStoreData {
		timeout = 60 * time.Second
	}
	if request.MessageType == FindNode {
		timeout = 30 * time.Second
	}
	if request.MessageType == BatchFindNode {
		timeout = 15 * time.Second
	}
	if isLong {
		timeout = 3 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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
		rawConn.SetDeadline(time.Now().UTC().Add(timeout))

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

func (s *Network) handleBatchFindValues(ctx context.Context, message *Message, reqID string) (res []byte, err error) {
	// Try to acquire the semaphore, wait up to 1 minute
	log.WithContext(ctx).Debug("Attempting to acquire semaphore immediately.")
	if !s.sem.TryAcquire(1) {
		log.WithContext(ctx).Info("Immediate acquisition failed. Waiting up to 1 minute.")
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := s.sem.Acquire(ctxWithTimeout, 1); err != nil {
			log.WithContext(ctx).Error("Failed to acquire semaphore within 1 minute.")
			// failed to acquire semaphore within 1 minute
			return s.generateResponseMessage(BatchFindValues, message.Sender, ResultFailed, errorBusy)
		}
		log.WithContext(ctx).Info("Semaphore acquired after waiting.")
	}

	// Add a defer function to recover from panic
	defer func() {
		s.sem.Release(1)

		if r := recover(); r != nil {
			// Log the error or handle it as you see fit
			log.WithContext(ctx).Errorf("HandleBatchFindValues Recovered from panic: %v", r)

			// Convert panic to error
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}

			res, _ = s.generateResponseMessage(BatchFindValues, message.Sender, ResultFailed, err.Error())
		}
	}()

	request, ok := message.Data.(*BatchFindValuesRequest)
	if !ok {
		return s.generateResponseMessage(BatchFindValues, message.Sender, ResultFailed, "invalid BatchFindValueRequest")
	}

	isDone, data, err := s.handleBatchFindValuesRequest(ctx, request, message.Sender.IP, reqID)
	if err != nil {
		return s.generateResponseMessage(BatchFindValues, message.Sender, ResultFailed, err.Error())
	}

	response := &BatchFindValuesResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
		Response: data,
		Done:     isDone,
	}

	resMsg := s.dht.newMessage(BatchFindValues, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleGetValuesRequest(ctx context.Context, message *Message, reqID string) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, BatchGetValues); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*BatchGetValuesRequest)
	if !ok {
		err := errors.New("invalid BatchGetValuesRequest")
		return s.generateResponseMessage(BatchGetValues, message.Sender, ResultFailed, err.Error())
	}

	log.WithContext(ctx).Infof("batch get values request received from %s", message.Sender.String())

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)
	keys := make([]string, len(request.Data))
	i := 0
	for key := range request.Data {
		keys[i] = key
		i++
	}

	values, count, err := s.dht.store.RetrieveBatchValues(ctx, keys)
	if err != nil {
		err = errors.Errorf("batch find values: %w", err)
		return s.generateResponseMessage(BatchGetValues, message.Sender, ResultFailed, err.Error())
	}

	log.WithContext(ctx).WithField("requested-keys", len(keys)).WithField("found", count).WithField("sender", message.Sender.String()).Info("batch get values request processed")

	for i, key := range keys {
		val := KeyValWithClosest{
			Value: values[i],
		}
		if len(val.Value) == 0 {
			decodedKey, err := hex.DecodeString(keys[i])
			if err != nil {
				err = errors.Errorf("batch find vals: decode key: %w - key %s", err, keys[i])
				return s.generateResponseMessage(BatchGetValues, message.Sender, ResultFailed, err.Error())
			}

			nodes, _ := s.dht.ht.closestContacts(Alpha, decodedKey, []*Node{message.Sender})
			val.Closest = nodes.Nodes
		}

		request.Data[key] = val
	}

	response := &BatchGetValuesResponse{
		Data: request.Data,
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}

	// new a response message
	resMsg := s.dht.newMessage(BatchGetValues, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleBatchFindValuesRequest(ctx context.Context, req *BatchFindValuesRequest, ip string, reqID string) (isDone bool, compressedData []byte, err error) {
	log.WithContext(ctx).WithField("p2p-req-id", reqID).WithField("keys", len(req.Keys)).WithField("from-ip", ip).Info("batch find values request received")
	if len(req.Keys) > 0 {
		log.WithContext(ctx).WithField("p2p-req-id", reqID).WithField("keys[0]", req.Keys[0]).WithField("keys[len]", req.Keys[len(req.Keys)-1]).
			WithField("from-ip", ip).Debug("first & last batch keys")
	}

	values, count, err := s.dht.store.RetrieveBatchValues(ctx, req.Keys)
	if err != nil {
		return false, nil, fmt.Errorf("failed to retrieve batch values: %w", err)
	}
	log.WithContext(ctx).WithField("p2p-req-id", reqID).WithField("values-len", len(values)).WithField("found", count).WithField("from-ip", ip).Info("batch find values request processed")

	isDone, count, compressedData, err = findOptimalCompression(count, req.Keys, values)
	if err != nil {
		return false, nil, fmt.Errorf("failed to find optimal compression: %w", err)
	}

	log.WithContext(ctx).WithField("p2p-req-id", reqID).WithField("compressed-data-len", utils.BytesToMB(uint64(len(compressedData)))).WithField("found", count).
		WithField("from-ip", ip).Info("batch find values response sent")

	return isDone, compressedData, nil
}

func findOptimalCompression(count int, keys []string, values [][]byte) (bool, int, []byte, error) {
	dataMap := make(map[string][]byte)
	for i, key := range keys {
		dataMap[key] = values[i]
	}

	compressedData, err := compressMap(dataMap)
	if err != nil {
		return true, 0, nil, err
	}

	// If the initial compressed data is under the threshold
	if utils.BytesIntToMB(len(compressedData)) < defaultMaxPayloadSize {
		log.WithField("compressed-data-len", utils.BytesToMB(uint64(len(compressedData)))).WithField("count", count).Debug("initial compression")
		return true, len(dataMap), compressedData, nil
	}

	iter := 0
	currentValuesCount := count
	for utils.BytesIntToMB(len(compressedData)) >= defaultMaxPayloadSize {
		size := utils.BytesIntToMB(len(compressedData))
		log.WithField("compressed-data-len", size).WithField("current-count", currentValuesCount).WithField("iter", iter).Debug("optimal compression")
		iter++
		// Find top 10 heaviest values and set their keys to nil in the map
		var heavyKeys []string
		currentValuesCount, heavyKeys = findTopHeaviestKeys(dataMap, size)
		for _, key := range heavyKeys {
			dataMap[key] = nil
		}

		// Recompress
		compressedData, err = compressMap(dataMap)
		if err != nil {
			return false, 0, nil, err
		}
	}

	// Calculate the count of non-nil keys
	counter := 0
	for _, v := range dataMap {
		if len(v) > 0 {
			counter++
		}
	}

	// if we were not able to fit even 1 key, there's nothing we can do at this point
	return counter == 0, counter, compressedData, nil
}

func compressMap(dataMap map[string][]byte) ([]byte, error) {
	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data map: %w", err)
	}

	compressedData, err := utils.Compress(dataBytes, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}

	return compressedData, nil
}

func findTopHeaviestKeys(dataMap map[string][]byte, size int) (int, []string) {
	type kv struct {
		Key string
		Len int
	}

	var sorted []kv
	count := 0
	for k, v := range dataMap {
		if len(v) > 0 { // Only consider non-nil values
			count++
			sorted = append(sorted, kv{k, len(v)})
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Len > sorted[j].Len
	})

	n := 10          // number of keys to remove from payload if payload is heavier than allowed size
	if count <= 50 { // if keys are less than 50, we'd wanna try a smaller decrement number
		n = 5
	}
	if count <= 10 { // if keys are less than 10, we'd wanna try a smaller decrement number
		n = 1
	}

	if size > (2 * defaultMaxPayloadSize) {
		log.Debug("find optimal compression decreasing payload by half")
		n = count / 2
	}

	topKeys := []string{}
	for i := 0; i < n && i < len(sorted); i++ {
		topKeys = append(topKeys, sorted[i].Key)
	}

	return count, topKeys
}

func (s *Network) handleBatchStoreData(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, BatchStoreData); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*BatchStoreDataRequest)
	if !ok {
		err := errors.New("invalid BatchStoreDataRequest")
		return s.generateResponseMessage(BatchStoreData, message.Sender, ResultFailed, err.Error())
	}

	log.P2P().WithContext(ctx).Info("handle batch store data request received")

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)

	if err := s.dht.store.StoreBatch(ctx, request.Data, 1, false); err != nil {
		err = errors.Errorf("batch store the data: %w", err)
		return s.generateResponseMessage(BatchStoreData, message.Sender, ResultFailed, err.Error())
	}

	response := &StoreDataResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}
	log.P2P().WithContext(ctx).Info("handle batch store data request processed")

	// new a response message
	resMsg := s.dht.newMessage(BatchStoreData, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handleBatchFindNode(ctx context.Context, message *Message) (res []byte, err error) {
	defer func() {
		if response, err := s.handlePanic(ctx, message.Sender, BatchFindNode); response != nil || err != nil {
			res = response
		}
	}()

	request, ok := message.Data.(*BatchFindNodeRequest)
	if !ok {
		err := errors.New("invalid FindNodeRequest")
		return s.generateResponseMessage(BatchFindNode, message.Sender, ResultFailed, err.Error())
	}

	// add the sender to queries hash table
	s.dht.addNode(ctx, message.Sender)

	response := &BatchFindNodeResponse{
		Status: ResponseStatus{
			Result: ResultOk,
		},
	}
	closestMap := make(map[string][]*Node)

	log.WithContext(ctx).WithField("sender", message.Sender.String()).Info("Batch Find Nodes Request Received")
	for _, hashedTargetID := range request.HashedTarget {
		closest, _ := s.dht.ht.closestContacts(K, hashedTargetID, []*Node{message.Sender})
		closestMap[base58.Encode(hashedTargetID)] = closest.Nodes
	}
	response.ClosestNodes = closestMap
	log.WithContext(ctx).WithField("sender", message.Sender.String()).Info("Batch Find Nodes Request Processed")

	// new a response message
	resMsg := s.dht.newMessage(BatchFindNode, message.Sender, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) generateResponseMessage(messageType int, receiver *Node, result ResultType, errMsg string) ([]byte, error) {
	responseStatus := ResponseStatus{
		Result: result,
		ErrMsg: errMsg,
	}

	var response interface{}

	switch messageType {
	case StoreData, BatchStoreData:
		response = &StoreDataResponse{Status: responseStatus}
	case FindNode, BatchFindNode:
		response = &BatchFindNodeResponse{Status: responseStatus}
	case FindValue, BatchFindValues:
		response = &BatchFindValuesResponse{Status: responseStatus}
	case Replicate:
		response = &ReplicateDataResponse{Status: responseStatus}
	case BatchGetValues:
		response = &BatchGetValuesResponse{Status: responseStatus}
	default:
		return nil, fmt.Errorf("unsupported message type %d", messageType)
	}

	resMsg := s.dht.newMessage(messageType, receiver, response)
	return s.encodeMesage(resMsg)
}

func (s *Network) handlePanic(ctx context.Context, sender *Node, messageType int) (res []byte, err error) {
	if r := recover(); r != nil {
		log.WithContext(ctx).Errorf("p2p network: recovered from panic: %v", r)

		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}

		if res, err := s.generateResponseMessage(messageType, sender, ResultFailed, err.Error()); err != nil {
			log.WithContext(ctx).Errorf("Error generating response message: %v", err)
		} else {
			return res, err
		}
	}

	return nil, nil
}
