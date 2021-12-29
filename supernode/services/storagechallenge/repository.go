package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
)

type repository interface {
	// ListSymbolFileKeysFromNFTTicket func
	ListSymbolFileKeysFromNFTTicket(ctx context.Context) ([]string, error)
	// GetSymbolFileByKey func
	GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error)
	// StoreSymbolFile func
	StoreSymbolFile(ctx context.Context, data []byte) (key string, err error)
	// RemoveSymbolFileByKey func
	RemoveSymbolFileByKey(ctx context.Context, key string) error
	// GetListOfMasternode func
	GetListOfMasternode(ctx context.Context) ([]string, error)
	// GetNClosestMasternodeIDsToComparisionString func
	GetNClosestMasternodeIDsToComparisionString(ctx context.Context, n int, comparisonString string, listMasternodes []string, ignores ...string) []string
	// GetNClosestMasternodesToAGivenFileUsingKademlia func
	GetNClosestMasternodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores ...string) []string
	// GetNClosestFileHashesToAGivenComparisonString func
	GetNClosestFileHashesToAGivenComparisonString(ctx context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string
}

func newRepository(p2p p2p.Client, pClient pastel.Client) repository {
	return &repo{
		p2p:     p2p,
		pClient: pClient,
	}
}

type repo struct {
	p2p     p2p.Client
	pClient pastel.Client
}

func (r *repo) ListSymbolFileKeysFromNFTTicket(ctx context.Context) ([]string, error) {
	var keys = make([]string, 0)
	regTickets, err := r.pClient.RegTickets(ctx)
	if err != nil {
		return keys, err
	}
	for _, regTicket := range regTickets {
		for _, key := range regTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs {
			keys = append(keys, string(key))
		}
	}

	return keys, nil
}

func (r *repo) GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error) {
	return r.p2p.Retrieve(ctx, key, getFromLocalOnly)
}

func (r *repo) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return r.p2p.Store(ctx, data)
}

func (r *repo) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return r.p2p.Delete(ctx, key)
}

func (r *repo) GetListOfMasternode(ctx context.Context) ([]string, error) {
	var ret = make([]string, 0)
	listMN, err := r.pClient.MasterNodesList(ctx)
	if err != nil {
		return ret, err
	}

	for _, node := range listMN {
		ret = append(ret, node.ExtKey)
	}

	return ret, nil
}

func (r *repo) GetNClosestMasternodeIDsToComparisionString(_ context.Context, n int, comparisonString string, listMasternodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listMasternodes, ignores...)
}

func (r *repo) GetNClosestMasternodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores ...string) []string {
	return r.p2p.NClosestNodes(ctx, n, comparisonString, ignores...)
}

func (r *repo) GetNClosestFileHashesToAGivenComparisonString(_ context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listFileHashes, ignores...)
}
