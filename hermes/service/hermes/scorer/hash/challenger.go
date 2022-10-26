package hash

import (
	"context"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/store"
	"github.com/pastelnetwork/gonode/mixins"

	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer/network"
	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
)

// ImageHashChallenger helps with thumbnail challenges
type ImageHashChallenger struct {
	meshHandler         *network.MeshHandler
	pastelClient        pastel.Client
	imageHandler        *mixins.NftImageHandler
	fingerprintsHandler *mixins.FingerprintsHandler
	store               store.ScoreStore

	dataHash           []byte
	creatorPastelID    string
	creatorPassphrase  string
	creatorBlockHeight int
	creatorBlockHash   string
	creationTimestamp  string
}

// NewImageHashChallenger returns a new instance of image hash Challenger
func NewImageHashChallenger(meshHandler *network.MeshHandler, pastelClient pastel.Client, store store.ScoreStore, pastelID string, pass string) *ImageHashChallenger {
	return &ImageHashChallenger{
		meshHandler:         meshHandler,
		pastelClient:        pastelClient,
		fingerprintsHandler: mixins.NewFingerprintsHandler(mixins.NewPastelHandler(pastelClient)),
		imageHandler:        mixins.NewImageHandler(mixins.NewPastelHandler(pastelClient)),
		store:               store,
		creatorPastelID:     pastelID,
		creatorPassphrase:   pass,
	}
}

// Run runs the challenger once
func (s *ImageHashChallenger) Run(ctx context.Context) error {
	if err := s.run(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("run image hash challenger failed")
	} else {
		log.WithContext(ctx).Info("image hash challenger run success")
	}

	return nil
}

func (s *ImageHashChallenger) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	creatorBlockHeight, creatorBlockHash, err := s.meshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	log.WithContext(ctx).Info("image hash challenger connected")

	s.creatorBlockHeight = creatorBlockHeight
	s.creatorBlockHash = creatorBlockHash
	//set timestamp for nft reg metadata
	s.creationTimestamp = time.Now().Format("YYYY-MM-DD hh:mm:ss")

	log.WithContext(ctx).WithField("block-height", creatorBlockHeight).Info("Image Hash challenge mesh of SNs created")
	nodesDone := s.meshHandler.ConnectionsSupervisor(ctx, cancel)

	// send registration metadata
	if err := s.sendRegMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}
	log.WithContext(ctx).Info("Image Hash challenge Reg Metadata sent")

	nftFile, err := newTestImageFile(creatorBlockHash)
	if err != nil {
		return errors.Errorf("creating test image: %w", err)
	}

	// probe ORIGINAL image for average rareness, nsfw and seen score
	if err := s.probeImage(ctx, nftFile, nftFile.Name()); err != nil {
		return errors.Errorf("probe image: %w", err)
	}
	log.WithContext(ctx).Info("Image Hash challenge Image probed")

	if err := s.fingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, 50); err != nil {
		return errors.Errorf("DD and/or Fingerprint ID error: %w", err)
	}
	log.WithContext(ctx).Info("Image Hash challenge DD And Fgs created")

	if err := s.imageHandler.CreateCopyWithEncodedFingerprint(ctx,
		s.creatorPastelID, s.creatorPassphrase,
		s.fingerprintsHandler.FinalFingerprints, nftFile); err != nil {
		return errors.Errorf("encode image with fingerprint: %w", err)
	}

	if s.dataHash, err = s.imageHandler.GetHash(); err != nil {
		return errors.Errorf("get image hash: %w", err)
	}
	log.WithContext(ctx).Info("Image Hash challenge Image Hash generated")

	if err := s.uploadImage(ctx); err != nil {
		return errors.Errorf("upload image: %w", err)
	}
	log.WithContext(ctx).Info("Image Hash challenge completed")

	// don't need SNs anymore
	_ = s.meshHandler.CloseSNsConnections(ctx, nodesDone)

	return nil
}

func (s *ImageHashChallenger) uploadImage(ctx context.Context) error {
	// Upload image with pqgsinganature and its thumb to supernodes
	if err := s.uploadImageWithThumbnail(ctx, s.imageHandler.ImageEncodedWithFingerprints, files.ThumbnailCoordinate{TopLeftX: 0, TopLeftY: 0,
		BottomRightX: 640, BottomRightY: 480}); err != nil {
		return errors.Errorf("upload encoded image and thumbnail coordinate: %w", err)
	}

	badNodes := s.imageHandler.GetUnmatchingHashPastelID()
	log.WithContext(ctx).WithField("bad-nodes", badNodes).Info("Bad Nodes rcvd")
	for _, bnode := range badNodes {
		for _, n := range s.meshHandler.Nodes {
			if bnode == n.PastelID() {
				score := &domain.SnScore{
					TxID:      n.TxID(),
					PastelID:  n.TxID(),
					IPAddress: n.Address(),
				}

				log.WithContext(ctx).WithField("txid", n.TxID()).WithField("ip-address", n.Address()).
					Info("Increment score of SN")
				updatedScore, err := s.store.IncrementScore(ctx, score, 1)
				if err != nil {
					return errors.Errorf("unable to increment score: %w , pastelID: %s", err, n.PastelID())
				}

				if updatedScore.Score >= domain.ScoreReportLimit {
					log.WithContext(ctx).Info("score limit exceeded, reporting SN to cNode")
					if err := s.pastelClient.IncrementPoseBanScore(ctx, strings.Split(n.TxID(), "-")[0], n.Idx()); err != nil {
						log.WithContext(ctx).WithField("txid", n.TxID()).WithError(err).Error("failed to increment pose-ban score")
					}
				}
			}
		}
	}

	return nil
}

// uploadImageWithThumbnail uploads the image with pqsignatured appended and thumbnail's coordinate to super nodes
func (s *ImageHashChallenger) uploadImageWithThumbnail(ctx context.Context, file *files.File, thumbnail files.ThumbnailCoordinate) error {
	group, gctx := errgroup.WithContext(ctx)
	s.imageHandler.ClearHashes()

	for _, someNode := range s.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*network.NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			hash1, hash2, hash3, err := nftRegNode.UploadImageWithThumbnail(gctx, file, thumbnail)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", someNode).Error("upload image with thumbnail failed")
				return err
			}
			s.imageHandler.AddNewHashes(hash1, hash2, hash3, someNode.PastelID())
			return nil
		})
	}

	return group.Wait()
}

func (s *ImageHashChallenger) probeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	s.fingerprintsHandler.Clear()

	// Send image to supernodes for probing.
	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range s.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*network.NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			// result is 4.B.5
			compress, stateOk, err := nftRegNode.ProbeImage(gctx, file)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
				return errors.Errorf("remote node %s: indicated processing error", someNode.String())
			}
			// Step 5.1
			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", someNode.String(), err)
			}
			s.fingerprintsHandler.AddNew(fingerprintAndScores, fingerprintAndScoresBytes, signature, someNode.PastelID())

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}
	// Step 5.2 - 5.3
	if err := s.fingerprintsHandler.Match(ctx); err != nil {
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}

	return nil
}

func (s *ImageHashChallenger) sendRegMetadata(ctx context.Context) error {
	if s.creatorBlockHash == "" {
		return errors.New("empty current block hash")
	}
	if s.creatorPastelID == "" {
		return errors.New("empty creator pastelID")
	}

	regMetadata := &types.NftRegMetadata{
		BlockHash:       s.creatorBlockHash,
		CreatorPastelID: s.creatorPastelID,
		BlockHeight:     strconv.Itoa(s.creatorBlockHeight),
		Timestamp:       s.creationTimestamp,
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range s.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*network.NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = nftRegNode.SendRegMetadata(gctx, regMetadata)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

func newTestImageFile(blockHash string) (*files.File, error) {
	imageStorage := files.NewStorage(fs.NewFileStorage(os.TempDir()))
	imgFile := imageStorage.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create storage file: %w", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	png.Encode(f, img)
	imgFile.SetFormat(1)

	encSig := qrsignature.New(qrsignature.BlockHash([]byte(blockHash)))
	if err := imgFile.Encode(encSig); err != nil {
		return nil, errors.Errorf("encode fingerprint into image: %w", err)
	}

	return imgFile, nil
}
