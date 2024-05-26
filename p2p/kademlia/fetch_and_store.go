package kademlia

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	json "github.com/json-iterator/go"

	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

const (
	maxBatchAttempts                     = 1
	oneMB                                = 1024 * 1024 // 1 MB in bytes
	totalMaxAttempts                     = 20
	maxSingleBatchIterations             = 10
	failedKeysClosestContactsLookupCount = 12
	fetchBatchSize                       = 400
)

// FetchAndStore fetches all keys from the queries TODO replicate list, fetches value from respective nodes and stores them in the queries store
func (s *DHT) FetchAndStore(ctx context.Context) error {
	log.WithContext(ctx).Info("Getting fetch and store keys")
	keys, err := s.store.GetAllToDoRepKeys(failedKeysClosestContactsLookupCount+maxBatchAttempts+1, totalMaxAttempts)
	if err != nil {
		return fmt.Errorf("get all keys error: %w", err)
	}
	log.WithContext(ctx).WithField("count", len(keys)).Info("got keys from queries store")

	if len(keys) == 0 {
		return nil
	}

	//wg := sync.WaitGroup{}
	//wg.Add(len(keys))        // Add count of all keys before spawning goroutines
	var successCounter int32 // Create a counter for successful operations

	for i := 0; i < len(keys); i++ {
		key := keys[i]

		func(info domain.ToRepKey) {
			//defer wg.Done()
			cctx, ccancel := context.WithTimeout(ctx, 30*time.Second)
			defer ccancel()

			sKey, err := hex.DecodeString(info.Key)
			if err != nil {
				log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("hex decode key failed")
				return
			}

			n := Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}

			b := backoff.WithMaxRetries(backoff.NewConstantBackOff(2*time.Second), 10)
			var value []byte // replace with the actual type of "value"
			err = backoff.Retry(func() error {
				val, err := s.GetValueFromNode(cctx, sKey, &n)
				if err != nil {
					return err
				}
				value = val
				return nil
			}, b)

			if err != nil {
				log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("fetch & store key failed")
				value, err = s.iterateFindValue(cctx, IterateFindValue, sKey)
				if err != nil {
					log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("iterate fetch for replication failed")
					return
				} else if len(value) == 0 {
					log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("iterate fetch for replication failed 0 val")
					return
				}

				log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).Info("iterate fetch for replication success")
			}

			if err := s.store.Store(cctx, sKey, value, 0, false); err != nil {
				log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("fetch & queries store key failed")
				return
			}

			if err := s.store.DeleteRepKey(info.Key); err != nil {
				log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).WithError(err).Error("delete key from todo list failed")
				return
			}

			atomic.AddInt32(&successCounter, 1) // Increment the counter atomically

			log.WithContext(cctx).WithField("key", info.Key).WithField("ip", info.IP).Info("fetch & store key success")
		}(key)

		time.Sleep(100 * time.Millisecond)
	}

	//wg.Wait()

	log.WithContext(ctx).WithField("todo-keys", len(keys)).WithField("successfully-added-keys", atomic.LoadInt32(&successCounter)).Infof("Successfully fetched & stored keys") // Log the final count

	return nil
}

// BatchFetchAndStoreFailedKeys fetches all failed keys from the queries TODO replicate list, fetches value from respective nodes and stores them in the queries store
func (s *DHT) BatchFetchAndStoreFailedKeys(ctx context.Context) error {
	log.WithContext(ctx).Debug("Getting failed batch fetch and store keys")
	keys, err := s.store.GetAllToDoRepKeys(maxBatchAttempts+1, failedKeysClosestContactsLookupCount+maxBatchAttempts+1) // 2 - 14
	if err != nil {
		return fmt.Errorf("get all keys error: %w", err)
	}
	log.WithContext(ctx).WithField("count", len(keys)).Info("read failed keys from store")

	if len(keys) == 0 {
		return nil
	}

	repKeys := make([]domain.ToRepKey, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		igList := s.ignorelist.ToNodeList()
		sKey, err := hex.DecodeString(keys[i].Key)
		if err != nil {
			log.WithContext(ctx).WithField("key", keys[i].Key).WithField("ip", keys[i].IP).WithError(err).Error("hex decode key failed")
			continue
		}

		nl, _ := s.ht.closestContacts(failedKeysClosestContactsLookupCount, sKey, igList)
		attempt := (keys[i].Attempts - maxBatchAttempts) + 1

		if len(nl.Nodes) > attempt {
			repKey := domain.ToRepKey{
				Key:  keys[i].Key,
				ID:   string(nl.Nodes[attempt].ID),
				IP:   nl.Nodes[attempt].IP,
				Port: nl.Nodes[attempt].Port,
			}

			repKeys = append(repKeys, repKey)
		}
	}
	log.WithField("count", len(repKeys)).Info("got 2nd tier replication keys from queries store")

	if err := s.GroupAndBatchFetch(ctx, repKeys, 0, false); err != nil {
		log.WithContext(ctx).WithError(err).Error("group and batch fetch failed-keys error")
		return fmt.Errorf("group and batch fetch failed keys error: %w", err)
	}

	return nil
}

// BatchFetchAndStore fetches all keys from the queries TODO replicate list, fetches value from respective nodes and stores them in the queries store
func (s *DHT) BatchFetchAndStore(ctx context.Context) error {
	log.WithContext(ctx).Debug("Getting batch fetch and store keys")
	keys, err := s.store.GetAllToDoRepKeys(0, maxBatchAttempts)
	if err != nil {
		return fmt.Errorf("get all keys error: %w", err)
	}
	log.WithContext(ctx).WithField("count", len(keys)).Info("got batch todo rep-keys from queries store")

	if len(keys) == 0 {
		return nil
	}

	if err := s.GroupAndBatchFetch(ctx, keys, 0, false); err != nil {
		log.WithContext(ctx).WithError(err).Error("group and batch fetch error")
		return fmt.Errorf("group and batch fetch error: %w", err)
	}

	return nil
}

// GroupAndBatchFetch gets values from nodes in batches and store them
func (s *DHT) GroupAndBatchFetch(ctx context.Context, repKeys []domain.ToRepKey, datatype int, isOriginal bool) error {
	nodeMap := make(map[string][]*domain.ToRepKey)

	// Group keys by Node
	for i := 0; i < len(repKeys); i++ {
		node := &Node{
			ID:   []byte(repKeys[i].ID),
			IP:   repKeys[i].IP,
			Port: repKeys[i].Port,
		}
		nodeKey := generateKeyFromNode(node)
		nodeMap[nodeKey] = append(nodeMap[nodeKey], &repKeys[i])
	}

	// Fetch from each Node and store directly
	for nodeKey, repKeyList := range nodeMap {
		node, err := getNodeFromKey(nodeKey)
		if err != nil {
			return fmt.Errorf("invalid nodeKey %s: %w", nodeKey, err)
		}

		// Fetch from node in batches
		for i := 0; i < len(repKeyList); i += fetchBatchSize {
			end := i + fetchBatchSize
			if end > len(repKeyList) {
				end = len(repKeyList)
			}

			// Convert repKeyList[i:end] to byteKeys
			stringKeys := make([]string, end-i)
			for j, key := range repKeyList[i:end] {
				stringKeys[j] = key.Key
			}

			iterations := 0
			totalKeysFound := 0
			for len(stringKeys) > 0 && iterations < maxSingleBatchIterations {
				iterations++
				log.WithContext(ctx).WithField("node-ip", node.IP).WithField("count", len(stringKeys)).WithField("keys[0]", stringKeys[0]).
					WithField("keys[len()]", stringKeys[len(stringKeys)-1]).Info("fetching batch values from node")

				isDone, retMap, failedKeys, err := s.GetBatchValuesFromNode(ctx, stringKeys, node)
				if err != nil {
					// Log the error but don't stop the process, continue to the next node
					log.WithContext(ctx).WithField("node-ip", node.IP).WithError(err).Info("failed to get batch values")
					continue
				}

				// Convert retMap to response
				stringDelKeys := make([]string, 0)
				var response [][]byte
				for key, value := range retMap {
					if len(value) > 0 {
						stringDelKeys = append(stringDelKeys, key)
						response = append(response, value)
						totalKeysFound++
					}
				}

				if len(stringDelKeys) > 0 {
					// Store the values directly
					err = s.store.StoreBatch(ctx, response, datatype, isOriginal)
					if err != nil {
						// Log the error but don't stop the process, continue to the next node
						log.WithContext(ctx).WithField("node-ip", node.IP).WithError(err).Info("failed to store batch values")
						continue
					}

					// Delete the keys that were successfully stored
					err = s.store.BatchDeleteRepKeys(stringDelKeys)
					if err != nil {
						// Log the error but don't stop the process, continue to the next node
						log.WithContext(ctx).WithField("node-ip", node.IP).WithError(err).Info("failed to delete rep keys")
						continue
					}
				} else {
					log.WithContext(ctx).WithField("node-ip", node.IP).Warn("no values found in batch fetch")
				}

				if isDone && len(failedKeys) > 0 {
					if err := s.store.IncrementAttempts(failedKeys); err != nil {
						log.WithContext(ctx).WithField("node-ip", node.IP).WithError(err).Info("failed to increment attempts")
						// not adding 'continue' here because we want to delete the keys from the todo list
					}
				} else if isDone {
					stringKeys = []string{}
				} else if !isDone {
					stringKeys = failedKeys
				}
			}

			log.WithContext(ctx).WithField("node-ip", node.IP).WithField("count", totalKeysFound).WithField("iterations", iterations).Info("fetch batch values from node successfully")
		}
	}

	return nil
}

// GetBatchValuesFromNode get values from node in bateches
func (s *DHT) GetBatchValuesFromNode(ctx context.Context, keys []string, n *Node) (bool, map[string][]byte, []string, error) {
	log.WithContext(ctx).WithField("node-ip", n.IP).WithField("keys", len(keys)).Info("sending batch fetch request")

	messageType := BatchFindValues

	data := &BatchFindValuesRequest{Keys: keys}
	request := s.newMessage(messageType, n, data)

	var response *Message

	operation := func() error {
		var err error
		response, err = s.network.Call(ctx, request, true)
		if err != nil {
			return fmt.Errorf("call error: %w", err)
		}

		if response == nil {
			return fmt.Errorf("response is nil")
		}

		v, ok := response.Data.(*BatchFindValuesResponse)
		if !ok {
			return fmt.Errorf("batch get request failure - %s - node: %s", response.String(), n.String())
		}

		if v == nil {
			return fmt.Errorf("response data is nil")
		}

		if v.Status.Result == ResultOk {
			return nil
		} else if v.Status.Result == ResultFailed {
			if v.Status.ErrMsg == errorBusy {
				return fmt.Errorf("batch get request failure - %s - node: %s", "server busy", n.String())
			}

			return fmt.Errorf("batch get request failure - %s - node: %s", response.String(), n.String())
		}

		return err
	}

	// Set up the backoff parameters
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 10 * time.Second
	bo.MaxElapsedTime = 30 * time.Second // max time before stop retrying
	bo.Multiplier = 1

	if err := backoff.Retry(operation, bo); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Debugf("network call request %s failed", request.String())
		return false, nil, nil, fmt.Errorf("network call request %s failed: %w", request.String(), err)
	}

	v, ok := response.Data.(*BatchFindValuesResponse)
	isDone := false
	if ok && v.Status.Result == ResultOk {
		// First, decompress the data
		decompressedData, err := utils.Decompress(v.Response)
		if err != nil {
			return isDone, nil, nil, fmt.Errorf("failed to decompress data: %w", err)
		}

		// Next, unmarshal the decompressed data back into a map
		var decompressedMap map[string][]byte
		err = json.Unmarshal(decompressedData, &decompressedMap)
		if err != nil {
			return isDone, nil, nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		retMap, failedKeys, err := VerifyAndFilter(decompressedMap)
		if err != nil {
			return isDone, nil, nil, fmt.Errorf("failed to verify and filter data: %w", err)
		}
		log.WithContext(ctx).WithField("node-ip", n.IP).WithField("received-keys", len(decompressedMap)).
			WithField("verified-keys", len(retMap)).WithField("failed-keys", len(failedKeys)).Info("batch fetch response rcvd and keys verified")

		return v.Done, retMap, failedKeys, nil
	}

	return false, nil, nil, fmt.Errorf("batch get request failure - %s - node: %s", response.String(), n.String())
}

// VerifyAndFilter verifies the data and filters out the failed keys
func VerifyAndFilter(decompressedMap map[string][]byte) (map[string][]byte, []string, error) {
	var retMap = make(map[string][]byte)
	var failedKeys []string

	for key, value := range decompressedMap {
		if len(value) == 0 {
			failedKeys = append(failedKeys, key)
			continue
		}

		// Compute the SHA256 hash of the value using the helper function
		hash, err := utils.Sha3256hash(value)
		if err != nil {
			failedKeys = append(failedKeys, key)
			log.WithError(err).Error("failed to compute hash")
			continue
		}

		// Encode the hash to a hex string
		hashHex := hex.EncodeToString(hash)

		// Compare the computed hash with the key
		if hashHex == key {
			retMap[key] = value
		} else {
			log.WithField("key", key).WithField("hash", hashHex).Error("hash mismatch")
			failedKeys = append(failedKeys, key)
		}
	}

	return retMap, failedKeys, nil
}
