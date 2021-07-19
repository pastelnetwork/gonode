package artworkdownload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	MinRQIDsFileLine = 4
	RQSymbolsDirName = "symbols"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	ResampledArtwork *artwork.File
	Artwork          *artwork.File

	RQSymbolsDir string
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
func (task *Task) Download(ctx context.Context, txid, timestamp, signature, ttxid string) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(StatusTaskStarted); err != nil {
		return nil, err
	}

	// Create directory for storing symbols of Artwork
	task.RQSymbolsDir = path.Join(task.Service.config.RqFilesDir, task.ID(), RQSymbolsDirName)
	err = os.MkdirAll(task.RQSymbolsDir, os.ModePerm)
	if err != nil {
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

		log.WithContext(ctx).Debug(fmt.Sprintf("Art ticket: %s", string(artRegTicket.RegTicketData.ArtTicket)))

		// Decode Art Ticket
		err = task.decodeRegTicket(&artRegTicket)
		if err != nil {
			return err
		}

		pastelID := artRegTicket.RegTicketData.ArtTicketData.Author

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

		// Get symbol identifiers files from Kademlia by using rq_ids - from Art Registration ticket
		// Get the list of "symbols/chunks" from Kademlia by using symbol identifiers from file
		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_oti)
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		for _, id := range artRegTicket.RegTicketData.ArtTicketData.AppTicketData.RQIDs {
			var rqIDsData []byte
			rqIDsData, err = task.p2pClient.Retrieve(ctx, id)
			if err != nil {
				return err
			}

			var rqIDs []string
			rqIDs, err = task.getRQSymbolIDs(rqIDsData)
			if err != nil {
				continue
			}

			symbols := make(map[string][]byte)
			for _, id := range rqIDs {
				var symbol []byte
				symbol, err = task.p2pClient.Retrieve(ctx, id)
				if err != nil {
					break
				}

				// Validate that the hash of each "symbol/chunk" matches its id
				h := sha3.Sum256(symbol)
				if string(h[:]) != id {
					err = errors.New("Symbol id mismatched")
					break
				}
				symbols[id] = symbol
			}
			if len(symbols) != len(rqIDs) {
				err = errors.New("Could not retrieve all symbols from Kademlia")
				continue
			}

			// Write all symbols to files and store in a directory and pass to rqservice for decoding
			err = task.writeSymbolsToFiles(symbols)
			if err != nil {
				task.removeAllSymbolFiles()
				continue
			}

			// Restore artwork
			var restoredFilePath string
			// FIXME: Restore after merge raptorQ interface
			// restoredFilePath, err = task.Service.raptorQClient.Decode(task.RQSymbolsDir)
			// if err != nil {
			// 	continue
			// }

			var restoredFile []byte
			restoredFile, err = ioutil.ReadFile(restoredFilePath)
			if err != nil {
				continue
			}
			// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
			fileHash := sha3.Sum256(restoredFile)
			if string(fileHash[:]) != string(artRegTicket.RegTicketData.ArtTicketData.AppTicketData.DataHash) {
				err = errors.New("File mismatched")
				continue
			}

			file = restoredFile
		}

		if len(file) == 0 {
			return err
		}

		return nil
	})

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
		return errors.New(fmt.Sprintf("Could not parse art ticket. %w", err))
	}
	log.Debug(fmt.Sprintf("App ticket: %s", string(artRegTicket.RegTicketData.ArtTicketData.AppTicket)))

	appTicketData := artRegTicket.RegTicketData.ArtTicketData.AppTicket
	// n := base64.StdEncoding.DecodedLen(len(artRegTicket.RegTicketData.ArtTicketData.AppTicket))
	// appTicketData := make([]byte, n)
	// _, err = base64.StdEncoding.Decode(appTicketData, artRegTicket.RegTicketData.ArtTicketData.AppTicket)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Could not decode app ticket. %w", err))
	// }

	err = json.Unmarshal(appTicketData, &artRegTicket.RegTicketData.ArtTicketData.AppTicketData)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse app ticket. %w", err))
	}

	return nil
}

func (task Task) getRQSymbolIDs(rqIDsData []byte) (rqIDs []string, err error) {
	lines := strings.Split(string(rqIDsData), "\n")
	// First line is RANDOM-GUID
	// Second line is BLOCK_HASH
	// Third line is PASTELID
	// All the rest are rq IDs

	if len(lines) < MinRQIDsFileLine {
		err = errors.New(fmt.Sprintf("Invalid RQ IDs file: %s", string(rqIDsData)))
		return
	}

	rqIDs = lines[3:]

	return
}

func (task *Task) writeSymbolsToFiles(symbols map[string][]byte) error {
	for id, symbol := range symbols {
		filePath := path.Join(task.RQSymbolsDir, id)
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(symbol)
		if err != nil {
			return err
		}
	}
	return nil
}

func (task *Task) removeAllSymbolFiles() error {
	d, err := os.Open(task.RQSymbolsDir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(1)
	if err != nil && err != io.EOF {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(path.Join(task.RQSymbolsDir, name))
		if err != nil {
			return err
		}
	}
	return nil
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
