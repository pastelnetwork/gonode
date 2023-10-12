package walletnode

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/pastelnetwork/gonode/common/storage/files"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/nftregister"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterNft represents grpc service for registration Nft.
type RegisterNft struct {
	pb.UnimplementedRegisterNftServer

	*common.RegisterNft
}

// Session implements walletnode.RegisterNftServer.Session()
func (service *RegisterNft) Session(stream pb.RegisterNft_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *nftregister.NftRegistrationTask

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewNftRegistrationTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)

	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := task.NetworkHandler.Session(ctx, req.IsPrimary); err != nil {
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

// AcceptedNodes implements walletnode.RegisterNftServer.AcceptedNodes()
func (service *RegisterNft) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.NetworkHandler.AcceptedNodes(ctx)
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

// GetDDDatabaseHash implements walletnode.RegisterNftServer.GetDupeDetectionDBHash()
func (service *RegisterNft) GetDDDatabaseHash(ctx context.Context, _ *pb.GetDBHashRequest) (*pb.DBHashReply, error) {
	hash, err := service.GetDupeDetectionDatabaseHash(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("GetDDDatabaseHash")
		return nil, fmt.Errorf("GetDDDatabaseHash: %w", err)
	}

	resp := &pb.DBHashReply{
		Hash: hash,
	}

	log.WithContext(ctx).WithField("hash len", len(hash)).Debug("DupeDetection DB hash returned")

	return resp, nil
}

// GetDDServerStats implements walletnode.RegisterNFTServer.GetDDServerStats()
func (service *RegisterNft) GetDDServerStats(ctx context.Context, _ *pb.DDServerStatsRequest) (*pb.DDServerStatsReply, error) {
	log.WithContext(ctx).Info("request for dd-server stats has been received")

	resp, err := service.GetDupeDetectionServerStats(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving dd-server stats from SN gRPC call")
	}

	stats := &pb.DDServerStatsReply{
		MaxConcurrent:  resp.MaxConcurrent,
		WaitingInQueue: resp.WaitingInQueue,
		Executing:      resp.Executing,
		Version:        resp.Version,
	}

	log.WithContext(ctx).WithField("stats", stats).Info("DD server stats returned")
	return stats, nil
}

// GetTopMNs implements walletnode.RegisterNFTServer.GetTopMNs()
func (service *RegisterNft) GetTopMNs(ctx context.Context, _ *pb.GetTopMNsRequest) (*pb.GetTopMNsReply, error) {
	log.WithContext(ctx).Debug("request for mn-top list has been received")

	mnTopList, err := service.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving mn-top list on the SN")
	}

	var mnList []string
	for _, mn := range mnTopList {
		mnList = append(mnList, mn.ExtAddress)
	}

	resp := &pb.GetTopMNsReply{
		MnTopList: mnList,
	}

	log.WithContext(ctx).WithField("mn-list", mnList).Debug("top mn-list has been returned")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterNftServer.ConnectTo()
func (service *RegisterNft) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.NetworkHandler.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")
	return resp, nil
}

// MeshNodes implements walletnode.RegisterNftServer.MeshNodes
func (service *RegisterNft) MeshNodes(ctx context.Context, req *pb.MeshNodesRequest) (*pb.MeshNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("MeshNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	meshedNodes := []types.MeshedSuperNode{}
	for _, node := range req.GetNodes() {
		meshedNodes = append(meshedNodes, types.MeshedSuperNode{
			NodeID: node.NodeID,
			SessID: node.SessID,
		})
	}

	err = task.NetworkHandler.MeshNodes(ctx, meshedNodes)
	return &pb.MeshNodesReply{}, err
}

// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
func (service *RegisterNft) SendRegMetadata(ctx context.Context, req *pb.SendRegMetadataRequest) (*pb.SendRegMetadataReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("SendRegMetadata request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	reqMetadata := &types.NftRegMetadata{
		BlockHash:       req.BlockHash,
		CreatorPastelID: req.CreatorPastelID,
		BlockHeight:     req.BlockHeight,
		Timestamp:       req.Timestamp,
		GroupID:         req.GroupId,
		CollectionTxID:  req.CollectionTxid,
	}

	err = task.SendRegMetadata(ctx, reqMetadata)

	return &pb.SendRegMetadataReply{}, err
}

// ProbeImage implements walletnode.RegisterNftServer.ProbeImage()
// As part of register NFT New Art Registration Workflow from
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step number 4.A begins here, with file reception
func (service *RegisterNft) ProbeImage(stream pb.RegisterNft_ProbeImageServer) (retErr error) {
	ctx := stream.Context()

	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenProbeImage")
		retErr = recErr
	})

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	image := service.Storage.NewFile()

	// Do file receiving
	receiveFunc := func() error {
		file, subErr := image.Create()
		if subErr != nil {
			return errors.Errorf("open file %q: %w", file.Name(), subErr)
		}
		defer file.Close()
		log.WithContext(ctx).WithField("filename", file.Name()).Debug("ProbeImage request")

		wr := bufio.NewWriter(file)
		defer wr.Flush()

		for {
			req, subErr := stream.Recv()
			if subErr != nil {
				if subErr == io.EOF {
					break
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}
				return errors.Errorf("receive ProbeImage: %w", subErr)
			}

			if _, subErr := wr.Write(req.Payload); subErr != nil {
				return errors.Errorf("write to file %q: %w", file.Name(), subErr)
			}
		}

		return nil
	}
	err = receiveFunc()
	if err != nil {
		return err
	}

	err = image.UpdateFormat()
	if err != nil {
		return errors.Errorf("add image format: %w", err)
	}

	exists, err := task.HashExists(ctx, image)
	if exists || err != nil {
		resp := &pb.ProbeImageReply{
			IsExisting:                        exists,
			CompressedSignedDDAndFingerprints: []byte{},
		}

		if err != nil {
			log.WithContext(ctx).WithError(err).Error("hash exist check failed")
			resp.ErrString = err.Error()
		}

		return stream.SendAndClose(resp)
	}

	//task.ProbeImage will begin Step 3,
	// Call dd-service to generate near duplicate fingerprints and dupe-detection info from re-sampled image (img1-r)
	compressedDDFingerAndScores, err := task.ProbeImage(ctx, image)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("probe image call failed")
		return err
	}

	// Ends up at Step 4.B.5
	// Return base64’ed (and compressed) dd_and_fingerprints and its signature back to the WalletNode over the open gRPC connection
	resp := &pb.ProbeImageReply{
		CompressedSignedDDAndFingerprints: compressedDDFingerAndScores,
	}

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send ProbeImage response: %w", err)
	}
	return nil
}

// UploadImage implements walletnode.RegisterNft.UploadImageWithThumbnail
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration Step 7
// All SuperNodes - Generate preview thumbnail from the image and calculate its hash
func (service *RegisterNft) UploadImage(stream pb.RegisterNft_UploadImageServer) (retErr error) {
	ctx := stream.Context()
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenUploadImage")
		retErr = recErr
	})

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
	thumbnail := files.ThumbnailCoordinate{}
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

					coordinates := metaData.GetThumbnail()
					if coordinates == nil {
						return errors.Errorf("no thumbnail coordinates")
					}
					thumbnail.TopLeftX = coordinates.TopLeftX
					thumbnail.TopLeftY = coordinates.TopLeftY
					thumbnail.BottomRightX = coordinates.BottomRightX
					thumbnail.BottomRightY = coordinates.BottomRightY

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
			log.WithContext(ctx).WithField("Filename", imageFile.Name()).Debugf("caculated from payload %s", base64.URLEncoding.EncodeToString(hashFromPayload))
			log.WithContext(ctx).WithField("Filename", imageFile.Name()).Debugf("sent by client %s", base64.URLEncoding.EncodeToString(hash))
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

// SendSignedNFTTicket implements walletnode.RegisterNft.SendSignedNFTTicket
// Step 11.B ALL SuperNode - Validate signature and IDs in the ticket
// Step 12. ALL SuperNode - Calculate final Registration Fee
func (service *RegisterNft) SendSignedNFTTicket(ctx context.Context, req *pb.SendSignedNFTTicketRequest) (retRes *pb.SendSignedNFTTicketReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedNFTTicket")
		retErr = recErr
	})

	log.WithContext(ctx).WithField("req", req).Debug("SignTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	registrationFee, err := task.GetNftRegistrationFee(ctx, req.NftTicket, req.CreatorSignature, req.Label, req.RqFiles, req.DdFpFiles, req.EncodeParameters.Oti)
	if err != nil {
		return nil, errors.Errorf("get total storage fee %w", err)
	}

	rsp := pb.SendSignedNFTTicketReply{
		RegistrationFee: registrationFee,
	}

	return &rsp, nil
}

// SendPreBurntFeeTxid implements walletnode.RegisterNft.SendPreBurntFeeTxid
// Receives Step 14, Starts Step 15
func (service *RegisterNft) SendPreBurntFeeTxid(ctx context.Context, req *pb.SendPreBurntFeeTxidRequest) (retRes *pb.SendPreBurntFeeTxidReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicSendPreBurntFeeTxid")
		retErr = recErr
	})

	log.WithContext(ctx).WithField("req", req).Debug("SendPreBurntFeeTxidRequest request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from meta data %w", err)
	}
	// Step 15-17
	err = task.ValidatePreBurnTransaction(ctx, req.Txid)
	if err != nil {
		return nil, errors.Errorf("validate preburn transaction %w", err)
	}

	// Step 18 - 19
	nftRegTxid, err := task.ActivateAndStoreNft(ctx)
	if err != nil {
		return nil, errors.Errorf("activate and store NFT %w", err)
	}

	rsp := pb.SendPreBurntFeeTxidReply{
		NFTRegTxid: nftRegTxid,
	}
	return &rsp, nil
}

// Desc returns a description of the service.
func (service *RegisterNft) Desc() *grpc.ServiceDesc {
	return &pb.RegisterNft_ServiceDesc
}

// NewRegisterNft returns a new RegisterNft instance.
func NewRegisterNft(service *nftregister.NftRegistrationService) *RegisterNft {
	return &RegisterNft{
		RegisterNft: common.NewRegisterNft(service),
	}
}
