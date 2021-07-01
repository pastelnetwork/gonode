package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/raptorq/node"
)

type raptorQ struct {
	conn   *clientConn
	client pb.RaptorQClient
}

// Encode data, and return a list of identifier of symbols
func (service *raptorQ) Encode(ctx context.Context, data []byte) ([][]byte, error) {
	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)
	output := [][]byte{}
	req := pb.UploadDataRequest{
		Data: data,
	}

	stream, err := service.client.Encode(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("failed to send request: %w", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errors.Errorf("failed to receive : %w", err)
		}

		output = append(output, res.GetSymbol())
	}

	return output, nil
}

func (service *raptorQ) EncodeInfo(ctx context.Context, data []byte) (*node.EncodeInfo, error) {
	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	req := pb.UploadDataRequest{
		Data: data,
	}

	res, err := service.client.EncoderInfo(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("failed to send request: %w", err)
	}

	output := &node.EncodeInfo{
		SymbolIds: res.Name,
		EncoderParam: node.EncoderParameters{
			TransferLength:  res.EncoderParams.TransferLength,
			SymbolSize:      res.EncoderParams.SymbolSize,
			NumSourceBlocks: res.EncoderParams.NumSourceBlocks,
			NumSubBlocks:    res.EncoderParams.NumSubBlocks,
			SymbolAlignment: res.EncoderParams.SymbolAlignment,
		},
	}

	return output, nil
}

func (service *raptorQ) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRaptorQ(conn *clientConn) node.RaptorQ {
	return &raptorQ{
		conn:   conn,
		client: pb.NewRaptorQClient(conn),
	}
}
