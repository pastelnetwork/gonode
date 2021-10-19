// Code generated by goa v3.5.2, DO NOT EDIT.
//
// external_dupe_detection_api client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package externaldupedetectionapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "external_dupe_detection_api" service client.
type Client struct {
	InitiateSubmissionEndpoint goa.Endpoint
}

// NewClient initializes a "external_dupe_detection_api" service client given
// the endpoints.
func NewClient(initiateSubmission goa.Endpoint) *Client {
	return &Client{
		InitiateSubmissionEndpoint: initiateSubmission,
	}
}

// InitiateSubmission calls the "initiate_submission" endpoint of the
// "external_dupe_detection_api" service.
func (c *Client) InitiateSubmission(ctx context.Context, p *ExternalDupeDetetionAPIInitiateSubmission) (res *Externaldupedetetionapiinitiatesubmissionresult, err error) {
	var ires interface{}
	ires, err = c.InitiateSubmissionEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Externaldupedetetionapiinitiatesubmissionresult), nil
}