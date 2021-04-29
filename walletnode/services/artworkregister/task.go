package artworkregister

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/walletnode/services/artworkregister/state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var taskID uint32

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID     int
	State  *state.State
	Ticket *Ticket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	log.Debugf("Start task %v", task.ID)
	defer log.Debugf("End task %v", task.ID)

	if err := task.run(ctx); err != nil {
		if err, ok := err.(*TaskError); ok {
			task.State.Update(state.NewMessage(state.StatusTaskRejected))
			log.WithField("error", err).Debugf("Task %d is rejected", task.ID)
			return nil
		}
		return err
	}

	task.State.Update(state.NewMessage(state.StatusTaskCompleted))
	log.Debugf("Task %d is completed", task.ID)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	superNodes, err := task.findSuperNodes(ctx)
	if err != nil {
		return err
	}

	for _, superNode := range superNodes {
		superNode.Address = "127.0.0.1:4444"
		if err := task.connect(ctx, superNode); err != nil {
			return err
		}
	}
	// if len(superNodes) < task.config.NumberSuperNodes {
	// 	task.State.Update(state.NewMessage(state.StatusErrorTooLowFee))
	// 	return NewTaskError(errors.Errorf("not found enough available SuperNodes with acceptable storage fee", task.config.NumberSuperNodes))
	// }

	return nil
}

func (task *Task) connect(ctx context.Context, superNode *SuperNode) error {
	log.Debugf("connect to %s", superNode.Address)

	ctxConnect, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	conn, err := grpc.DialContext(ctxConnect, superNode.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return NewTaskError(errors.Errorf("fail to dial: %v", err).WithField("address", superNode.Address))
	}
	defer conn.Close()

	log.Debugf("connected to %s", superNode.Address)

	client := pb.NewWalletNodeClient(conn)
	stream, err := client.RegisterArtowrk(ctx)
	if err != nil {
		return errors.New(err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
			req := &pb.WalletNodeRegisterArtworkRequest{
				TestOneof: &pb.WalletNodeRegisterArtworkRequest_AcceptSecondaryNodes{
					AcceptSecondaryNodes: &pb.WalletNodeRegisterArtworkRequest_AcceptSecondaryNodesRequest{
						ConnID: "12345",
					},
				},
			}

			if err := stream.Send(req); err != nil {
				if status.Code(err) == codes.Canceled {
					log.Debug("stream closed (context cancelled)")
					return nil
				}

				return NewTaskError(errors.Errorf("could not perform connect command: %v", err).WithField("address", superNode.Address))
			}
		}
	}
}

func (task *Task) findSuperNodes(ctx context.Context) (SuperNodes, error) {
	var superNodes SuperNodes

	mns, err := task.pastel.TopMasterNodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Ticket.MaximumFee {
			continue
		}
		superNodes = append(superNodes, &SuperNode{
			Address:   mn.ExtAddress,
			PastelKey: mn.ExtKey,
			Fee:       mn.Fee,
		})
	}

	return superNodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	return &Task{
		Service: service,
		ID:      int(atomic.AddUint32(&taskID, 1)),
		Ticket:  Ticket,
		State:   state.New(state.NewMessage(state.StatusTaskStarted)),
	}
}
