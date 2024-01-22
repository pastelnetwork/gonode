package common

import (
	"context"
	"encoding/json"
	"fmt"
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
	checkMinBalance               bool
	minBalance                    int64
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

// DDStats contains dd-service stats
type DDStats struct {
	Version        string
	WaitingInQueue int32
	Executing      int32
	MaxConcurrency int32
	IsReady        bool
}

// SNData contains data required for supernode
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
	CheckMinBalance               bool
	MinBalance                    int64
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
		checkMinBalance:               opts.Configs.CheckMinBalance,
		minBalance:                    opts.Configs.MinBalance,
	}
}

// SetupMeshOfNSupernodesNodes sets Mesh
func (m *MeshHandler) SetupMeshOfNSupernodesNodes(ctx context.Context, sortKey []byte) (int, string, error) {
	log.WithContext(ctx).Info("SetupMeshOfNSupernodesNodes Starting...")

	// Get current block height & hash
	blockNum, blockHash, err := m.pastelHandler.GetBlock(ctx)
	if err != nil {
		return -1, "", errors.Errorf("get block: %v", err)
	}
	log.WithContext(ctx).Infof("Current block is %d", blockNum)

	connectedNodes, err := m.findNValidTopSuperNodes(ctx, m.minNumberSuperNodes, []string{}, sortKey, blockNum)
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

	connectedNodes, err := m.findNValidTopSuperNodes(ctx, n, skipNodes, []byte{}, 0)
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
func (m *MeshHandler) findNValidTopSuperNodes(ctx context.Context, n int, skipNodes []string, sortKey []byte, block int) (SuperNodeList, error) {

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

	candidateNodes, err := m.GetCandidateNodes(ctx, WNTopNodes, n, block)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting candidate nodes")
	}

	if len(candidateNodes) < n {
		if err == nil {
			err = errors.New("not enough candidate nodes found with required parameters")
		}
		log.WithContext(ctx).WithField("count", len(candidateNodes)).WithError(err)
		return nil, fmt.Errorf("unable to find enough Supernodes: required: %d - found: %d - err: %w", n, len(candidateNodes), err)
	}

	// Connect to top nodes to find 3SN and validate their info
	foundNodes, err := m.connectToAndValidateSuperNodes(ctx, candidateNodes, n, m.HashCheckMaxRetries, skipNodes)
	if err != nil {
		m.task.UpdateStatus(StatusErrorFindRespondingSNs)
		log.WithContext(ctx).WithError(err).Errorf("unable to validate MN")
		return nil, errors.Errorf("validate MNs info: %v", err)
	}
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("Found %d valid Supernodes", len(foundNodes))

	if len(sortKey) > 0 {
		log.WithContext(ctx).WithField("before sort", foundNodes).Info("nodes before sorting")
		foundNodes.Sort(sortKey)
		log.WithContext(ctx).WithField("after sort", foundNodes).Info("nodes after sorting")
	} else {
		log.WithContext(ctx).Info("no sort key provided")
	}

	return foundNodes, nil
}

// GetCandidateNodes returns list of candidate nodes
func (m *MeshHandler) GetCandidateNodes(ctx context.Context, candidatesNodes SuperNodeList, n int, block int) (SuperNodeList, error) {
	if !m.checkDDDatabaseHashes && !m.requireSNAgreementOnMNTopList {
		return candidatesNodes, nil
	}

	log.WithContext(ctx).WithField("given nodes", len(candidatesNodes)).Info("GetCandidateNodes Starting...")

	type result struct {
		node    *SuperNodeClient
		balance int64
		top     []string
		hash    string
		data    DDStats
		err     error
	}

	// Creating a buffered channel to collect results
	resultsChan := make(chan *result, len(candidatesNodes))

	WNTopNodesList := SuperNodeList{}
	topMap := make(map[string][]string)
	hashes := make(map[string]string)
	dataMap := make(map[string]DDStats)
	balances := make(map[string]int64)

	var wg sync.WaitGroup
	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  pastel.SignAlgorithmED448,
	}

	// Start a goroutine for each node
	for _, someNode := range candidatesNodes {
		wg.Add(1)
		go func(node *SuperNodeClient) {
			defer wg.Done()
			balance, top, hash, data, err := m.GetRequiredDataFromSN(ctx, *node, secInfo, block)
			resultsChan <- &result{node, balance, top, hash, data, err}
		}(someNode)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process results as they come in
	for res := range resultsChan {
		if res.err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(res.err).Errorf("failed to get required data from supernode - address: %s; pastelID: %s ", res.node.String(), res.node.PastelID())
			continue
		}

		WNTopNodesList.Add(res.node)
		topMap[res.node.Address()] = res.top
		hashes[res.node.Address()] = res.hash
		dataMap[res.node.Address()] = res.data
		balances[res.node.Address()] = res.balance
	}

	if len(WNTopNodesList) < n {
		return WNTopNodesList, errors.New("failed to get required data from enough candidate nodes")
	}

	var err error
	if m.requireSNAgreementOnMNTopList {

		// Filter out nodes with different top MNs list
		given := WNTopNodesList
		WNTopNodesList, err = m.filterByTopNodesList(ctx, WNTopNodesList, topMap)
		if err != nil {
			return nil, fmt.Errorf("filter by top nodes list: %w", err)
		}

		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("given nodes count", len(given)).
			WithField("matching top 10 count", len(WNTopNodesList)).Info("Matched top nodes list for nodes")

		if len(WNTopNodesList) < len(given) {
			topJSON, _ := json.Marshal(topMap)

			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("given nodes", given.String()).
				WithField("matched list", WNTopNodesList).WithField("block", block).WithField("top 10 lists", string(topJSON)).Warn("some nodes left out due to mismatching top 10 list")
		}
	}

	if len(WNTopNodesList) < n {
		return WNTopNodesList, errors.New("failed to get enough nodes with matching top 10 list")
	}

	if m.checkMinBalance {
		WNTopNodesList = m.filterByMinBalanceReq(ctx, WNTopNodesList, balances, m.minBalance)
		if len(WNTopNodesList) < n {
			return WNTopNodesList, errors.New("failed to get enough nodes with minimum balance")
		}
	}
	log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("nodes count", len(candidatesNodes)).
		Info("given nodes count after checking min req balance")

	if m.checkDDDatabaseHashes {
		WNTopNodesList, err = m.filterByDDHash(ctx, WNTopNodesList, hashes)
		if err != nil {
			return nil, fmt.Errorf("filter by top nodes list: %w", err)
		}

		if len(WNTopNodesList) < n {
			return WNTopNodesList, errors.New("failed to get enough nodes with matching fingerprints database hash")
		}

		WNTopNodesList, err = m.filterDDServerRequestsInWaitingQueue(ctx, WNTopNodesList, dataMap)
		if err != nil {
			return nil, fmt.Errorf("filter by dd-service stats: %w", err)
		}

		if len(WNTopNodesList) < n {
			return WNTopNodesList, errors.New("failed to get required number of nodes with available dd-service")
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
			log.WithContext(ctx).WithError(err).Warn("Connecting to primary and secondary failed. Disconnecting any still connected")
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

func (m *MeshHandler) connectToAndValidateSuperNodes(ctx context.Context, candidatesNodes SuperNodeList, n int, _ int, skipNodes []string) (SuperNodeList, error) {
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

func (m *MeshHandler) connectToNodes(ctx context.Context, nodesToConnect SuperNodeList, _ int, secInfo *alts.SecInfo, skipNodes []string) SuperNodeList {
	var nodes SuperNodeList
	var mutex sync.Mutex
	var wg sync.WaitGroup

	shouldSkip := func(node *SuperNodeClient) bool {
		for _, skipNode := range skipNodes {
			if skipNode == node.PastelID() {
				return true
			}
		}
		return false
	}

	for _, someNode := range nodesToConnect {
		if shouldSkip(someNode) {
			continue
		}

		wg.Add(1)
		go func(node *SuperNodeClient) {
			defer wg.Done()

			if err := node.Connect(ctx, m.connectToNodeTimeout, secInfo, m.logRequestID); err != nil {
				log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("Failed to connect to Supernodes - address: %s; pastelID: %s ", node.String(), node.PastelID())
			} else {
				mutex.Lock()
				nodes = append(nodes, node)
				mutex.Unlock()
			}
		}(someNode)
	}

	wg.Wait()

	return nodes
}

func (m *MeshHandler) filterByMinBalanceReq(ctx context.Context, WNTopNodesList SuperNodeList, balances map[string]int64, reqBalance int64) (retList SuperNodeList) {
	retList = SuperNodeList{}
	for _, someNode := range WNTopNodesList {
		balance := balances[someNode.Address()]
		if balance >= reqBalance {
			retList = append(retList, someNode)
		} else {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("balance", balance).WithField("required balance", reqBalance).
				Warnf("not enough balance %s", someNode.Address())
		}
	}

	return retList
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

	log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("given nodes count", len(nodesList)).
		WithField("matching hash nodes count", len(retList)).WithField("matching hash", matchingKey).Infof("DD Database hashes matched")

	return retList, nil
}

func (m *MeshHandler) filterDDServerRequestsInWaitingQueue(ctx context.Context, nodesList SuperNodeList, dataMap map[string]DDStats) (retList SuperNodeList, err error) {
	log.WithContext(ctx).WithField("req-id", m.logRequestID).Infof("checking dd-server no of requests waiting in a list of %d nodes ", len(nodesList))
	retList = SuperNodeList{}
	// Connect to top nodes to find 3SN and validate their info
	for _, someNode := range nodesList {
		waiting, executing, maxConcurrent := dataMap[someNode.Address()].WaitingInQueue, dataMap[someNode.Address()].Executing, dataMap[someNode.Address()].MaxConcurrency
		if maxConcurrent-executing-waiting > 0 {
			retList = append(retList, someNode)
		} else {

			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("executing", executing).WithField("max-concurrent", maxConcurrent).
				WithField("waiting tasks", waiting).Warnf("dd-service not ready %s", someNode.Address())
		}
	}

	log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("give nodes count", len(nodesList)).
		WithField("dd-service ready nodes", len(retList)).Infof("dd-service ready nodes matched")

	versionMap := make(map[string]SuperNodeList)

	for _, someNode := range retList {
		version := dataMap[someNode.Address()].Version
		versionMap[version] = append(versionMap[version], someNode)
	}

	// Find the key (version) with the highest count of nodes
	var maxVersion string
	var maxCount int

	for version, nodeList := range versionMap {
		if len(nodeList) > maxCount {
			maxCount = len(nodeList)
			maxVersion = version
		}
	}

	log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("version", maxVersion).WithField("count", maxCount).
		Infof("dd-service version matched")

	return versionMap[maxVersion], nil
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
	if err := primary.Connect(ctx, m.connectToNodeTimeout, secInfo, m.logRequestID); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to connect to primary node %q", primary)
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := primary.Session(ctx, true); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to open session with primary node %q", primary)
		return nil, fmt.Errorf("session: %w", err)
	}
	log.WithContext(ctx).WithField("task", m.logRequestID).WithField("remote-sess-id", primary.SessID()).Infof("Connected to primary node %q", primary.Address())

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

					if err := sameNode.Connect(ctx, m.connectToNodeTimeout, secInfo, m.logRequestID); err != nil {
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
					log.WithContext(ctx).WithField("task", m.logRequestID).WithField("sameNode", sameNode.SessID()).Infof("Connected to secondry node %q", sameNode.Address())
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
		log.WithContext(ctx).Infof("Primary accepted %q secondary node", pastelID)

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
	for _, someNode := range candidateNodes {
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
func (m MeshHandler) GetRequiredDataFromSN(ctx context.Context, someNode SuperNodeClient, secInfo *alts.SecInfo, block int) (balance int64, top []string, hash string, data DDStats, err error) {
	logger := log.WithContext(ctx).WithField("address", someNode.Address())

	if err := someNode.Connect(ctx, m.connectToNodeTimeout, secInfo, m.logRequestID); err != nil {
		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("Failed to connect to Supernodes - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
		return balance, top, hash, data, err
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
		resp, err := someNode.GetTopMNs(ctx, block)
		if err != nil {
			logger.WithError(err).Error("error getting top mns through gRPC call")
			return balance, top, hash, data, fmt.Errorf("error getting top mns through gRPC call: node: %s - err: %w", someNode.Address(), err)
		}

		top = resp.MnTopList
		balance = resp.CurrBalance
	}

	if m.checkDDDatabaseHashes {
		h, err := someNode.GetDupeDetectionDBHash(ctx)
		if err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("failed to get dd database hash - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			return balance, top, hash, data, fmt.Errorf("failed to get dd database hash - address: %s; pastelID: %s err: %s", someNode.Address(), someNode.PastelID(), err.Error())
		}

		hash = h

		s, err := someNode.GetDDServerStats(ctx)
		if err != nil {
			log.WithContext(ctx).WithField("req-id", m.logRequestID).WithError(err).Errorf("failed to get dd server stats - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			return balance, top, hash, data, fmt.Errorf("failed to get dd server stats - address: %s; pastelID: %s err: %s", someNode.Address(), someNode.PastelID(), err.Error())
		}

		data = DDStats{
			Version:        s.Version,
			WaitingInQueue: s.WaitingInQueue,
			Executing:      s.Executing,
			MaxConcurrency: s.MaxConcurrent,
			IsReady:        true,
		}

		log.WithContext(ctx).WithField("req-id", m.logRequestID).WithField("hash", hash).WithField("node", someNode.Address()).WithField("stats", data).Info("DD-service stats & hash for node has been received")
	}

	return balance, top, hash, data, nil
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

// NewMeshHandlerSimple returns new MeshHandlerSimple
func NewMeshHandlerSimple(nodeClient node.ClientInterface, nodeMaker node.RealNodeMaker, pastelHandler *mixins.PastelHandler) *MeshHandler {
	return &MeshHandler{
		nodeMaker:     nodeMaker,
		nodeClient:    nodeClient,
		pastelHandler: pastelHandler,
	}
}
