package userdataprocess

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// UserDataTask is the task of userdata processing.
type UserDataTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler

	service *UserDataService

	// information of nodes process to set userdata
	resultChan chan *userdata.ProcessResult
	request    *userdata.ProcessRequest

	// information user pastelid to retrieve userdata
	userpastelid  string // user pastelid
	resultChanGet chan *userdata.ProcessRequest

	//userdata []userdata.ProcessResult

	err error
}

// Run starts the task
func (task *UserDataTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *UserDataTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if task.request == nil {
		// This is process to retrieve userdata
		maxNode := 1 // Get data from 1 supernode only, currently we choose the 1st ranked supernode, but this may change later

		if err := task.MeshHandler.ConnectToNSuperNodes(ctx, maxNode); err != nil {
			return errors.Errorf("connect to top rank nodes: %w", err)
		}
	} else {
		_, _, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
		if err != nil {
			return errors.Errorf("connect to top rank nodes: %w", err)
		}

	}
	task.MeshHandler.DisconnectInactiveNodes(ctx)

	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	if task.request == nil {
		// PROCESS TO RETRIEVE USERDATA FROM METADATA LAYER
		if err := task.RetrieveUserdata(ctx, task.userpastelid); err != nil {
			return errors.Errorf("retrieve userdata error: %w", err)
		}
		// Post on result channel
		node0 := task.MeshHandler.Nodes[0]
		userDataNode, ok := node0.SuperNodeAPIInterface.(*UserDataNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegisterNode", node0.String())
		}
		if userDataNode.ResultGet != nil {
			task.resultChanGet <- userDataNode.ResultGet
		} else {
			return errors.Errorf("Error, retrieved nil from userdata node.")
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
		passphrase := task.request.UserPastelIDPassphrase
		task.request.UserPastelIDPassphrase = ""
		// Marshal task.request to byte array for signing
		js, err := json.Marshal(task.request)
		if err != nil {
			return errors.Errorf("marshal request: %w", err)
		}

		// Hash the request
		hashvalue, err := utils.Sha3256hash(js)
		if err != nil {
			return errors.Errorf("hash request %w", err)
		}

		// Sign request with Wallet Node's pastelID and passphrase
		signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, hashvalue, task.request.UserPastelID, passphrase, "ed448")
		if err != nil {
			return errors.Errorf("sign ticket %w", err)
		}

		userdata := &userdata.ProcessRequestSigned{
			Userdata:     task.request,
			UserdataHash: hex.EncodeToString(hashvalue),
			Signature:    hex.EncodeToString(signature),
		}

		// Send userdata to supernodes for storing in MDL's rqlite db.
		if err := task.SendUserdata(ctx, userdata); err != nil {
			return err
		}

		res, err := task.AggregateResult(ctx)
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
			userDataNode.Result = res
			return err
		})
	}
	return group.Wait()
}

// RetrieveUserdata retrieve the userdata from Metadata layer
func (task *UserDataTask) RetrieveUserdata(ctx context.Context, pasteluserid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		userDataNode, ok := someNode.SuperNodeAPIInterface.(*UserDataNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			res, err := userDataNode.RetrieveUserdata(ctx, pasteluserid)
			userDataNode.ResultGet = res
			return err
		})
	}
	return group.Wait()
}

// AggregateResult aggregate all results return by all supernode, and consider it valid or not
func (task *UserDataTask) AggregateResult(_ context.Context) (userdata.ProcessResult, error) {
	// There is following common scenarios when supernodes response:
	// 1. Secondary node and primary node both response userdata validation result error
	// 2. Secondary node response userdata validation result success and primary node provide further processing result
	// 3. Some node fail to response, or not in the 2 case above, then we need to aggregate result and consider what happen

	// This part is for case 1 or 2 above, and we trust the primary node so we use its response
	for _, someNode := range task.MeshHandler.Nodes {
		userDataNode, ok := someNode.SuperNodeAPIInterface.(*UserDataNode)
		if !ok {
			//TODO: use assert here?
			return userdata.ProcessResult{}, errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
		}
		if someNode.IsPrimary() {
			result := userDataNode.Result
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
	// userpastelid is empty on createUserdata and updateUserdata; and not empty on getUserdata
	task := &UserDataTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		request:        request,
		userpastelid:   userpastelid,
		resultChan:     make(chan *userdata.ProcessResult),
		resultChanGet:  make(chan *userdata.ProcessRequest),
	}

	task.MeshHandler = common.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &UserDataNodeMaker{},
		service.pastelHandler,
		request.UserPastelID, request.UserPastelIDPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	return task
}
