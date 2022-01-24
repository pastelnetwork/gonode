package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

// List represents multiple Node.
type List []*NftRegisterNodeClient

// SendPreBurntFeeTxid send txid of transaction in which 10% of registration fee is preburnt
func (nodes *List) SendPreBurntFeeTxid(ctx context.Context, txid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			ticketTxid, err := node.SendPreBurntFeeTxid(ctx, txid)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("send pre-burnt fee txid failed")
				return err
			}
			if !node.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, node.PastelID())
			}

			node.regNFTTxid = ticketTxid
			return nil
		})
	}
	return group.Wait()
}

// RegArtTicketID return txid of RegArt ticket
func (nodes *List) RegArtTicketID() string {
	for i := range *nodes {
		if (*nodes)[i].IsPrimary() {
			return (*nodes)[i].regNFTTxid
		}
	}
	return string("")
}
