package storagechallenge

import (
	"context"

	appcontext "github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"google.golang.org/grpc"
)

// StorageChallenge represents grpc service for storage challenge.
type StorageChallenge struct {
	service storagechallenge.StorageChallenge
	pb.UnimplementedStorageChallengeServer
}

// Desc returns a description of the service.
func (service *StorageChallenge) Desc() *grpc.ServiceDesc {
	return &pb.StorageChallenge_ServiceDesc
}

func (service *StorageChallenge) GenerateStorageChallenges(ctx context.Context, req *pb.GenerateStorageChallengeRequest) (resp *pb.GenerateStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "GenerateStorageChallenges").Debug("handled generate storage challenge action")
	// validate request body
	es := validateGenerateStorageChallengeData(req)
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "GenerateStorageChallenges").Errorf("[GenerateStorageChallenge][Validation Error] %s", es)
		return &pb.GenerateStorageChallengeReply{}, err
	}

	// calling domain service to process bussiness logics
	err = service.service.GenerateStorageChallenges(appcontext.FromContext(ctx), req.GetCurrentBlockHash(), req.GetChallengingMasternodeId(), int(req.GetChallengesPerMasternodePerBlock()))
	return &pb.GenerateStorageChallengeReply{}, err
}

func (service *StorageChallenge) ProcessStorageChallenge(ctx context.Context, req *pb.StorageChallengeRequest) (resp *pb.StorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "ProcessStorageChallenge").Debug("handled process storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "ProcessStorageChallenge").Errorf("[ProcessStorageChallenge][Validation Error] %s", es)
		return &pb.StorageChallengeReply{Data: req.GetData()}, err
	}

	// calling domain service to process bussiness logics
	err = service.service.ProcessStorageChallenge(appcontext.FromContext(ctx), mapChallengeMessage(req.GetData()))
	return &pb.StorageChallengeReply{Data: req.GetData()}, err
}

func (service *StorageChallenge) VerifyStorageChallenge(ctx context.Context, req *pb.VerifyStorageChallengeRequest) (resp *pb.VerifyStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "VerifyStorageChallenge").Debug("handled verify storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "VerifyStorageChallenge").Errorf("[VerifyStorageChallenge][Validation Error] %s", es)
		return &pb.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}
	// calling domain service to process bussiness logics
	err = service.service.VerifyStorageChallenge(appcontext.FromContext(ctx), mapChallengeMessage(req.GetData()))
	return &pb.VerifyStorageChallengeReply{Data: req.GetData()}, err
}

// NewStorageChallenge returns a new StorageChallenge instance.
func NewStorageChallenge(service storagechallenge.StorageChallenge) *StorageChallenge {
	return &StorageChallenge{
		service: service,
	}
}
