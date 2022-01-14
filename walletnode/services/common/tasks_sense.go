package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/pastel"
	"time"
)

type ImageHandler struct {
	FileStorage *files.Storage
	FileDb      storage.KeyValue
	imageTTL    time.Duration
}

func NewImageHandler(
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
	imageTTL time.Duration) *ImageHandler {
	return &ImageHandler{
		FileStorage: files.NewStorage(fileStorage),
		FileDb:      db,
		imageTTL:    imageTTL,
	}
}

// StoreImage stores the image in the FileDb and storage and return image_id
func (st *ImageHandler) StoreFileNameIntoStorage(ctx context.Context, fileName *string) (string, string, error) {

	id, err := random.String(8, random.Base62Chars)
	if err := st.FileDb.Set(id, []byte(*fileName)); err != nil {
		return "", "0", err
	}

	file, err := st.FileStorage.File(*fileName)
	if err != nil {
		return "", "0", err
	}
	file.RemoveAfter(st.imageTTL)

	return id, time.Now().Add(st.imageTTL).Format(time.RFC3339), nil
}

// GetImgData returns the image data from the storage
func (st *ImageHandler) GetImgData(imageID string) ([]byte, error) {
	// get image filename from storage based on image_id
	filename, err := st.FileDb.Get(imageID)
	if err != nil {
		return nil, errors.Errorf("get image filename from FileDb: %w", err)
	}

	// get image data from storage
	file, err := st.FileStorage.File(string(filename))
	if err != nil {
		return nil, errors.Errorf("get image fd: %v", err)
	}

	imgData, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("read image data: %w", err)
	}

	return imgData, nil
}

type PastelHandler struct {
	PastelClient pastel.Client
}

func NewPastelHandler(pastelClient pastel.Client) *PastelHandler {
	return &PastelHandler{
		PastelClient: pastelClient,
	}
}

// VerifySignature verifies the signature of the image
func (pt *PastelHandler) VerifySignature(ctx context.Context, data []byte, signature string, pastelID string) error {
	ok, err := pt.PastelClient.Verify(ctx, data, signature, pastelID, pastel.SignAlgorithmED448)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("signature verification failed")
	}

	return nil
}

// GetEstimatedFee returns the estimated sense fee for the given image
func (pt *PastelHandler) GetEstimatedFee(ctx context.Context, ImgSizeInMb int64) (float64, error) {
	actionFees, err := pt.PastelClient.GetActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return 0, err
	}
	return actionFees.SenseFee, nil
}
func (pt *PastelHandler) pastelTopNodes(ctx context.Context, nodeClient node.Client, maximumFee float64) (node.List, error) {
	var nodes node.List

	mns, err := pt.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > maximumFee {
			continue
		}
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = pt.PastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes = append(nodes, node.NewNode(nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

func (pt *PastelHandler) WaitTxidValid(ctx context.Context, txID string, expectedConfirms int64, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s", expectedConfirms, txID)
	blockTracker := blocktracker.New(pt.PastelClient)
	baseBlkCnt, err := blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := pt.PastelClient.GetRawTransactionVerbose1(subCtx, txID)
				if err != nil {
					return errors.Errorf("get transaction: %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("not enough confirmations: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("check confirmations failed")
			} else {
				return nil
			}

			currentBlkCnt, err := blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				continue
			}

			if currentBlkCnt-baseBlkCnt >= int32(expectedConfirms)+2 {
				return errors.Errorf("timeout when wating for confirmation of transaction %s", txID)
			}
		}
	}
}
