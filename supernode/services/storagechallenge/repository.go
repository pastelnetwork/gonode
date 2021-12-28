package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
)

type repository interface {
	// ListKeys func
	ListKeys(ctx context.Context) ([]string, error)
	// GetSymbolFileByKey func
	GetSymbolFileByKey(ctx context.Context, key string) ([]byte, error)
	// StoreSymbolFile func
	StoreSymbolFile(ctx context.Context, data []byte) (key string, err error)
	// RemoveSymbolFileByKey func
	RemoveSymbolFileByKey(ctx context.Context, key string) error
	// GetNClosestXORDistanceMasternodesToComparisionString func
	GetNClosestXORDistanceMasternodesToComparisionString(ctx context.Context, n int, comparisonString string, ignores ...string) []string
	// GetNClosestXORDistanceFileHashesToComparisonString func
	GetNClosestXORDistanceFileHashesToComparisonString(ctx context.Context, n int, comparisonString string, symbolFileKeys []string, ignores ...string) []string
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

func (r *repo) ListKeys(ctx context.Context) ([]string, error) {
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
	// mock integration test, this should be removed after testing
	// resp, err := httpClient.Get("http://helper:8088/store/keys")
	// if err != nil {
	// 	return keys, err
	// }
	// defer resp.Body.Close()
	// return keys, json.NewDecoder(resp.Body).Decode(&keys)

	return keys, nil
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

func (r *repo) GetNClosestXORDistanceMasternodesToComparisionString(ctx context.Context, n int, comparisonString string, ignores ...string) []string {
	return r.p2p.NClosestNodes(ctx, n, comparisonString, ignores...)
}

func (r *repo) GetNClosestXORDistanceFileHashesToComparisonString(_ context.Context, n int, comparisonString string, symbolFileKeys []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, symbolFileKeys, ignores...)
}
