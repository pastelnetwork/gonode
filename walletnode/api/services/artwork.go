package services

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

const (
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// NftApiHandler represents services for artworks endpoints.
type NftApiHandler struct {
	*Common
	register *artworkregister.NftRegisterService
	search   *artworksearch.NftSearchService
	download *artworkdownload.NftDownloadService
	db       storage.KeyValue
	imageTTL time.Duration
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *NftApiHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// Mount configures the mux to serve the artworks endpoints.
func (service *NftApiHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := artworks.NewEndpoints(service)
	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		NftRegUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// --- Register NFT ---
// UploadImage uploads an image and return unique image id.
func (service *NftApiHandler) UploadImage(ctx context.Context, p *artworks.UploadImagePayload) (res *artworks.Image, err error) {
	if p.Filename == nil {
		return nil, artworks.MakeBadRequest(errors.New("file not specified"))
	}

	id, expiry, err := service.register.StoreFile(ctx, p.Filename)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	res = &artworks.Image{
		ImageID:   id,
		ExpiresIn: expiry,
	}
	return res, nil
}

// Register runs registers process for the new NFT.
func (service *NftApiHandler) Register(_ context.Context, p *artworks.RegisterPayload) (res *artworks.RegisterResult, err error) {
	taskID, err := service.register.AddTask(p)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	res = &artworks.RegisterResult{
		TaskID: taskID,
	}
	return res, nil
}

// RegisterTaskState streams the state of the registration process.
func (service *NftApiHandler) RegisterTaskState(ctx context.Context, p *artworks.RegisterTaskStatePayload, stream artworks.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		return artworks.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &artworks.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return artworks.MakeInternalServerError(err)
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// RegisterTask returns a single task.
func (service *NftApiHandler) RegisterTask(_ context.Context, p *artworks.RegisterTaskPayload) (res *artworks.Task, err error) {
	task := service.register.GetTask(p.TaskID)
	if task == nil {
		return nil, artworks.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	res = &artworks.Task{
		ID:     p.TaskID,
		Status: task.Status().String(),
		Ticket: artworkregister.ToNftRegisterTicket(task.Request),
		States: toArtworkStates(task.StatusHistory()),
	}
	return res, nil
}

// RegisterTasks returns list of all tasks.
func (service *NftApiHandler) RegisterTasks(_ context.Context) (res artworks.TaskCollection, err error) {
	tasks := service.register.Tasks()
	for _, task := range tasks {
		res = append(res, &artworks.Task{
			ID:     task.ID(),
			Status: task.Status().String(),
			Ticket: artworkregister.ToNftRegisterTicket(task.Request),
		})
	}
	return res, nil
}

// --- Download registered NFT ---
func (service *NftApiHandler) Download(ctx context.Context, p *artworks.ArtworkDownloadPayload) (res *artworks.DownloadResult, err error) {
	log.WithContext(ctx).Info("Start downloading")
	defer log.WithContext(ctx).Info("Finished downloading")
	taskID := service.download.AddTask(p)
	task := service.download.GetTask(taskID)
	defer task.Cancel()

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil, artworks.MakeBadRequest(errors.Errorf("context done: %w", ctx.Err()))
		case status := <-sub():
			if status.IsFailure() {
				if strings.Contains(utils.SafeErrStr(task.Error()), "ticket ownership") {
					return nil, artworks.MakeBadRequest(errors.New("failed to verify ownership"))
				}

				return nil, artworks.MakeInternalServerError(task.Error())
			}

			if status.IsFinal() {
				log.WithContext(ctx).WithField("size", fmt.Sprintf("%d bytes", len(task.File))).Info("NFT downloaded")
				res = &artworks.DownloadResult{
					File: task.File,
				}

				return res, nil
			}
		}
	}
}

// --- Search registered NFTs and return details ---
// ArtSearch searches for NFT & streams the result based on filters
func (service *NftApiHandler) ArtSearch(ctx context.Context, p *artworks.ArtSearchPayload, stream artworks.ArtSearchServerStream) error {
	defer stream.Close()

	taskID := service.search.AddTask(p)
	task := service.search.GetTask(taskID)

	resultChan := task.SubscribeSearchResult()
	for {
		select {
		case <-ctx.Done():
			return nil
		case search, ok := <-resultChan:
			if !ok {
				if task.Status().IsFailure() {
					return artworks.MakeInternalServerError(task.Error())
				}

				return nil
			}

			res := toArtSearchResult(search)
			if err := stream.Send(res); err != nil {
				return artworks.MakeInternalServerError(err)
			}
		}
	}
}

// ArtworkGet returns NFT detail
func (service *NftApiHandler) ArtworkGet(ctx context.Context, p *artworks.ArtworkGetPayload) (res *artworks.ArtworkDetail, err error) {
	ticket, err := service.search.RegTicket(ctx, p.Txid)
	if err != nil {
		return nil, artworks.MakeBadRequest(err)
	}

	res = toArtworkDetail(ticket)
	data, err := service.search.GetThumbnail(ctx, ticket, p.UserPastelID, p.UserPassphrase)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	res.PreviewThumbnail = data

	return res, nil
}

// NewNftApiHandler returns the artworks NftApiHandler implementation.
func NewNftApiHandler(register *artworkregister.NftRegisterService, search *artworksearch.NftSearchService, download *artworkdownload.NftDownloadService) *NftApiHandler {
	return &NftApiHandler{
		Common:   NewCommon(),
		register: register,
		search:   search,
		download: download,
		db:       memory.NewKeyValue(),
		imageTTL: defaultImageTTL,
	}
}
