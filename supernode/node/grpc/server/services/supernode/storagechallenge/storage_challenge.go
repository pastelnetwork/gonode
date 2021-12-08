package storagechallenge

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	appcontext "github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"google.golang.org/grpc"
)

// AppSvc interface
type AppSvc interface {
	pb.StorageChallengeServer
	Desc() *grpc.ServiceDesc
}

// appSc represents grpc service for storage challenge.
type appSc struct {
	actor       messaging.Actor
	appActorPID *actor.PID
	pb.UnimplementedStorageChallengeServer
}

// Desc returns a description of the service.
func (service *appSc) Desc() *grpc.ServiceDesc {
	return &pb.StorageChallenge_ServiceDesc
}

// GenerateStorageChallenges gRPC handler
func (service *appSc) GenerateStorageChallenges(ctx context.Context, req *pb.GenerateStorageChallengeRequest) (resp *pb.GenerateStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "GenerateStorageChallenges").Debug("handled generate storage challenge action")
	// validate request body
	es := validateGenerateStorageChallengeData(req)
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "GenerateStorageChallenges").Errorf("[GenerateStorageChallenge][Validation Error] %s", es)
		return &pb.GenerateStorageChallengeReply{}, err
	}
	appCtx := appcontext.FromContext(ctx)
	// calling async actor to process storage challenge
	err = service.actor.Send(appCtx, service.appActorPID, newGenerateStorageChallengeMsg(appCtx, req.GetMerkleroot(), req.GetChallengingMasternodeId(), req.GetChallengesPerMasternodePerBlock()))
	return &pb.GenerateStorageChallengeReply{}, err
}

// ProcessStorageChallenge gRPC handler
func (service *appSc) ProcessStorageChallenge(ctx context.Context, req *pb.ProcessStorageChallengeRequest) (resp *pb.ProcessStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "ProcessStorageChallenge").Debug("handled process storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "ProcessStorageChallenge").Errorf("[ProcessStorageChallenge][Validation Error] %s", es)
		return &pb.ProcessStorageChallengeReply{Data: req.GetData()}, err
	}

	appCtx := appcontext.FromContext(ctx)
	// calling async actor to process storage challenge
	err = service.actor.Send(appCtx, service.appActorPID, newProcessStorageChallengeMsg(appCtx, mapChallengeMessage(req.GetData())))
	return &pb.ProcessStorageChallengeReply{Data: req.GetData()}, err
}

// VerifyStorageChallenge gRPC handler
func (service *appSc) VerifyStorageChallenge(ctx context.Context, req *pb.VerifyStorageChallengeRequest) (resp *pb.VerifyStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("grpc-server", "VerifyStorageChallenge").Debug("handled verify storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("grpc-server", "VerifyStorageChallenge").Errorf("[VerifyStorageChallenge][Validation Error] %s", es)
		return &pb.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}

	appCtx := appcontext.FromContext(ctx)
	// calling async actor to process verify storage challenge
	err = service.actor.Send(appCtx, service.appActorPID, newVerifyStorageChallengeMsg(appCtx, mapChallengeMessage(req.GetData())))
	return &pb.VerifyStorageChallengeReply{Data: req.GetData()}, err
}

// NewStorageChallenge returns a new StorageChallenge instance.
func NewStorageChallenge(domainService storagechallenge.StorageChallenge) (appSvc AppSvc, stopActor func()) {
	logger := log.DefaultLogger
	localActor := messaging.NewActor(actor.NewActorSystem())
	pid, err := localActor.RegisterActor(&applicationActor{
		domainService: domainService,
		logger:        logger,
	}, "application_actor")
	if err != nil {
		panic(err)
	}

	return &appSc{
		actor:       localActor,
		appActorPID: pid,
	}, localActor.Stop
}
