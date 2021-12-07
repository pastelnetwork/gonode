package artworkdownload

import (
	"bytes"
	"context"
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
		log.WithContext(ctx).WithField("status", status.String()).Debug("Status updated")
	})

	return task.RunAction(ctx)
}

// DownloadThumbnail downloads thumbnail of given hash.
func (task *Task) DownloadThumbnail(ctx context.Context, key []byte) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}

	var file []byte
	<-task.NewAction(func(ctx context.Context) error {
		base58Key := base58.Encode(key)
		file, err = task.p2pClient.Retrieve(ctx, base58Key)
		if err != nil {
			err = errors.Errorf("fetch p2p key : %s, error: %w", string(base58Key), err)
			task.UpdateStatus(StatusKeyNotFound)
		}
		return nil
	})

	return file, err
}

// Download downloads image and return the image.
func (task *Task) Download(ctx context.Context, txid, timestamp, signature, ttxid string) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}

	var file []byte
	var nftRegTicket pastel.RegTicket

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
		nftRegTicket, err = task.pastelClient.RegTicket(ctx, txid)
		if err != nil {
			err = errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			task.UpdateStatus(StatusArtRegGettingFailed)
			return nil
		}

		log.WithContext(ctx).Debugf("Art ticket: %s", string(nftRegTicket.RegTicketData.NFTTicket))

		// Decode Art Ticket
		err = task.decodeRegTicket(&nftRegTicket)
		if err != nil {
			task.UpdateStatus(StatusArtRegDecodingFailed)
			return nil
		}

		pastelID := nftRegTicket.RegTicketData.NFTTicketData.Author

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
				err = errors.New("could not find any available trade tickets")
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
				log.WithContext(ctx).WithField("ttxid", ttxid).Errorf("could not find trade ticket of transaction")
				err = errors.Errorf("could not find trade ticket of transaction %s", ttxid)
				task.UpdateStatus(StatusTradeTicketMismatched)
				return nil
			}
		}
		// Validate timestamp signature with PastelID from Trade ticket
		// by calling command `pastelid verify timestamp-string signature PastelID`
		var isValid bool
		isValid, err = task.pastelClient.Verify(ctx, []byte(timestamp), signature, pastelID, "ed448")
		if err != nil {
			err = errors.Errorf("timestamp signature verify: %w", err)
			task.UpdateStatus(StatusTimestampVerificationFailed)
			return nil
		}

		if !isValid {
			err = errors.New("invalid signature timestamp")
			task.UpdateStatus(StatusTimestampInvalid)
			return nil
		}

		// Get symbol identifiers files from Kademlia by using rq_ids - from Art Registration ticket
		// Get the list of "symbols/chunks" from Kademlia by using symbol identifiers from file
		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_oti)
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		file, err = task.restoreFile(ctx, &nftRegTicket)
		if err != nil {
			err = errors.Errorf("restore file: %w", err)
			return nil
		}

		if len(file) == 0 {
			err = errors.New("nil restored file")
			task.UpdateStatus(StatusFileEmpty)
		}

		return nil
	})

	return file, err
}

func (task *Task) restoreFile(ctx context.Context, nftRegTicket *pastel.RegTicket) ([]byte, error) {
	var file []byte
	var lastErr error
	var err error

	if len(nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs) == 0 {
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

	for _, id := range nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs {
		var rqIDsData []byte
		rqIDsData, err = task.p2pClient.Retrieve(ctx, id)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Retrieve compressed symbol IDs file from P2P failed")
			lastErr = errors.Errorf("retrieve compressed symbol IDs file: %w", err)
			task.UpdateStatus(StatusSymbolFileNotFound)
			continue
		}

		fileContent, err := zstd.Decompress(nil, rqIDsData)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Decompress compressed symbol IDs file failed")
			lastErr = errors.Errorf("decompress symbol IDs file: %w", err)
			task.UpdateStatus(StatusSymbolFileInvalid)
			continue
		}

		log.WithContext(ctx).WithField("Content", string(fileContent)).Debugf("symbol IDs file")
		var rqIDs []string
		rqIDs, err = task.getRQSymbolIDs(fileContent)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Parse symbol IDs failed")
			lastErr = errors.Errorf("parse symbol IDs: %w", err)
			task.UpdateStatus(StatusSymbolFileInvalid)
			continue
		}

		log.WithContext(ctx).Debugf("Symbol IDs: %v", rqIDs)

		symbols := make(map[string][]byte)
		for _, id := range rqIDs {
			var symbol []byte
			symbol, err = task.p2pClient.Retrieve(ctx, id)
			if err != nil {
				log.WithContext(ctx).WithField("SymbolID", id).Warn("Could not retrieve symbol")
				task.UpdateStatus(StatusSymbolNotFound)
				break
			}

			// Validate that the hash of each "symbol/chunk" matches its id
			h := sha3.Sum256(symbol)
			storedID := base58.Encode(h[:])
			if storedID != id {
				log.WithContext(ctx).Warnf("Symbol ID mismatched, expect %v, got %v", id, storedID)
				task.UpdateStatus(StatusSymbolMismatched)
				break
			}
			symbols[id] = symbol
		}
		if len(symbols) != len(rqIDs) {
			log.WithContext(ctx).WithField("SymbolIDsFileId", id).Warn("Could not retrieve all symbols")
			lastErr = errors.New("could not retrieve all symbols from Kademlia")
			task.UpdateStatus(StatusSymbolsNotEnough)
			continue
		}

		// Restore artwork
		var decodeInfo *rqnode.Decode
		encodeInfo := rqnode.Encode{
			Symbols: symbols,
			EncoderParam: rqnode.EncoderParameters{
				Oti: nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQOti,
			},
		}

		decodeInfo, err = rqService.Decode(ctx, &encodeInfo)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Restore file with rqserivce")
			lastErr = errors.Errorf("restore file with rqserivce: %w", err)
			task.UpdateStatus(StatusFileDecodingFailed)
			continue
		}
		task.UpdateStatus(StatusFileDecoded)

		// log.WithContext(ctx).Debugf("Restored file path: %s", decodeInfo.Path)
		// log.WithContext(ctx).Debugf("Restored file: %s", string(restoredFile))
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		fileHash := sha3.Sum256(decodeInfo.File)

		if !bytes.Equal(fileHash[:], nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash) {
			log.WithContext(ctx).WithField("SymbolIDsFileId", id).Warn("hash file mismatched")
			lastErr = errors.New("hash file mismatched")
			task.UpdateStatus(StatusFileMismatched)
			continue
		}

		return decodeInfo.File, nil
	}

	return file, lastErr
}

func (task *Task) decodeRegTicket(nftRegTicket *pastel.RegTicket) error {
	articketData, err := pastel.DecodeNFTTicket(nftRegTicket.RegTicketData.NFTTicket)
	if err != nil {
		return errors.Errorf("convert NFT ticket: %w", err)
	}
	nftRegTicket.RegTicketData.NFTTicketData = *articketData

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
