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
		pastelID   = flag.String("pastelid", "", "creator or owner's Pastel ID")
		passphrase = flag.String("passphrase", "", "passphrase associated with Pastel ID")
	)
	//flag.Usage = usage
	flag.Parse()

	if len(*regNFTTxid) == 0 || len(*pastelID) == 0 || len(*passphrase) == 0 {
		usage()
		return
	}

	client := doHTTP("http", "localhost:8080", 100, true)

	endpoint := client.Download()
	// create payload
	payload, err := cli.BuildDownloadPayload(*regNFTTxid, *pastelID, *passphrase)
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

	result, ok := res.(*artworks.DownloadResult)
	if !ok {
		printError(errors.New("invalid data"))
		return
	}

	ioutil.WriteFile(*regNFTTxid, result.File, os.ModePerm)
}

func prettyPrint(s interface{}) {
	m, _ := json.MarshalIndent(s, "", "    ")
	fmt.Println(string(m))
}

func usage() {
	fmt.Fprintf(os.Stderr, `%s is a command line client to simulate ui

 Usage:
     %s -regnfttxid <RegNFTTXID> -pastelid <PastelID> -passphrase <passphrase>

     -regnfttxid    string:  reg NFT tx id
	 -pastelid      string:  creator or owner's Pastel ID
	 -passphrase    string:  passphrase associated with Pastel ID

 Commands:
 %s
 Additional help:
    None

 Example:
 %s -regnfttxid 7602bf7c95487521c8a8c56a0a45a03438fcfd607a0089942382afd5fd3867c5 -pastelid jXXsFY5CMSFy5T1cHFiqTsK2B5kMv93zP48G2FdfPWuPP6p2jH6c1JrNGL3VJHDxFvq6eZSsWppC3191sbG3oD -passphrase passphrase
 `, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
