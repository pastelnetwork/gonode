package walletnode

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb/database"
	pbsn "github.com/pastelnetwork/gonode/proto/supernode"
	pbwn "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ProcessUserdata represents grpc service for processing userdata.
type ProcessUserdata struct {
	pbwn.UnimplementedProcessUserdataServer

	*common.ProcessUserdata
	databaseOps *database.Ops
}

// Session implements walletnode.ProcessUserdataServer.Session()
func (service *ProcessUserdata) Session(stream pbwn.ProcessUserdata_SessionServer) error {
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

	if err := task.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pbwn.MDLSessionReply{
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

// AcceptedNodes implements walletnode.ProcessUserdataServer.AcceptedNodes()
func (service *ProcessUserdata) AcceptedNodes(ctx context.Context, req *pbwn.MDLAcceptedNodesRequest) (*pbwn.MDLAcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.AcceptedNodes(ctx)
	if err != nil {
		return nil, err
	}

	var peers []*pbwn.MDLAcceptedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pbwn.MDLAcceptedNodesReply_Peer{
			NodeID: node.ID,
		})
	}

	resp := &pbwn.MDLAcceptedNodesReply{
		Peers: peers,
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.ProcessUserdataServer.ConnectTo()
func (service *ProcessUserdata) ConnectTo(ctx context.Context, req *pbwn.MDLConnectToRequest) (*pbwn.MDLConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pbwn.MDLConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")
	return resp, nil
}

// SendUserdata implements walletnode.ProcessUserdataServer.SendUserdata()
func (service *ProcessUserdata) SendUserdata(ctx context.Context, req *pbwn.UserdataRequest) (*pbwn.UserdataReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdata request")

	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.Errorf("Request receive from walletnode is nil")
	}
	// Convert protobuf request to UserdataProcessRequest
	request := userdata.ProcessRequestSigned{
		Userdata: &userdata.ProcessRequest{
			RealName:        req.RealName,
			FacebookLink:    req.FacebookLink,
			TwitterLink:     req.TwitterLink,
			NativeCurrency:  req.NativeCurrency,
			Location:        req.Location,
			PrimaryLanguage: req.PrimaryLanguage,
			Categories:      req.Categories,
			Biography:       req.Biography,
			AvatarImage: userdata.UserImageUpload{
				Content:  req.AvatarImage.Content,
				Filename: req.AvatarImage.Filename,
			},
			CoverPhoto: userdata.UserImageUpload{
				Content:  req.CoverPhoto.Content,
				Filename: req.CoverPhoto.Filename,
			},
			ArtistPastelID:    req.ArtistPastelID,
			Timestamp:         req.Timestamp,
			PreviousBlockHash: req.PreviousBlockHash,
			Command:           req.Command,
			Data:              req.Data,
		},
		UserdataHash: req.UserdataHash,
		Signature:    req.Signature,
	}

	processResult, err := task.SupernodeProcessUserdata(ctx, &request)
	if err != nil {
		return nil, errors.Errorf("SupernodeProcessUserdata can not process %w", err)
	}
	if processResult.ResponseCode == userdata.ErrorOnContent {
		return &pbwn.UserdataReply{
			ResponseCode:    processResult.ResponseCode,
			Detail:          processResult.Detail,
			RealName:        processResult.RealName,
			FacebookLink:    processResult.FacebookLink,
			TwitterLink:     processResult.TwitterLink,
			NativeCurrency:  processResult.NativeCurrency,
			Location:        processResult.Location,
			PrimaryLanguage: processResult.PrimaryLanguage,
			Categories:      processResult.Categories,
			Biography:       processResult.Biography,
			AvatarImage:     processResult.AvatarImage,
			CoverPhoto:      processResult.CoverPhoto,
		}, nil
	}

	if task.ConnectedTo != nil {
		// This is secondary node, we just response to walletnode here
		return &pbwn.UserdataReply{
			ResponseCode: processResult.ResponseCode,
			Detail:       processResult.Detail,
		}, nil
	}

	// Process actual write to rqlite db happen here
	var actionErr error
	<-task.NewAction(func(ctx context.Context) error {
		if processResult.ResponseCode == userdata.SuccessVerifyAllSignature {
			if service.databaseOps == nil {
				processResult.ResponseCode = userdata.ErrorRQLiteDBNotFound
				processResult.Detail = userdata.Description[userdata.ErrorRQLiteDBNotFound]
				actionErr = errors.Errorf("databaseOps service object is empty")
				return nil
			}
			// Send data to SN contain the leader rqlite
			if !service.databaseOps.IsLeader() {
				if task.ConnectedToLeader != nil {
					if _, err := task.ConnectedToLeader.ProcessUserdata.SendUserdataToLeader(ctx, request); err != nil {
						processResult.ResponseCode = userdata.ErrorWriteToRQLiteDBFail
						processResult.Detail = userdata.Description[userdata.ErrorWriteToRQLiteDBFail]
						actionErr = errors.Errorf("failed to send or write userdata to leader rqlite %w", err)
						return nil
					}
					// Write success:
					processResult.ResponseCode = userdata.SuccessProcess
					processResult.Detail = userdata.Description[userdata.SuccessProcess]

				} else {
					processResult.ResponseCode = userdata.ErrorRQLiteDBNotFound
					processResult.Detail = userdata.Description[userdata.ErrorRQLiteDBNotFound]
					actionErr = errors.Errorf("leader rqlite node object is empty")
					return nil
				}
			} else {
				// This supernode contain rqlite leader, write to db here

				if len((*req).Data) == 0 {
					// This is user specified data (user profile data)
					reqsn := pbsn.UserdataRequest{
						RealName:        (*req).RealName,
						FacebookLink:    (*req).FacebookLink,
						TwitterLink:     (*req).TwitterLink,
						NativeCurrency:  (*req).NativeCurrency,
						Location:        (*req).Location,
						PrimaryLanguage: (*req).PrimaryLanguage,
						Categories:      (*req).Categories,
						Biography:       (*req).Biography,
						AvatarImage: &pbsn.UserdataRequest_UserImageUpload{
							Content:  (*req).AvatarImage.Content,
							Filename: (*req).AvatarImage.Filename,
						},
						CoverPhoto: &pbsn.UserdataRequest_UserImageUpload{
							Content:  (*req).CoverPhoto.Content,
							Filename: (*req).CoverPhoto.Filename,
						},
						ArtistPastelID:    (*req).ArtistPastelID,
						Timestamp:         (*req).Timestamp,
						Signature:         (*req).Signature,
						PreviousBlockHash: (*req).PreviousBlockHash,
						Command:           (*req).Command,
						Data:              (*req).Data,
					}

					err := service.databaseOps.WriteUserData(ctx, &reqsn)
					if err != nil {
						processResult.ResponseCode = userdata.ErrorWriteToRQLiteDBFail
						processResult.Detail = userdata.Description[userdata.ErrorWriteToRQLiteDBFail]
						actionErr = errors.Errorf("error while writting data: %w", err)
						return nil
					}

					// If can go to here, all process of setting userdata have passed
					// We can consider to pass userdata.SuccessProcess or userdata.SuccessWriteToRQLiteDB (both have same success meaning)
					processResult.ResponseCode = userdata.SuccessProcess
					processResult.Detail = userdata.Description[userdata.SuccessProcess]
				} else {
					// This is walletnode metric (for both get/set data)
					request := pbsn.Metric{
						Command: (*req).Command,
						Data:    (*req).Data,
						// No need to pass Signature of the data or PastelID to the database operation
					}

					var result interface{}
					var err error
					if result, err = service.databaseOps.ProcessCommand(ctx, &request); err != nil {
						log.WithContext(ctx).Debugf("Error ProcessCommand:%s, Data %s, err:%s", request.Command, string(request.Data), err.Error())
						processResult.ResponseCode = userdata.ErrorProcessMetric
						processResult.Detail = userdata.Description[userdata.ErrorProcessMetric]
						actionErr = errors.Errorf("error while processing command: %w", err)
						return nil
					}

					// Marshal the response from MetadataLayer
					js, err := json.Marshal(result)
					if err != nil {
						log.WithContext(ctx).Debugf("Error Marshal result: %s")
						processResult.ResponseCode = userdata.ErrorProcessMetric
						processResult.Detail = userdata.Description[userdata.ErrorProcessMetric]
						actionErr = errors.Errorf("error while marshaling command: %w", err)
						return nil
					}

					processResult.ResponseCode = userdata.SuccessProcess
					processResult.Detail = userdata.Description[userdata.SuccessProcess]
					processResult.Data = js
					return nil
				}
			}
		}
		return nil
	})

	return &pbwn.UserdataReply{
		ResponseCode: processResult.ResponseCode,
		Detail:       processResult.Detail,
		Data:         processResult.Data,
	}, actionErr
}

// ReceiveUserdata implements walletnode.ProcessUserdataServer.ReceiveUserdata()
func (service *ProcessUserdata) ReceiveUserdata(ctx context.Context, req *pbwn.RetrieveRequest) (*pbwn.UserdataRequest, error) {
	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}
	log.WithContext(ctx).WithField("userpastelid", req.Userpastelid).Debugf("ReceiveUserdata request")

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	userpastelid := req.Userpastelid
	result, err := task.ReceiveUserdata(ctx, userpastelid)
	if err != nil {
		return nil, err
	}

	// Generate protobuf response respProto
	respProto := &pbwn.UserdataRequest{
		RealName:        result.RealName,
		FacebookLink:    result.FacebookLink,
		TwitterLink:     result.TwitterLink,
		NativeCurrency:  result.NativeCurrency,
		Location:        result.Location,
		PrimaryLanguage: result.PrimaryLanguage,
		Categories:      result.Categories,
		Biography:       result.Biography,
		AvatarImage: &pbwn.UserdataRequest_UserImageUpload{
			Content:  result.AvatarImage.Content,
			Filename: result.AvatarImage.Filename,
		},
		CoverPhoto: &pbwn.UserdataRequest_UserImageUpload{
			Content:  result.CoverPhoto.Content,
			Filename: result.CoverPhoto.Filename,
		},
		ArtistPastelID:    result.ArtistPastelID,
		Timestamp:         result.Timestamp,
		PreviousBlockHash: result.PreviousBlockHash,
		Command:           result.Command,
		Data:              result.Data,
	}

	return respProto, nil
}

// Desc returns a description of the service.
func (service *ProcessUserdata) Desc() *grpc.ServiceDesc {
	return &pbwn.ProcessUserdata_ServiceDesc
}

// NewProcessUserdata returns a new ProcessUserdata instance.
func NewProcessUserdata(service *userdataprocess.Service, databaseOps *database.Ops) *ProcessUserdata {
	return &ProcessUserdata{
		ProcessUserdata: common.NewProcessUserdata(service),
		databaseOps:     databaseOps,
	}
}
