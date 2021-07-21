package userdataprocess

import (
	"context"
	"fmt"
	"sync"

	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/pastel"
)

// Task is the task of searching for artwork.
type Task struct {
	task.Task
	*Service

	resultChan chan *UserdataProcessResult
	err        error
	request    *UserdataProcessRequest
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")
	defer close(task.resultChan)

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

	// TODO (1): Make this init and connect to super nodes to generic reusable function to avoid code duplication (1)
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return err
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughMasterNode)
		return errors.New("unable to find enough Supernodes to send userdata to")
	}

	// Try to create mesh of supernodes, connecting to all supernodes in a different sequences.
	var nodes node.List
	var errs error
	for primaryRank := range topNodes {
		nodes, err = task.meshNodes(ctx, topNodes, primaryRank)
		if err != nil {
			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warnf("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
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

	// End TODO (1) ---- 

	// Get the previous block hash
	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return nil, errors.Errorf("failed to get block info with given block num %d: %w", blockNum, err)
	}

	// Decode hash string to byte
	task.request.PreviousBlockHash  = blockInfo.Hash
	if err != nil {
		return nil, errors.Errorf("failed to convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	// Get the value of task.Request.ArtistPastelIDPassphrase for sign data, then empty it in the request to make sure it not sent to supernodes
	passphrase := task.Request.ArtistPastelIDPassphrase
	task.Request.ArtistPastelIDPassphrase = ""
	// Marshal task.request to byte array for signing
	js, err := json.Marshal(task.request)
	if err != nil {
		return nil, errors.Errorf("failed to encode ticket %w", err)
	}

	// Hash the request
	hashvalue = userdata.sha3256hash(js)

	// Sign request with Wallet Node's pastelID and passphrase
	signature, err := task.pastelClient.Sign(ctx, hashvalue, task.Request.ArtistPastelID, passphrase)
	if err != nil {
		return nil, errors.Errorf("failed to sign ticket %w", err)
	}

	userdata := &UserdataProcessRequestSigned{
		Userdata: 		&(task.Request),
		UserdataHash:	hashvalue,
		Signature:		signature,
	}

	// Send userdata to supernodes for storing in MDL's rqlite db.
	if err := nodes.SendUserdata(ctx, userdata) ; err != nil {
		return err
	} else {

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
func (task *Task) AggregateResult(ctx context.Context,nodes node.List) (UserdataProcessResult, error) {
	aggregate := make(map[string][]int)
	count := 0
	for i, node := range nodes {
		node := node
		if node.result != nil {
			count++
			// Marshal result to get the hash value
			js, err := json.Marshal(*(node.result))
			if err != nil {
				errors.Errorf("failed marshal the result %w", err)
			}
			log.WithContext(ctx).Debugf("Node result print:%s", string(js))

			// Hash the result
			hashvalue = userdata.sha3256hash(js)
			aggregate[hashvalue] = append(aggregate[hashvalue],i)
		}
	}

	if count < task.config.MinimalNodeConfirmSuccess {
		// If there is not enough reponse from supernodes
		return &UserdataProcessResult{
			ResponseCode: userdata.ERROR_NOT_ENOUGH_SUPERNODE_RESPONSE
			Detail 		: userdata.Description[userdata.ERROR_NOT_ENOUGH_SUPERNODE_RESPONSE]
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
		return &UserdataProcessResult{
			ResponseCode: userdata.ERROR_NOT_ENOUGH_SUPERNODE_CONFIRM
			Detail 		: userdata.Description[userdata.ERROR_NOT_ENOUGH_SUPERNODE_CONFIRM]
		}, nil
	}

	// There is enough supernodes verified our request, so we return that response
	if (len(aggregate[finalHashKey]) > 0 && nodes[(aggregate[finalHashKey])[0]] != nil) {
		result := *(nodes[(aggregate[finalHashKey])[0]].result)
		return result, nil
	}
	
	return nil, errors.Errorf("failed to Aggregate Result")
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
	meshNodes.SetPrimary(primaryIndex)
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


// pastelTopNodes retrieve the top super nodes we want to send userdata to
func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}


// Error returns task err
func (task *Task) Error() error {
	return task.err
}

// SubscribeProcessResult returns the result state of userdata process
func (task *Task) SubscribeProcessResult() <-chan *UserdataProcessResult {
	return task.resultChan
}

// NewTask returns a new Task instance.
func NewTask(service *Service, request *UserdataProcessRequest) *Task {
	return &Task{
		Task:       task.New(StatusTaskStarted),
		Service:    service,
		request:    request,
		resultChan: make(chan *UserdataProcessResult),
	}
}
