package healthcheck

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/p2p"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"google.golang.org/grpc"
)

const (
	defaultInputMaxLength  = 64
	defaultExpiresDuration = 10 * time.Minute
)

// StatsMngr is interface of StatsManger, return stats of system
type StatsMngr interface {
	// Stats returns stats of system
	Stats(ctx context.Context) (map[string]interface{}, error)
}

// HealthCheck represents grpc service for supernode healthcheck
type HealthCheck struct {
	pb.UnimplementedHealthCheckServer
	StatsMngr
	p2pService    p2p.P2P
	rqliteKVStore *RqliteKVStore
	cleanTracker  CleanTracker
}

func (service *HealthCheck) p2pClean(ctx context.Context, key string) {
	err := service.p2pService.Delete(ctx, key)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("key", key).Error("P2PRemoveKeyFailed")
	} else {
		log.WithContext(ctx).WithField("key", key).Info("P2PRemoveKeySuccessfully")
	}
}

func (service *HealthCheck) rqliteClean(ctx context.Context, key string) {
	err := service.rqliteKVStore.Delete(ctx, key)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("key", key).Error("RqliteRemoveKeyFailed")
	} else {
		log.WithContext(ctx).WithField("key", key).Info("RqliteRemoveKeySuccessfully")
	}
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
	if len(value) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}
	key, err := service.p2pService.Store(ctx, []byte(value))

	if err != nil {
		return nil, err
	}

	// track to delete key soon
	service.cleanTracker.Track(key, service.p2pClean, time.Now().Add(defaultExpiresDuration))

	return &pb.P2PStoreReply{Key: key}, nil
}

// P2PRetrieve get data of given key
func (service *HealthCheck) P2PRetrieve(ctx context.Context, in *pb.P2PRetrieveRequest) (*pb.P2PRetrieveReply, error) {
	key := in.GetKey()
	value, err := service.p2pService.Retrieve(ctx, key)

	if err != nil {
		return &pb.P2PRetrieveReply{Value: ""}, err
	}
	return &pb.P2PRetrieveReply{Value: string(value)}, nil
}

// RqliteStore do store a value - and return a key - is hash of string
func (service *HealthCheck) RqliteStore(ctx context.Context, in *pb.RqliteStoreRequest) (*pb.RqliteStoreReply, error) {
	if len(in.GetValue()) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}

	value := in.GetValue()
	key, err := service.rqliteKVStore.Store(ctx, []byte(value))
	if err != nil {
		return nil, err
	}

	// track to delete key soon
	service.cleanTracker.Track(key, service.rqliteClean, time.Now().Add(defaultExpiresDuration))

	return &pb.RqliteStoreReply{Key: key}, nil
}

// RqliteRetrieve get data of given key
func (service *HealthCheck) RqliteRetrieve(ctx context.Context, in *pb.RqliteRetrieveRequest) (*pb.RqliteRetrieveReply, error) {
	if len(in.GetKey()) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}

	key := in.GetKey()
	value, err := service.rqliteKVStore.Retrieve(ctx, key)

	if err != nil {
		return nil, err
	}

	return &pb.RqliteRetrieveReply{Value: string(value)}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck(mngr StatsMngr, p2pService p2p.P2P, metaDb metadb.MetaDB, cleanTracker CleanTracker) *HealthCheck {
	return &HealthCheck{
		StatsMngr:     mngr,
		p2pService:    p2pService,
		rqliteKVStore: NewRqliteKVStore(metaDb),
		cleanTracker:  cleanTracker,
	}
}
