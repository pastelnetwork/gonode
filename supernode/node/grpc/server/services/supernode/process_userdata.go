package supernode

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

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

const defaultGRPCPort = 4444

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

	// Primary node will connect to leader node here.
	if !service.databaseOps.IsLeader() {
		var extAddress string

		leaderAddress := service.databaseOps.LeaderAddress()
		parts := strings.Split(leaderAddress, ":")
		//temporary use this to test on local machine
		localEnv := false
		if localEnv {
			port := parts[len(parts)-1]
			portnum, err := strconv.Atoi(port)
			if err != nil {
				return nil, errors.Errorf("cannot convert port to num %s / err: %w", port, err)
			}
			//4041 -> 4444, 4051 -> 4445
			portnum = (portnum/10)%10 + 4440
			extAddress = fmt.Sprintf("127.0.0.1:%d", portnum)
		} else {
			ipAddress := parts[0]
			extAddress = fmt.Sprintf("%s:%d", ipAddress, defaultGRPCPort)
		}

		if err := task.ConnectToLeader(ctx, extAddress); err != nil {
			return nil, err
		}
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
	if len(req.Data) == 0 {
		// This is user specified data (user profile data)
		if err := service.databaseOps.WriteUserData(ctx, req); err != nil {
			return nil, errors.Errorf("error occurs while writting to database: %w", err)
		}

		return &pb.SuperNodeReply{
			ResponseCode: userdata.SuccessWriteToRQLiteDB,
			Detail:       userdata.Description[userdata.SuccessWriteToRQLiteDB],
		}, nil
	}

	// This is walletnode metric (for both get/set data)
	request := pb.Metric{
		Command: req.Command,
		Data:    req.Data,
		// No need to pass Signature of the data or PastelID to the database operation
	}

	var result interface{}
	var err error
	if result, err = service.databaseOps.ProcessCommand(ctx, &request); err != nil {
		return nil, errors.Errorf("error occurs while process metric to database: %w", err)
	}

	// Marshal the response from MetadataLayer
	js, err := json.Marshal(result)
	if err != nil {
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorProcessMetric,
			Detail:       userdata.Description[userdata.ErrorProcessMetric],
		}, nil
	}

	return &pb.SuperNodeReply{
		ResponseCode: userdata.SuccessProcessMetric,
		Detail:       userdata.Description[userdata.SuccessProcessMetric],
		Data:         js,
	}, nil

}

// StoreMetric implements supernode.ProcessUserdataServer.StoreMetric()
func (service *ProcessUserdata) StoreMetric(ctx context.Context, req *pb.Metric) (*pb.SuperNodeReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log.WithContext(ctx).WithField("req", req).Debugf("StoreMetric request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	task := service.NewTask()
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	// Get hash of the data to verify
	hashvalue, err := userdata.Sha3256hash(req.Data)
	if err != nil {
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorSupernodeVerifyMetricFail,
			Detail:       userdata.Description[userdata.ErrorSupernodeVerifyMetricFail],
		}, nil
	}

	// Validate user signature
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		log.WithContext(ctx).Debugf("failed to decode signature %s of node %s", req.Signature, req.PastelID)
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, nil
	}

	ok, err := task.Service.PastelClient.Verify(ctx, hashvalue, string(signature), req.PastelID, "ed448")
	if err != nil || !ok {
		log.WithContext(ctx).Debugf("failed to verify signature %s of node %s", req.Signature, req.PastelID)
		return &pb.SuperNodeReply{
			ResponseCode: userdata.ErrorVerifyUserdataFail,
			Detail:       userdata.Description[userdata.ErrorVerifyUserdataFail],
		}, nil
	}
	if service.databaseOps.IsLeader() {
		var result interface{}
		if result, err = service.databaseOps.ProcessCommand(ctx, req); err != nil {
			return nil, errors.Errorf("error occurs while writting to database: %w", err)
		}

		js, err := json.Marshal(result)
		if err != nil {
			return nil, errors.Errorf("error occurs while Marshal MetadataLayer result: %w", err)
		}

		return &pb.SuperNodeReply{
			ResponseCode: userdata.SuccessWriteToRQLiteDB,
			Detail:       userdata.Description[userdata.SuccessWriteToRQLiteDB],
			Data:         js,
		}, nil
	}

	if err := task.ConnectToLeader(ctx, service.databaseOps.LeaderAddress()); err != nil {
		return nil, err
	}

	metric := userdata.Metric{
		Command: req.Command,
		Data:    req.Data,
	}

	if task.ConnectedToLeader != nil {
		resp := userdata.SuperNodeReply{}
		if resp, err = task.ConnectedToLeader.ProcessUserdata.StoreMetric(ctx, metric); err != nil {
			return nil, errors.Errorf("error occurs while writting to database: %w", err)
		}
		return &pb.SuperNodeReply{
			ResponseCode: userdata.SuccessWriteToRQLiteDB,
			Detail:       userdata.Description[userdata.SuccessWriteToRQLiteDB],
			Data:         resp.Data,
		}, nil
	}

	log.WithContext(ctx).Debugf("Process Metric ConnectedToLeader is null")
	return &pb.SuperNodeReply{
		ResponseCode: userdata.ErrorProcessMetric,
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
