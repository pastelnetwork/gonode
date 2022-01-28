package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"sync"
	"time"
)

type MeshHandler struct {
	task          *WalletNodeTask
	pastelHandler *mixins.PastelHandler

	nodeMaker  node.NodeMaker
	nodeClient node.ClientInterface

	// TODO: make config interface and use it here instead of individual items
	minNumberSuperNodes    int
	connectToNodeTimeout   time.Duration
	acceptNodesTimeout     time.Duration
	connectToNextNodeDelay time.Duration

	callersPastelID string
	passphrase      string

	Nodes SuperNodeList
}

func NewMeshHandler(task *WalletNodeTask,
	nodeClient node.ClientInterface, nodeMaker node.NodeMaker,
	pastelHandler *mixins.PastelHandler,
	pastelID string, passphrase string,
	minSNs int, connectToNodeTimeout time.Duration, acceptNodesTimeout time.Duration, connectToNextNodeDelay time.Duration,
) *MeshHandler {
	return &MeshHandler{
		task:                   task,
		nodeMaker:              nodeMaker,
		pastelHandler:          pastelHandler,
		nodeClient:             nodeClient,
		callersPastelID:        pastelID,
		passphrase:             passphrase,
		minNumberSuperNodes:    minSNs,
		connectToNodeTimeout:   connectToNodeTimeout,
		acceptNodesTimeout:     acceptNodesTimeout,
		connectToNextNodeDelay: connectToNextNodeDelay,
	}
}

// SetupMeshOfNSupernodesNodes
func (m *MeshHandler) SetupMeshOfNSupernodesNodes(ctx context.Context) (int, string, error) {

	// Get current block height & hash
	blockNum, blockHash, err := m.pastelHandler.GetBlock(ctx)
	if err != nil {
		return -1, "", errors.Errorf("get block: %v", err)
	}

	candidatesNodes, err := m.findNValidTopSuperNodes(ctx, m.minNumberSuperNodes)
	if err != nil {
		return 0, "", err
	}

	// Close all connected connections - setMesh will Connect again
	m.disconnectNodes(ctx, candidatesNodes)

	meshedNodes, err := m.setMesh(ctx, candidatesNodes, m.minNumberSuperNodes)
	if err != nil {
		return 0, "", err
	}
	m.Nodes = meshedNodes

	return blockNum, blockHash, nil
}

func (m MeshHandler) ConnectToNSuperNodes(ctx context.Context, n int) error {

	connectedNodes, err := m.findNValidTopSuperNodes(ctx, n)
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
func (m *MeshHandler) findNValidTopSuperNodes(ctx context.Context, n int) (SuperNodeList, error) {

	// Retrieve supernodes with the highest ranks.
	candidatesNodes, err := m.getTopNodes(ctx)
	if err != nil {
		return nil, errors.Errorf("call masternode top: %w", err)
	}

	if len(candidatesNodes) < n {
		m.task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		return nil, errors.New("unable to find enough Supernodes")
	}

	// Connect to top nodes to find 3SN and validate their info
	candidatesNodes, err = m.connectToAndValidateSuperNodes(ctx, candidatesNodes, n)
	if err != nil {
		m.task.UpdateStatus(StatusErrorFindRespondingSNs)
		return nil, errors.Errorf("validate MNs info: %v", err)
	}
	return candidatesNodes, nil
}

func (m *MeshHandler) setMesh(ctx context.Context, candidatesNodes SuperNodeList, n int) (SuperNodeList, error) {
	var errs error
	var err error

	var meshedNodes SuperNodeList
	for primaryRank := range candidatesNodes {
		meshedNodes, err = m.connectToPrimarySecondary(ctx, candidatesNodes, primaryRank)
		if err != nil {
			// close connected connections
			m.disconnectNodes(ctx, candidatesNodes)

			if errors.IsContextCanceled(err) {
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

// validateMNsInfo - validate MNs info, until found at least N valid MNs
func (m *MeshHandler) connectToAndValidateSuperNodes(ctx context.Context, candidatesNodes SuperNodeList, n int) (SuperNodeList, error) {
	count := 0

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  "ed448",
	}

	var nodes SuperNodeList
	for _, someNode := range candidatesNodes {
		if err := someNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
			continue
		}

		count++
		nodes = append(nodes, someNode)
		if count == n {
			break
		}
	}

	if count < n {
		return nil, errors.Errorf("validate %d Supernodes from pastel network", n)
	}

	return nodes, nil
}

// meshNodes establishes communication between supernodes.
func (m *MeshHandler) connectToPrimarySecondary(ctx context.Context, candidatesNodes SuperNodeList, primaryIndex int) (SuperNodeList, error) {

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  "ed448",
	}

	primary := candidatesNodes[primaryIndex]
	log.WithContext(ctx).Debugf("Trying to connect to primary node %q", primary)
	if err := primary.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
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
						return
					}
					if err := sameNode.Session(ctx, false); err != nil {
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
			return nil, errors.New("not found accepted node")
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

// ValidBurnTxID returns whether the burn txid is valid at ALL SNs
func (m *MeshHandler) CheckSNReportedState() bool {
	for _, someNode := range m.Nodes {
		if !someNode.IsRemoteState() {
			return false
		}
	}
	return true
}

// DisconnectInactive disconnects nodes which were not marked as activated.
func (m *MeshHandler) DisconnectInactiveNodes(ctx context.Context) {
	log.WithContext(ctx).Debug("close connections to inactive supernodes")

	for _, someNode := range m.Nodes {
		_ = m.disconnectFromNode(ctx, someNode, false)
	}
}

// disconnectNodes disconnects all passed nodes
func (m *MeshHandler) disconnectNodes(ctx context.Context, nodes SuperNodeList) {
	for _, someNode := range nodes {
		_ = m.disconnectFromNode(ctx, someNode, false)
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

func (m *MeshHandler) AddNewNode(address string, pastelID string) {
	m.Nodes.AddNewNode(m.nodeClient, address, pastelID, m.nodeMaker)
}

func NewMeshHandlerSimple(nodeClient node.ClientInterface, nodeMaker node.NodeMaker) *MeshHandler {
	return &MeshHandler{
		nodeMaker:  nodeMaker,
		nodeClient: nodeClient,
	}
}
