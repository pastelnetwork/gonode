package userdataprocess

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess/node"
)

// Task is the task of userdata processing.
type Task struct {
	task.Task
	*Service

	// information of nodes process to set userdata
	nodes      node.List
	resultChan chan *userdata.UserdataProcessResult
	err        error
	request    *userdata.UserdataProcessRequest

	// information user pastelid to retrieve userdata
	userpastelid  string // user pastelid
	resultChanGet chan *userdata.UserdataProcessRequest
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")
	defer close(task.resultChan)
	defer close(task.resultChanGet)

	if err := task.run(ctx); err != nil {
		task.err = err
		task.UpdateStatus(StatusTaskFailure)
		log.WithContext(ctx).WithError(err).Warnf("Task failed")

		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)

	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	maxNode := task.config.NumberSuperNodes
	if task.request == nil {
		// This is process to retrieve userdata
		maxNode = 1 // Get data from 1 supernode only, currently we choose the 1st ranked supernode, but this may change later
	}

	// TODO: Make this init and connect to super nodes to generic reusable function to avoid code duplication (1)
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx, maxNode)
	if err != nil {
		return err
	}
	if len(topNodes) < maxNode {
		task.UpdateStatus(StatusErrorNotEnoughMasterNode)
		return errors.New("unable to find enough Supernodes to send userdata to")
	}

	// Try to create mesh of supernodes, connecting to all supernodes in a different sequences.
	var nodes node.List
	var errs error
	nodes, err = task.meshNodes(ctx, topNodes, 0) // Connect a mesh node with primary is 1st ranked SN
	if err != nil {
		if errors.IsContextCanceled(err) {
			return err
		}
		errs = errors.Append(errs, err)
		log.WithContext(ctx).WithError(err).Warnf("Could not create a mesh of the nodes")
	}

	if len(nodes) < maxNode {
		return errors.Errorf("Could not create a mesh of %d nodes: %w", task.config.NumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	nodes.Activate()
	// Disconnect supernodes that are not involved in the process.
	topNodes.DisconnectInactive()

	// Cancel context when any connection is broken.
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return nodes.WaitConnClose(ctx)
	})
	task.UpdateStatus(StatusConnected)
	task.nodes = nodes

	if task.request == nil {
		// PROCESS TO RETRIEVE USERDATA FROM METADATA LAYER
		if err := nodes.ReceiveUserdata(ctx, task.userpastelid); err != nil {
			return err
		} else {
			// Post on result channel
			node := nodes[0]
			if node.ResultGet != nil {
				task.resultChanGet <- node.ResultGet
			} else {
				return errors.Errorf("failed to receive userdata")
			}

			log.WithContext(ctx).Debug("Finished retrieve userdata")
		}
	} else {
		// PROCESS TO SET/UPDATE USERDATA TO METADATA LAYER
		// Get the previous block hash
		// Get block num

		// TODO: Unblock this part after pastelClient support it
		/*
		blockNum, err := task.pastelClient.GetBlockCount(ctx)
		if err != nil {
			return errors.Errorf("failed to get block num: %w", err)
		}

		// Get block hash string
		blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
		if err != nil {
			return errors.Errorf("failed to get block info with given block num %d: %w", blockNum, err)
		} */

		// Decode hash string to byte
		task.request.PreviousBlockHash = "" /* blockInfo.Hash */
		if err != nil {
			return errors.Errorf("failed to convert hash string %s to bytes: %w", "" /*blockInfo.Hash*/, err)
		}

		// Get the value of task.request.ArtistPastelIDPassphrase for sign data, then empty it in the request to make sure it not sent to supernodes
		passphrase := task.request.ArtistPastelIDPassphrase
		task.request.ArtistPastelIDPassphrase = ""
		// Marshal task.request to byte array for signing
		js, err := json.Marshal(task.request)
		if err != nil {
			return errors.Errorf("failed to encode request %w", err)
		}

		// Hash the request
		hashvalue, err := userdata.Sha3256hash(js)
		if err != nil {
			return errors.Errorf("failed to hash request %w", err)
		}

		// Sign request with Wallet Node's pastelID and passphrase
		signature, err := task.pastelClient.Sign(ctx, hashvalue, task.request.ArtistPastelID, passphrase, "ed448")
		if err != nil {
			return errors.Errorf("failed to sign ticket %w", err)
		}

		userdata := &userdata.UserdataProcessRequestSigned{
			Userdata:     task.request,
			UserdataHash: hex.EncodeToString(hashvalue),
			Signature:    hex.EncodeToString(signature),
		}

		// Send userdata to supernodes for storing in MDL's rqlite db.
		if err := nodes.SendUserdata(ctx, userdata); err != nil {

			return err
		} else {
			res, err := task.AggregateResult(ctx, nodes)
			if err != nil {
				return err
			}
			// Post on result channel
			task.resultChan <- &res
			log.WithContext(ctx).WithField("userdata_result", res).Debug("Posted userdata result")
		}
	}

	// close the connections
	for i := range nodes {
		if err := nodes[i].Connection.Close(); err != nil {
			return errors.Errorf("failed to close connection to node %s %w", nodes[i].PastelID(), err)
		}
	}

	// Wait for all connections to disconnect.
	return groupConnClose.Wait()
}

// AggregateResult aggregate all results return by all supernode, and consider it valid or not
func (task *Task) AggregateResult(ctx context.Context, nodes node.List) (userdata.UserdataProcessResult, error) {
	// There is following common scenarios when supernodes response:
	// 1. Secondary node and primary node both response userdata validation result error
	// 2. Secondary node response userdata validation result success and primary node provide further processing result
	// 3. Some node fail to response, or not in the 2 case above, then we need to aggregate result and consider what happen

	// This part is for case 1 or 2 above, and we trust the primary node so we use its response
	for _, node := range nodes {
		node := node
		if node.IsPrimary() {
			result := node.Result
			if result == nil {
				return userdata.UserdataProcessResult{}, errors.Errorf("Primary node have empty result")
			}
			return *result, nil
		}
	}

	// This part of aggregate the response for case 3 and for future use
	// This is for in case we want to do descrepancy check between all result of nodes
	// for node reputation score
	/* aggregate := make(map[string][]int)
	count := 0
	for i, node := range nodes {
		node := node
		if node.Result != nil {
			count++
			// Marshal result to get the hash value
			js, err := json.Marshal(*(node.Result))
			if err != nil {
				errors.Errorf("failed marshal the result %w", err)
			}
			log.WithContext(ctx).Debugf("Node result print:%s", string(js))

			// Hash the result
			hashvalue, err := userdata.Sha3256hash(js)
			if err != nil {
				errors.Errorf("failed hash the result %w", err)
			}
			aggregate[hex.EncodeToString(hashvalue)] = append(aggregate[hex.EncodeToString(hashvalue)],i)
		}
	}

	if count < task.config.MinimalNodeConfirmSuccess {
		// If there is not enough reponse from supernodes
		return userdata.UserdataProcessResult{
			ResponseCode: userdata.ErrorNotEnoughSupernodeResponse,
			Detail:       userdata.Description[userdata.ErrorNotEnoughSupernodeResponse],
		}, nil
	}

	max := 0
	var finalHashKey string
	for key, element := range aggregate {
		if max < len(element) {
			max = len(element)
			finalHashKey = key
		}
	}

	if len(aggregate[finalHashKey]) < task.config.MinimalNodeConfirmSuccess {
		return userdata.UserdataProcessResult{
			ResponseCode: userdata.ErrorNotEnoughSupernodeConfirm,
			Detail:       userdata.Description[userdata.ErrorNotEnoughSupernodeConfirm],
		}, nil
	}

	// There is enough supernodes verified our request, so we return that response
	if len(aggregate[finalHashKey]) > 0 && nodes[(aggregate[finalHashKey])[0]] != nil {
		result := *(nodes[(aggregate[finalHashKey])[0]].Result)
		return result, nil
	} */

	return userdata.UserdataProcessResult{}, errors.Errorf("failed to Aggregate Result")
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List

	primary := nodes[primaryIndex]
	if err := primary.Connect(ctx, task.config.connectTimeout); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
	}
	primary.SetPrimary(true)

	if len(nodes) == 1 {
		// If the number of nodes only have 1 node, we use this primary node and return directly
		meshNodes.Add(primary)
		return meshNodes, nil
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	var secondaries node.List
	go func() {
		for i, node := range nodes {
			node := node

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(task.config.connectToNextNodeDelay):
				go func() {
					defer errors.Recover(log.Fatal)

					if err := node.Connect(ctx, task.config.connectTimeout); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					secondaries.Add(node)

					if err := node.ConnectTo(ctx, primary.PastelID(), primary.SessID()); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Seconary %q connected to primary", node)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, task.config.acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	meshNodes.Add(primary)
	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		node := secondaries.FindByPastelID(pastelID)
		if node == nil {
			return nil, errors.New("not found accepted node")
		}
		meshNodes.Add(node)
	}
	return meshNodes, nil
}

// pastelTopNodes retrieve the top super nodes we want to send userdata to, limit by maxNode
func (task *Task) pastelTopNodes(ctx context.Context, maxNode int) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	count := 0
	for _, mn := range mns {
		count++
		if count <= maxNode {
			nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
		} else {
			break
		}
	}

	return nodes, nil
}

// Error returns task err
func (task *Task) Error() error {
	return task.err
}

// SubscribeProcessResult returns the result state of userdata process
func (task *Task) SubscribeProcessResult() <-chan *userdata.UserdataProcessResult {
	return task.resultChan
}

// SubscribeProcessResultGet returns the result state of userdata process
func (task *Task) SubscribeProcessResultGet() <-chan *userdata.UserdataProcessRequest {
	return task.resultChanGet
}

// NewTask returns a new Task instance.
func NewTask(service *Service, request *userdata.UserdataProcessRequest, userpastelid string) *Task {
	return &Task{
		Task:          task.New(StatusTaskStarted),
		Service:       service,
		request:       request,
		userpastelid:  userpastelid,
		resultChan:    make(chan *userdata.UserdataProcessResult),
		resultChanGet: make(chan *userdata.UserdataProcessRequest),
	}
}
