package transport

import "google.golang.org/grpc/credentials"

type ed448AuthInfo struct {
}

func newAuth() credentials.AuthInfo {
	return &ed448AuthInfo{}
}

func (authInfo *ed448AuthInfo) AuthType() string {
	return "ed448"
}
