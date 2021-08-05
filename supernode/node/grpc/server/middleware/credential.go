package middleware

import (
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"google.golang.org/grpc"
)

// AltsCredential returns a ServerOption that sets the Creds for the server.
func AltsCredential(secClient alts.SecClient, secInfo *alts.SecInfo) grpc.ServerOption {
	altsTCServer := credentials.NewServerCreds(secClient, secInfo)
	return grpc.Creds(altsTCServer)
}
