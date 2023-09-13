// Code generated by goa v3.13.0, DO NOT EDIT.
//
// sense client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package sense

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "sense" service client.
type Client struct {
	UploadImageEndpoint       goa.Endpoint
	StartProcessingEndpoint   goa.Endpoint
	RegisterTaskStateEndpoint goa.Endpoint
	GetTaskHistoryEndpoint    goa.Endpoint
	DownloadEndpoint          goa.Endpoint
}

// NewClient initializes a "sense" service client given the endpoints.
func NewClient(uploadImage, startProcessing, registerTaskState, getTaskHistory, download goa.Endpoint) *Client {
	return &Client{
		UploadImageEndpoint:       uploadImage,
		StartProcessingEndpoint:   startProcessing,
		RegisterTaskStateEndpoint: registerTaskState,
		GetTaskHistoryEndpoint:    getTaskHistory,
		DownloadEndpoint:          download,
	}
}

// UploadImage calls the "uploadImage" endpoint of the "sense" service.
// UploadImage may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) UploadImage(ctx context.Context, p *UploadImagePayload) (res *Image, err error) {
	var ires any
	ires, err = c.UploadImageEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Image), nil
}

// StartProcessing calls the "startProcessing" endpoint of the "sense" service.
// StartProcessing may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) StartProcessing(ctx context.Context, p *StartProcessingPayload) (res *StartProcessingResult, err error) {
	var ires any
	ires, err = c.StartProcessingEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StartProcessingResult), nil
}

// RegisterTaskState calls the "registerTaskState" endpoint of the "sense"
// service.
// RegisterTaskState may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) RegisterTaskState(ctx context.Context, p *RegisterTaskStatePayload) (res RegisterTaskStateClientStream, err error) {
	var ires any
	ires, err = c.RegisterTaskStateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(RegisterTaskStateClientStream), nil
}

// GetTaskHistory calls the "getTaskHistory" endpoint of the "sense" service.
// GetTaskHistory may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) GetTaskHistory(ctx context.Context, p *GetTaskHistoryPayload) (res []*TaskHistory, err error) {
	var ires any
	ires, err = c.GetTaskHistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.([]*TaskHistory), nil
}

// Download calls the "download" endpoint of the "sense" service.
// Download may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) Download(ctx context.Context, p *DownloadPayload) (res *DownloadResult, err error) {
	var ires any
	ires, err = c.DownloadEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadResult), nil
}
