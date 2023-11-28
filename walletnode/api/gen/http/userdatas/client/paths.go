// Code generated by goa v3.14.0, DO NOT EDIT.
//
// HTTP request path constructors for the userdatas service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"fmt"
)

// CreateUserdataUserdatasPath returns the URL path to the userdatas service createUserdata HTTP endpoint.
func CreateUserdataUserdatasPath() string {
	return "/userdatas/create"
}

// UpdateUserdataUserdatasPath returns the URL path to the userdatas service updateUserdata HTTP endpoint.
func UpdateUserdataUserdatasPath() string {
	return "/userdatas/update"
}

// GetUserdataUserdatasPath returns the URL path to the userdatas service getUserdata HTTP endpoint.
func GetUserdataUserdatasPath(pastelid string) string {
	return fmt.Sprintf("/userdatas/%v", pastelid)
}
