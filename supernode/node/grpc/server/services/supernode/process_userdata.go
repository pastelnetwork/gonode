package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/service/userdata"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/database"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)
// ProcessUserdata represents grpc service for processing userdata.
type ProcessUserdata struct {
	pb.UnimplementedProcessUserdataServer

	*common.ProcessUserdata
	databaseOps *database.Ops
}

// Session implements supernode.ProcessUserdataServer.Session()
func (service *ProcessUserdata) Session(stream pb.ProcessUserdata_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *userdataprocess.Task

	if sessID, ok := service.SessID(ctx); ok {
		task = service.Task(sessID)
		if task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("failed to receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := task.SessionNode(ctx, req.NodeID); err != nil {
		return err
	}

	resp := &pb.MDLSessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("failed to send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")

	for {
		if _, err := stream.Recv(); err != nil {
			if err == io.EOF {
				return errors.Errorf("handshake stream closed: %w", err)
			}
		}
	}
}

// SendUserdataToPrimary implements supernode.ProcessUserdataServer.SendUserdataToPrimary()
func (service *ProcessUserdata) SendUserdataToPrimary(ctx context.Context, req *pb.SuperNodeRequest) (*pb.SuperNodeReply, error) {
	// This code run in primary supernode
	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdataToPrimary request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	snrequest := userdata.SuperNodeRequest{
		UserdataHash:       req.UserdataHash,
		UserdataResultHash: req.UserdataResultHash,
		HashSignature:      req.HashSignature,
		NodeID:             req.NodeID,
	}

	if err := task.AddPeerSNDataSigned(ctx, snrequest); err != nil {
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorPrimarySupernodeFailToProcess,
			Detail:       userdata.Description[userdata.ErrorPrimarySupernodeFailToProcess],
		}, nil
	}

	return &pb.SuperNodeReply{
		ResponseCode: userdata.SuccessAddDataToPrimarySupernode,
		Detail:       userdata.Description[userdata.SuccessAddDataToPrimarySupernode],
	}, nil
}

// SendUserdataToLeader implements supernode.ProcessUserdataServer.SendUserdataToLeader()
func (service *ProcessUserdata) SendUserdataToLeader(ctx context.Context, req *pb.UserdataRequest) (*pb.SuperNodeReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdataToLeader request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	// This code run in supernode contain leader rqlite db
	// Process write the data to rqlite happen here
	if err := service.databaseOps.WriteUserData(ctx, req); err != nil {
		return nil, errors.Errorf("error occurs while writting to database: %w", err)
	}

	return &pb.SuperNodeReply{
		ResponseCode: userdata.SuccessWriteToRQLiteDB,
		Detail:       userdata.Description[userdata.SuccessWriteToRQLiteDB],
	}, nil
}

// Desc returns a description of the service.
func (service *ProcessUserdata) Desc() *grpc.ServiceDesc {
	return &pb.ProcessUserdata_ServiceDesc
}

// NewProcessUserdata returns a new ProcessUserdata instance.
func NewProcessUserdata(service *userdataprocess.Service, databaseOps *database.Ops) *ProcessUserdata {
	return &ProcessUserdata{
		ProcessUserdata: common.NewProcessUserdata(service),
		databaseOps:     databaseOps,
	}
}
