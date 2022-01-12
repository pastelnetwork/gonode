package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister/node"
)

type SuperNodesTasks struct {
}

func (task *SuperNodesTasks) connectToTopRankNodes(ctx context.Context) error {

	/* Select 3 SNs */
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("call masternode top: %w", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorFindTopNodes)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// Get current block height & hash
	if err := task.getBlock(ctx); err != nil {
		return errors.Errorf("get block: %v", err)
	}

	// Connect to top nodes to find 3SN and validate their infor
	err = task.validateMNsInfo(ctx, topNodes)
	if err != nil {
		task.UpdateStatus(StatusErrorFindResponsdingSNs)
		return errors.Errorf("validate MNs info: %v", err)
	}

	/* Establish mesh with 3 SNs */
	var nodes node.List
	var errs error

	for primaryRank := range task.nodes {
		nodes, err = task.meshNodes(ctx, task.nodes, primaryRank)
		if err != nil {
			// close connected connections
			task.nodes.DisconnectAll()

			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warn("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
		// close connected connections
		topNodes.DisconnectAll()
		return errors.Errorf("Could not create a mesh of %d nodes: %w", task.config.NumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	nodes.Activate()

	// Cancel context when any connection is broken.
	task.UpdateStatus(StatusConnected)

	// Send all meshed supernode info to nodes - that will be used to node send info to other nodes
	meshedSNInfo := []types.MeshedSuperNode{}
	for _, node := range nodes {
		meshedSNInfo = append(meshedSNInfo, types.MeshedSuperNode{
			NodeID: node.PastelID(),
			SessID: node.SessID(),
		})
	}

	for _, node := range nodes {
		err = node.MeshNodes(ctx, meshedSNInfo)
		if err != nil {
			nodes.DisconnectAll()
			return errors.Errorf("could not send info of meshed nodes: %w", err)
		}
	}
	return nil
}

func (task *SuperNodesTasks) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = task.pastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}
