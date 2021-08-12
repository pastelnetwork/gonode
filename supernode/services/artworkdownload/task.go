package artworkdownload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

const (
	minRQIDsFileLine = 4
	rqSymbolsDirName = "symbols"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	RQSymbolsDir string
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = task.context(ctx)
	defer task.UpdateStatus(StatusTaskCompleted)
	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	return task.RunAction(ctx)
}

// Download downloads image and return the image.
func (task *Task) Download(ctx context.Context, txid, timestamp, signature, ttxid string) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}

	var file []byte
	var artRegTicket pastel.RegTicket

	<-task.NewAction(func(ctx context.Context) error {
		// Validate timestamp is not older than 10 minutes
		now := time.Now()
		lastTenMinutes := now.Add(time.Duration(-10) * time.Minute)
		requestTime, _ := time.Parse(time.RFC3339, timestamp)
		if lastTenMinutes.After(requestTime) {
			err = errors.New("request time is older than 10 minutes")
			task.UpdateStatus(StatusRequestTooLate)
			return nil
		}

		// Get Art Registration ticket by txid
		artRegTicket, err = task.pastelClient.RegTicket(ctx, txid)
		if err != nil {
			err = errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			task.UpdateStatus(StatusArtRegGettingFailed)
			return nil
		}

		log.WithContext(ctx).Debugf("Art ticket: %s", string(artRegTicket.RegTicketData.ArtTicket))

		// Decode Art Ticket
		err = task.decodeRegTicket(&artRegTicket)
		if err != nil {
			task.UpdateStatus(StatusArtRegDecodingFailed)
			return nil
		}

		pastelID := base58.Encode(artRegTicket.RegTicketData.ArtTicketData.Author)

		if len(ttxid) > 0 {
			// Get list of non sold Trade ticket owened by the owner of the PastelID from request
			// by calling command `tickets list trade available`
			var tradeTickets []pastel.TradeTicket
			tradeTickets, err = task.pastelClient.ListAvailableTradeTickets(ctx)
			if err != nil {
				err = errors.Errorf("could not get available trade tickets: %w", err)
				task.UpdateStatus(StatusListTradeTicketsFailed)
				return nil
			}

			// Validate that Trade ticket with ttxid is in the list
			if len(tradeTickets) == 0 {
				err = errors.New("not found any available trade tickets")
				task.UpdateStatus(StatusTradeTicketsNotFound)
				return nil
			}
			isTXIDValid := false
			for _, t := range tradeTickets {
				if t.TXID == ttxid {
					isTXIDValid = true
					pastelID = t.Ticket.PastelID
					break
				}
			}
			if !isTXIDValid {
				log.WithContext(ctx).WithField("ttxid", ttxid).Errorf("not found trade ticket of transaction")
				err = errors.Errorf("not found trade ticket of transaction %s", ttxid)
				task.UpdateStatus(StatusTradeTicketMismatched)
				return nil
			}
		}
		// Validate timestamp signature with PastelID from Trade ticket
		// by calling command `pastelid verify timestamp-string signature PastelID`
		var isValid bool
		isValid, err = task.pastelClient.Verify(ctx, []byte(timestamp), signature, pastelID, "ed448")
		if err != nil {
			task.UpdateStatus(StatusTimestampVerificationFailed)
			return nil
		}
		if !isValid {
			err = errors.New("timestamp verification failed")
			task.UpdateStatus(StatusTimestampInvalid)
			return nil
		}

		// Get symbol identifiers files from Kademlia by using rq_ids - from Art Registration ticket
		// Get the list of "symbols/chunks" from Kademlia by using symbol identifiers from file
		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_oti)
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		file, err = task.restoreFile(ctx, &artRegTicket)
		if err != nil {
			err = errors.Errorf("failed to restore file: %w", err)
		}

		if len(file) == 0 {
			err = errors.Errorf("empty file")
			task.UpdateStatus(StatusFileEmpty)
		}

		return nil
	})

	return file, err
}

func (task *Task) restoreFile(ctx context.Context, artRegTicket *pastel.RegTicket) ([]byte, error) {
	var file []byte
	var err error

	if len(artRegTicket.RegTicketData.ArtTicketData.AppTicketData.RQIDs) == 0 {
		task.UpdateStatus(StatusArtRegTicketInvalid)
		return file, errors.Errorf("ticket has empty symbol identifier files")
	}

	var rqConnection rqnode.Connection
	rqConnection, err = task.Service.raptorQClient.Connect(ctx, task.Service.config.RaptorQServiceAddress)
	if err != nil {
		task.UpdateStatus(StatusRQServiceConnectionFailed)
		return file, errors.Errorf("could not connect to rqservice: %w", err)
	}
	defer rqConnection.Done()
	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)

	for _, id := range artRegTicket.RegTicketData.ArtTicketData.AppTicketData.RQIDs {
		var rqIDsData []byte
		rqIDsData, err = task.p2pClient.Retrieve(ctx, id)
		if err != nil {
			err = errors.Errorf("could not retrieve compressed symbol file from Kademlia: %w", err)
			task.UpdateStatus(StatusSymbolFileNotFound)
			continue
		}

		fileContent, err := zstd.Decompress(nil, rqIDsData)
		if err != nil {
			err = errors.Errorf("could not decompress symbol file: %w", err)
			task.UpdateStatus(StatusSymbolFileInvalid)
			continue
		}

		log.WithContext(ctx).Debugf("symbolFile: %s", string(fileContent))
		var rqIDs []string
		rqIDs, err = task.getRQSymbolIDs(fileContent)
		if err != nil {
			task.UpdateStatus(StatusSymbolFileInvalid)
			continue
		}
		log.WithContext(ctx).Debugf("rqIDs: %v", rqIDs)

		symbols := make(map[string][]byte)
		for _, id := range rqIDs {
			var symbol []byte
			symbol, err = task.p2pClient.Retrieve(ctx, id)
			if err != nil {
				log.WithContext(ctx).Debugf("Could not retrieve symbol of key: %s", id)
				task.UpdateStatus(StatusSymbolNotFound)
				break
			}

			// Validate that the hash of each "symbol/chunk" matches its id
			h := sha3.Sum256(symbol)
			storedID := base58.Encode(h[:])
			if storedID != id {
				err = errors.New("symbol id mismatched")
				log.WithContext(ctx).Debugf("Symbol id mismatched, expect %v, got %v", id, storedID)
				task.UpdateStatus(StatusSymbolMismatched)
				break
			}
			symbols[id] = symbol
		}
		if len(symbols) != len(rqIDs) {
			err = errors.New("could not retrieve all symbols from Kademlia")
			task.UpdateStatus(StatusSymbolsNotEnough)
			continue
		}

		// Restore artwork
		var decodeInfo *rqnode.Decode
		encodeInfo := rqnode.Encode{
			Symbols: symbols,
			EncoderParam: rqnode.EncoderParameters{
				Oti: artRegTicket.RegTicketData.ArtTicketData.AppTicketData.RQOti,
			},
		}

		decodeInfo, err = rqService.Decode(ctx, &encodeInfo)
		if err != nil {
			err = errors.Errorf("failed to restore file with rqserivce: %w", err)
			task.UpdateStatus(StatusFileDecodingFailed)
			continue
		}
		task.UpdateStatus(StatusFileDecoded)

		// log.WithContext(ctx).Debugf("Restored file path: %s", decodeInfo.Path)
		// log.WithContext(ctx).Debugf("Restored file: %s", string(restoredFile))
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		fileHash := sha3.Sum256(decodeInfo.File)

		if !bytes.Equal(fileHash[:], artRegTicket.RegTicketData.ArtTicketData.AppTicketData.DataHash) {
			err = errors.New("file mismatched")
			task.UpdateStatus(StatusFileMismatched)
			continue
		}

		file = decodeInfo.File
		break
	}

	return file, err
}

func (task *Task) decodeRegTicket(artRegTicket *pastel.RegTicket) error {
	// n := base64.StdEncoding.DecodedLen(len(artRegTicket.RegTicketData.ArtTicket))
	// artTicketData := make([]byte, n)
	// _, err := base64.StdEncoding.Decode(artTicketData, artRegTicket.RegTicketData.ArtTicket)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Could not decode art ticket. %w", err))
	// }

	err := json.Unmarshal(artRegTicket.RegTicketData.ArtTicket, &artRegTicket.RegTicketData.ArtTicketData)
	if err != nil {
		return errors.Errorf("Could not parse art ticket. %w", err)
	}
	// log.Debug(fmt.Sprintf("App ticket: %s", string(artRegTicket.RegTicketData.ArtTicketData.AppTicket)))

	appTicketData := artRegTicket.RegTicketData.ArtTicketData.AppTicket
	// n := base64.StdEncoding.DecodedLen(len(artRegTicket.RegTicketData.ArtTicketData.AppTicket))
	// appTicketData := make([]byte, n)
	// _, err = base64.StdEncoding.Decode(appTicketData, artRegTicket.RegTicketData.ArtTicketData.AppTicket)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Could not decode app ticket. %w", err))
	// }

	err = json.Unmarshal(appTicketData, &artRegTicket.RegTicketData.ArtTicketData.AppTicketData)
	if err != nil {
		return errors.Errorf("Could not parse app ticket. %w", err)
	}
	return nil
}

func (task Task) getRQSymbolIDs(rqIDsData []byte) (rqIDs []string, err error) {
	lines := strings.Split(string(rqIDsData), "\n")
	// First line is RANDOM-GUID
	// Second line is BLOCK_HASH
	// Third line is PASTELID
	// All the rest are rq IDs

	if len(lines) < minRQIDsFileLine {
		err = errors.Errorf("Invalid symbol identifiers file: %s", string(rqIDsData))
		return
	}

	l := len(lines)
	rqIDs = lines[3 : l-1]

	return
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
