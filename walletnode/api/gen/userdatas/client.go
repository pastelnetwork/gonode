// Code generated by goa v3.4.3, DO NOT EDIT.
//
// userdatas client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package userdatas

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "userdatas" service client.
type Client struct {
	CreateUserdataEndpoint        goa.Endpoint
	UpdateUserdataEndpoint        goa.Endpoint
	UserdataGetEndpoint           goa.Endpoint
	SetUserFollowRelationEndpoint goa.Endpoint
	GetFollowersEndpoint          goa.Endpoint
	GetFolloweesEndpoint          goa.Endpoint
	GetFriendsEndpoint            goa.Endpoint
	SetUserLikeArtEndpoint        goa.Endpoint
	GetUsersLikeArtEndpoint       goa.Endpoint
}

// NewClient initializes a "userdatas" service client given the endpoints.
func NewClient(createUserdata, updateUserdata, userdataGet, setUserFollowRelation, getFollowers, getFollowees, getFriends, setUserLikeArt, getUsersLikeArt goa.Endpoint) *Client {
	return &Client{
		CreateUserdataEndpoint:        createUserdata,
		UpdateUserdataEndpoint:        updateUserdata,
		UserdataGetEndpoint:           userdataGet,
		SetUserFollowRelationEndpoint: setUserFollowRelation,
		GetFollowersEndpoint:          getFollowers,
		GetFolloweesEndpoint:          getFollowees,
		GetFriendsEndpoint:            getFriends,
		SetUserLikeArtEndpoint:        setUserLikeArt,
		GetUsersLikeArtEndpoint:       getUsersLikeArt,
	}
}

// CreateUserdata calls the "createUserdata" endpoint of the "userdatas"
// service.
func (c *Client) CreateUserdata(ctx context.Context, p *CreateUserdataPayload) (res *UserdataProcessResult, err error) {
	var ires interface{}
	ires, err = c.CreateUserdataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserdataProcessResult), nil
}

// UpdateUserdata calls the "updateUserdata" endpoint of the "userdatas"
// service.
func (c *Client) UpdateUserdata(ctx context.Context, p *UpdateUserdataPayload) (res *UserdataProcessResult, err error) {
	var ires interface{}
	ires, err = c.UpdateUserdataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserdataProcessResult), nil
}

// UserdataGet calls the "userdataGet" endpoint of the "userdatas" service.
func (c *Client) UserdataGet(ctx context.Context, p *UserdataGetPayload) (res *UserSpecifiedData, err error) {
	var ires interface{}
	ires, err = c.UserdataGetEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserSpecifiedData), nil
}

// SetUserFollowRelation calls the "setUserFollowRelation" endpoint of the
// "userdatas" service.
func (c *Client) SetUserFollowRelation(ctx context.Context, p *SetUserFollowRelationPayload) (res *SetUserFollowRelationResult, err error) {
	var ires interface{}
	ires, err = c.SetUserFollowRelationEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SetUserFollowRelationResult), nil
}

// GetFollowers calls the "getFollowers" endpoint of the "userdatas" service.
func (c *Client) GetFollowers(ctx context.Context, p *GetFollowersPayload) (res *GetFollowersResult, err error) {
	var ires interface{}
	ires, err = c.GetFollowersEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*GetFollowersResult), nil
}

// GetFollowees calls the "getFollowees" endpoint of the "userdatas" service.
func (c *Client) GetFollowees(ctx context.Context, p *GetFolloweesPayload) (res *GetFolloweesResult, err error) {
	var ires interface{}
	ires, err = c.GetFolloweesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*GetFolloweesResult), nil
}

// GetFriends calls the "getFriends" endpoint of the "userdatas" service.
func (c *Client) GetFriends(ctx context.Context, p *GetFriendsPayload) (res *GetFriendsResult, err error) {
	var ires interface{}
	ires, err = c.GetFriendsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*GetFriendsResult), nil
}

// SetUserLikeArt calls the "setUserLikeArt" endpoint of the "userdatas"
// service.
func (c *Client) SetUserLikeArt(ctx context.Context, p *SetUserLikeArtPayload) (res *SetUserLikeArtResult, err error) {
	var ires interface{}
	ires, err = c.SetUserLikeArtEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SetUserLikeArtResult), nil
}

// GetUsersLikeArt calls the "getUsersLikeArt" endpoint of the "userdatas"
// service.
func (c *Client) GetUsersLikeArt(ctx context.Context, p *GetUsersLikeArtPayload) (res *GetUsersLikeArtResult, err error) {
	var ires interface{}
	ires, err = c.GetUsersLikeArtEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*GetUsersLikeArtResult), nil
}
