package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/node"
)

// BlockMessage is block message
type BlockMessage struct {
	Block int
}

const (
	//P2PConnectionCloseError represents the p2p connection close error
	P2PConnectionCloseError = "grpc: the client connection is closing"
	//P2PServiceNotRunningError represents the p2p service not running error
	P2PServiceNotRunningError = "p2p service is not running"
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
	return strings.Contains(errString, P2PConnectionCloseError)
}

// IsP2PServiceNotRunningError checks if the passed errString is of p2p not running
func IsP2PServiceNotRunningError(errString string) bool {
	return strings.Contains(errString, P2PServiceNotRunningError)
}
