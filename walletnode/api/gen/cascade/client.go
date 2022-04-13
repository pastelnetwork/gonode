// Code generated by goa v3.6.2, DO NOT EDIT.
//
// cascade client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package cascade

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "cascade" service client.
type Client struct {
	UploadImageEndpoint       goa.Endpoint
	StartProcessingEndpoint   goa.Endpoint
	RegisterTaskStateEndpoint goa.Endpoint
	DownloadEndpoint          goa.Endpoint
}

// NewClient initializes a "cascade" service client given the endpoints.
func NewClient(uploadImage, startProcessing, registerTaskState, download goa.Endpoint) *Client {
	return &Client{
		UploadImageEndpoint:       uploadImage,
		StartProcessingEndpoint:   startProcessing,
		RegisterTaskStateEndpoint: registerTaskState,
		DownloadEndpoint:          download,
	}
}

// UploadImage calls the "uploadImage" endpoint of the "cascade" service.
func (c *Client) UploadImage(ctx context.Context, p *UploadImagePayload) (res *Image, err error) {
	var ires interface{}
	ires, err = c.UploadImageEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Image), nil
}

// StartProcessing calls the "startProcessing" endpoint of the "cascade"
// service.
func (c *Client) StartProcessing(ctx context.Context, p *StartProcessingPayload) (res *StartProcessingResult, err error) {
	var ires interface{}
	ires, err = c.StartProcessingEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StartProcessingResult), nil
}

// RegisterTaskState calls the "registerTaskState" endpoint of the "cascade"
// service.
func (c *Client) RegisterTaskState(ctx context.Context, p *RegisterTaskStatePayload) (res RegisterTaskStateClientStream, err error) {
	var ires interface{}
	ires, err = c.RegisterTaskStateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(RegisterTaskStateClientStream), nil
}

// Download calls the "download" endpoint of the "cascade" service.
func (c *Client) Download(ctx context.Context, p *DownloadPayload) (res *DownloadResult, err error) {
	var ires interface{}
	ires, err = c.DownloadEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadResult), nil
}
