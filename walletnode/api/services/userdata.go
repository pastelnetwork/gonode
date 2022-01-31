package services

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"

	goahttp "goa.design/goa/v3/http"
)

// UserdataAPIHandler represents services for userdatas endpoints.
type UserdataAPIHandler struct {
	*Common
	process *userdataprocess.UserDataService
}

// Mount configures the mux to serve the nfts endpoints.
func (service *UserdataAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
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
func (service *UserdataAPIHandler) CreateUserdata(ctx context.Context, req *userdatas.CreateUserdataPayload) (*userdatas.UserdataProcessResult, error) {
	return service.processUserdata(ctx, fromUserdataCreateRequest(req))
}

// UpdateUserdata update the userdata in rqlite db
func (service *UserdataAPIHandler) UpdateUserdata(ctx context.Context, req *userdatas.UpdateUserdataPayload) (*userdatas.UserdataProcessResult, error) {
	return service.processUserdata(ctx, fromUserdataUpdateRequest(req))
}

// ProcessUserdata will send userdata to Super Nodes to store in Metadata layer
func (service *UserdataAPIHandler) processUserdata(ctx context.Context, request *userdata.ProcessRequest) (*userdatas.UserdataProcessResult, error) {
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
func (service *UserdataAPIHandler) GetUserdata(ctx context.Context, req *userdatas.GetUserdataPayload) (*userdatas.UserSpecifiedData, error) {
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

// NewUserdataAPIHandler returns the UserdataAPIHandler implementation.
func NewUserdataAPIHandler(process *userdataprocess.UserDataService) *UserdataAPIHandler {
	return &UserdataAPIHandler{
		Common:  NewCommon(),
		process: process,
	}
}
