package services

import (
	"context"

	// "github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"

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
	srv := server.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, api.ErrorHandler, nil, UserdatasProcessUserdataDecoderFunc(ctx, service))
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// ProcessUserdata will send userdata to Super Nodes to store in Metadata layer 
func (service *Userdata) ProcessUserdata(ctx context.Context, req *userdatas.ProcessUserdataPayload) (*userdatas.UserSpecifiedData, error) {
	request := fromUserdataProcessRequest(req)
	taskID := service.process.AddTask(request)
	task := service.process.Task(taskID)

	resultChanGet := task.SubscribeProcessResultGet()
	for {
		select {
		case <-ctx.Done():
			return nil,nil
		case response, ok := <-resultChanGet:
			if !ok {
				if task.Status().IsFailure() {
					return nil,userdatas.MakeInternalServerError(task.Error())
				}

				return nil,nil
			}

			res := toUserSpecifiedData(response)
			return res,nil
		}
	}
}

// UserdataGet will get userdata from Super Nodes to store in Metadata layer 
func (service *Userdata) UserdataGet (ctx context.Context, pastelid string) (*userdatas.ProcessUserdataPayload, error) {
	taskID := service.process.AddTask(nil, request)
	task := service.process.Task(taskID)

	resultChan := task.SubscribeProcessResult()
	for {
		select {
		case <-ctx.Done():
			return nil,nil
		case response, ok := <-resultChan:
			if !ok {
				if task.Status().IsFailure() {
					return nil,userdatas.MakeInternalServerError(task.Error())
				}

				return nil,nil
			}

			res := toUserdataProcessResult(response)
			return res,nil
		}
	}
}

// NewUserdata returns the Userdata implementation.
func NewUserdata(process *userdataprocess.Service) *Userdata {
	return &Userdata{
		Common:   NewCommon(),
		process: process,
	}
}
