package services

import (
	"context"
	"time"

	"goa.design/goa/v3/security"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
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
		&websocket.Upgrader{},
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

// RegisterTaskState - Registers a task state
func (service *CollectionAPIHandler) RegisterTaskState(ctx context.Context, p *collection.RegisterTaskStatePayload, stream collection.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		log.Error("unable to get task")
		return collection.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &collection.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return collection.MakeInternalServerError(err)
			}

			if status.IsFailure() {
				if task.Error() != nil {
					errStr := task.Error()
					log.WithContext(ctx).WithError(errStr).Errorf("error registering cascade")
				}

			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// GetTaskHistory - Gets a task's history
func (service *CollectionAPIHandler) GetTaskHistory(ctx context.Context, p *collection.GetTaskHistoryPayload) (history []*collection.TaskHistory, err error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, collection.MakeInternalServerError(errors.New("error retrieving status"))
	}
	defer store.CloseHistoryDB(ctx)

	statuses, err := store.QueryTaskHistory(p.TaskID)
	if err != nil {
		return nil, collection.MakeNotFound(errors.New("task not found"))
	}

	for _, entry := range statuses {
		timestamp := entry.CreatedAt.String()
		historyItem := &collection.TaskHistory{
			Timestamp: &timestamp,
			Status:    entry.Status,
		}

		if entry.Details != nil {
			historyItem.Details = &collection.Details{
				Message: &entry.Details.Message,
				Fields:  entry.Details.Fields,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *CollectionAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// NewCollectionAPIIHandler returns the swagger OpenAPI implementation.
func NewCollectionAPIIHandler(config *Config, register *collectionregister.CollectionRegistrationService) *CollectionAPIHandler {
	return &CollectionAPIHandler{
		Common:   NewCommon(config),
		register: register,
	}
}
