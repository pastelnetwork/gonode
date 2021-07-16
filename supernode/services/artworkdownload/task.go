package artworkdownload

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	ResampledArtwork *artwork.File
	Artwork          *artwork.File

	// acceptedMu sync.Mutex
	// accpeted   Nodes

	// connectedTo *Node
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = task.context(ctx)
	defer log.WithContext(ctx).Debug("Task canceled")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	return task.RunAction(ctx)
}

// Download downloads image and return the image.
func (task *Task) Download(_ context.Context, txid, timestamp, signature, ttxid string) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(StatusTaskStarted); err != nil {
		return nil, err
	}

	var file []byte
	var artRegTicket pastel.RegTicket

	<-task.NewAction(func(ctx context.Context) error {
		// Validate timestamp is not older than 10 minutes
		now := time.Now()
		lastTenMinutes := now.Add(time.Duration(-10) * time.Minute)
		requestTime, _ := time.Parse(time.RFC3339, timestamp)
		if lastTenMinutes.After(requestTime) {
			return errors.New("Request time is older than 10 minutes")
		}

		// Get Art Registration ticket by txid
		artRegTicket, err = task.pastelClient.RegTicket(ctx, txid)
		if err != nil {
			return err
		}

		pastelID := artRegTicket.RegTicketData.ArtTicketData.AppTicketData.AuthorPastelID

		if len(ttxid) > 0 {
			// Get list of non sold Trade ticket owened by the owner of the PastelID from request
			// by calling command `tickets list trade available`
			var tradeTickets []pastel.TradeTicket
			tradeTickets, err = task.pastelClient.ListAvailableTradeTickets(ctx)
			if err != nil {
				return err
			}

			// Validate that Trade ticket with ttxid is in the list
			if len(tradeTickets) == 0 {
				return errors.New("Could not get any available trade tickets")
			}
			isTXIDValid := false
			for _, t := range tradeTickets {
				if t.ArtTXID == ttxid {
					isTXIDValid = true
					pastelID = t.PastelID
					break
				}
			}
			if !isTXIDValid {
				return errors.New(fmt.Sprintf("Not found trade ticket of transaction %s", ttxid))
			}
		}
		// Validate timestamp signature with PastelID from Trade ticket
		// by calling command `pastelid verify timestamp-string signature PastelID`
		var isValid bool
		isValid, err = task.pastelClient.Verify(ctx, []byte(timestamp), signature, pastelID)
		if err != nil {
			return err
		}
		if !isValid {
			return errors.New("Invalid timestamp")
		}

		// FIXME: Get the list of “symbols/chunks” - rq_ids - from Art Registration ticket and request them from Kademlia
		for _, id := range artRegTicket.RegTicketData.ArtTicketData.AppTicketData.RQIDs {
			var data []byte
			data, err = task.p2pClient.Retrieve(ctx, id)
			if err != nil {
				return err
			}
			// Validate that the hash of each "symbol/chunk" matches its id
			h := sha3.Sum256(data)
			if string(h[:]) != id {
				err = errors.New("Symbol id mismatched")
				return err
			}
		}

		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_coti and rq_ssoti)

		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		fileHash := sha3.Sum256(file)
		if string(fileHash[:]) != string(artRegTicket.RegTicketData.ArtTicketData.AppTicketData.DataHash) {
			err = errors.New("File mismatched")
			return err
		}
		return nil
	})

	return file, err
}

func (task *Task) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
	}
}
