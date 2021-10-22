package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	cli "github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/client"
	"github.com/pkg/errors"
	goahttp "goa.design/goa/v3/http"
)

func doHTTP(scheme, host string, timeout int, debug bool) *cli.Client {
	var (
		doer goahttp.Doer
	)
	{
		doer = &http.Client{Timeout: time.Duration(timeout) * time.Second}
		if debug {
			doer = goahttp.NewDebugDoer(doer)
		}
	}

	var (
		dialer *websocket.Dialer
	)
	{
		dialer = websocket.DefaultDialer
	}

	return cli.NewClient(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		true,
		dialer,
		nil,
	)
}

func printError(err error) {
	fmt.Println(err.Error())
}

func main() {

	var (
		regNFTTxid = flag.String("regnfttxid", "", "reg NFT txid")
		pastelID   = flag.String("pastelid", "", "valid user pastelid")
		passphrase = flag.String("passphase", "", "valid user passhrase")
	)
	//flag.Usage = usage
	flag.Parse()

	if *regNFTTxid == "" {
		usage()
		return
	}

	client := doHTTP("http", "localhost:8080", 100, true)

	endpoint := client.ArtworkGet()
	// create payload
	/*
		body := `{
			"user_passphrase": "passphrase",
			"user_pastelid": "jXYX6qAEiQKvTmpLZvNKQuphrzKPLACkx73zo9mE3B1kRQ7sjvnJMit3fHVPGNCY7REwdnVB2H42FaZoG8keAi"
		  }`
	*/
	bodyPayload := &cli.ArtworkGetRequestBody{
		UserPastelID:   *pastelID,
		UserPassphrase: *passphrase,
	}
	jsonBody, err := json.Marshal(bodyPayload)
	if err != nil {
		printError(errors.Wrap(err, "error marshal request body"))
		return
	}

	payload, err := cli.BuildArtworkGetPayload(string(jsonBody), *regNFTTxid)
	if err != nil {
		printError(errors.Wrap(err, "error creating payload"))
		return
	}
	if err != nil {
		printError(errors.Wrap(err, "error creating payload"))
		return
	}

	// send request
	res, err := endpoint(context.Background(), payload)
	if err != nil {
		printError(errors.Wrap(err, "error send data"))
		return
	}

	result, ok := res.(*artworks.ArtworkDetail)
	if !ok {
		printError(errors.New("invalid data"))
		return
	}

	ioutil.WriteFile("outputThumbnail", result.Thumbnail1, os.ModePerm)
}

func prettyPrint(s interface{}) {
	m, _ := json.MarshalIndent(s, "", "    ")
	fmt.Println(string(m))
}

func usage() {
	flag.PrintDefaults()
}
