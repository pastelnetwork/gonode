package services

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/collection/server"
	"github.com/pastelnetwork/gonode/walletnode/services/collectionregister"
	goahttp "goa.design/goa/v3/http"
)

// CollectionAPIHandler - CollectionAPIHandler service
type CollectionAPIHandler struct {
	*Common
	register *collectionregister.CollectionRegistrationService
}

// Mount configures the mux to serve the OpenAPI enpoints.
func (service *CollectionAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := collection.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// RegisterCollection - Starts a processing for a collection registration
func (service *CollectionAPIHandler) RegisterCollection(ctx context.Context, p *collection.RegisterCollectionPayload) (res *collection.RegisterCollectionResponse, err error) {
	taskID, err := service.register.AddTask(p)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to add task")
		return nil, cascade.MakeInternalServerError(err)
	}

	res = &collection.RegisterCollectionResponse{
		TaskID: taskID,
	}

	log.Infof("task has been added: %s", taskID)

	return res, nil
}
