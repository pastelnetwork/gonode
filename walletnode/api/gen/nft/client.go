// Code generated by goa v3.6.2, DO NOT EDIT.
//
// nft client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package nft

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "nft" service client.
type Client struct {
	RegisterEndpoint          goa.Endpoint
	RegisterTaskStateEndpoint goa.Endpoint
	GetTaskHistoryEndpoint    goa.Endpoint
	RegisterTaskEndpoint      goa.Endpoint
	RegisterTasksEndpoint     goa.Endpoint
	UploadImageEndpoint       goa.Endpoint
	NftSearchEndpoint         goa.Endpoint
	NftGetEndpoint            goa.Endpoint
	DownloadEndpoint          goa.Endpoint
}

// NewClient initializes a "nft" service client given the endpoints.
func NewClient(register, registerTaskState, getTaskHistory, registerTask, registerTasks, uploadImage, nftSearch, nftGet, download goa.Endpoint) *Client {
	return &Client{
		RegisterEndpoint:          register,
		RegisterTaskStateEndpoint: registerTaskState,
		GetTaskHistoryEndpoint:    getTaskHistory,
		RegisterTaskEndpoint:      registerTask,
		RegisterTasksEndpoint:     registerTasks,
		UploadImageEndpoint:       uploadImage,
		NftSearchEndpoint:         nftSearch,
		NftGetEndpoint:            nftGet,
		DownloadEndpoint:          download,
	}
}

// Register calls the "register" endpoint of the "nft" service.
func (c *Client) Register(ctx context.Context, p *RegisterPayload) (res *RegisterResult, err error) {
	var ires interface{}
	ires, err = c.RegisterEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RegisterResult), nil
}

// RegisterTaskState calls the "registerTaskState" endpoint of the "nft"
// service.
func (c *Client) RegisterTaskState(ctx context.Context, p *RegisterTaskStatePayload) (res RegisterTaskStateClientStream, err error) {
	var ires interface{}
	ires, err = c.RegisterTaskStateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(RegisterTaskStateClientStream), nil
}

// GetTaskHistory calls the "getTaskHistory" endpoint of the "nft" service.
func (c *Client) GetTaskHistory(ctx context.Context, p *GetTaskHistoryPayload) (res []*TaskHistory, err error) {
	var ires interface{}
	ires, err = c.GetTaskHistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.([]*TaskHistory), nil
}

// RegisterTask calls the "registerTask" endpoint of the "nft" service.
func (c *Client) RegisterTask(ctx context.Context, p *RegisterTaskPayload) (res *Task, err error) {
	var ires interface{}
	ires, err = c.RegisterTaskEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// RegisterTasks calls the "registerTasks" endpoint of the "nft" service.
func (c *Client) RegisterTasks(ctx context.Context) (res TaskCollection, err error) {
	var ires interface{}
	ires, err = c.RegisterTasksEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(TaskCollection), nil
}

// UploadImage calls the "uploadImage" endpoint of the "nft" service.
func (c *Client) UploadImage(ctx context.Context, p *UploadImagePayload) (res *ImageRes, err error) {
	var ires interface{}
	ires, err = c.UploadImageEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ImageRes), nil
}

// NftSearch calls the "nftSearch" endpoint of the "nft" service.
func (c *Client) NftSearch(ctx context.Context, p *NftSearchPayload) (res NftSearchClientStream, err error) {
	var ires interface{}
	ires, err = c.NftSearchEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(NftSearchClientStream), nil
}

// NftGet calls the "nftGet" endpoint of the "nft" service.
func (c *Client) NftGet(ctx context.Context, p *NftGetPayload) (res *NftDetail, err error) {
	var ires interface{}
	ires, err = c.NftGetEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*NftDetail), nil
}

// Download calls the "download" endpoint of the "nft" service.
func (c *Client) Download(ctx context.Context, p *DownloadPayload) (res *DownloadResult, err error) {
	var ires interface{}
	ires, err = c.DownloadEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadResult), nil
}
