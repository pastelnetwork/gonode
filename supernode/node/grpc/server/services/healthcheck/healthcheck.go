package healthcheck

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/p2p"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"google.golang.org/grpc"
)

const (
	defaultP2PDataMaxLength     = 32
	defaultP2PExpiresDuraction  = 15 * time.Minute
	defaultRqliteQueryMaxLength = 1 * 1024 // 1k
)

// StatsMngr is interface of StatsManger, return stats of system
type StatsMngr interface {
	// Stats returns stats of system
	Stats(ctx context.Context) (map[string]interface{}, error)
}

// P2PTracker is interface of lifetime tracker that will periodically check to remove expired temporary p2p's (key, value)
type P2PTracker interface {
	// Track monitor to delete temporary key
	Track(key string, expires time.Time)
}

// HealthCheck represents grpc service for supernode healthcheck
type HealthCheck struct {
	pb.UnimplementedHealthCheckServer
	StatsMngr
	p2pService p2p.P2P
	metaDb     metadb.MetaDB
	p2pTracker P2PTracker
}

// Status will send a message to and get back a reply from supernode
func (service *HealthCheck) Status(ctx context.Context, _ *pb.StatusRequest) (*pb.StatusReply, error) {
	stats, err := service.Stats(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to Stats(): %w", err)
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		return nil, errors.Errorf("failed to Marshal(): %w", err)
	}

	// echos received message
	return &pb.StatusReply{StatusInJson: string(jsonData)}, nil
}

// P2PStore store a data, and return given key
func (service *HealthCheck) P2PStore(ctx context.Context, in *pb.P2PStoreRequest) (*pb.P2PStoreReply, error) {
	value := in.GetValue()
	if len(value) > defaultP2PDataMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}
	key, err := service.p2pService.Store(ctx, value)

	if err != nil {
		return nil, err
	}

	// track to delete key soon
	service.p2pTracker.Track(key, time.Now().Add(defaultP2PExpiresDuraction))

	return &pb.P2PStoreReply{Key: key}, nil
}

// P2PRetrieve get data of given key
func (service *HealthCheck) P2PRetrieve(ctx context.Context, in *pb.P2PRetrieveRequest) (*pb.P2PRetrieveReply, error) {
	key := in.GetKey()
	value, err := service.p2pService.Retrieve(ctx, key)

	if err != nil {
		return &pb.P2PRetrieveReply{Value: nil}, err
	}
	return &pb.P2PRetrieveReply{Value: value}, nil
}

// QueryRqlite do a query
func (service *HealthCheck) QueryRqlite(ctx context.Context, in *pb.QueryRqliteRequest) (*pb.QueryRqliteReply, error) {
	if len(in.GetQuery()) > defaultRqliteQueryMaxLength {
		return nil, errors.New("length too big (>1kbytes")
	}

	// do query
	queryResult, err := service.metaDb.Query(ctx, in.GetQuery(), "nono")
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return &pb.QueryRqliteReply{Result: "{}"}, nil
	}

	// Get all rows of output, and fill into map to convert to json string
	row := 0
	resultMap := map[string]interface{}{}

	for queryResult.Next() {
		rowMap, err := queryResult.Map()
		if err != nil {
			return nil, errors.Errorf("error while extracting result: %w", err)
		}
		resultMap[strconv.Itoa(row)] = rowMap
		row = row + 1
	}

	resultJSON, err := json.Marshal(resultMap)
	if err != nil {
		return nil, errors.Errorf("error while marshal result: %w", err)
	}

	return &pb.QueryRqliteReply{Result: string(resultJSON)}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck(mngr StatsMngr, p2pService p2p.P2P, metaDb metadb.MetaDB, p2pTracker P2PTracker) *HealthCheck {
	return &HealthCheck{
		StatsMngr:  mngr,
		p2pService: p2pService,
		metaDb:     metaDb,
		p2pTracker: p2pTracker,
	}
}
