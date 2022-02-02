package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
)

func contextWithMDSessID(ctx context.Context, sessID string) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func contextWithLogPrefix(ctx context.Context, connId string) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, connId))
}
