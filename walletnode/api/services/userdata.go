package services

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"

	goahttp "goa.design/goa/v3/http"
)

// Userdatas represents services for userdatas endpoints.
type Userdata struct {
	*Common
	process *userdataprocess.Service
}

// Mount configures the mux to serve the artworks endpoints.
func (service *Userdata) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := userdatas.NewEndpoints(service)
	srv := server.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, api.ErrorHandler, nil, &websocket.Upgrader{}, nil, UserdatasProcessUserdataDecoderFunc(ctx, service))
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// ProcessUserdata will send userdata to Super Nodes to store in Metadata layer 
func (service *Userdata) ProcessUserdata(ctx context.Context, req *userdatas.UserdataProcessPayload, stream userdatas.ProcessUserdataServerStream) error {
	defer stream.Close()
	request := fromUserdataProcessRequest(req)
	taskID := service.process.AddTask(request)
	task := service.process.Task(taskID)

	resultChan := task.SubscribeProcessResult()
	for {
		select {
		case <-ctx.Done():
			return nil
		case response, ok := <-resultChan:
			if !ok {
				if task.Status().IsFailure() {
					return artworks.MakeInternalServerError(task.Error())
				}

				return nil
			}

			res := toUserdataProcessResult(response)
			if err := stream.Send(res); err != nil {
				return userdatas.MakeInternalServerError(err)
			}
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
