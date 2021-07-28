package userdataprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)

const (
	// TODO: Move this to config pass from app.go later.
	defaultNumberSuperNodes   = 10
	minimalNodeConfirmSuccess = 8
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	acceptedMu sync.Mutex
	accepted   Nodes

	// userdata signed by this supernode
	ownSNData userdata.SuperNodeRequest

	// valid only for a task run as primary
	peersSNDataSignedMtx    *sync.Mutex
	peersSNDataSigned       map[string]userdata.SuperNodeRequest
	allPeersSNDatasReceived chan struct{}

	// valid only for secondary node
	connectedTo *Node

	// valid only for primary node
	connectedToLeader *Node
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = task.context(ctx)
	defer log.WithContext(ctx).Debug("Task canceled")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	return task.RunAction(ctx)
}

// Session is handshake wallet to supernode
func (task *Task) Session(_ context.Context, isPrimary bool) error {
	if err := task.RequiredStatus(StatusTaskStarted); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if isPrimary {
			log.WithContext(ctx).Debugf("Acts as primary node")
			task.UpdateStatus(StatusPrimaryMode)
			return nil
		}

		log.WithContext(ctx).Debugf("Acts as secondary node")
		task.UpdateStatus(StatusSecondaryMode)

		return nil
	})
	return nil
}

// AcceptedNodes waits for connection supernodes, as soon as there is the required amount returns them.
func (task *Task) AcceptedNodes(serverCtx context.Context) (Nodes, error) {
	if err := task.RequiredStatus(StatusPrimaryMode); err != nil {
		return nil, err
	}

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("Waiting for supernodes to connect")

		sub := task.SubscribeStatus()
		for {
			select {
			case <-serverCtx.Done():
				return nil
			case <-ctx.Done():
				return nil
			case status := <-sub():
				if status.Is(StatusConnected) {
					return nil
				}
			}
		}
	})
	return task.accepted, nil
}

// SessionNode accepts secondary node
func (task *Task) SessionNode(_ context.Context, nodeID string) error {
	task.acceptedMu.Lock()
	defer task.acceptedMu.Unlock()

	if err := task.RequiredStatus(StatusPrimaryMode); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if node := task.accepted.ByID(nodeID); node != nil {
			return errors.Errorf("node %q is already registered", nodeID)
		}

		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			return errors.Errorf("failed to get node by extID %s %w", nodeID, err)
		}
		task.accepted.Add(node)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debugf("Accept secondary node")

		if len(task.accepted) >= task.config.NumberConnectedNodes {
			task.UpdateStatus(StatusConnected)
		}
		return nil
	})
	return nil
}

// ConnectTo connects to primary node
func (task *Task) ConnectTo(_ context.Context, nodeID, sessID string, nodetype int) error {
	if err := task.RequiredStatus(StatusSecondaryMode); err != nil {
		return err
	}

	task.NewAction(func(ctx context.Context) error {
		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			return err
		}

		if err := node.connect(ctx); err != nil {
			return err
		}

		if err := node.Session(ctx, task.config.PastelID, sessID); err != nil {
			return err
		}

		if nodetype == userdata.NodeTypePrimary {
			task.connectedTo = node
		} else if nodetype == userdata.NodeTypeLeader {
			task.connectedToLeader = node
		} else {
			return errors.Errorf("Invalid NodeType")
		}
		task.UpdateStatus(StatusConnected)
		return nil
	})
	return nil
}

// SupernodeProcessUserdata process the userdata send from Walletnode
func (task *Task) SupernodeProcessUserdata(ctx context.Context, req *userdata.UserdataProcessRequestSigned) (userdata.UserdataProcessResult, error) {
	log.WithContext(ctx).Debugf("supernodeProcessUserdata on user PastelID: %s", req.Userdata.ArtistPastelID)

	validateResult, err := task.validateUserdata(req.Userdata)
	if err != nil {
		return userdata.UserdataProcessResult{}, errors.Errorf("failed to validateUserdata")
	}
	if validateResult.ResponseCode == userdata.ErrorOnContent {
		// If the request from Walletnode fail the validation, return the response to Walletnode and stop process further
		return validateResult, nil
	}

	// Validation the request from Walletnode successful, we continue to get acknowledgement and confirmation process from multiple SNs
	snRequest := userdata.SuperNodeRequest{}

	// Marshal validateResult to byte array for signing
	js, err := json.Marshal(validateResult)
	if err != nil {
		return userdata.UserdataProcessResult{}, errors.Errorf("failed to encode validateResult %w", err)
	}
	// Hash the validateResult
	hashvalue, err := userdata.Sha3256hash(js)
	if err != nil {
		return userdata.UserdataProcessResult{}, errors.Errorf("failed to hash userdata %w", err)
	}
	snRequest.UserdataResultHash = string(hashvalue)
	snRequest.UserdataHash = req.UserdataHash

	task.ownSNData = snRequest // At this step this SuperNodeRequest only contain UserdataResultHash and UserdataHash

	<-task.NewAction(func(ctx context.Context) error {
		// sign the data if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %d", task.connectedTo == nil)
		if err := task.signAndSendSNDataSigned(ctx, task.ownSNData, task.connectedTo == nil); err != nil {
			return errors.Errorf("failed to signed and send SuperNodeRequest")
		}
		return nil
	})

	// only primary node start this action
	var processResult userdata.UserdataProcessResult
	if task.connectedTo == nil {
		<-task.NewAction(func(ctx context.Context) error {
			log.WithContext(ctx).Debug("waiting for signature from peers")
			for {
				select {
				case <-ctx.Done():
					err := ctx.Err()
					if err != nil {
						log.WithContext(ctx).Debug("waiting for SuperNodeRequest from peers cancelled or timeout")
					}
					return err
				case <-task.allPeersSNDatasReceived:
					log.WithContext(ctx).Debugf("all SuperNodeRequest received so start validation")

					if resp, err := task.verifyPeersUserdata(ctx); err != nil {
						return errors.Errorf("fail to verifyPeersUserdata: %w", err)
					} else {
						processResult = resp
						return nil
					}
				}
			}
		})
	}

	return processResult, nil
}

// retrieveUserdata get the userdata from database
func (task *Task) retrieveUserdata(ctx context.Context, userpastelid string) (userdata.UserdataProcessRequest, error) {
	log.WithContext(ctx).Debugf("retrieveUserdata on user PastelID: %s", userpastelid)

	// only primary node start this action
	return task.Service.databaseOps.ReadUserData(ctx, userpastelid)
}

// Sign and send SNDataSigned if not primary
func (task *Task) signAndSendSNDataSigned(ctx context.Context, sndata userdata.SuperNodeRequest, isPrimary bool) error {
	log.WithContext(ctx).Debugf("signAndSendSNDataSigned begin to sign SuperNodeRequest")
	signature, err := task.pastelClient.Sign(ctx, []byte(sndata.UserdataHash+sndata.UserdataResultHash), task.config.PastelID, task.config.PassPhrase)
	if err != nil {
		return errors.Errorf("failed to sign sndata %w", err)
	}
	if !isPrimary {
		sndata.HashSignature = string(signature)
		sndata.NodeID = task.config.PastelID
		log.WithContext(ctx).Debug("send signed sndata to primary node")
		if _, err := task.connectedTo.ProcessUserdata.SendUserdataToPrimary(ctx, sndata); err != nil {
			return errors.Errorf("failed to send signature to primary node %s at address %s %w", task.connectedTo.ID, task.connectedTo.Address, err)
		}
	}
	return nil
}

// AddPeerSNDataSigned gather all the Userdata with signature from other Supernodes
func (task *Task) AddPeerSNDataSigned(ctx context.Context, snrequest userdata.SuperNodeRequest) error {
	log.WithContext(ctx).Debugf("addPeerSNDataSigned begin to aggregate SuperNodeRequest")

	task.peersSNDataSignedMtx.Lock()
	defer task.peersSNDataSignedMtx.Unlock()

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive supernode userdata result signed from node %s", snrequest.NodeID)
		if node := task.accepted.ByID(snrequest.NodeID); node == nil {
			return errors.Errorf("node %s not in accepted list", snrequest.NodeID)
		}

		task.peersSNDataSigned[snrequest.NodeID] = snrequest
		if len(task.peersSNDataSigned) == len(task.accepted) {
			log.WithContext(ctx).Debug("all signature received")
			go func() {
				close(task.allPeersSNDatasReceived)
			}()
		}
		return nil
	})
	return nil
}

func (task *Task) verifyPeersUserdata(ctx context.Context) (userdata.UserdataProcessResult, error) {
	log.WithContext(ctx).Debugf("all supernode data signed received so start validation")
	successCount := 0
	dataMatchingCount := 0
	for _, sndata := range task.peersSNDataSigned {
		// Verify the data
		if ok, err := task.pastelClient.Verify(ctx, []byte(sndata.UserdataHash+sndata.UserdataResultHash), string(sndata.HashSignature), sndata.NodeID); err != nil {
			errors.Errorf("failed to verify signature %s of node %s", string(sndata.HashSignature), sndata.NodeID)
			continue
		} else {
			if !ok {
				errors.Errorf("signature of node %s mistmatch", sndata.NodeID)
				continue
			} else {
				successCount++
				if task.ownSNData.UserdataHash == sndata.UserdataHash &&
					task.ownSNData.UserdataResultHash == sndata.UserdataResultHash {
					dataMatchingCount++
				}
			}
		}
	}

	if successCount < minimalNodeConfirmSuccess-1 {
		// If there is not enough userdata signed from other supernodes
		return userdata.UserdataProcessResult{
			ResponseCode: userdata.ErrorNotEnoughSupernodeConfirm,
			Detail:       userdata.Description[userdata.ErrorNotEnoughSupernodeConfirm],
		}, nil
	} else if dataMatchingCount < minimalNodeConfirmSuccess-1 {
		// If the userdata between supernodes is not matching with each other
		return userdata.UserdataProcessResult{
			ResponseCode: userdata.ErrorUserdataMismatchBetweenSupernode,
			Detail:       userdata.Description[userdata.ErrorUserdataMismatchBetweenSupernode],
		}, nil
	} else {
		// Success
		return userdata.UserdataProcessResult{
			ResponseCode: userdata.SuccessVerifyAllSignature,
			Detail:       userdata.Description[userdata.SuccessVerifyAllSignature],
		}, nil
	}
}

func (task *Task) validateUserdata(req *userdata.UserdataProcessRequest) (userdata.UserdataProcessResult, error) {
	result := userdata.UserdataProcessResult{}

	textLengthLimit := 1000

	contentValidation := userdata.SuccessValidateContent

	// Biography validation
	if len(req.Biography) > textLengthLimit {
		result.Biography = "Biography text length is greater than the limit " + string(textLengthLimit) + " characters"
		contentValidation = userdata.ErrorOnContent
	}

	// Primary Language validation
	// TODO: Or the “Primary Language” field might be limited to a list of language names that we can validate against
	// (this way we can avoid the problem of users writing in “Russian” and “русский” which creates confusion and data fragmentation).

	// FacebookLink validation
	if !(strings.Contains(strings.ToLower(req.FacebookLink), "facebook.com") || strings.Contains(strings.ToLower(req.FacebookLink), "fb.com")) {
		result.Biography = "Facebook Link is not valid"
		contentValidation = userdata.ErrorOnContent
	}

	// Image validation
	// Image Extension validation
	allowExtension := []string{".png", ".jpeg", ".jpg"}
	isAvatarMatchExtention := false
	isCoverPhotoMatchExtention := false
	for _, extension := range allowExtension {
		if strings.Contains(strings.ToLower(req.AvatarImage.Filename), extension) {
			isAvatarMatchExtention = true
		}
		if strings.Contains(strings.ToLower(req.CoverPhoto.Filename), extension) {
			isCoverPhotoMatchExtention = true
		}
	}
	if !isAvatarMatchExtention {
		result.AvatarImage = "Avatar extension must be in the following: " + strings.Join(allowExtension, ",")
		contentValidation = userdata.ErrorOnContent
	}
	if !isCoverPhotoMatchExtention {
		result.CoverPhoto = "CoverPhoto extension must be in the following: " + strings.Join(allowExtension, ",")
		contentValidation = userdata.ErrorOnContent
	}

	// Image Size validation
	minSizeLimit := 10 * 1024   // 10 kb
	maxSizeLimit := 1000 * 1024 // 1000 kb
	if len(req.AvatarImage.Content) < minSizeLimit || len(req.AvatarImage.Content) > maxSizeLimit {
		result.AvatarImage = "Avatar size must be in the range: [" + string(minSizeLimit) + "-" + string(maxSizeLimit) + "]"
		contentValidation = userdata.ErrorOnContent
	}
	if len(req.CoverPhoto.Content) < minSizeLimit || len(req.CoverPhoto.Content) > maxSizeLimit {
		result.CoverPhoto = "Cover Photo size must be in the range: [" + string(minSizeLimit) + "-" + string(maxSizeLimit) + "]"
		contentValidation = userdata.ErrorOnContent
	}

	// Image Resolution validation
	// TODO: Requirement not specified yet, will do later

	// Other field Validation
	// TODO: Requirement not specified yet, will do later

	result.ResponseCode = int32(contentValidation)
	result.Detail = userdata.Description[contentValidation]
	return result, nil
}

func (task *Task) pastelNodeByExtKey(ctx context.Context, nodeID string) (*Node, error) {
	masterNodes, err := task.pastelClient.MasterNodesTop(ctx)
	log.WithContext(ctx).Debugf("master node %s", masterNodes)

	if err != nil {
		return nil, err
	}

	for _, masterNode := range masterNodes {
		if masterNode.ExtKey != nodeID {
			continue
		}
		node := &Node{
			client:  task.Service.nodeClient,
			ID:      masterNode.ExtKey,
			Address: masterNode.ExtAddress,
		}
		return node, nil
	}

	return nil, errors.Errorf("node %q not found", nodeID)
}

func (task *Task) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	return &Task{
		Task:                    task.New(StatusTaskStarted),
		Service:                 service,
		peersSNDataSignedMtx:    &sync.Mutex{},
		peersSNDataSigned:       make(map[string]userdata.SuperNodeRequest),
		allPeersSNDatasReceived: make(chan struct{}),
	}
}
