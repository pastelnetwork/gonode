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
}

// NewDupeDetectionTaskHelper creates instance of DupeDetectionHandler
func NewDupeDetectionTaskHelper(task *SuperNodeTask,
	ddClient ddclient.DDServerClient,
	pastelID string, passPhrase string,
	network *NetworkHandler,
	pastelClient pastel.Client,
	preburntTxMinConfirmations int,
) *DupeDetectionHandler {
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
//  Implementing https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration starting with step 4.A.3
//  Call dd-service to generate near duplicate fingerprints and dupe-detection info from re-sampled image (img1-r)
func (h *DupeDetectionHandler) ProbeImage(_ context.Context, file *files.File, blockHash string, creatorPastelID string, tasker Tasker) ([]byte, error) {
	if err := h.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	var err error

	<-h.NewAction(func(ctx context.Context) error {
		defer errors.Recover(func(recErr error) {
			log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).WithError(recErr).Error("PanicWhenProbeImage")
		})
		h.UpdateStatus(StatusImageProbed)

		//SuperNode makes ImageRarenessScore gRPC call to dd-service
		h.myDDAndFingerprints, err = h.GenFingerprintsData(ctx, file, blockHash, creatorPastelID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate fingerprints data")
			err = errors.Errorf("generate fingerprints data: %w", err)
			return nil
		}

		// Begin send signed DDAndFingerprints to other SNs
		// Send base64’ed (and compressed) dd_and_fingerprints and its signature to the 2 OTHER SNs
		// step 4.A.5
		if h.NetworkHandler.meshedNodes == nil || len(h.NetworkHandler.meshedNodes) != 3 {
			log.WithContext(ctx).Error("Not enough meshed SuperNodes")
			err = errors.New("not enough meshed SuperNodes")
			return nil
		}

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
				return nil
			}

			if err = node.Connect(ctx); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("connect to node")
				return nil
			}

			//supernode/services/nftregister/task.go
			if err = tasker.SendDDFBack(ctx, node.SuperNodePeerAPIInterface, &nodeInfo, h.ServerPastelID, h.myDDAndFingerprints.ZstdCompressedFingerprint); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"sessID":  nodeInfo.SessID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("send signed DDAndFingerprints failed")
				return nil
			}
		}

		// wait for other SNs shared their signed DDAndFingerprints
		// start to implement 4.B here, waiting for SN's to return dd and fingerprints and signatures
		select {
		case <-ctx.Done():
			err = ctx.Err()
			log.WithContext(ctx).WithError(err).Error("ctx.Done()")
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
					return nil
				}
				dDAndFingerprintsList = append(dDAndFingerprintsList, v)
			}

			// calculate final result from DdDAndFingerprints from all SNs node
			if len(dDAndFingerprintsList) != 3 {
				err = errors.Errorf("not enough DDAndFingerprints, len: %d", len(dDAndFingerprintsList))
				log.WithContext(ctx).WithFields(log.Fields{
					"list": dDAndFingerprintsList,
				}).Errorf("not enough DDAndFingerprints")
				return nil
			}

			h.calculatedDDAndFingerprints, err = pastel.CombineFingerPrintAndScores(
				dDAndFingerprintsList[0],
				dDAndFingerprintsList[1],
				dDAndFingerprintsList[2],
			)

			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("call CombineFingerPrintAndScores() failed")
				return nil
			}

			// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
			var compressed []byte
			compressed, err = h.compressAndSignDDAndFingerprints(ctx, h.calculatedDDAndFingerprints)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("compress combine DDAndFingerPrintAndScore failed")

				err = errors.Errorf("compress combine DDAndFingerPrintAndScore failed: %w", err)
				return nil
			}
			h.calculatedDDAndFingerprints.ZstdCompressedFingerprint = compressed
			return nil
		case <-time.After(30 * time.Second):
			log.WithContext(ctx).Error("waiting for DDAndFingerprints from peers timeout")
			err = errors.New("waiting for DDAndFingerprints timeout")
			return nil
		}
	})

	if err != nil {
		return nil, err
	}

	return h.calculatedDDAndFingerprints.ZstdCompressedFingerprint, nil
}

// GenFingerprintsData calls DD server to get DD and FP data
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Call dd-service to generate near duplicate fingerprints and dupe-detection info from re-sampled image (img1-r)
// Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
// Step 4.A.3 - 4.A.4
func (h *DupeDetectionHandler) GenFingerprintsData(ctx context.Context, file *files.File, blockHash string, creatorPastelID string) (*pastel.DDAndFingerprints, error) {
	img, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("get content of image %s: %w", file.Name(), err)
	}

	// Get DDAndFingerprints
	// SuperNode makes ImageRarenessScore gRPC call to dd-service
	ddAndFingerprints, err := h.DdClient.ImageRarenessScore(
		ctx,
		img,
		file.Format().String(),
		blockHash,
		creatorPastelID,
	)

	if err != nil {
		return nil, errors.Errorf("call ImageRarenessScore(): %w", err)
	}

	// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
	//Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
	compressed, err := h.compressAndSignDDAndFingerprints(ctx, ddAndFingerprints)
	if err != nil {
		return nil, errors.Errorf("call compressSignedDDAndFingerprints failed: %w", err)
	}

	ddAndFingerprints.ZstdCompressedFingerprint = compressed
	return ddAndFingerprints, nil
}

//  https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
//  Step 4.A.4
//  Sign dd_and_fingerprints with SN own PastelID (private key) using cNode API
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
//  https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
//  Step 4.A.5 - Send base64’ed (and compressed) dd_and_fingerprints and its signature to the 2 OTHER SNs
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
