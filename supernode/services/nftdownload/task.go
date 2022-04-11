package nftdownload

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/supernode/services/common"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"

	"github.com/pastelnetwork/gonode/common/b85"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// NftDownloadingTask is the task of registering new Nft.
type NftDownloadingTask struct {
	*common.SuperNodeTask

	*NftDownloaderService

	RQSymbolsDir string
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
	nftTicket := &pastel.NFTTicket{}
	json.Unmarshal(regTicket.RegTicketData.NFTTicket, nftTicket)

	//decode the appticketdata
	appTicketData, err := b85.Decode(nftTicket.AppTicket)
	if err != nil {
		log.Warnf("b85 decoding failed, trying to base64 decode - err: %v", err)
		appTicketData, err = base64.StdEncoding.DecodeString(regTicket.RegTicketData.NFTTicketData.AppTicket)
		if err != nil {
			return nil, fmt.Errorf("b64 decode: %v", err)
		}
	}

	appTicket := &pastel.AppTicket{}
	json.Unmarshal(appTicketData, appTicket)

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
	resMap[1] = file2
	return resMap, err
}

// DownloadDDAndFingerprints gets dd and fp file from ticket based on id and returns the file.
func (task *NftDownloadingTask) DownloadDDAndFingerprints(ctx context.Context, txid string) ([]byte, error) {
	log.WithContext(ctx).WithField("txid", txid).Println("Getting dd and fingerprints for txid")
	var err error
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
	nftTicket := &pastel.NFTTicket{}
	json.Unmarshal(regTicket.RegTicketData.NFTTicket, nftTicket)

	//decode the appticketdata
	appTicketData, err := b85.Decode(nftTicket.AppTicket)
	if err != nil {
		log.Warnf("b85 decoding failed, trying to base64 decode - err: %v", err)
		appTicketData, err = base64.StdEncoding.DecodeString(regTicket.RegTicketData.NFTTicketData.AppTicket)
		if err != nil {
			return nil, fmt.Errorf("b64 decode: %v", err)
		}
	}

	appTicket := &pastel.AppTicket{}
	json.Unmarshal(appTicketData, appTicket)

	DDAndFingerprintsIDs := appTicket.DDAndFingerprintsIDs
	log.WithContext(ctx).WithField("ddandfpids", DDAndFingerprintsIDs).WithField("NFTTicket", nftTicket).Println("Found dd and fp ids")

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
		file, err := task.P2PClient.Retrieve(ctx, DDAndFingerprintsIDs[i])
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails tried to get this file and failed. ")
			continue
		}
		log.WithContext(ctx).WithField("file", file).Println("Got the file")
		decompressedData, err := zstd.Decompress(nil, file)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails failed to decompress this file. ")
			continue
		}
		log.WithContext(ctx).Println("Decompressed the file")
		//base64 dataset doesn't contain periods, so we just find the first index of period and chop it and everything else
		firstIndexOfSeparator := bytes.IndexByte(decompressedData, pastel.SeparatorByte)
		if firstIndexOfSeparator < 1 {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails got a bad separator index. ")
			continue
		}
		dataToBase64Decode := decompressedData[:firstIndexOfSeparator]
		log.WithContext(ctx).WithField("first index of separator", firstIndexOfSeparator).WithField("decompd data", string(decompressedData)).WithField("datatobase64decode", string(dataToBase64Decode)).Println("About to base64 decode file")
		dataToJSONDecode, err := utils.B64Decode(dataToBase64Decode)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails could not base64 decode. ")
			continue
		}
		log.WithContext(ctx).Println("base64 decoded the file")
		ddAndFingerprintsStruct := &pastel.DDAndFingerprints{}
		err = json.Unmarshal(dataToJSONDecode, ddAndFingerprintsStruct)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails could not JSON unmarshal. ")
			continue
		}
		log.WithContext(ctx).WithField("ddfpstruct", ddAndFingerprintsStruct).Println("Returning this file in byte form, this was json unmarshallable")
		//dataToJSONDecode is just the DDAndFingerprints file we'd like to return at this point
		return dataToJSONDecode, nil

	}
	return nil, errors.Errorf("could not get dd and fingerprints for any file tested")
}

//utility functions to download dupe detection and fingerprint files
// //
// func (service *NftSearchingService) findAndConnectToTopValidSupernode(ctx context.Context, pastelID string, passphrase string) (node.ConnectionInterface, error) {
// 	mns, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	topAddress := ""
// 	gotANode := false
// 	for _, mn := range mns {
// 		if !gotANode && mn.ExtKey == "" || mn.ExtAddress == "" {
// 			continue
// 		}

// 		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
// 		_, err = service.pastelHandler.PastelClient.FindTicketByID(ctx, mn.ExtKey)
// 		if err != nil {
// 			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
// 			continue
// 		}
// 		gotANode = true
// 		topAddress = mn.ExtAddress
// 	}
// 	alts := &alts.SecInfo{
// 		PastelID:   pastelID,
// 		PassPhrase: passphrase,
// 		Algorithm:  "ed448",
// 	}

// 	conn, err := service.nodeClient.Connect(ctx, topAddress, alts)
// 	return conn, err
// }

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
// func (service *NftSearchingService) GetDDAndFingerprintDetailsFromAppTicketDetails(ctx context.Context, DDAndFingerprintsIDs []string) (*pastel.DDAndFingerprints, error) {
//connect to a node
// conn, err := service.findAndConnectToTopValidSupernode(ctx, opts.Configs.PastelID, opts.Configs.PassPhrase)
// if err != nil {
// 	log.WithContext(ctx).Warn("Failed to connect to supernode for ")
// }

// successfulDecoding := false
// for i := 0; i < len(DDAndFingerprintsIDs) && !successfulDecoding; i++ {
// 	file, err := conn.DownloadFile().DownloadFile(ctx, []byte(DDAndFingerprintsIDs[i]))
// 	if err != nil {
// 		log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails tried to get this file and failed. ")
// 		continue
// 	}
// 	decompressedData, err := zstd.Decompress(nil, file)
// 	if err != nil {
// 		log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails failed to decompress this file. ")
// 		continue
// 	}

// }
// }

func (task *NftDownloadingTask) getTicketInfo(ctx context.Context, txid string, ttype string) (pastelID string, rqIDs []string,
	rqOti []byte, dataHash []byte) {

	switch ttype {
	case pastel.ActionTypeCascade:
		ticket, err := task.PastelClient.ActionRegTicket(ctx, txid)
		if err != nil {
			err = errors.Errorf("could not get action registered ticket: %w, txid: %s", err, txid)
			task.UpdateStatus(common.StatusNftRegGettingFailed)
			return
		}
		fmt.Println("got ticket: ", ticket)
		actionTicket, err := pastel.DecodeActionTicket([]byte(ticket.ActionTicketData.ActionTicket))
		if err != nil {
			err = errors.Errorf("cloud not decode action ticket: %w", err)
			task.UpdateStatus(common.StatusNftRegDecodingFailed)
			return
		}
		ticket.ActionTicketData.ActionTicketData = *actionTicket

		log.WithContext(ctx).Debugf("Art ticket: %s", string(ticket.ActionTicketData.ActionTicket))
		cTicket, err := ticket.ActionTicketData.ActionTicketData.APICascadeTicket()
		if err != nil {
			err = errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			task.UpdateStatus(common.StatusNftRegDecodingFailed)
			return
		}

		pastelID = ticket.ActionTicketData.ActionTicketData.Caller
		rqIDs = cTicket.RQIDs
		rqOti = cTicket.RQOti
		dataHash = cTicket.DataHash
	default:
		nftRegTicket, err := task.PastelClient.RegTicket(ctx, txid)
		if err != nil {
			err = errors.Errorf("could not get registered ticket: %w, txid: %s", err, txid)
			task.UpdateStatus(common.StatusNftRegGettingFailed)
			return
		}

		log.WithContext(ctx).Debugf("Art ticket: %s", string(nftRegTicket.RegTicketData.NFTTicket))

		// Decode Art Request
		err = task.decodeRegTicket(&nftRegTicket)
		if err != nil {
			task.UpdateStatus(common.StatusNftRegDecodingFailed)
			return
		}
		pastelID = nftRegTicket.RegTicketData.NFTTicketData.Author
		rqIDs = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs
		rqOti = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.RQOti
		dataHash = nftRegTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash
	}

	return
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
		// Validate timestamp is not older than 10 minutes
		now := time.Now()
		lastTenMinutes := now.Add(time.Duration(-10) * time.Minute)
		requestTime, _ := time.Parse(time.RFC3339, timestamp)
		if lastTenMinutes.After(requestTime) {
			err = errors.New("request time is older than 10 minutes")
			task.UpdateStatus(common.StatusRequestTooLate)
			return nil
		}

		pastelID, rqIDs, rqOti, dataHash := task.getTicketInfo(ctx, txid, ttype)
		if pastelID == "" {
			// err in retrieval
			err = errors.New("getTicketInfo failed")
			return nil
		}

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

		// Get symbol identifiers files from Kademlia by using rq_ids - from Art Registration ticket
		// Get the list of "symbols/chunks" from Kademlia by using symbol identifiers from file
		// Pass all symbols/chunks to the raptorq service to decode (also passing encoder parameters: rq_oti)
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		file, err = task.restoreFile(ctx, rqIDs, rqOti, dataHash)
		if err != nil {
			err = errors.Errorf("restore file: %w", err)
			task.UpdateStatus(common.StatusFileRestoreFailed)
			return nil
		}

		if len(file) == 0 {
			err = errors.New("nil restored file")
			task.UpdateStatus(common.StatusFileEmpty)
		}

		return nil
	})

	return file, err
}

func (task *NftDownloadingTask) restoreFile(ctx context.Context, rqIDs []string, rqOti []byte, dataHash []byte) ([]byte, error) {
	var file []byte
	var lastErr error
	var err error

	if len(rqIDs) == 0 {
		task.UpdateStatus(common.StatusNftRegTicketInvalid)
		return file, errors.Errorf("ticket has empty symbol identifier files")
	}

	var rqConnection rqnode.Connection
	rqConnection, err = task.NftDownloaderService.RQClient.Connect(ctx, task.NftDownloaderService.config.RaptorQServiceAddress)
	if err != nil {
		task.UpdateStatus(common.StatusRQServiceConnectionFailed)
		return file, errors.Errorf("could not connect to rqservice: %w", err)
	}
	defer rqConnection.Done()
	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.NftDownloaderService.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)

	for _, id := range rqIDs {
		var rqIDsData []byte
		rqIDsData, err = task.P2PClient.Retrieve(ctx, id)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Retrieve compressed symbol IDs file from P2P failed")
			lastErr = errors.Errorf("retrieve compressed symbol IDs file: %w", err)
			task.UpdateStatus(common.StatusSymbolFileNotFound)
			continue
		}

		log.WithContext(ctx).WithField("rqIDsData", string(rqIDsData)).Debugf("rqIDs Data")

		var rqIDs []string
		rqIDs, err = task.getRQSymbolIDs(ctx, id, rqIDsData)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Parse symbol IDs failed")
			lastErr = errors.Errorf("parse symbol IDs: %w", err)
			task.UpdateStatus(common.StatusSymbolFileInvalid)
			continue
		}

		log.WithContext(ctx).Debugf("Symbol IDs: %v", rqIDs)

		symbols := make(map[string][]byte)
		for _, id := range rqIDs {
			var symbol []byte
			symbol, err = task.P2PClient.Retrieve(ctx, id)
			if err != nil {
				log.WithContext(ctx).WithField("SymbolID", id).Warn("Could not retrieve symbol")
				task.UpdateStatus(common.StatusSymbolNotFound)
				break
			}

			// Validate that the hash of each "symbol/chunk" matches its id
			h := sha3.Sum256(symbol)

			storedID := base58.Encode(h[:])
			if storedID != id {
				log.WithContext(ctx).Warnf("Symbol ID mismatched, expect %v, got %v", id, storedID)
				task.UpdateStatus(common.StatusSymbolMismatched)
				break
			}
			symbols[id] = symbol
		}
		if len(symbols) != len(rqIDs) {
			log.WithContext(ctx).WithField("SymbolIDsFileId", id).Warn("Could not retrieve all symbols")
			lastErr = errors.New("could not retrieve all symbols from Kademlia")
			task.UpdateStatus(common.StatusSymbolsNotEnough)
			continue
		}

		// Restore Nft
		var decodeInfo *rqnode.Decode
		encodeInfo := rqnode.Encode{
			Symbols: symbols,
			EncoderParam: rqnode.EncoderParameters{
				Oti: rqOti,
			},
		}

		decodeInfo, err = rqService.Decode(ctx, &encodeInfo)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Restore file with rqserivce")
			lastErr = errors.Errorf("restore file with rqserivce: %w", err)
			task.UpdateStatus(common.StatusFileDecodingFailed)
			continue
		}
		task.UpdateStatus(common.StatusFileDecoded)

		// log.WithContext(ctx).Debugf("Restored file path: %s", decodeInfo.Path)
		// log.WithContext(ctx).Debugf("Restored file: %s", string(restoredFile))
		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		fileHash := sha3.Sum256(decodeInfo.File)

		if !bytes.Equal(fileHash[:], dataHash) {
			log.WithContext(ctx).WithField("SymbolIDsFileId", id).Warn("hash file mismatched")
			lastErr = errors.New("hash file mismatched")
			task.UpdateStatus(common.StatusFileMismatched)
			continue
		}

		return decodeInfo.File, nil
	}

	return file, lastErr
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
	fileContent, err := zstd.Decompress(nil, rqIDsData)
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
		return rqIDs, errors.Errorf("parsing file: %s - err: %w", string(rqIDsData), err)
	}

	return file.SymbolIdentifiers, nil
}

func (task *NftDownloadingTask) removeArtifacts() {
}

// NewNftDownloadingTask returns a new Task instance.
func NewNftDownloadingTask(service *NftDownloaderService) *NftDownloadingTask {
	return &NftDownloadingTask{
		SuperNodeTask:        common.NewSuperNodeTask(logPrefix),
		NftDownloaderService: service,
	}
}
