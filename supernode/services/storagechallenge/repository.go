package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
)

type repository interface {
	// ListKeys func
	ListKeys(ctx context.Context) []string
	// GetSymbolFileByKey func
	GetSymbolFileByKey(ctx context.Context, key string) ([]byte, error)
	// StoreSymbolFile func
	StoreSymbolFile(ctx context.Context, data []byte) (key string, err error)
	// RemoveSymbolFileByKey func
	RemoveSymbolFileByKey(ctx context.Context, key string) error
	// GetNClosestXORDistanceMasternodesToComparisionString func
	GetNClosestXORDistanceMasternodesToComparisionString(ctx context.Context, n int, comparisonString string) []*kademlia.Node
	// GetNClosestXORDistanceFileHashesToComparisonString func
	GetNClosestXORDistanceFileHashesToComparisonString(ctx context.Context, n int, comparisonString string, symbolFileKeys []string) []string
}

func newRepository(client p2p.Client) repository {
	return &repo{
		p2p: client,
	}
}

type repo struct {
	p2p p2p.Client
}

func (r *repo) ListKeys(ctx context.Context) []string {
	var offset, limit = 0, 30
	var ret = make([]string, 0)
	for {
		keys := r.p2p.Keys(ctx, offset, limit)
		ret = append(ret, keys...)
		offset += limit
		if len(keys) < limit {
			break
		}
	}
	return ret
}

func (r *repo) GetSymbolFileByKey(ctx context.Context, key string) ([]byte, error) {
	return r.p2p.Retrieve(ctx, key, true)
}

func (r *repo) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return r.p2p.Store(ctx, data)
}

func (r *repo) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return r.p2p.Delete(ctx, key)
}

func (r *repo) GetNClosestXORDistanceMasternodesToComparisionString(ctx context.Context, n int, comparisonString string) []*kademlia.Node {
	return r.p2p.NClosestNodes(ctx, n, comparisonString)
}

func (r *repo) GetNClosestXORDistanceFileHashesToComparisonString(_ context.Context, n int, comparisonString string, symbolFileKeys []string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, symbolFileKeys)
}
