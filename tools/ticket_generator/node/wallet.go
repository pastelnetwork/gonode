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
	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"

	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goaclient "github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/client"
	goahttp "goa.design/goa/v3/http"
)

type Wallet struct {
	config config.WalletNode

	goaClient *goaclient.Client
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

	log.WithContext(ctx).Info("calling upload image endpoint")
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
// Returns art
func (wallet *Wallet) RegisterNFT(ctx context.Context, imgID string) (string, error) {

}

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
	return &Wallet{
		config:    config,
		goaClient: goaClient,
	}
}
