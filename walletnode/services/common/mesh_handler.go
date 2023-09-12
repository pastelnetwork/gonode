package common

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	//DDServerPendingRequestsThreshold is the threshold for DD-server pending requests
	DDServerPendingRequestsThreshold = 0
	// DDServerExecutingRequestsThreshold is the threshold for DD-server executing requests
	DDServerExecutingRequestsThreshold = 2
)

// MeshHandler provides networking functionality, including mesh
type MeshHandler struct {
	task          *WalletNodeTask
	pastelHandler *mixins.PastelHandler

	nodeMaker  node.RealNodeMaker
	nodeClient node.ClientInterface

	// TODO: make config interface and use it here instead of individual items
	minNumberSuperNodes    int
	connectToNodeTimeout   time.Duration
	acceptNodesTimeout     time.Duration
	connectToNextNodeDelay time.Duration

	callersPastelID string
	passphrase      string

	Nodes                         SuperNodeList
	UseMaxNodes                   bool
	checkDDDatabaseHashes         bool
	HashCheckMaxRetries           int
	requireSNAgreementOnMNTopList bool
	logRequestID                  string
}

// MeshHandlerOpts set of options to pass to NewMeshHandler
type MeshHandlerOpts struct {
	Task          *WalletNodeTask
	NodeClient    node.ClientInterface
	NodeMaker     node.RealNodeMaker
	PastelHandler *mixins.PastelHandler
	Configs       *MeshHandlerConfig
	LogRequestID  string
}

type DDStats struct {
	WaitingInQueue int32
	Executing      int32
	IsReady        bool
}

type SNData struct {
	TopMNsList map[string][]string
	Hashes     map[string]string
	DDStats    map[string]DDStats
}

// MeshHandlerConfig config subset used by MeshHandler
type MeshHandlerConfig struct {
	PastelID                      string
	Passphrase                    string
	MinSNs                        int
	ConnectToNodeTimeout          time.Duration
	AcceptNodesTimeout            time.Duration
	ConnectToNextNodeDelay        time.Duration
	UseMaxNodes                   bool
	CheckDDDatabaseHashes         bool
	HashCheckMaxRetries           int
	RequireSNAgreementOnMNTopList bool
}

// NewMeshHandler returns new NewMeshHandler
func NewMeshHandler(opts MeshHandlerOpts) *MeshHandler {
	return &MeshHandler{
		task:                          opts.Task,
		nodeMaker:                     opts.NodeMaker,
		pastelHandler:                 opts.PastelHandler,
		nodeClient:                    opts.NodeClient,
		callersPastelID:               opts.Configs.PastelID,
		passphrase:                    opts.Configs.Passphrase,
		minNumberSuperNodes:           opts.Configs.MinSNs,
		connectToNodeTimeout:          opts.Configs.ConnectToNodeTimeout,
		acceptNodesTimeout:            opts.Configs.AcceptNodesTimeout,
		connectToNextNodeDelay:        opts.Configs.ConnectToNextNodeDelay,
		UseMaxNodes:                   opts.Configs.UseMaxNodes,
		checkDDDatabaseHashes:         opts.Configs.CheckDDDatabaseHashes,
		HashCheckMaxRetries:           opts.Configs.HashCheckMaxRetries,
		requireSNAgreementOnMNTopList: opts.Configs.RequireSNAgreementOnMNTopList,
		logRequestID:                  opts.LogRequestID,
	}
}

// SetupMeshOfNSupernodesNodes sets Mesh
func (m *MeshHandler) SetupMeshOfNSupernodesNodes(ctx context.Context) (int, string, error) {
	log.WithContext(ctx).Info("SetupMeshOfNSupernodesNodes Starting...")

	// Get current block height & hash
	blockNum, blockHash, err := m.pastelHandler.GetBlock(ctx)
	if err != nil {
		return -1, "", errors.Errorf("get block: %v", err)
	}
	log.WithContext(ctx).Infof("Current block is %d", blockNum)

	connectedNodes, err := m.findNValidTopSuperNodes(ctx, m.minNumberSuperNodes, []string{})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("failed validating of %d SNs", m.minNumberSuperNodes)
		return 0, "", errors.Errorf("failed validating of %d SNs: %w", m.minNumberSuperNodes, err)
	}

	meshedNodes, err := m.setMesh(ctx, connectedNodes, m.minNumberSuperNodes)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("req-id", m.logRequestID).Errorf("failed setting up Mesh of %d SNs from %d SNs", m.minNumberSuperNodes, len(connectedNodes))
		return 0, "", errors.Errorf("failed setting up Mesh of %d SNs from %d SNs %w", m.minNumberSuperNodes, len(connectedNodes), err)
	}
	m.Nodes = meshedNodes

	return blockNum, blockHash, nil
}

// ConnectToNSuperNodes sets single simple connection to N SNs
func (m *MeshHandler) ConnectToNSuperNodes(ctx context.Context, n int, skipNodes []string) error {

	connectedNodes, err := m.findNValidTopSuperNodes(ctx, n, skipNodes)
	if err != nil {
		return err
	}

	// Activate supernodes that are in the mesh.
	connectedNodes.Activate()
	// Cancel context when any connection is broken.
	m.task.UpdateStatus(StatusConnected)

	m.Nodes = connectedNodes

	return nil
}

// ConnectToTopRankNodes - find N supernodes and create mesh of 3 nodes
func (m *MeshHandler) findNValidTopSuperNodes(ctx context.Context, n int, skipNodes []string) (SuperNodeList, error) {

	// Retrieve supernodes with the highest ranks.
	WNTopNodes, err := m.getTopNodes(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("req-id", m.logRequestID).Errorf("Could not get top masternodes")
		return nil, errors.Errorf("call masternode top: %w", err)
	}

	if len(WNTopNodes) < n {
		m.task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		log.WithContext(ctx).WithError(err).WithField("req-id", m.logRequestID).Errorf("unable to find enough Supernodes: %d", n)
		return nil, fmt.Errorf("unable to find enough Supernodes: %d", n)
	}
	log.WithContext(ctx).Infof("Found %d Supernodes", len(WNTopNodes))

	candidateNodes, err := m.GetCandidateNodes(ctx, WNTopNodes, n)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting candidate nodes")
	}

	if len(candidateNodes) < n {
		err := errors.New("not enough candidate nodes found with required parameters")
		log.WithContext(ctx).WithField("count", len(candidateNodes)).WithError(err)
		return nil, err
	}

	// Connect to top nodes to find 3SN and validate their info
	foundNodes, err := m.connectToAndValidateSuperNodes(ctx, candidateNodes, n, m.HashCheckMaxRetries, skipNodes)
	if err != nil {
		m.task.UpdateStatus(StatusErrorFindRespondingSNs)
		log.WithContext(ctx).WithError(err).Errorf("unable to validate MN")
		return nil, errors.Errorf("validate MNs info: %v", err)
	}
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("Found %d valid Supernodes", len(foundNodes))
	return foundNodes, nil
}

func (m *MeshHandler) GetCandidateNodes(ctx context.Context, candidatesNodes SuperNodeList, n int) (SuperNodeList, error) {
	if !m.checkDDDatabaseHashes && !m.requireSNAgreementOnMNTopList {
		return candidatesNodes, nil
	}

	log.WithContext(ctx).WithField("given nodes", len(candidatesNodes)).Info("GetCandidateNodes Starting...")

	WNTopNodesList := SuperNodeList{}
	topMap := make(map[string][]string)
	hashes := make(map[string]string)
	dataMap := make(map[string]DDStats)
	var err error

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  pastel.SignAlgorithmED448,
	}

	// Get WN top nodes list from all the candidate nodes
	for _, someNode := range candidatesNodes {
		top, hash, data, err := m.GetRequiredDataFromSN(ctx, *someNode, secInfo)
		if err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("failed to get required data from supernode - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			continue
		}

		WNTopNodesList.Add(someNode)
		topMap[someNode.Address()] = top
		hashes[someNode.Address()] = hash
		dataMap[someNode.Address()] = data
	}

	if len(WNTopNodesList) < n {
		return nil, errors.New("failed to get required data from all candidate nodes")
	}

	if m.requireSNAgreementOnMNTopList {

		// Filter out nodes with different top MNs list
		WNTopNodesList, err = m.filterByTopNodesList(ctx, WNTopNodesList, topMap)
		if err != nil {
			return nil, fmt.Errorf("filter by top nodes list: %w", err)
		}

		log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("Matching top nodes list for nodes %d", len(WNTopNodesList))
	}

	if len(WNTopNodesList) < n {
		return nil, errors.New("failed to get enough nodes with matching top 10 list")
	}

	if m.checkDDDatabaseHashes {
		WNTopNodesList, err = m.filterByDDHash(ctx, WNTopNodesList, hashes)
		if err != nil {
			return nil, fmt.Errorf("filter by top nodes list: %w", err)
		}

		if len(WNTopNodesList) < n {
			return nil, errors.New("failed to get enough nodes with matching fingerprints database hash")
		}

		WNTopNodesList, err = m.filterDDServerRequestsInWaitingQueue(ctx, WNTopNodesList, dataMap)
		if err != nil {
			return nil, fmt.Errorf("filter by dd-service stats: %w", err)
		}

	}

	return WNTopNodesList, nil
}

func (m *MeshHandler) setMesh(ctx context.Context, candidatesNodes SuperNodeList, n int) (SuperNodeList, error) {
	log.WithContext(ctx).Debug("setMesh Starting...")

	var errs error
	var err error

	var meshedNodes SuperNodeList
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Info("connected nodes: ", candidatesNodes)
	for primaryRank := range candidatesNodes {
		meshedNodes, err = m.connectToPrimarySecondary(ctx, candidatesNodes, primaryRank)
		if err != nil {
			log.WithContext(ctx).WithError(err).Warn("Connecting to primary ans secondary failed. Disconnecting any still connected")
			// close connected connections
			m.disconnectNodes(ctx, candidatesNodes)

			if errors.IsContextCanceled(err) {
				log.WithContext(ctx).Info("context cancelled")
				return nil, err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warn("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(meshedNodes) < n {
		// close connected connections
		log.WithContext(ctx).WithError(err).WithField("req-id", m.logRequestID).Warnf("Could not create a mesh of %d nodes. Disconnecting any still connected", n)
		m.disconnectNodes(ctx, meshedNodes)
		return nil, errors.Errorf("Could not create a mesh of %d nodes: %w", n, errs)
	}

	// Send all meshed supernode info to nodes - that will be used by each remote node to send info to other nodes
	var meshedSNInfo []types.MeshedSuperNode
	for _, someNode := range meshedNodes {
		meshedSNInfo = append(meshedSNInfo, types.MeshedSuperNode{
			NodeID: someNode.PastelID(),
			SessID: someNode.SessID(),
		})
	}

	for _, someNode := range meshedNodes {
		err := someNode.MeshNodes(ctx, meshedSNInfo)
		if err != nil {
			m.disconnectNodes(ctx, meshedNodes)
			return nil, errors.Errorf("could not send info of meshed nodes: %w", err)
		}
	}

	// Activate supernodes that are in the mesh.
	meshedNodes.Activate()
	// Cancel context when any connection is broken.
	m.task.UpdateStatus(StatusConnected)

	return meshedNodes, nil
}

// GetTopNodes get list of current top masternodes and validate them (in the order received from pasteld)
func (m *MeshHandler) getTopNodes(ctx context.Context) (SuperNodeList, error) {

	var nodes SuperNodeList

	mns, err := m.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = m.pastelHandler.PastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes.AddNewNode(m.nodeClient, mn.ExtAddress, mn.ExtKey, m.nodeMaker)
	}
	return nodes, nil
}

func (m *MeshHandler) connectToAndValidateSuperNodes(ctx context.Context, candidatesNodes SuperNodeList, n int, maxRetries int, skipNodes []string) (SuperNodeList, error) {
	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  pastel.SignAlgorithmED448,
	}

	return m.connectToAndValidateSuperNodesWithoutHashCheck(ctx, candidatesNodes, n, secInfo, skipNodes)
}

func (m *MeshHandler) connectToAndValidateSuperNodesWithoutHashCheck(ctx context.Context, candidatesNodes SuperNodeList, n int, secInfo *alts.SecInfo, skipNodes []string) (SuperNodeList, error) {
	nodes := m.connectToNodes(ctx, candidatesNodes, n, secInfo, skipNodes)

	if len(nodes) < n {
		log.WithContext(ctx).WithField("req-id", m.logRequestID).Errorf("Failed to connect with required number of supernodes - only found %d good", n)
		return nil, errors.Errorf("validate %d Supernodes from pastel network", n)
	}

	return nodes, nil
}

func (m *MeshHandler) connectToNodes(ctx context.Context, nodesToConnect SuperNodeList, n int, secInfo *alts.SecInfo, skipNodes []string) SuperNodeList {
	nodes := SuperNodeList{}

	for _, someNode := range nodesToConnect {
		skip := false
		for i := 0; i < len(skipNodes); i++ {
			if skipNodes[i] == someNode.PastelID() {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if err := someNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("Failed to connect to Supernodes - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			continue
		}

		nodes = append(nodes, someNode)
		if len(nodes) == n && !m.UseMaxNodes {
			break
		}
	}

	return nodes
}

// filterByDDHash matches database hashes
func (m *MeshHandler) filterByDDHash(ctx context.Context, nodesList SuperNodeList, hashes map[string]string) (retList SuperNodeList, err error) {
	retList = SuperNodeList{}
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("Matching database hash for nodes %d", len(nodesList))
	counter := make(map[string]SuperNodeList)
	// Connect to top nodes to find 3SN and validate their info
	for _, someNode := range nodesList {
		hash := hashes[someNode.Address()]

		if counter[hash] == nil {
			counter[hash] = SuperNodeList{}
		}

		counter[hash] = append(counter[hash], someNode)
	}

	matchingKey := ""
	for key, val := range counter {
		if len(val) > len(retList) {
			retList = val
			matchingKey = key
		}
	}

	log.WithContext(ctx).WithField("given nodes count", len(nodesList)).WithField("matching hash nodes count", len(retList)).WithField("matching hash", matchingKey).Infof("DD Database hashes matched")

	return retList, nil
}

func (m *MeshHandler) filterDDServerRequestsInWaitingQueue(ctx context.Context, nodesList SuperNodeList, dataMap map[string]DDStats) (retList SuperNodeList, err error) {
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("checking dd-server no of requests waiting in a list of %d nodes ", len(nodesList))
	retList = SuperNodeList{}
	// Connect to top nodes to find 3SN and validate their info
	for _, someNode := range nodesList {
		waiting, executing := dataMap[someNode.Address()].WaitingInQueue, dataMap[someNode.Address()].Executing
		if waiting > DDServerPendingRequestsThreshold {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("waiting tasks", waiting).Warnf("pending requests in queue exceeds the threshold for node %s", someNode.Address())
			continue
		}

		if executing > DDServerExecutingRequestsThreshold {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("executing tasks", executing).Warnf("executing tasks exceeds the threshold for node %s", someNode.Address())
			continue
		}

		retList = append(retList, someNode)
	}

	log.WithContext(ctx).WithField("give nodes count", len(nodesList)).WithField("dd-service ready nodes", len(retList)).Infof("dd-service ready nodes matched")

	return retList, nil
}

// meshNodes establishes communication between supernodes.
func (m *MeshHandler) connectToPrimarySecondary(ctx context.Context, candidatesNodes SuperNodeList, primaryIndex int) (SuperNodeList, error) {
	log.WithContext(ctx).Debug("connectToPrimarySecondary Starting...")

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  "ed448",
	}

	primary := candidatesNodes[primaryIndex]
	log.WithContext(ctx).Debugf("Trying to connect to primary node %q", primary)
	if err := primary.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to connect to primary node %q", primary)
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := primary.Session(ctx, true); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to open session with primary node %q", primary)
		return nil, fmt.Errorf("session: %w", err)
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	// FIXME: ugly hack here
	secondariesMtx := &sync.Mutex{}
	var secondaries SuperNodeList
	go func() {
		for i, someNode := range candidatesNodes {
			sameNode := someNode

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(m.connectToNextNodeDelay):
				go func() {
					defer errors.Recover(log.Fatal)

					if err := sameNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("Failed to connect to secondary node %q", sameNode.String())
						return
					}
					if err := sameNode.Session(ctx, false); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("Failed to open session with secondary node %q", sameNode.String())
						return
					}
					// Should not run this code in go routine
					func() {
						secondariesMtx.Lock()
						defer secondariesMtx.Unlock()
						secondaries.Add(sameNode)
					}()

					if err := sameNode.ConnectTo(ctx, types.MeshedSuperNode{
						NodeID: primary.PastelID(),
						SessID: primary.SessID(),
					}); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Secondary %q connected to primary", sameNode)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, m.acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	var meshedNodes SuperNodeList
	primary.SetPrimary(true)
	meshedNodes.Add(primary)

	secondariesMtx.Lock()
	defer secondariesMtx.Unlock()

	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		someNode := secondaries.FindByPastelID(pastelID)
		if someNode == nil {
			return nil, fmt.Errorf("not found accepted node: %s", pastelID)
		}
		meshedNodes.Add(someNode)
	}

	return meshedNodes, nil
}

func (m *MeshHandler) filterByTopNodesList(ctx context.Context, WNTopNodesList SuperNodeList, topMap map[string][]string) (SuperNodeList, error) {
	maxOutliers := 3
	minRequiredNodes := 3
	minRequiredCommonIPs := 5
	itemCount := make(map[string]int)
	peerNodeList := make(map[string][]string)
	for _, localNode := range WNTopNodesList {
		itemCount[localNode.Address()] = 0
	}

	badPeers := 0
	log.WithContext(ctx).Info("getting candidate nodes")
	for _, someNode := range WNTopNodesList {
		someNodeIP := someNode.Address()
		log.WithContext(ctx).WithField("address", someNodeIP)

		peerNodeList[someNodeIP] = topMap[someNodeIP]
		for _, topNodeIP := range peerNodeList[someNodeIP] {
			itemCount[topNodeIP]++
		}
	}

	// SNs DDoS protection
	if badPeers > maxOutliers {
		return nil, fmt.Errorf("failed to get mn-top lists from %d SNs", badPeers)
	}

	candidateNodes := SuperNodeList{}
	for _, someNode := range WNTopNodesList {
		someNodeIP := someNode.Address()
		count := itemCount[someNodeIP]
		if count >= len(WNTopNodesList)-maxOutliers {
			candidateNodes = append(candidateNodes, someNode)
		}
	}

	finalCandidateNodes := SuperNodeList{}
	for _, someNode := range WNTopNodesList {
		someNodeIP := someNode.Address()
		commonCount := 0
		for _, remoteNodeIP := range peerNodeList[someNodeIP] {
			// check if the IP from the remote node is in the candidate list and often enough (>=7)
			if itemCount[remoteNodeIP] >= len(WNTopNodesList)-maxOutliers {
				commonCount++
			}
		}

		if commonCount >= minRequiredCommonIPs {
			finalCandidateNodes = append(finalCandidateNodes, someNode)
		}
	}

	if len(finalCandidateNodes) < minRequiredNodes {
		return nil, fmt.Errorf("failed to find at least 3 common SNs in mn-top lists from current top MNs")
	}

	return finalCandidateNodes, nil
}

// GetRequiredDataFromSN retrieves the MN top list from the given SN, dd databse hash and dd-service stats if required
func (m MeshHandler) GetRequiredDataFromSN(ctx context.Context, someNode SuperNodeClient, secInfo *alts.SecInfo) (top []string, hash string, data DDStats, err error) {
	logger := log.WithContext(ctx).WithField("address", someNode.Address())

	if err := someNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("Failed to connect to Supernodes - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
		return top, hash, data, err
	}

	defer func() {
		someNode.RLock()
		defer someNode.RUnlock()

		if someNode.ConnectionInterface != nil {
			if err = someNode.ConnectionInterface.Close(); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"pastelId": someNode.PastelID(),
					"addr":     someNode.String(),
				}).WithError(err).Debugf("connection already closed. SN - %s", someNode.String())
			}
		}
	}()

	if m.requireSNAgreementOnMNTopList {
		resp, err := someNode.GetTopMNs(ctx)
		if err != nil {
			logger.WithError(err).Error("error getting top mns through gRPC call")
			return top, hash, data, fmt.Errorf("error getting top mns through gRPC call: node: %s - err: %w", someNode.Address(), err)
		}

		top = resp.MnTopList
	}

	if m.checkDDDatabaseHashes {
		h, err := someNode.GetDupeDetectionDBHash(ctx)
		if err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("failed to get dd database hash - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			return top, hash, data, fmt.Errorf("failed to get dd database hash - address: %s; pastelID: %s err: %s", someNode.Address(), someNode.PastelID(), err.Error())
		}

		hash = h
		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("node", someNode.Address()).WithField("hash", hash).Info("DD Database hash for node received")

		s, err := someNode.GetDDServerStats(ctx)
		if err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("failed to get dd server stats - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			return top, hash, data, fmt.Errorf("failed to get dd server stats - address: %s; pastelID: %s err: %s", someNode.Address(), someNode.PastelID(), err.Error())
		}

		data = DDStats{
			WaitingInQueue: s.WaitingInQueue,
			Executing:      s.Executing,
			IsReady:        true,
		}

		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("node", someNode.Address()).WithField("stats", data).Info("DD-server stats for node has been received")
	}

	return top, hash, data, nil
}

// CloseSNsConnections closes connections to all nodes
func (m *MeshHandler) CloseSNsConnections(ctx context.Context, nodesDone chan struct{}) error {
	if nodesDone != nil { // Check for nil before closing
		close(nodesDone)
	}

	log.WithContext(ctx).Debug("close connections to ALL supernodes")

	ok := true
	for _, someNode := range m.Nodes {
		if err := m.disconnectFromNode(ctx, someNode, false); err != nil {
			ok = false
		}
	}
	if !ok {
		log.WithContext(ctx).Debug("failed to close connections to ALL supernodes")
	}

	return nil
}

// ConnectionsSupervisor supervises the connection to top rank nodes
// and cancels any ongoing context if the connections are broken
func (m *MeshHandler) ConnectionsSupervisor(ctx context.Context, cancel context.CancelFunc) chan struct{} {
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return m.Nodes.WaitConnClose(ctx, nodesDone)
	})

	return nodesDone
}

// CheckSNReportedState returns SNs state
func (m *MeshHandler) CheckSNReportedState() bool {
	for _, someNode := range m.Nodes {
		if !someNode.IsRemoteState() {
			return false
		}
	}
	return true
}

// DisconnectInactiveNodes disconnects nodes which were not marked as activated.
func (m *MeshHandler) DisconnectInactiveNodes(ctx context.Context) {
	log.WithContext(ctx).Debug("close connections to inactive supernodes")

	for _, someNode := range m.Nodes {
		_ = m.disconnectFromNode(ctx, someNode, true)
	}
}

// disconnectNodes disconnects all passed nodes
func (m *MeshHandler) disconnectNodes(ctx context.Context, nodes SuperNodeList) {
	for _, someNode := range nodes {
		if err := m.disconnectFromNode(ctx, someNode, false); err != nil {
			log.WithContext(ctx).WithError(err).Error("disconnect node failure")
		}
	}
}

func (m *MeshHandler) disconnectFromNode(ctx context.Context, someNode *SuperNodeClient, _ bool) error {
	someNode.RLock()
	defer someNode.RUnlock()

	/*
		if onlyIfInactive && someNode.IsActive() {
			fmt.Println("closing connection now: ", someNode.String())
			return nil
		}*/

	var err error
	if someNode.ConnectionInterface != nil {
		if err = someNode.ConnectionInterface.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": someNode.PastelID(),
				"addr":     someNode.String(),
			}).WithError(err).Debugf("connection already closed. SN - %s", someNode.String())
		}
	} else {
		log.WithContext(ctx).WithFields(log.Fields{
			"pastelId": someNode.PastelID(),
			"addr":     someNode.String(),
		}).Infof("supernode connection is closed already: SN - %s", someNode.String())
	}
	return err
}

// AddNewNode adds new SN to the (internal) collection
func (m *MeshHandler) AddNewNode(address string, pastelID string) {
	m.Nodes.AddNewNode(m.nodeClient, address, pastelID, m.nodeMaker)
}

// Reset resets meshhandler
func (m *MeshHandler) Reset() {
	m.Nodes = SuperNodeList{}
}

func combinations(candidatesNodes SuperNodeList) []SuperNodeList {
	n := len(candidatesNodes)
	var result []SuperNodeList

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				result = append(result, SuperNodeList{candidatesNodes[i], candidatesNodes[j], candidatesNodes[k]})
			}
		}
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1) // Generate a random index from 0 to i
		// Swap elements at index i and j
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func getWNTopNodesHash(WNTopNodesList SuperNodeList) (string, error) {
	var WNTopNodesIPs []string
	for _, node := range WNTopNodesList {
		WNTopNodesIPs = append(WNTopNodesIPs, node.Address())
	}

	WNTopNodesHash, err := sortIPsAndGenerateHash(WNTopNodesIPs)
	if err != nil {
		return "", err
	}

	return WNTopNodesHash, nil
}

func sortIPsAndGenerateHash(ips []string) (string, error) {
	realIPs := make([]net.IP, 0, len(ips))

	for _, ip := range ips {
		host, _, err := net.SplitHostPort(ip)
		if err != nil {
			return "", err
		}

		realIPs = append(realIPs, net.ParseIP(host))
	}

	sort.Slice(realIPs, func(i, j int) bool {
		return bytes.Compare(realIPs[i], realIPs[j]) < 0
	})

	sortedIPs := make([]string, 0, len(realIPs))
	for _, ip := range realIPs {
		sortedIPs = append(sortedIPs, ip.String())
	}

	// Concatenate the sorted IPs into a single string
	concatenated := strings.Join(sortedIPs, "")

	// Create a new hash
	h := sha256.New()

	// Write the concatenated string into the hash
	h.Write([]byte(concatenated))

	// Compute the final hash and print it
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

// NewMeshHandlerSimple returns new MeshHandlerSimple
func NewMeshHandlerSimple(nodeClient node.ClientInterface, nodeMaker node.RealNodeMaker, pastelHandler *mixins.PastelHandler) *MeshHandler {
	return &MeshHandler{
		nodeMaker:     nodeMaker,
		nodeClient:    nodeClient,
		pastelHandler: pastelHandler,
	}
}
