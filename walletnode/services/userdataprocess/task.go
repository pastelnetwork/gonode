package userdataprocess

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
)

// Task is the task of userdata processing.
type UserDataTask struct {
	*common.WalletNodeTask

	MeshHandler *mixins.MeshHandler

	service *UserDataService

	// information of nodes process to set userdata
	resultChan chan *userdata.ProcessResult
	request    *userdata.ProcessRequest

	// information user pastelid to retrieve userdata
	userpastelid  string // user pastelid
	resultChanGet chan *userdata.ProcessRequest

	userdata []userdata.ProcessResult
}

// Run starts the task
func (task *UserDataTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *UserDataTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	maxNode := task.service.config.NumberSuperNodes
	if task.request == nil {
		// This is process to retrieve userdata
		maxNode = 1 // Get data from 1 supernode only, currently we choose the 1st ranked supernode, but this may change later

		if err := task.MeshHandler.ConnectToNSuperNodes(ctx, maxNode); err != nil {
			return errors.Errorf("connect to top rank nodes: %w", err)
		}
	} else {
		_, _, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
		if err != nil {
			return errors.Errorf("connect to top rank nodes: %w", err)
		}

	}
	_ = task.MeshHandler.DisconnectInactiveNodes(ctx)

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.MeshHandler.Nodes.WaitConnClose(ctx, nodesDone)
	})

	if task.request == nil {
		// PROCESS TO RETRIEVE USERDATA FROM METADATA LAYER
		if err := task.ReceiveUserdata(ctx, task.userpastelid); err != nil {
			return errors.Errorf("receive userdata: %w", err)
		}
		// Post on result channel
		node := nodes[0]
		if node.ResultGet != nil {
			task.resultChanGet <- node.ResultGet
		} else {
			return errors.Errorf("receive userdata")
		}

		log.WithContext(ctx).Debug("Finished retrieve userdata")

	} else {
		// PROCESS TO SET/UPDATE USERDATA TO METADATA LAYER
		// Get the previous block hash
		// Get block num
		blockHash := ""
		blockNum, err := task.service.pastelHandler.PastelClient.GetBlockCount(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Debug("Get block num failed")
		} else {
			// Get block hash string
			blockInfo, err := task.service.pastelHandler.PastelClient.GetBlockVerbose1(ctx, blockNum)
			if err != nil {
				log.WithContext(ctx).WithError(err).Debug("Get block info failed")
			} else {
				blockHash = blockInfo.Hash
			}
		}

		task.request.PreviousBlockHash = blockHash

		// Get the value of task.request.ArtistPastelIDPassphrase for sign data, then empty it in the request to make sure it not sent to supernodes
		passphrase := task.request.ArtistPastelIDPassphrase
		task.request.ArtistPastelIDPassphrase = ""
		// Marshal task.request to byte array for signing
		js, err := json.Marshal(task.request)
		if err != nil {
			return errors.Errorf("marshal request: %w", err)
		}

		// Hash the request
		hashvalue, err := userdata.Sha3256hash(js)
		if err != nil {
			return errors.Errorf("hash request %w", err)
		}

		// Sign request with Wallet Node's pastelID and passphrase
		signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, hashvalue, task.request.ArtistPastelID, passphrase, "ed448")
		if err != nil {
			return errors.Errorf("sign ticket %w", err)
		}

		userdata := &userdata.ProcessRequestSigned{
			Userdata:     task.request,
			UserdataHash: hex.EncodeToString(hashvalue),
			Signature:    hex.EncodeToString(signature),
		}

		// Send userdata to supernodes for storing in MDL's rqlite db.
		if err := nodes.SendUserdata(ctx, userdata); err != nil {
			return err
		}

		res, err := task.AggregateResult(ctx, nodes)
		// Post on result channel
		task.resultChan <- &res
		if err != nil {
			return err
		}

		log.WithContext(ctx).WithField("userdata_result", res).Debug("Posted userdata result")
	}

	// close the connections
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	return nil
}

// SendUserdata sends the userdata to supernodes for verification and store in MDL
func (task *UserDataTask) SendUserdata(ctx context.Context, req *userdata.ProcessRequestSigned) error {
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		userDataNode, ok := someNode.SuperNodeAPIInterface.(*UserDataNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			res, err := userDataNode.SendUserdata(ctx, req)
			node.Result = res
			return err
		})
	}
	return group.Wait()
}

// ReceiveUserdata retrieve the userdata from Metadata layer
func (task *UserDataTask) ReceiveUserdata(ctx context.Context, pasteluserid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			res, err := task.ReceiveUserdata(ctx, pasteluserid)
			node.ResultGet = res
			return err
		})
	}
	return group.Wait()
}

// AggregateResult aggregate all results return by all supernode, and consider it valid or not
func (task *UserDataTask) AggregateResult(_ context.Context, nodes node.List) (userdata.ProcessResult, error) {
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
				return userdata.ProcessResult{}, errors.Errorf("Primary node have empty result")
			}
			return *result, nil
		}
	}

	return userdata.ProcessResult{}, errors.Errorf("aggregate Result failed")
}

// Error returns task err
func (task *UserDataTask) Error() error {
	return task.err
}

// SubscribeProcessResult returns the result state of userdata process
func (task *UserDataTask) SubscribeProcessResult() <-chan *userdata.ProcessResult {
	return task.resultChan
}

// SubscribeProcessResultGet returns the result state of userdata process
func (task *UserDataTask) SubscribeProcessResultGet() <-chan *userdata.ProcessRequest {
	return task.resultChanGet
}

func (task *UserDataTask) removeArtifacts() {
}

// NewUserDataTask returns a new Task instance.
func NewUserDataTask(service *UserDataService, request *userdata.ProcessRequest, userpastelid string) *UserDataTask {
	task := &UserDataTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		request:        request,
		userpastelid:   userpastelid,
		resultChan:     make(chan *userdata.ProcessResult),
		resultChanGet:  make(chan *userdata.ProcessRequest),
	}

	task.MeshHandler = mixins.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &UserDataNodeMaker{},
		service.pastelHandler,
		request.ArtistPastelID, request.ArtistPastelIDPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.acceptNodesTimeout, service.config.connectToNextNodeDelay,
	)

	return task
}
