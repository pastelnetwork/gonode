package userdataprocess

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"strings"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)

var imageAllowExtension = []string{".png", ".jpeg", ".jpg"}

const (
	// TODO: Move this to config pass from app.go later.
	imageMinSizeLimit        = 10 * 1024   //10kB
	imageMaxSizeLimit        = 1000 * 1024 //1000kB
	biographyTextLengthLimit = 1000
	facebookLongURL          = "facebook.com"
	facebookShortURL         = "fb.com"
)

// UserDataTask is the task of registering new User.
type UserDataTask struct {
	*common.SuperNodeTask

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
	ConnectedTo *Node

	// valid only for primary node
	ConnectedToLeader *Node
}

// Run starts the task
func (task *UserDataTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// Session is handshake wallet to supernode
func (task *UserDataTask) Session(_ context.Context, isPrimary bool) error {
	if err := task.RequiredStatus(common.StatusTaskStarted); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if isPrimary {
			log.WithContext(ctx).Debug("Acts as primary node")
			task.UpdateStatus(common.StatusPrimaryMode)
			return nil
		}

		log.WithContext(ctx).Debug("Acts as secondary node")
		task.UpdateStatus(common.StatusSecondaryMode)

		return nil
	})
	return nil
}

// AcceptedNodes waits for connection supernodes, as soon as there is the required amount returns them.
func (task *UserDataTask) AcceptedNodes(serverCtx context.Context) (Nodes, error) {
	if err := task.RequiredStatus(common.StatusPrimaryMode); err != nil {
		return nil, err
	}

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debug("Waiting for supernodes to connect")

		sub := task.SubscribeStatus()
		for {
			select {
			case <-serverCtx.Done():
				return nil
			case <-ctx.Done():
				return nil
			case status := <-sub():
				if status.Is(common.StatusConnected) {
					return nil
				}
			}
		}
	})
	return task.accepted, nil
}

// SessionNode accepts secondary node
func (task *UserDataTask) SessionNode(_ context.Context, nodeID string) error {
	task.acceptedMu.Lock()
	defer task.acceptedMu.Unlock()

	err := task.RequiredStatus(common.StatusPrimaryMode)
	if err != nil {
		return err
	}

	var actionErr error
	<-task.NewAction(func(ctx context.Context) error {
		if node := task.accepted.ByID(nodeID); node != nil {
			actionErr = errors.Errorf("node %q is already registered", nodeID)
			return nil
		}

		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			actionErr = errors.Errorf("get node by extID %s %w", nodeID, err)
			return nil
		}
		task.accepted.Add(node)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debug("Accept secondary node")

		if len(task.accepted) >= task.config.NumberSuperNodes-1 {
			task.UpdateStatus(common.StatusConnected)
		}
		return nil
	})
	return actionErr
}

// ConnectTo connects to primary node
func (task *UserDataTask) ConnectTo(_ context.Context, nodeID, sessID string) error {
	err := task.RequiredStatus(common.StatusSecondaryMode)
	if err != nil {
		return err
	}

	var actionErr error
	task.NewAction(func(ctx context.Context) error {
		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			actionErr = err
			return nil
		}

		if err := node.connect(ctx); err != nil {
			actionErr = err
			return nil
		}

		if err := node.Session(ctx, task.config.PastelID, sessID); err != nil {
			actionErr = err
			return nil
		}

		task.ConnectedTo = node

		task.UpdateStatus(common.StatusConnected)
		return nil
	})
	return actionErr
}

// SupernodeProcessUserdata process the userdata send from Walletnode
func (task *UserDataTask) SupernodeProcessUserdata(ctx context.Context, req *userdata.ProcessRequestSigned) (userdata.ProcessResult, error) {
	log.WithContext(ctx).Debugf("supernodeProcessUserdata on user PastelID: %s", req.Userdata.UserPastelID)

	validateResult, err := task.validateUserdata(req.Userdata)
	if err != nil {
		return userdata.ProcessResult{}, errors.Errorf("validateUserdata")
	}
	if validateResult.ResponseCode == userdata.ErrorOnContent {
		// If the request from Walletnode fail the validation, return the response to Walletnode and stop process further
		return validateResult, nil
	}

	// Validate user signature
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		log.WithContext(ctx).Debugf("failed to decode signature %s of user %s", req.Signature, req.Userdata.UserPastelID)
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, nil
	}

	userdatahash, err := hex.DecodeString(req.UserdataHash)
	if err != nil {
		log.WithContext(ctx).Debugf("failed to decode userdata hash %s of user %s", req.UserdataHash, req.Userdata.UserPastelID)
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, nil
	}

	ok, err := task.pastelClient.Verify(ctx, userdatahash, string(signature), req.Userdata.UserPastelID, "ed448")
	if err != nil || !ok {
		log.WithContext(ctx).Debugf("failed to verify signature %s of user %s", req.Userdata.UserPastelID, req.Userdata.UserPastelID)
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, nil
	}

	// Validation the request from Walletnode successful, we continue to get acknowledgement and confirmation process from multiple SNs
	snRequest := userdata.SuperNodeRequest{}

	// Marshal validateResult to byte array for signing
	js, err := json.Marshal(validateResult)
	if err != nil {
		return userdata.ProcessResult{}, errors.Errorf("encode validateResult %w", err)
	}
	// Hash the validateResult
	hashvalue, err := utils.Sha3256hash(js)
	if err != nil {
		return userdata.ProcessResult{}, errors.Errorf("hash userdata %w", err)
	}
	snRequest.UserdataResultHash = hex.EncodeToString(hashvalue)
	snRequest.UserdataHash = req.UserdataHash

	task.ownSNData = snRequest // At this step this SuperNodeRequest only contain UserdataResultHash and UserdataHash

	var actionErr error
	<-task.NewAction(func(ctx context.Context) error {
		// sign the data if not primary node
		isPrimary := false
		if task.ConnectedTo == nil {
			isPrimary = true
		}
		log.WithContext(ctx).Debugf("isPrimary: %t", isPrimary)
		if err := task.signAndSendSNDataSigned(ctx, task.ownSNData, isPrimary); err != nil {
			actionErr = errors.Errorf("signed and send SuperNodeRequest:%w", err)
		}
		return nil
	})
	if actionErr != nil {
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, errors.Errorf("error on action: %w", actionErr)
	}

	// only primary node start this action
	var processResult userdata.ProcessResult
	if task.ConnectedTo == nil {
		<-task.NewAction(func(ctx context.Context) error {
			log.WithContext(ctx).Debug("waiting for signature from peers")
			for {
				select {
				case <-ctx.Done():
					err := ctx.Err()
					if err != nil {
						actionErr = errors.Errorf("waiting for SuperNodeRequest from peers cancelled or timeout: %w", err)
					}
					return nil
				case <-task.allPeersSNDatasReceived:
					log.WithContext(ctx).Debug("all SuperNodeRequest received so start validation")
					resp, err := task.verifyPeersUserdata(ctx)
					if err != nil {
						actionErr = errors.Errorf("fail to verifyPeersUserdata: %w", err)
					}
					processResult = resp
					return nil
				}
			}
		})
	}

	return processResult, actionErr
}

// ReceiveUserdata get the userdata from database
func (task *UserDataTask) ReceiveUserdata(ctx context.Context, userpastelid string) (r userdata.ProcessRequest, err error) {
	log.WithContext(ctx).Debugf("ReceiveUserdata on user PastelID: %s", userpastelid)

	if task.Service.databaseOps == nil {
		return r, errors.New("databaseOps is nil")
	}

	// only primary node start this action
	return task.Service.databaseOps.ReadUserData(ctx, userpastelid)
}

// Sign and send SNDataSigned if not primary
func (task *UserDataTask) signAndSendSNDataSigned(ctx context.Context, sndata userdata.SuperNodeRequest, isPrimary bool) error {
	log.WithContext(ctx).Debugf("signAndSendSNDataSigned begin to sign SuperNodeRequest")
	signature, err := task.pastelClient.Sign(ctx, []byte(sndata.UserdataHash+sndata.UserdataResultHash), task.config.PastelID, task.config.PassPhrase, "ed448")
	if err != nil {
		return errors.Errorf("sign sndata %w", err)
	}
	if !isPrimary {
		sndata.HashSignature = hex.EncodeToString(signature)
		sndata.NodeID = task.config.PastelID
		log.WithContext(ctx).Debug("send signed sndata to primary node")
		if _, err := task.ConnectedTo.ProcessUserdataInterface.SendUserdataToPrimary(ctx, sndata); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s %w", task.ConnectedTo.ID, task.ConnectedTo.Address, err)
		}
	}
	return nil
}

// AddPeerSNDataSigned gather all the Userdata with signature from other Supernodes
func (task *UserDataTask) AddPeerSNDataSigned(ctx context.Context, snrequest userdata.SuperNodeRequest) error {
	log.WithContext(ctx).Debugf("addPeerSNDataSigned begin to aggregate SuperNodeRequest")

	task.peersSNDataSignedMtx.Lock()
	defer task.peersSNDataSignedMtx.Unlock()

	var actionErr error
	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive supernode userdata result signed from node %s", snrequest.NodeID)
		if node := task.accepted.ByID(snrequest.NodeID); node == nil {
			actionErr = errors.Errorf("node %s not in accepted list", snrequest.NodeID)
			return nil
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
	return actionErr
}

func (task *UserDataTask) verifyPeersUserdata(ctx context.Context) (userdata.ProcessResult, error) {
	log.WithContext(ctx).Debugf("all supernode data signed received so start validation")
	successCount := 0
	dataMatchingCount := 0
	for _, sndata := range task.peersSNDataSigned {
		// Verify the data
		signature, err := hex.DecodeString(sndata.HashSignature)
		if err != nil {
			log.WithContext(ctx).Debugf("failed to decode signature %s of node %s", sndata.HashSignature, sndata.NodeID)
			continue
		}

		ok, err := task.pastelClient.Verify(ctx, []byte(sndata.UserdataHash+sndata.UserdataResultHash), string(signature), sndata.NodeID, "ed448")
		if err != nil {
			log.WithContext(ctx).Debugf("failed to verify signature %s of node %s", sndata.HashSignature, sndata.NodeID)
			continue
		}

		if !ok {
			log.WithContext(ctx).Debugf("signature of node %s mistmatch", sndata.NodeID)
			continue
		}

		successCount++
		if task.ownSNData.UserdataHash == sndata.UserdataHash &&
			task.ownSNData.UserdataResultHash == sndata.UserdataResultHash {
			dataMatchingCount++
		}
	}

	if successCount < task.config.MinimalNodeConfirmSuccess-1 {
		// If there is not enough signed data from other supernodes
		// or signed data cannot be verify
		// or signed data have signature mismatch
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorNotEnoughSupernodeConfirm,
			Detail:       userdata.Description[userdata.ErrorNotEnoughSupernodeConfirm],
		}, nil
	} else if dataMatchingCount < task.config.MinimalNodeConfirmSuccess-1 {
		// If the userdata between supernodes is not matching with each other
		return userdata.ProcessResult{
			ResponseCode: userdata.ErrorUserdataMismatchBetweenSupernode,
			Detail:       userdata.Description[userdata.ErrorUserdataMismatchBetweenSupernode],
		}, nil
	}
	// Success
	return userdata.ProcessResult{
		ResponseCode: userdata.SuccessVerifyAllSignature,
		Detail:       userdata.Description[userdata.SuccessVerifyAllSignature],
	}, nil
}

func (task *UserDataTask) validateUserdata(req *userdata.ProcessRequest) (userdata.ProcessResult, error) {
	result := userdata.ProcessResult{}

	contentValidation := userdata.SuccessValidateContent

	// Biography validation
	if len(req.Biography) > biographyTextLengthLimit {
		result.Biography = fmt.Sprintf("Biography text length is greater than the limit %d characters", biographyTextLengthLimit)
		contentValidation = userdata.ErrorOnContent
	}

	// Primary Language validation
	// TODO: The “Primary Language” field might be limited to a list of language names that we can validate against
	// (this way we can avoid the problem of users writing in “Russian” and “русский” which creates confusion and data fragmentation).

	// FacebookLink validation
	if len(req.FacebookLink) > 0 && !(strings.Contains(strings.ToLower(req.FacebookLink), facebookLongURL) || strings.Contains(strings.ToLower(req.FacebookLink), facebookShortURL)) {
		result.Biography = "Facebook Link is not valid"
		contentValidation = userdata.ErrorOnContent
	}

	// Image validation
	// Image Extension validation
	isAvatarMatchExtention := false
	isCoverPhotoMatchExtention := false
	for _, extension := range imageAllowExtension {
		if strings.Contains(strings.ToLower(req.AvatarImage.Filename), extension) {
			isAvatarMatchExtention = true
		}
		if strings.Contains(strings.ToLower(req.CoverPhoto.Filename), extension) {
			isCoverPhotoMatchExtention = true
		}
	}
	if len(req.AvatarImage.Filename) > 0 && !isAvatarMatchExtention {
		result.AvatarImage = fmt.Sprintf("Avatar extension must be in the following: %s", strings.Join(imageAllowExtension, ","))
		contentValidation = userdata.ErrorOnContent
	}
	if len(req.CoverPhoto.Filename) > 0 && !isCoverPhotoMatchExtention {
		result.CoverPhoto = fmt.Sprintf("CoverPhoto extension must be in the following: %s", strings.Join(imageAllowExtension, ","))
		contentValidation = userdata.ErrorOnContent
	}

	// Image Size validation
	if len(req.AvatarImage.Filename) > 0 && (len(req.AvatarImage.Content) < imageMinSizeLimit || len(req.AvatarImage.Content) > imageMaxSizeLimit) {
		result.AvatarImage = fmt.Sprintf("Avatar Image size must be in the range: [%d-%d]", imageMinSizeLimit, imageMaxSizeLimit)
		contentValidation = userdata.ErrorOnContent
	}
	if len(req.CoverPhoto.Filename) > 0 && (len(req.CoverPhoto.Content) < imageMinSizeLimit || len(req.CoverPhoto.Content) > imageMaxSizeLimit) {
		result.CoverPhoto = fmt.Sprintf("Cover Photo size must be in the range: [%d-%d]", imageMinSizeLimit, imageMaxSizeLimit)
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

func (task *UserDataTask) pastelNodeByExtKey(ctx context.Context, nodeID string) (*Node, error) {
	masterNodes, err := task.pastelClient.MasterNodesTop(ctx)
	// log.WithContext(ctx).Debugf("master node %s", masterNodes)

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

func (task *UserDataTask) getRQliteLeaderNode(ctx context.Context, extAddress string) (*Node, error) {
	log.WithContext(ctx).Debugf("getRQliteLeaderNode node %s", extAddress)

	node := &Node{
		client:  task.Service.nodeClient,
		ID:      "", // We don't need to care about ID of the leader node.
		Address: extAddress,
	}
	return node, nil
}

// ConnectToLeader connects to RQLite Leader node
func (task *UserDataTask) ConnectToLeader(ctx context.Context, extAddress string) error {
	log.WithContext(ctx).Debugf("ConnectToLeader on address %s", extAddress)
	var actionErr error
	task.NewAction(func(ctx context.Context) error {
		node, err := task.getRQliteLeaderNode(ctx, extAddress)
		if err != nil {
			actionErr = errors.Errorf("error while acquiring rqlite leader node: %w", err)
			return nil
		}

		if err := node.connect(ctx); err != nil {
			actionErr = errors.Errorf("error while connecting to rqlite leader node: %w", err)
			return nil
		}

		task.ConnectedToLeader = node

		return nil
	})
	return actionErr
}
func (task *UserDataTask) removeArtifacts() {
}

// NewUserDataTask returns a new Task instance.
func NewUserDataTask(service *Service) *UserDataTask {
	return &UserDataTask{
		SuperNodeTask:           common.NewSuperNodeTask(logPrefix),
		Service:                 service,
		peersSNDataSignedMtx:    &sync.Mutex{},
		peersSNDataSigned:       make(map[string]userdata.SuperNodeRequest),
		allPeersSNDatasReceived: make(chan struct{}),
	}
}
