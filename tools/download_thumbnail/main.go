package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	json "github.com/json-iterator/go"

	"github.com/gorilla/websocket"
	cli "github.com/pastelnetwork/gonode/walletnode/api/gen/http/nft/client"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
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

	endpoint := client.NftGet()
	// create payload
	/*
		body := `{
			"user_passphrase": "passphrase",
			"user_pastelid": "jXYX6qAEiQKvTmpLZvNKQuphrzKPLACkx73zo9mE3B1kRQ7sjvnJMit3fHVPGNCY7REwdnVB2H42FaZoG8keAi"
		  }`
	*/
	bodyPayload := &cli.NftGetRequestBody{
		UserPastelID:   *pastelID,
		UserPassphrase: *passphrase,
	}
	jsonBody, err := json.Marshal(bodyPayload)
	if err != nil {
		printError(errors.Wrap(err, "error marshal request body"))
		return
	}

	payload, err := cli.BuildNftGetPayload(string(jsonBody), *regNFTTxid)
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

	result, ok := res.(*nft.NftDetail)
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
