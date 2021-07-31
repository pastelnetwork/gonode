package supernode

import (
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
	"github.com/pastelnetwork/gonode/metadb/network/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"github.com/pastelnetwork/gonode/metadb/database"
)

// ProcessUserdata represents grpc service for processing userdata.
type ProcessUserdata struct {
	pb.UnimplementedProcessUserdataServer

	*common.ProcessUserdata
	databaseOps *database.DatabaseOps
}

// Session implements supernode.ProcessUserdataServer.Session()
func (service *ProcessUserdata) Session(stream pb.ProcessUserdata_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *userdataprocess.Task

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
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
				return nil
			}
			switch status.Code(err) {
			case codes.Canceled, codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}


// SendUserdataToPrimary implements supernode.ProcessUserdataServer.SendUserdataToPrimary()
func (service *ProcessUserdata) SendUserdataToPrimary(ctx context.Context, req *pb.SuperNodeRequest) (*pb.SuperNodeReply, error) {
	// This code run in primary supernode

	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdataToPrimary request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	snrequest := userdata.SuperNodeRequest{
		UserdataHash		: req.UserdataHash,
		UserdataResultHash	: req.UserdataResultHash,
		HashSignature		: req.HashSignature,
		NodeID				: req.NodeID,
	}

	if err := task.AddPeerSNDataSigned(ctx, snrequest); err != nil {
		errors.Errorf("failed to add peer signature %w", err)
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorPrimarySupernodeFailToProcess,
			Detail: userdata.Description[userdata.ErrorPrimarySupernodeFailToProcess],
		}, nil
	}

	return &pb.SuperNodeReply{
		ResponseCode: userdata.SuccessAddDataToPrimarySupernode,
		Detail: userdata.Description[userdata.SuccessAddDataToPrimarySupernode],
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
	err := service.databaseOps.WriteUserData(ctx, *req)
	if err != nil {
		return nil, err
	}

	return &pb.SuperNodeReply{
		ResponseCode: userdata.SuccessWriteToRQLiteDB,
		Detail: userdata.Description[userdata.SuccessWriteToRQLiteDB],
	}, nil
}

// Desc returns a description of the service.
func (service *ProcessUserdata) Desc() *grpc.ServiceDesc {
	return &pb.ProcessUserdata_ServiceDesc
}

// NewProcessUserdata returns a new ProcessUserdata instance.
func NewProcessUserdata(service *userdataprocess.Service, databaseOps *database.DatabaseOps) *ProcessUserdata {
	return &ProcessUserdata{
		ProcessUserdata: common.NewProcessUserdata(service),
		databaseOps: databaseOps,
	}
}
