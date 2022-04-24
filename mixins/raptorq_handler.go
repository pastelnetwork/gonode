package mixins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// RQHandler handles rq stuff
type RQHandler struct {
	pastelHandler *PastelHandler
	rqClient      rqnode.ClientInterface

	raptorQServiceAddress string
	rqFilesDir            string
	numberRQIDSFiles      uint32
	maxRQIDs              uint32

	RQIDs          []string
	RQEncodeParams rqnode.EncoderParameters
	RQIDsFile      []byte
	RQIDsIc        uint32
}

// NewRQHandler creates new rq handler
func NewRQHandler(rqClient rqnode.ClientInterface, pastelHandler *PastelHandler, raptorQServiceAddress string, rqFilesDir string,
	numberRQIDSFiles uint32, maxRQIDs uint32) *RQHandler {
	return &RQHandler{
		rqClient:              rqClient,
		raptorQServiceAddress: raptorQServiceAddress,
		rqFilesDir:            rqFilesDir,
		numberRQIDSFiles:      numberRQIDSFiles,
		maxRQIDs:              maxRQIDs,
		pastelHandler:         pastelHandler,
	}
}

// GenRQIdentifiersFiles generates rq id files
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step 9
func (h *RQHandler) GenRQIdentifiersFiles(ctx context.Context, file *files.File, operationBlockHash string, callerPastelID string, callerPassphrase string) error {
	log.Debugf("Connect to %s", h.raptorQServiceAddress)
	conn, err := h.rqClient.Connect(ctx, h.raptorQServiceAddress)
	if err != nil {
		return errors.Errorf("connect to raptorQ: %w", err)
	}
	defer conn.Close()

	content, err := file.Bytes()
	if err != nil {
		return errors.Errorf("read data file content: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: h.rqFilesDir,
	})

	// FIXME :
	// - check format of block hash - should be base58 or not
	encodeInfo, err := rqService.EncodeInfo(ctx, content, h.numberRQIDSFiles, operationBlockHash, callerPastelID)
	if err != nil {
		return errors.Errorf("generate RaptorQ symbols identifiers: %w", err)
	}

	var rqIDsFilesCount uint32
	for _, rawSymbolIDFile := range encodeInfo.SymbolIDFiles {
		if len(rawSymbolIDFile.SymbolIdentifiers) == 0 {
			return errors.Errorf("empty symbol identifiers - rawFile : %s", rawSymbolIDFile)
		}

		err := h.generateRQIDs(ctx, rawSymbolIDFile, callerPastelID, callerPassphrase)
		if err != nil {
			return errors.Errorf("create RQIDs file :%w", err)
		}
		rqIDsFilesCount++
		break
	}

	if rqIDsFilesCount != h.numberRQIDSFiles {
		return errors.Errorf("number of RaptorQ symbol identifiers files must be %d, most probably old version of rq-service is installed", h.numberRQIDSFiles)
	}
	h.RQEncodeParams = encodeInfo.EncoderParam

	return nil
}

func (h *RQHandler) generateRQIDs(ctx context.Context, rawFile rqnode.RawSymbolIDFile, callerPastelID string, callerPassphrase string) error {
	rqIDsfile, err := json.Marshal(rawFile)
	if err != nil {
		return fmt.Errorf("marshal rqID file")
	}

	signature, err := h.pastelHandler.PastelClient.Sign(ctx, rqIDsfile, callerPastelID,
		callerPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign identifiers file: %w", err)
	}

	encRqIDsfile := utils.B64Encode(rqIDsfile)

	var buffer bytes.Buffer
	buffer.Write(encRqIDsfile)
	buffer.WriteString(".")
	buffer.Write(signature)
	rqIDFile := buffer.Bytes()

	h.RQIDsIc = rand.Uint32()
	h.RQIDs, _, err = pastel.GetIDFiles(rqIDFile, h.RQIDsIc, h.maxRQIDs)
	if err != nil {
		return fmt.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, rqIDFile, 22)
	if err != nil {
		return errors.Errorf("compress: %w", err)
	}
	h.RQIDsFile = utils.B64Encode(comp)

	return nil
}

// IsEmpty() ...
func (h RQHandler) IsEmpty() bool {
	return len(h.RQIDs) == 0 || len(h.RQIDsFile) == 0 || len(h.RQEncodeParams.Oti) == 0
}
