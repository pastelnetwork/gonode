package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/node"
)

const (
	p2pConnectionCloseError = "grpc: the client connection is closing"
)

// CreateNewP2PConnection creates a new p2p connection using sn host and port
func CreateNewP2PConnection(host string, port int, sn node.SNClientInterface) (node.HermesP2PInterface, error) {
	snAddr := fmt.Sprintf("%s:%d", host, port)
	log.WithField("sn-addr", snAddr).Info("connecting with SN-Service")

	conn, err := sn.Connect(context.Background(), snAddr)
	if err != nil {
		return nil, errors.Errorf("unable to connect with SN service: %w", err)
	}
	log.Info("connection established with SN-Service")

	p2p := conn.HermesP2P()
	log.Info("connection established with SN-P2P-Service")

	return p2p, nil
}

// IsP2PConnectionCloseError checks if the passed errString is of gRPC connection closing
func IsP2PConnectionCloseError(errString string) bool {
	return strings.Contains(errString, p2pConnectionCloseError)
}
