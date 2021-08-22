package services

import (
	"context"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"

	goahttp "goa.design/goa/v3/http"
)

// Userdata represents services for userdatas endpoints.
type Userdata struct {
	*Common
	process *userdataprocess.Service
}

// Mount configures the mux to serve the artworks endpoints.
func (service *Userdata) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := userdatas.NewEndpoints(service)
	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		UserdatasCreateUserdataDecoderFunc(ctx, service),
		UserdatasUpdateUserdataDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// CreateUserdata create the userdata in rqlite db
func (service *Userdata) CreateUserdata(ctx context.Context, req *userdatas.CreateUserdataPayload) (*userdatas.UserdataProcessResult, error) {
	return service.processUserdata(ctx, fromUserdataCreateRequest(req))
}

// UpdateUserdata update the userdata in rqlite db
func (service *Userdata) UpdateUserdata(ctx context.Context, req *userdatas.UpdateUserdataPayload) (*userdatas.UserdataProcessResult, error) {
	return service.processUserdata(ctx, fromUserdataUpdateRequest(req))
}

// ProcessUserdata will send userdata to Super Nodes to store in Metadata layer
func (service *Userdata) processUserdata(ctx context.Context, request *userdata.ProcessRequest) (*userdatas.UserdataProcessResult, error) {
	taskID := service.process.AddTask(request, "")
	task := service.process.Task(taskID)

	resultChan := task.SubscribeProcessResult()
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case response, ok := <-resultChan:
			if !ok {
				if task.Status().IsFailure() {
					return nil, userdatas.MakeInternalServerError(task.Error())
				}

				return nil, userdatas.MakeInternalServerError(errors.New("no info retrieve"))
			}

			res := toUserdataProcessResult(response)
			return res, nil
		}
	}
}

// GetUserdata will get userdata from Super Nodes to store in Metadata layer
func (service *Userdata) GetUserdata(ctx context.Context, req *userdatas.GetUserdataPayload) (*userdatas.UserSpecifiedData, error) {
	userpastelid := req.Pastelid

	taskID := service.process.AddTask(nil, userpastelid)
	task := service.process.Task(taskID)

	resultChanGet := task.SubscribeProcessResultGet()
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case response, ok := <-resultChanGet:
			if !ok {
				if task.Status().IsFailure() {
					return nil, userdatas.MakeInternalServerError(task.Error())
				}

				return nil, nil
			}

			res := toUserSpecifiedData(response)
			return res, nil
		}
	}
}

// SetUserFollowRelation Set a follower, followee relationship to metadb
func (service *Userdata) SetUserFollowRelation(ctx context.Context, req *userdatas.SetUserFollowRelationPayload) (*userdatas.SetUserFollowRelationResult, error) {
	// Generalize the data to be get/set by marshaling it
	data := userdata.UserFollow{
		FollowerPastelID: req.FollowerPastelID,
		FolloweePastelID: req.FolloweePastelID,
	}
	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandUserFollowWrite,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.SetUserFollowRelationResult{
		ResponseCode: result.ResponseCode,
		Detail:       result.Detail,
	}, nil
}

// GetFollowers Get followers of a user
func (service *Userdata) GetFollowers(ctx context.Context, req *userdatas.GetFollowersPayload) (*userdatas.GetFollowersResult, error) {
	// Generalize the data to be get/set by marshaling it
	data := userdata.PaginationIDStringQuery{
		ID: req.Pastelid,
	}
	if req.Limit != nil {
		data.Limit = *req.Limit
	}
	if req.Offset != nil {
		data.Offset = *req.Offset
	}

	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandGetFollowers,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	var relationResult userdata.UserRelationshipQueryResult
	if err := json.Unmarshal(result.Data, &relationResult); err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.GetFollowersResult{
		TotalCount: relationResult.TotalCount,
		Result:     toRelationshipInfoArray(relationResult.Items),
	}, nil
}

// GetFollowees Get followees of a user
func (service *Userdata) GetFollowees(ctx context.Context, req *userdatas.GetFolloweesPayload) (*userdatas.GetFolloweesResult, error) {
	data := userdata.PaginationIDStringQuery{
		ID: req.Pastelid,
	}
	if req.Limit != nil {
		data.Limit = *req.Limit
	}
	if req.Offset != nil {
		data.Offset = *req.Offset
	}

	// Generalize the data to be get/set by marshaling it
	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandGetFollowees,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	var relationResult userdata.UserRelationshipQueryResult
	if err := json.Unmarshal(result.Data, &relationResult); err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.GetFolloweesResult{
		TotalCount: relationResult.TotalCount,
		Result:     toRelationshipInfoArray(relationResult.Items),
	}, nil
}

// GetFriends Get friends of a user
func (service *Userdata) GetFriends(ctx context.Context, req *userdatas.GetFriendsPayload) (*userdatas.GetFriendsResult, error) {
	data := userdata.PaginationIDStringQuery{
		ID: req.Pastelid,
	}
	if req.Limit != nil {
		data.Limit = *req.Limit
	}
	if req.Offset != nil {
		data.Offset = *req.Offset
	}

	// Generalize the data to be get/set by marshaling it
	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandGetFriend,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	var relationResult userdata.UserRelationshipQueryResult
	if err := json.Unmarshal(result.Data, &relationResult); err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.GetFriendsResult{
		TotalCount: relationResult.TotalCount,
		Result:     toRelationshipInfoArray(relationResult.Items),
	}, nil
}

// SetUserLikeArt Notify a new like event of an user to an art
func (service *Userdata) SetUserLikeArt(ctx context.Context, req *userdatas.SetUserLikeArtPayload) (*userdatas.SetUserLikeArtResult, error) {
	data := userdata.ArtLike{
		ArtID:    req.ArtPastelID,
		PastelID: req.UserPastelID,
	}
	// Generalize the data to be get/set by marshaling it
	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandArtLikeWrite,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.SetUserLikeArtResult{
		ResponseCode: result.ResponseCode,
		Detail:       result.Detail,
	}, nil
}

// GetUsersLikeArt Get users that liked an art
func (service *Userdata) GetUsersLikeArt(ctx context.Context, req *userdatas.GetUsersLikeArtPayload) (*userdatas.GetUsersLikeArtResult, error) {
	data := userdata.PaginationIDStringQuery{
		ID: req.ArtID,
	}
	if req.Limit != nil {
		data.Limit = *req.Limit
	}
	if req.Offset != nil {
		data.Offset = *req.Offset
	}

	// Generalize the data to be get/set by marshaling it
	js, err := json.Marshal(data)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Generate the request
	request := userdata.ProcessRequest{
		Command: userdata.CommandUsersLikeNft,
		Data:    js,
	}

	// Send the request to set it in Metadata Layer
	result, err := service.processUserdata(ctx, &request)
	if err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	var relationResult userdata.UserRelationshipQueryResult
	if err := json.Unmarshal(result.Data, &relationResult); err != nil {
		return nil, userdatas.MakeInternalServerError(err)
	}

	// Return the result of Metadata Layer process this request
	return &userdatas.GetUsersLikeArtResult{
		TotalCount: relationResult.TotalCount,
		Result:     toRelationshipInfoArray(relationResult.Items),
	}, nil
}

// NewUserdata returns the Userdata implementation.
func NewUserdata(process *userdataprocess.Service) *Userdata {
	return &Userdata{
		Common:  NewCommon(),
		process: process,
	}
}
