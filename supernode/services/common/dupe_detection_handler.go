package common

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
)

// Tasker interface to call back concrete Task for specific actions
type Tasker interface {
	SendDDFBack(ctx context.Context, node node.SuperNodePeerAPIInterface, nodeInfo *types.MeshedSuperNode, pastelID string, data []byte) error
}

// DupeDetectionHandler common operations related to Dupe Detection and Fingerprints processing
type DupeDetectionHandler struct {
	*SuperNodeTask
	*RegTaskHelper

	DdClient ddclient.DDServerClient

	ddMtx sync.Mutex

	myDDAndFingerprints         *pastel.DDAndFingerprints
	calculatedDDAndFingerprints *pastel.DDAndFingerprints

	// DDAndFingerprints created by node its self
	allDDAndFingerprints                  map[string]*pastel.DDAndFingerprints
	allSignedDDAndFingerprintsReceivedChn chan struct{}
	isOpenAPI                             bool
	groupID                               string
	collectionName                        string
}

// SetDDFields sets isOpenAPI, groupID, collectionName
func (h *DupeDetectionHandler) SetDDFields(isOpenAPI bool, groupID, collectionName string) {
	h.isOpenAPI = isOpenAPI
	h.groupID = groupID
	h.collectionName = collectionName
}

// NewDupeDetectionTaskHelper creates instance of DupeDetectionHandler
func NewDupeDetectionTaskHelper(task *SuperNodeTask,
	ddClient ddclient.DDServerClient, pastelID string, passPhrase string, network *NetworkHandler, pastelClient pastel.Client,
	preburntTxMinConfirmations int) *DupeDetectionHandler {
	return &DupeDetectionHandler{
		SuperNodeTask:                         task,
		DdClient:                              ddClient,
		allDDAndFingerprints:                  map[string]*pastel.DDAndFingerprints{},
		allSignedDDAndFingerprintsReceivedChn: make(chan struct{}),
		RegTaskHelper: NewRegTaskHelper(task, pastelID, passPhrase, network, pastelClient,
			preburntTxMinConfirmations),
	}
}

// ProbeImage uploads the resampled image compute and return a compression of pastel.DDAndFingerprints
//
//	Implementing https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration starting with step 4.A.3
//	Call dd-service to generate near duplicate fingerprints and dupe-detection info from re-sampled image (img1-r)
func (h *DupeDetectionHandler) ProbeImage(_ context.Context, file *files.File, blockHash string, blockHeight string, timestamp string, creatorPastelID string, tasker Tasker) ([]byte, error) {
	if err := h.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	var err error
	var retCompressed, compressed []byte

	<-h.NewAction(func(ctx context.Context) error {
		defer errors.Recover(func(recErr error) {
			log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).WithError(recErr).Error("PanicWhenProbeImage")
		})

		// Begin send signed DDAndFingerprints to other SNs
		// Send base64’ed (and compressed) dd_and_fingerprints and its signature to the 2 OTHER SNs
		// step 4.A.5
		if len(h.NetworkHandler.meshedNodes) != 3 {
			log.WithContext(ctx).Error("Not enough meshed SuperNodes")
			err = errors.New("not enough meshed SuperNodes")
			return err
		}

		registeringSupernode1 := h.NetworkHandler.meshedNodes[0].NodeID
		registeringSupernode2 := h.NetworkHandler.meshedNodes[1].NodeID
		registeringSupernode3 := h.NetworkHandler.meshedNodes[2].NodeID

		//SuperNode makes ImageRarenessScore gRPC call to dd-service
		//This is the original python logic
		// if valid_image:
		// 	sense_output_json_file_path = sense_api_json_results_path + sha256_hash_of_image_file + '.' + 'json'
		// 	image_filepath = output_path_for_image_file.replace('./','/home/ubuntu/sense_demo_api/')
		// 	pastel_block_hash_when_request_submitted = '00000000819d1ddc92d6a2f78d59c2554c3d47f4ec443c286fde29f84523a607'
		// 	pastel_block_height_when_request_submitted = '240669'
		// 	utc_timestamp_when_request_submitted = datetime.strftime(datetime.utcnow(), '%Y-%m-%d %H:%M:%S')
		// 	pastel_id_of_submitter = 'jXIGc0SbVWHFRKemuDSpkTIrCybP3tFbT2c5LOZZuKgBZf95eFZ62QkBP4mmFoZrQkVB0bmn70Oran3vg5fW1X'
		// 	pastel_id_of_registeringSupernode1 = 'jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis'
		// 	pastel_id_of_registeringSupernode2 = 'jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F'
		// 	pastel_id_of_registeringSupernode3 = 'jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ'
		// 	isPastelOpenapiRequest = 'False'

		// 	command_string = 'python3 /home/ubuntu/pastel/dd-service/tools/grpc-client.py 50052 ' + \
		// 	image_filepath + ' ' + \
		// 	pastel_block_hash_when_request_submitted + ' ' + \
		// 	pastel_block_height_when_request_submitted + ' ' + \
		// 	utc_timestamp_when_request_submitted + ' ' + \
		// 	pastel_id_of_submitter + ' ' + \
		// 	pastel_id_of_registeringSupernode1 + ' ' + \
		// 	pastel_id_of_registeringSupernode2 + ' ' + \
		// 	pastel_id_of_registeringSupernode3 + ' ' + \
		// 	isPastelOpenapiRequest + ' ' + \
		//
		log.WithContext(ctx).Info("asking dd server to process image")
		compressed, err = h.GenFingerprintsData(ctx, file, blockHash, blockHeight, timestamp, creatorPastelID, registeringSupernode1,
			registeringSupernode2, registeringSupernode3, h.isOpenAPI, h.groupID, h.collectionName)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate fingerprints data")
			err = errors.Errorf("generate fingerprints data: %w", err)
			return err
		}

		h.UpdateStatus(StatusImageProbed)

		for _, nodeInfo := range h.NetworkHandler.meshedNodes {
			// Don't send to itself
			if nodeInfo.NodeID == h.ServerPastelID {
				continue
			}

			var node *SuperNodePeer
			node, err = h.NetworkHandler.PastelNodeByExtKey(ctx, nodeInfo.NodeID)
			if err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("get node by extID")

				err = errors.Errorf("get node by extID: %w", err)
				return err
			}

			if err = node.Connect(ctx); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("connect to node")
				err = errors.Errorf("connect to node: %w", err)

				return err
			}
			log.WithContext(ctx).Debugf("sending dd_fp to other SN: %s", node.ID)

			//supernode/services/nftregister/task.go
			if err = tasker.SendDDFBack(ctx, node.SuperNodePeerAPIInterface, &nodeInfo, h.ServerPastelID, compressed); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"sessID":  nodeInfo.SessID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("send signed DDAndFingerprints failed")
				err = errors.Errorf("send signed DDAndFingerprints failed: %w", err)

				return err
			}
		}

		// wait for other SNs shared their signed DDAndFingerprints
		// start to implement 4.B here, waiting for SN's to return dd and fingerprints and signatures
		select {
		case <-ctx.Done():
			err = ctx.Err()
			log.WithContext(ctx).WithError(err).Error("ctx.Done() error from gen fingerprints data")
			if err != nil {
				log.WithContext(ctx).Error("waiting for DDAndFingerprints from peers context cancelled")
			}
			return nil
		case <-h.allSignedDDAndFingerprintsReceivedChn:
			//verification has already been completed in the Add method
			log.WithContext(ctx).Debug("all DDAndFingerprints received so start calculate final DDAndFingerprints")
			h.allDDAndFingerprints[h.ServerPastelID] = h.myDDAndFingerprints

			// get list of DDAndFingerprints in order of node rank
			dDAndFingerprintsList := []*pastel.DDAndFingerprints{}
			for _, node := range h.NetworkHandler.meshedNodes {
				v, ok := h.allDDAndFingerprints[node.NodeID]
				if !ok {
					err = errors.Errorf("not found DDAndFingerprints of node : %s", node.NodeID)
					log.WithContext(ctx).WithFields(log.Fields{
						"nodeID": node.NodeID,
					}).Errorf("DDAndFingerprints of node not found")
					return err
				}
				dDAndFingerprintsList = append(dDAndFingerprintsList, v)
			}

			// calculate final result from DdDAndFingerprints from all SNs node
			if len(dDAndFingerprintsList) != 3 {
				err = errors.Errorf("not enough DDAndFingerprints, len: %d", len(dDAndFingerprintsList))
				log.WithContext(ctx).WithFields(log.Fields{
					"list": dDAndFingerprintsList,
				}).Errorf("not enough DDAndFingerprints")
				return err
			}

			h.calculatedDDAndFingerprints, err = pastel.CombineFingerPrintAndScores(
				dDAndFingerprintsList[0],
				dDAndFingerprintsList[1],
				dDAndFingerprintsList[2],
			)

			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("call CombineFingerPrintAndScores() failed")
				err = errors.Errorf("call CombineFingerPrintAndScores() failed")
				return err
			}

			// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
			retCompressed, err = h.compressAndSignDDAndFingerprints(ctx, h.calculatedDDAndFingerprints)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("compress combine DDAndFingerPrintAndScore failed")

				err = errors.Errorf("compress combine DDAndFingerPrintAndScore failed: %w", err)
				return err
			}

			log.WithContext(ctx).Debug("DDAndFingerprints combined and compressed")
			return nil
		case <-time.After(25 * time.Minute):
			log.WithContext(ctx).Error("waiting for DDAndFingerprints from peers timeout")
			err = errors.New("waiting for DDAndFingerprints timeout")
			return err
		}
	})

	if err != nil {
		return nil, err
	}

	return retCompressed, nil
}

// GenFingerprintsData calls DD server to get DD and FP data
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Call dd-service to generate near duplicate fingerprints and dupe-detection info from re-sampled image (img1-r)
// Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
// Step 4.A.3 - 4.A.4
// GenFingerprintsData(ctx, file, blockHash, blockHeight, timestamp, creatorPastelID, registeringSupernode1, registeringSupernode2, registeringSupernode3, false, "")
func (h *DupeDetectionHandler) GenFingerprintsData(ctx context.Context, file *files.File, blockHash string, blockHeight string, timestamp string,
	creatorPastelID string, registeringSupernode1 string, registeringSupernode2 string, registeringSupernode3 string,
	isPastelOpenapiRequest bool, groupID string, collectionName string) ([]byte, error) {
	img, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("get content of image %s: %w", file.Name(), err)
	}

	// Get DDAndFingerprints
	// SuperNode makes ImageRarenessScore gRPC call to dd-service
	//   ctx context.Context, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string, sn_1 string, sn_2 string, sn_3 string, openapi_request bool,
	h.myDDAndFingerprints, err = h.DdClient.ImageRarenessScore(
		ctx,
		img,
		file.Format().String(),
		blockHash,
		blockHeight,
		timestamp,
		creatorPastelID,
		registeringSupernode1,
		registeringSupernode2,
		registeringSupernode3,
		isPastelOpenapiRequest,
		groupID,
		collectionName,
	)

	if err != nil {
		return nil, errors.Errorf("call ImageRarenessScore(): %w", err)
	}

	// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
	//Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
	compressed, err := h.compressAndSignDDAndFingerprints(ctx, h.myDDAndFingerprints)
	if err != nil {
		return nil, errors.Errorf("call compressSignedDDAndFingerprints failed: %w", err)
	}

	return compressed, nil
}

// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step 4.A.4
// Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
func (h *DupeDetectionHandler) compressAndSignDDAndFingerprints(ctx context.Context, ddData *pastel.DDAndFingerprints) ([]byte, error) {
	// JSON marshal the ddData
	ddDataBytes, err := json.Marshal(ddData)
	if err != nil {
		return nil, errors.Errorf("marshal DDAndFingerprints: %w", err)
	}

	// sign it
	signature, err := h.PastelHandler.PastelClient.Sign(ctx, ddDataBytes, h.ServerPastelID, h.serverPassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, errors.Errorf("sign DDAndFingerprints: %w", err)
	}

	// Compress the data and signature, there is a separator byte
	compressed, err := pastel.ToCompressSignedDDAndFingerprints(ddData, signature)
	if err != nil {
		return nil, errors.Errorf("compress SignedDDAndFingerprints: %w", err)
	}

	return compressed, nil
}

// AddSignedDDAndFingerprints adds signed dd and fp
//
//	https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
//	Step 4.A.5 - Send base64’ed (and compressed) dd_and_fingerprints and its signature to the 2 OTHER SNs
func (h *DupeDetectionHandler) AddSignedDDAndFingerprints(nodeID string, compressedSignedDDAndFingerprints []byte) error {
	h.ddMtx.Lock()
	defer h.ddMtx.Unlock()

	var err error

	<-h.NewAction(func(ctx context.Context) error {
		defer errors.Recover(func(recErr error) {
			log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).WithError(recErr).Error("PanicWhenProbeImage")
		})
		log.WithContext(ctx).Debugf("receive compressedSignedDDAndFingerprints from node %s", nodeID)
		// Check nodeID should in task.meshedNodes
		err = h.NetworkHandler.CheckNodeInMeshedNodes(nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("check if node existed in meshedNodes failed")
			return nil
		}

		// extract data
		var ddAndFingerprints *pastel.DDAndFingerprints
		var ddAndFingerprintsBytes []byte
		var signature []byte

		ddAndFingerprints, ddAndFingerprintsBytes, signature, err = pastel.ExtractCompressSignedDDAndFingerprints(compressedSignedDDAndFingerprints)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("extract compressedSignedDDAndFingerprints failed")
			return nil
		}

		var ok bool
		// Step 4.B.2
		// Verify signatures of dd_and_fingerprints returned by 2 other SuperNodes are correct:
		//If some are not the same - SuperNode must stop processing and return Error
		ok, err = h.PastelHandler.PastelClient.Verify(ctx, ddAndFingerprintsBytes, string(signature), nodeID, pastel.SignAlgorithmED448)
		if err != nil || !ok {
			err = errors.New("signature verification failed")
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("verify signature of ddAndFingerprintsBytes failed")
			return nil
		}

		h.allDDAndFingerprints[nodeID] = ddAndFingerprints

		// if enough ddAndFingerprints received (not include from itself)
		if len(h.allDDAndFingerprints) == 2 {
			log.WithContext(ctx).Debug("all ddAndFingerprints received")
			go func() {
				close(h.allSignedDDAndFingerprintsReceivedChn)
			}()
		}
		return nil
	})

	return err
}
