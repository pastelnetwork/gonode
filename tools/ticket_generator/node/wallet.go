package node

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"path"
	"time"

	"mime/multipart"

	"github.com/go-errors/errors"
	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"

	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goaclient "github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/client"
	goahttp "goa.design/goa/v3/http"
)

type Wallet struct {
	config config.WalletNode

	goaClient    *goaclient.Client
	pastelClient pastel.Client
}

// encode multipart request
func UploadImageEncoderFunc(file string) func(*multipart.Writer, *artworks.UploadImagePayload) error {
	return func(writer *multipart.Writer, p *artworks.UploadImagePayload) error {
		// &{map[Content-Disposition:[form-data; name="file"; filename="2cd43b_b462137ddab34c21a3741c5713850d2b~mv2.png"] Content-Type:[image/png]] 0xc000391380  map[] {0xc0003cf730} 0 0 <nil> <nil>}

		filename := path.Base(file)
		extension := path.Ext(file)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "image/"+extension)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))

		part, err := writer.CreatePart(h)
		if err != nil {
			return errors.Errorf("failed to create file part: %w", err)
		}
		if _, err := part.Write(p.Bytes); err != nil {
			return errors.Errorf("failed  to write file filed: %w", err)
		}
		if p.Filename != nil {
			if err := writer.WriteField("filename", *p.Filename); err != nil {
				return errors.Errorf("failed to write filename filed: %w", err)
			}
		}
		return nil
	}
}

// UploadNFT upload the image and return imagei d
func (wallet *Wallet) UploadNFT(ctx context.Context, nft string) (string, error) {
	file, err := ioutil.ReadFile(nft)
	if err != nil {
		return "", errors.Errorf("failed to open image %s: %w", nft, err)
	}
	body := &artworks.UploadImagePayload{
		Bytes: file,
	}

	ep := wallet.goaClient.UploadImage(UploadImageEncoderFunc(nft))

	rsp, err := ep(ctx, body)
	if err != nil {
		log.WithContext(ctx).Debugf("%v", err)
		return "", errors.Errorf("failed to upload image %s: %w", nft, err)
	}

	result, ok := rsp.(*artworks.Image)
	if !ok {
		return "", errors.Errorf("invalid response")
	}

	return result.ImageID, nil
}

// RegisterArt reigsters the art that already uploaded and identified by imgID
func (wallet *Wallet) RegisterNFT(ctx context.Context, imgID string) (string, error) {
	body := &artworks.RegisterPayload{
		ImageID:                  imgID,
		Name:                     "NFTName",
		Description:              ptrFromStr("Description"),
		Keywords:                 ptrFromStr("Keywords"),
		SeriesName:               ptrFromStr("SerSeriesName"),
		IssuedCopies:             5,
		YoutubeURL:               ptrFromStr("https://www.youtube.com/watch?v=0xl6Ufo4ZX0"),
		ArtistPastelID:           wallet.config.Artist.PastelID,
		ArtistPastelIDPassphrase: wallet.config.Artist.Passphrase,
		ArtistName:               "ArtistName",
		ArtistWebsiteURL:         ptrFromStr("https://awesomenft.com"),
		SpendableAddress:         wallet.config.Artist.SpendableAddr,
		MaximumFee:               100,
		Royalty:                  12,
		Green:                    false,
		ThumbnailCoordinate: &artworks.Thumbnailcoordinate{
			TopLeftX:     640,
			TopLeftY:     480,
			BottomRightX: 0,
			BottomRightY: 0,
		},
	}

	ep := wallet.goaClient.Register()
	rsp, err := ep(ctx, body)
	if err != nil {
		return "", errors.Errorf("failed to register task: %w", err)
	}

	result, ok := rsp.(*artworks.RegisterResult)
	if !ok {
		return "", errors.Errorf("invalid response")
	}

	return result.TaskID, nil
}

// GetBlockCount returns the current total number of block in the chain
func (wallet *Wallet) GetBlockCount(ctx context.Context) (int32, error) {
	blockCount, err := wallet.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return 0, errors.Errorf("failed to call getblockcount: %w", err)
	}
	return blockCount, nil
}

// GetActTicketByCreatorHeigh return txid of reg-act with creator-height
func (wallet *Wallet) FindActTicketByCreatorHeight(ctx context.Context, creatorHeight int32) (*pastel.ActTicket, error) {
	tickets, err := wallet.pastelClient.FindActTicketByCreatorHeight(ctx, creatorHeight)
	if err != nil {
		return nil, errors.Errorf("failed to call tickets find act %d: %w", creatorHeight, err)
	}
	if len(tickets) == 0 {
		return nil, errors.Errorf("act-ticket at creator-height %d not found", creatorHeight)
	}
	return &tickets[0], nil
}

// SubscribeTaskState returns register's task state
func (wallet *Wallet) SubscribeTaskState(ctx context.Context, taskID string) (*TaskStateReceiver, error) {
	// client := newGoaClient(wallet.config.PastelAPI.Hostname, wallet.config.PastelAPI.Port)
	// ep := client.RegisterTaskState()
	ep := wallet.goaClient.RegisterTaskState()

	body := &artworks.RegisterTaskStatePayload{
		TaskID: taskID,
	}

	rsp, err := ep(ctx, body)
	if err != nil {
		return nil, errors.Errorf("failed to call register/task/%s/state: %w", taskID, err)
	}

	stream, ok := rsp.(artworks.RegisterTaskStateClientStream)
	if !ok {
		return nil, errors.Errorf("invalid rsp")
	}

	return NewTaskStateReceiver(ctx, stream), nil
}

// NewWallet returns a client to communicate with both wallet's goNode and wallet's cNode
func NewWallet(config config.WalletNode) *Wallet {
	doer := goahttp.NewDebugDoer(
		&http.Client{Timeout: time.Minute},
	)

	goaClient := goaclient.NewClient(
		"http",
		fmt.Sprintf("%s:%d", config.RestAPI.Hostname, config.RestAPI.Port),
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		true,
		websocket.DefaultDialer,
		nil,
	)

	pastelCfg := &pastel.Config{
		Hostname: ptrFromStr(config.PastelAPI.Hostname),
		Port:     ptrFromInt(config.PastelAPI.Port),
		Username: ptrFromStr(config.PastelAPI.Username),
		Password: ptrFromStr(config.PastelAPI.Passphrase),
	}
	pastelClient := pastel.NewClient(pastelCfg)

	return &Wallet{
		config:       config,
		goaClient:    goaClient,
		pastelClient: pastelClient,
	}
}
