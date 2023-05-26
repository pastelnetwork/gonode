package common

import (
	"context"
	"fmt"
	"math/rand"
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

	Nodes                 SuperNodeList
	UseMaxNodes           bool
	checkDDDatabaseHashes bool
	HashCheckMaxRetries   int
}

// MeshHandlerOpts set of options to pass to NewMeshHandler
type MeshHandlerOpts struct {
	Task          *WalletNodeTask
	NodeClient    node.ClientInterface
	NodeMaker     node.RealNodeMaker
	PastelHandler *mixins.PastelHandler
	Configs       *MeshHandlerConfig
}

// MeshHandlerConfig config subset used by MeshHandler
type MeshHandlerConfig struct {
	PastelID               string
	Passphrase             string
	MinSNs                 int
	ConnectToNodeTimeout   time.Duration
	AcceptNodesTimeout     time.Duration
	ConnectToNextNodeDelay time.Duration
	UseMaxNodes            bool
	CheckDDDatabaseHashes  bool
	HashCheckMaxRetries    int
}

// NewMeshHandler returns new NewMeshHandler
func NewMeshHandler(opts MeshHandlerOpts) *MeshHandler {
	return &MeshHandler{
		task:                   opts.Task,
		nodeMaker:              opts.NodeMaker,
		pastelHandler:          opts.PastelHandler,
		nodeClient:             opts.NodeClient,
		callersPastelID:        opts.Configs.PastelID,
		passphrase:             opts.Configs.Passphrase,
		minNumberSuperNodes:    opts.Configs.MinSNs,
		connectToNodeTimeout:   opts.Configs.ConnectToNodeTimeout,
		acceptNodesTimeout:     opts.Configs.AcceptNodesTimeout,
		connectToNextNodeDelay: opts.Configs.ConnectToNextNodeDelay,
		UseMaxNodes:            opts.Configs.UseMaxNodes,
		checkDDDatabaseHashes:  opts.Configs.CheckDDDatabaseHashes,
		HashCheckMaxRetries:    opts.Configs.HashCheckMaxRetries,
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
		return 0, "", errors.Errorf("failed validating of %d SNs", m.minNumberSuperNodes)
	}

	meshedNodes, err := m.setMesh(ctx, connectedNodes, m.minNumberSuperNodes)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("failed setting up Mesh of %d SNs from %d SNs", m.minNumberSuperNodes, len(connectedNodes))
		return 0, "", errors.Errorf("failed setting up Mesh of %d SNs from %d SNs", m.minNumberSuperNodes, len(connectedNodes))
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
	candidatesNodes, err := m.getTopNodes(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not get top masternodes")
		return nil, errors.Errorf("call masternode top: %w", err)
	}

	if len(candidatesNodes) < n {
		m.task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		log.WithContext(ctx).WithError(err).Errorf("unable to find enough Supernodes: %d", n)
		return nil, fmt.Errorf("unable to find enough Supernodes: %d", n)
	}
	log.WithContext(ctx).Infof("Found %d Supernodes", len(candidatesNodes))

	// Connect to top nodes to find 3SN and validate their info
	foundNodes, err := m.connectToAndValidateSuperNodes(ctx, candidatesNodes, n, m.HashCheckMaxRetries, skipNodes)
	if err != nil {
		m.task.UpdateStatus(StatusErrorFindRespondingSNs)
		log.WithContext(ctx).WithError(err).Errorf("unable to validate MN")
		return nil, errors.Errorf("validate MNs info: %v", err)
	}
	log.WithContext(ctx).Infof("Found %d valid Supernodes", len(foundNodes))
	return foundNodes, nil
}

func (m *MeshHandler) setMesh(ctx context.Context, candidatesNodes SuperNodeList, n int) (SuperNodeList, error) {
	log.WithContext(ctx).Debug("setMesh Starting...")

	var errs error
	var err error

	var meshedNodes SuperNodeList
	log.WithContext(ctx).Info("connected nodes: ", candidatesNodes)
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
		log.WithContext(ctx).WithError(err).Warnf("Could not create a mesh of %d nodes. Disconnecting any still connected", n)
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

	if m.checkDDDatabaseHashes {
		return m.connectToAndValidateSuperNodesWithHashCheck(ctx, candidatesNodes, n, maxRetries, secInfo, skipNodes)
	}

	return m.connectToAndValidateSuperNodesWithoutHashCheck(ctx, candidatesNodes, n, secInfo, skipNodes)
}

func (m *MeshHandler) connectToAndValidateSuperNodesWithHashCheck(ctx context.Context, candidatesNodes SuperNodeList, n int, maxRetries int, secInfo *alts.SecInfo, skipNodes []string) (SuperNodeList, error) {
	combinations := combinations(candidatesNodes)

	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		for _, set := range combinations {
			nodes := m.connectToNodes(ctx, set, n, secInfo, skipNodes)

			if len(nodes) != n {
				continue
			}

			if err := m.matchDatabaseHash(ctx, nodes); err != nil {
				continue
			}

			return nodes, nil
		}

		log.WithContext(ctx).Infof("Could not match database hash, retrying in 30s... (attempt %d/%d)", retryCount+1, maxRetries)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while waiting to retry: %w", ctx.Err())
		case <-time.After(30 * time.Second):
		}
	}

	return nil, errors.Errorf("nodes not found, could not match database hash after %d attempts", maxRetries)
}

func (m *MeshHandler) connectToAndValidateSuperNodesWithoutHashCheck(ctx context.Context, candidatesNodes SuperNodeList, n int, secInfo *alts.SecInfo, skipNodes []string) (SuperNodeList, error) {
	nodes := m.connectToNodes(ctx, candidatesNodes, n, secInfo, skipNodes)

	if len(nodes) < n {
		log.WithContext(ctx).Errorf("Failed to validate enough Supernodes - only found %d good", n)
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
			log.WithContext(ctx).WithError(err).Errorf("Failed to connect to Supernodes - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			continue
		}

		nodes = append(nodes, someNode)
		if len(nodes) == n && !m.UseMaxNodes {
			break
		}
	}

	return nodes
}

// matchDatabaseHash matches database hashes
func (m *MeshHandler) matchDatabaseHash(ctx context.Context, nodesList SuperNodeList) error {

	log.WithContext(ctx).Infof("Matching database hash for nodes %d", len(nodesList))
	hashes := make(map[string]string)
	matcher := ""
	// Connect to top nodes to find 3SN and validate their info
	for i, someNode := range nodesList {
		hash, err := someNode.GetDupeDetectionDBHash(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("failed to get dd database hash - address: %s; pastelID: %s ", someNode.String(), someNode.PastelID())
			return fmt.Errorf("failed to get dd database hash - address: %s; pastelID: %s err: %s", someNode.Address(), someNode.PastelID(), err.Error())
		}

		log.WithContext(ctx).WithField("node", someNode.Address()).WithField("hash", hash).Info("DD Database hash for node received")
		hashes[someNode.Address()] = hash

		if i == 0 {
			matcher = hash
		}
	}

	for key, val := range hashes {
		if val != matcher {
			log.WithContext(ctx).Errorf("database hash mismatch for node %s", key)
			return fmt.Errorf("database hash mismatch for node %s", key)
		}
	}

	log.WithContext(ctx).Infof("DD Database hashes matched")

	return nil
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

// CloseSNsConnections closes connections to all nodes
func (m *MeshHandler) CloseSNsConnections(ctx context.Context, nodesDone chan struct{}) error {
	close(nodesDone)

	log.WithContext(ctx).Debug("close connections to ALL supernodes")

	ok := true
	for _, someNode := range m.Nodes {
		if err := m.disconnectFromNode(ctx, someNode, false); err != nil {
			ok = false
		}
	}
	if !ok {
		return errors.New("errors while disconnecting from some SNs")
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

func (m *MeshHandler) disconnectFromNode(ctx context.Context, someNode *SuperNodeClient, onlyIfInactive bool) error {
	someNode.RLock()
	defer someNode.RUnlock()

	if onlyIfInactive && someNode.IsActive() {
		return nil
	}

	var err error
	if someNode.ConnectionInterface != nil {
		if err = someNode.ConnectionInterface.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": someNode.PastelID(),
				"addr":     someNode.String(),
			}).WithError(err).Errorf("close supernode connection failed: SN - %s", someNode.String())
		}
	}
	return err
}

// AddNewNode adds new SN to the (internal) collection
func (m *MeshHandler) AddNewNode(address string, pastelID string) {
	m.Nodes.AddNewNode(m.nodeClient, address, pastelID, m.nodeMaker)
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

// NewMeshHandlerSimple returns new MeshHandlerSimple
func NewMeshHandlerSimple(nodeClient node.ClientInterface, nodeMaker node.RealNodeMaker, pastelHandler *mixins.PastelHandler) *MeshHandler {
	return &MeshHandler{
		nodeMaker:     nodeMaker,
		nodeClient:    nodeClient,
		pastelHandler: pastelHandler,
	}
}
