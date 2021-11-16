package walletnode

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterArtwork represents grpc service for registration artwork.
type RegisterArtwork struct {
	pb.UnimplementedRegisterArtworkServer

	*common.RegisterArtwork
}

// Session implements walletnode.RegisterArtworkServer.Session()
func (service *RegisterArtwork) Session(stream pb.RegisterArtwork_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *artworkregister.Task

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := task.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")

	for {
		if _, err := stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			switch status.Code(err) {
			case codes.Canceled, codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}

// AcceptedNodes implements walletnode.RegisterArtworkServer.AcceptedNodes()
func (service *RegisterArtwork) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.AcceptedNodes(ctx)
	if err != nil {
		return nil, err
	}

	var peers []*pb.AcceptedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pb.AcceptedNodesReply_Peer{
			NodeID: node.ID,
		})
	}

	resp := &pb.AcceptedNodesReply{
		Peers: peers,
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterArtworkServer.ConnectTo()
func (service *RegisterArtwork) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")
	return resp, nil
}

// ProbeImage implements walletnode.RegisterArtworkServer.ProbeImage()
func (service *RegisterArtwork) ProbeImage(stream pb.RegisterArtwork_ProbeImageServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	image := service.Storage.NewFile()
	file, err := image.Create()
	if err != nil {
		return errors.Errorf("open file %q: %w", file.Name(), err)
	}
	defer file.Close()
	log.WithContext(ctx).WithField("filename", file.Name()).Debug("ProbeImage request")

	wr := bufio.NewWriter(file)
	defer wr.Flush()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return errors.Errorf("receive ProbeImage: %w", err)
		}

		if _, err := wr.Write(req.Payload); err != nil {
			return errors.Errorf("write to file %q: %w", file.Name(), err)
		}
	}

	err = image.UpdateFormat()
	if err != nil {
		return errors.Errorf("add image format: %w", err)
	}

	fingerAndScores, err := task.ProbeImage(ctx, image)
	if err != nil {
		return err
	}

	resp := &pb.ProbeImageReply{
		DupeDetectionVersion:      fingerAndScores.DupeDectectionSystemVersion,
		HashOfCandidateImg:        fingerAndScores.HashOfCandidateImageFile,
		AverageRarenessScore:      fingerAndScores.OverallAverageRarenessScore,
		IsRareOnInternet:          fingerAndScores.IsRareOnInternet,
		MatchesFoundOnFirstPage:   fingerAndScores.MatchesFoundOnFirstPage,
		NumberOfPagesOfResults:    fingerAndScores.NumberOfPagesOfResults,
		UrlOfFirstMatchInPage:     fingerAndScores.URLOfFirstMatchInPage,
		OpenNsfwScore:             fingerAndScores.OpenNSFWScore,
		ZstdCompressedFingerprint: fingerAndScores.ZstdCompressedFingerprint,
		AlternativeNsfwScore: &pb.ProbeImageReply_AlternativeNSFWScore{
			Drawing: fingerAndScores.AlternativeNSFWScore.Drawing,
			Hentai:  fingerAndScores.AlternativeNSFWScore.Hentai,
			Neutral: fingerAndScores.AlternativeNSFWScore.Neutral,
			Porn:    fingerAndScores.AlternativeNSFWScore.Porn,
			Sexy:    fingerAndScores.AlternativeNSFWScore.Sexy,
		},
		ImageHashes: &pb.ProbeImageReply_ImageHashes{
			PerceptualHash: fingerAndScores.ImageHashes.PerceptualHash,
			AverageHash:    fingerAndScores.ImageHashes.AverageHash,
			DifferenceHash: fingerAndScores.ImageHashes.DifferenceHash,
			PDQHash:        fingerAndScores.ImageHashes.PDQHash,
			NeuralHash:     fingerAndScores.ImageHashes.NeuralHash,
		},
		PerceptualHashOverlapCount:                    fingerAndScores.PerceptualHashOverlapCount,
		NumberOfFingerprintsRequiringFurtherTesting_1: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting1,
		NumberOfFingerprintsRequiringFurtherTesting_2: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting2,
		NumberOfFingerprintsRequiringFurtherTesting_3: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting3,
		NumberOfFingerprintsRequiringFurtherTesting_4: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting4,
		NumberOfFingerprintsRequiringFurtherTesting_5: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting5,
		NumberOfFingerprintsRequiringFurtherTesting_6: fingerAndScores.NumberOfFingerprintsRequiringFurtherTesting6,
		NumberOfFingerprintsOfSuspectedDupes:          fingerAndScores.NumberOfFingerprintsOfSuspectedDupes,
		PearsonMax:                                    fingerAndScores.PearsonMax,
		SpearmanMax:                                   fingerAndScores.SpearmanMax,
		KendallMax:                                    fingerAndScores.KendallMax,
		HoeffdingMax:                                  fingerAndScores.HoeffdingMax,
		MutualInformationMax:                          fingerAndScores.MutualInformationMax,
		HsicMax:                                       fingerAndScores.HsicMax,
		XgbimportanceMax:                              fingerAndScores.XgbimportanceMax,
		PearsonTop_1BpsPercentile:                     fingerAndScores.PearsonTop1BpsPercentile,
		SpearmanTop_1BpsPercentile:                    fingerAndScores.SpearmanTop1BpsPercentile,
		KendallTop_1BpsPercentile:                     fingerAndScores.KendallTop1BpsPercentile,
		HoeffdingTop_10BpsPercentile:                  fingerAndScores.HoeffdingTop10BpsPercentile,
		MutualInformationTop_100BpsPercentile:         fingerAndScores.MutualInformationTop100BpsPercentile,
		HsicTop_100BpsPercentile:                      fingerAndScores.HsicTop100BpsPercentile,
		XgbimportanceTop_100BpsPercentile:             fingerAndScores.XgbimportanceTop100BpsPercentile,
		CombinedRarenessScore:                         fingerAndScores.CombinedRarenessScore,
		XgboostPredictedRarenessScore:                 fingerAndScores.XgboostPredictedRarenessScore,
		NnPredictedRarenessScore:                      fingerAndScores.NnPredictedRarenessScore,
		OverallAverageRarenessScore:                   fingerAndScores.OverallAverageRarenessScore,
		IsLikelyDupe:                                  fingerAndScores.IsLikelyDupe,
	}

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send ProbeImage response: %w", err)
	}
	return nil
}

// UploadImage implements walletnode.RegisterArtwork.UploadImageWithThumbnail
func (service *RegisterArtwork) UploadImage(stream pb.RegisterArtwork_UploadImageServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return errors.Errorf("task not found %w", err)
	}

	image := service.Storage.NewFile()
	imageFile, err := image.Create()
	if err != nil {
		return errors.Errorf("open image file %q: %w", imageFile.Name(), err)
	}
	log.WithContext(ctx).WithField("filename", imageFile.Name()).Debug("UploadImageWithThumbnail request")

	imageWriter := bufio.NewWriter(imageFile)

	imageSize := int64(0)
	thumbnail := artwork.ThumbnailCoordinate{}
	hash := make([]byte, 0)

	err = func() error {
		defer imageFile.Close()
		defer imageWriter.Flush()

		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}
				return errors.Errorf("receive UploadImageWithThumbnail: %w", err)
			}

			if imagePiece := req.GetImagePiece(); imagePiece != nil {
				n, err := imageWriter.Write(imagePiece)
				imageSize += int64(n)
				if err != nil {
					return errors.Errorf("write to file %q: %w", imageFile.Name(), err)
				}
			} else {
				if metaData := req.GetMetaData(); metaData != nil {
					if metaData.Size != imageSize {
						return errors.Errorf("incomplete payload, send = %d receive=%d", metaData.Size, imageSize)
					}

					cordinates := metaData.GetThumbnail()
					if cordinates == nil {
						return errors.Errorf("no thumbnail coordinates")
					}
					thumbnail.TopLeftX = cordinates.TopLeftX
					thumbnail.TopLeftY = cordinates.TopLeftY
					thumbnail.BottomRightX = cordinates.BottomRightX
					thumbnail.BottomRightY = cordinates.BottomRightY

					if metaData.Hash == nil {
						return errors.Errorf("empty hash")
					}
					hash = metaData.Hash

					if metaData.Format != "" {
						if err := image.SetFormatFromExtension(metaData.Format); err != nil {
							return errors.Errorf("set format %s for file %s", metaData.Format, image.Name())
						}
					}
				}
			}
		}

		return nil
	}()

	if err != nil {
		return err
	}

	err = func() error {
		hasher := sha3.New256()
		f, fileErr := image.Open()
		if fileErr != nil {
			return errors.Errorf("open file %w", fileErr)
		}
		defer f.Close()

		if _, err := io.Copy(hasher, f); err != nil {
			return errors.Errorf("hash file failed %w", err)
		}
		hashFromPayload := hasher.Sum(nil)
		if !bytes.Equal(hashFromPayload, hash) {
			log.WithField("Filename", imageFile.Name()).Debugf("caculated from payload %s", base64.URLEncoding.EncodeToString(hashFromPayload))
			log.WithField("Filename", imageFile.Name()).Debugf("sent by client %s", base64.URLEncoding.EncodeToString(hash))
			return errors.Errorf("wrong hash")
		}

		return nil
	}()

	if err != nil {
		return err
	}

	previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err := task.UploadImageWithThumbnail(ctx, image, thumbnail)
	if err != nil {
		return err
	}

	resp := &pb.UploadImageReply{
		PreviewThumbnailHash: previewThumbnailHash[:],
		MediumThumbnailHash:  mediumThumbnailHash[:],
		SmallThumbnailHash:   smallThumbnailHash[:],
	}
	log.WithContext(ctx).Debugf("preview thumbnail hash: %x\n", previewThumbnailHash)
	log.WithContext(ctx).Debugf("medium thumbnail hash: %x\n", mediumThumbnailHash)
	log.WithContext(ctx).Debugf("small thumbnail hash: %x\n", smallThumbnailHash)

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send UploadImageAndThumbnail response: %w", err)
	}

	return nil
}

// SendSignedNFTTicket implements walletnode.RegisterArtwork.SendSignedNFTTicket
func (service *RegisterArtwork) SendSignedNFTTicket(ctx context.Context, req *pb.SendSignedNFTTicketRequest) (*pb.SendSignedNFTTicketReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("SignTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	registrationFee, err := task.GetRegistrationFee(ctx, req.NftTicket, req.CreatorSignature, req.Key1, req.Key2, req.EncodeFiles, req.EncodeParameters.Oti)
	if err != nil {
		return nil, errors.Errorf("get total storage fee %w", err)
	}

	rsp := pb.SendSignedNFTTicketReply{
		RegistrationFee: registrationFee,
	}

	return &rsp, nil
}

// SendPreBurntFeeTxid implements walletnode.RegisterArtwork.SendPreBurntFeeTxid
func (service *RegisterArtwork) SendPreBurntFeeTxid(ctx context.Context, req *pb.SendPreBurntFeeTxidRequest) (*pb.SendPreBurntFeeTxidReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("SendPreBurntFeeTxidRequest request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from meta data %w", err)
	}

	nftRegTxid, err := task.ValidatePreBurnTransaction(ctx, req.Txid)
	if err != nil {
		return nil, errors.Errorf("validate preburn transaction %w", err)
	}

	rsp := pb.SendPreBurntFeeTxidReply{
		NFTRegTxid: nftRegTxid,
	}
	return &rsp, nil
}

// Desc returns a description of the service.
func (service *RegisterArtwork) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtwork_ServiceDesc
}

// NewRegisterArtwork returns a new RegisterArtwork instance.
func NewRegisterArtwork(service *artworkregister.Service) *RegisterArtwork {
	return &RegisterArtwork{
		RegisterArtwork: common.NewRegisterArtwork(service),
	}
}
