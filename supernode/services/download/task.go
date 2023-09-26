package download

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/supernode/services/common"

	"github.com/btcsuite/btcutil/base58"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
)

// NftDownloadingTask is the task of registering new Nft.
type NftDownloadingTask struct {
	*common.SuperNodeTask
	RqClient rqnode.ClientInterface
	*NftDownloaderService

	RQSymbolsDir string
	ttype        string
}

// Run starts the task
func (task *NftDownloadingTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// DownloadThumbnail gets thumbnail file from ticket based on id and returns the thumbnail.
func (task *NftDownloadingTask) DownloadThumbnail(ctx context.Context, txid string, numnails int32) (map[int][]byte, error) {
	var err, err2 error
	if err = task.RequiredStatus(common.StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}
	regTicket, err := task.PastelClient.RegTicket(ctx, txid)
	if err != nil {
		log.WithContext(ctx).WithField("txid", txid).Error("Could not find regticket with txid")
		return nil, errors.Errorf("Bad txid: %s", err)
	}

	//decode the nft ticket
	nftTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", txid).Error("Could not decode ticket")
		return nil, errors.Errorf("Could not decode: %s", err)
	}
	appTicket := nftTicket.AppTicketData

	thumbnailHash := appTicket.Thumbnail1Hash
	var file1, file2 []byte
	<-task.NewAction(func(ctx context.Context) error {
		base58Key := base58.Encode(thumbnailHash)
		file1, err = task.P2PClient.Retrieve(ctx, base58Key)
		if err != nil {
			err = errors.Errorf("fetch p2p key : %s, base58 key: %s,error: %w", thumbnailHash, string(base58Key), err)
			task.UpdateStatus(common.StatusKeyNotFound)
		}
		return nil
	})
	if numnails > 1 {
		thumbnailHash = appTicket.Thumbnail2Hash
		<-task.NewAction(func(ctx context.Context) error {
			base58Key := base58.Encode(thumbnailHash)
			file2, err2 = task.P2PClient.Retrieve(ctx, base58Key)
			if err != nil {
				err2 = errors.Errorf("fetch p2p key2 : %s, base58 key2: %s,error: %w", thumbnailHash, string(base58Key), err)
				task.UpdateStatus(common.StatusKeyNotFound)
			}
			return nil
		})
	}
	if err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = errors.Errorf(err.Error() + " " + err2.Error())
		}
	}
	resMap := make(map[int][]byte)
	resMap[0] = file1
	if numnails > 1 {
		resMap[1] = file2
	}

	return resMap, err
}

func (task *NftDownloadingTask) downloadDDFPFile(ctx context.Context, ddFpID string) ([]byte, error) {
	file, err := task.P2PClient.Retrieve(ctx, ddFpID)
	if err != nil {
		return nil, errors.Errorf("DDAndFingerPrintDetails tried to get this file and failed: %w", err)
	}

	if len(file) == 0 {
		return nil, errors.Errorf("DDAndFingerPrintDetails empty file received: %w", err)
	}

	decompressedData, err := utils.Decompress(file)
	if err != nil {
		return nil, errors.Errorf("DDAndFingerPrintDetails failed to decompress this file: %w", err)
	}

	//base64 dataset doesn't contain periods, so we just find the first index of period and chop it and everything else
	firstIndexOfSeparator := bytes.IndexByte(decompressedData, pastel.SeparatorByte)
	if firstIndexOfSeparator < 1 {
		return nil, errors.Errorf("DDAndFingerPrintDetails got a bad separator index: %w", err)
	}
	dataToBase64Decode := decompressedData[:firstIndexOfSeparator]
	dataToJSONDecode, err := utils.B64Decode(dataToBase64Decode)
	if err != nil {
		return nil, errors.Errorf("DDAndFingerPrintDetails could not base64 decode: %w", err)
	}

	/*ddAndFingerprintsStruct := &pastel.DDAndFingerprints{}
	err = json.Unmarshal(dataToJSONDecode, ddAndFingerprintsStruct)
	if err != nil {
		return nil, errors.Errorf("DDAndFingerPrintDetails could not JSON unmarshal: %w", err)
	}*/ // Is this required?

	return dataToJSONDecode, nil
}

// DownloadDDAndFingerprints gets dd and fp file from ticket based on id and returns the file.
func (task *NftDownloadingTask) DownloadDDAndFingerprints(ctx context.Context, txid string) ([]byte, error) {
	log.WithContext(ctx).WithField("txid", txid).Info("Getting dd and fingerprints for txid")
	var err error
	if err = task.RequiredStatus(common.StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}

	info, err := task.getTicketInfo(ctx, txid)
	if err != nil {
		return nil, errors.Errorf("Could not get ticket info: %w", err)
	}

	DDAndFingerprintsIDs := info.ddAndFpIDs
	log.WithContext(ctx).WithField("ddandfpids", info.ddAndFpIDs).WithField("txid", txid).Info("Found dd and fp ids")

	//utility function for getting the DD and Fingerprint Details given
	//	a list of IDs (presumably from AppTicketData's DDAndFingerprintIDs), fingerprint IC, and fingerprint max
	//1) iterate over DDAndFingerPrintsIDs to try to get the file from the p2p network
	//2) try to decompress the file
	//3) find the 4th to last period in the file. This will be where the separator bytes for signatures and counter were added
	//4) remove the period and everything after
	//5) try to base64 decode this
	//6) try to JSON decode into a pastel.DDAndFingerprints
	//7) if all these are successful, return this struct
	//8) else, something got messed up somewhere so keep iterating through the files until successful
	for i := 0; i < len(DDAndFingerprintsIDs); i++ {
		decKey := base58.Decode(DDAndFingerprintsIDs[i])
		dbKey := hex.EncodeToString(decKey)

		log.WithContext(ctx).WithField("hash", decKey).WithField("db_key", dbKey).WithField("txid", txid).
			WithField("dd_fp_id", DDAndFingerprintsIDs[i]).Info("DDAndFingerPrintDetails trying to fetch & decode this file")

		retData, err := task.downloadDDFPFile(ctx, DDAndFingerprintsIDs[i])
		if err != nil {
			log.WithContext(ctx).WithField("db_key", dbKey).WithField("txid", txid).WithField("dd_fp_id", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails failed to download this file. ")
			continue
		}

		log.WithContext(ctx).WithField("txid", txid).WithField("db_key", dbKey).
			WithField("Hash", DDAndFingerprintsIDs[i]).Info("DDAndFingerPrintDetails successfully downloaded this file. ")

		return retData, nil
	}

	return nil, errors.Errorf("could not get dd and fingerprints for any file tested")
}

//utility function for getting the DD and Fingerprint Details given
//	a list of IDs (presumably from AppTicketData's DDAndFingerprintIDs), fingerprint IC, and fingerprint max
//1) iterate over DDAndFingerPrintsIDs to try to get the file from the p2p network
//2) try to decompress the file
//3) find the 4th to last period in the file. This will be where the separator bytes for signatures and counter were added
//4) remove the period and everything after
//5) try to base64 decode this
//6) try to JSON decode into a pastel.DDAndFingerprints
//7) if all these are successful, return this struct
//8) else, something got messed up somewhere so keep iterating through the files until successful

type restoreInfo struct {
	pastelID   string
	rqIDs      []string
	rqOti      []byte
	dataHash   []byte
	ddAndFpIDs []string
	isPublic   bool
}

func (task *NftDownloadingTask) getTicketInfo(ctx context.Context, txid string) (info restoreInfo, err error) {

	switch task.ttype {
	case pastel.ActionTypeCascade, pastel.ActionTypeSense:
		ticket, err := task.PastelClient.ActionRegTicket(ctx, txid)
		if err != nil {
			task.UpdateStatus(common.StatusNftRegGettingFailed)
			return info, errors.Errorf("could not get action registered ticket: %w, txid: %s", err, txid)
		}

		actionTicket, err := pastel.DecodeActionTicket([]byte(ticket.ActionTicketData.ActionTicket))
		if err != nil {
			task.UpdateStatus(common.StatusNftRegDecodingFailed)
			return info, errors.Errorf("cloud not decode action ticket: %w", err)
		}
		ticket.ActionTicketData.ActionTicketData = *actionTicket

		if task.ttype == pastel.ActionTypeCascade {
			if ticket.ActionTicketData.ActionTicketData.ActionType != pastel.ActionTypeCascade {
				return info, errors.Errorf("ticket type mismatch - ticket is not cascade but %s: %w, txid: %s",
					ticket.ActionTicketData.ActionTicketData.ActionType, err, txid)
			}

			cTicket, err := ticket.ActionTicketData.ActionTicketData.APICascadeTicket()
			if err != nil {
				task.UpdateStatus(common.StatusNftRegDecodingFailed)
				return info, errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			}

			info.pastelID = ticket.ActionTicketData.ActionTicketData.Caller
			info.rqIDs = cTicket.RQIDs
			info.rqOti = cTicket.RQOti
			info.dataHash = cTicket.DataHash
			info.isPublic = cTicket.MakePubliclyAccessible
		} else {
			if ticket.ActionTicketData.ActionTicketData.ActionType != pastel.ActionTypeSense {
				return info, errors.Errorf("ticket type mismatch - ticket is not sense but %s: %w, txid: %s",
					ticket.ActionTicketData.ActionTicketData.ActionType, err, txid)
			}

			sTicket, err := ticket.ActionTicketData.ActionTicketData.APISenseTicket()
			if err != nil {
				task.UpdateStatus(common.StatusNftRegDecodingFailed)
				return info, errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			}
			info.ddAndFpIDs = sTicket.DDAndFingerprintsIDs
			info.dataHash = sTicket.DataHash
		}
	default:
		nftRegTicket, err := task.PastelClient.RegTicket(ctx, txid)
		if err != nil {
			task.UpdateStatus(common.StatusNftRegGettingFailed)
			return info, errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
		}

		// Decode Art Request
		err = task.decodeRegTicket(&nftRegTicket)
		if err != nil {
			task.UpdateStatus(common.StatusNftRegDecodingFailed)
			return info, errors.Errorf("could not decode registered ticket: %w, txid: %s", err, txid)
		}
		info.pastelID = nftRegTicket.RegTicketData.NFTTicketData.Author
		info.rqIDs = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs
		info.rqOti = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQOti
		info.dataHash = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash
		info.ddAndFpIDs = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs
		info.isPublic = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.MakePubliclyAccessible
	}

	return info, nil
}

// Download downloads image and return the image.
func (task *NftDownloadingTask) Download(ctx context.Context, txid, timestamp, signature, ttxid, ttype string) ([]byte, error) {
	var err error
	if err = task.RequiredStatus(common.StatusTaskStarted); err != nil {
		log.WithContext(ctx).WithField("status", task.Status().String()).Error("Wrong task status")
		return nil, errors.Errorf("wrong status: %w", err)
	}

	var file []byte

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).WithField("txid", txid).WithField("ttype", ttype).Info("Downloading File request received")
		// Validate timestamp is not older than 10 minutes
		now := time.Now().UTC()
		lastTenMinutes := now.Add(time.Duration(-10) * time.Minute)
		requestTime, _ := time.Parse(time.RFC3339, timestamp)
		if lastTenMinutes.After(requestTime) {
			err = errors.New("request time is older than 10 minutes")
			task.UpdateStatus(common.StatusRequestTooLate)
			return nil
		}
		task.ttype = ttype

		if ttype == pastel.ActionTypeSense {
			data, err := task.DownloadDDAndFingerprints(ctx, txid)
			if err != nil {
				err = errors.Errorf("downloa dd & fingerprints file: %w", err)
				task.UpdateStatus(common.StatusFileRestoreFailed)
				return nil
			}

			file = data

			if len(file) == 0 {
				log.WithContext(ctx).WithField("txid", txid).Info("sense file nil downloaded")
				err = errors.New("nil restored file")
				task.UpdateStatus(common.StatusFileEmpty)
			}

			log.WithContext(ctx).WithField("txid", txid).Info("sense file downloaded successfully")
			return nil
		}

		info, err := task.getTicketInfo(ctx, txid)
		if err != nil {
			return errors.Errorf("Could not get ticket info: %w", err)
		}

		pastelID := info.pastelID
		if info.pastelID == "" {
			// err in retrieval
			err = errors.New("getTicketInfo failed")
			return nil
		}

		log.WithContext(ctx).WithField("txid", txid).Info("file download check owner")
		if len(ttxid) > 0 {
			// Get list of non sold Trade ticket owened by the owner of the PastelID from request
			// by calling command `tickets list trade available`
			var tradeTickets []pastel.TradeTicket
			tradeTickets, err = task.PastelClient.ListAvailableTradeTickets(ctx)
			if err != nil {
				err = errors.Errorf("could not get available trade tickets: %w", err)
				task.UpdateStatus(common.StatusListTradeTicketsFailed)
				return nil
			}

			// Validate that Trade ticket with ttxid is in the list
			if len(tradeTickets) == 0 {
				err = errors.New("not found any available trade tickets")
				task.UpdateStatus(common.StatusTradeTicketsNotFound)
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
				task.UpdateStatus(common.StatusTradeTicketMismatched)
				return nil
			}
		}
		// Validate timestamp signature with PastelID from Trade ticket
		// by calling command `pastelid verify timestamp-string signature PastelID`
		if !info.isPublic {
			var isValid bool
			isValid, err = task.PastelClient.Verify(ctx, []byte(timestamp), signature, pastelID, "ed448")
			if err != nil {
				err = errors.Errorf("timestamp signature verify: %w", err)
				task.UpdateStatus(common.StatusTimestampVerificationFailed)
				return nil
			}

			if !isValid {
				err = errors.New("invalid signature timestamp")
				task.UpdateStatus(common.StatusTimestampInvalid)
				return nil
			}
		}

		log.WithContext(ctx).WithField("txid", txid).Info("file download begin")
		// Get symbol identifiers files from Kademlia by using rq_ids - from Art Registration ticket
		// Get the list of "symbols/chunks" from Kademlia by using symbol identifiers from file
		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_oti)
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		file, err = task.restoreFile(ctx, info.rqIDs, info.rqOti, info.dataHash, txid)
		if err != nil {
			log.WithContext(ctx).WithField("txid", txid).Error("restore file failed")
			err = errors.Errorf("restore file: %w", err)
			task.UpdateStatus(common.StatusFileRestoreFailed)
			return nil
		}

		if len(file) == 0 {
			log.WithContext(ctx).WithField("txid", txid).Error("nil file downloaded")
			err = errors.New("nil restored file")
			task.UpdateStatus(common.StatusFileEmpty)
		}

		log.WithContext(ctx).WithField("txid", txid).Info("file downloaded successfully")
		return nil
	})

	return file, err
}

func (task *NftDownloadingTask) decodeRegTicket(nftRegTicket *pastel.RegTicket) error {
	articketData, err := pastel.DecodeNFTTicket(nftRegTicket.RegTicketData.NFTTicket)
	if err != nil {
		return errors.Errorf("convert NFT ticket: %w", err)
	}
	nftRegTicket.RegTicketData.NFTTicketData = *articketData

	return nil
}

func (task *NftDownloadingTask) getRQSymbolIDs(ctx context.Context, id string, rqIDsData []byte) (rqIDs []string, err error) {
	fileContent, err := utils.Decompress(rqIDsData)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Decompress compressed symbol IDs file failed")
		task.UpdateStatus(common.StatusSymbolFileInvalid)
	}

	log.WithContext(ctx).WithField("Content", string(fileContent)).Debugf("symbol IDs file")

	var rqData []byte
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] == pastel.SeparatorByte {
			rqData = fileContent[:i]
			if i+1 >= len(fileContent) {
				return rqIDs, errors.New("invalid rqIDs data")
			}
			break
		}
	}

	rqDataJSON, err := utils.B64Decode(rqData)
	if err != nil {
		return rqIDs, errors.Errorf("decode %s failed: %w", string(rqData), err)
	}

	file := rqnode.RawSymbolIDFile{}
	if err := json.Unmarshal(rqDataJSON, &file); err != nil {
		log.WithContext(ctx).WithError(err).WithField("Content", string(fileContent)).
			WithField("file", string(rqIDsData)).Error("rq: parsing symbolID file failure")

		return rqIDs, errors.Errorf("parsing file: %s - file content: %s - err: %w", string(rqIDsData), string(fileContent), err)
	}

	return file.SymbolIdentifiers, nil
}

func (task *NftDownloadingTask) removeArtifacts() {
}

// NewNftDownloadingTask returns a new Task instance.
func NewNftDownloadingTask(service *NftDownloaderService) *NftDownloadingTask {
	return &NftDownloadingTask{
		SuperNodeTask:        common.NewSuperNodeTask(logPrefix, service.historyDB),
		NftDownloaderService: service,
		RqClient:             rqgrpc.NewClient(),
	}
}
