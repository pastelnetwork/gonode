package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func registerTicket(client *cli.Client, imageId string) (string, error) {
	endpoint := client.Register()
	payloadJson := fmt.Sprintf(
		`{
		  "artist_name": "Leonardo da Vinci",
		  "artist_pastelid": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS",
		  "artist_pastelid_passphrase": "passphrase",
		  "artist_website_url": "https://www.leonardodavinci.net",
		  "description": "The Mona Lisa is an oil painting by Italian artist, inventor, and writer Leonardo da Vinci. Likely completed in 1506, the piece features a portrait of a seated woman set against an imaginary landscape.",
		  "green": false,
		  "image_id": "%s",
		  "issued_copies": 1,
		  "keywords": "Renaissance, sfumato, portrait",
		  "maximum_fee": 100,
		  "name": "Mona Lisa",
		  "royalty": 12,
		  "series_name": "Famous artist",
		  "spendable_address": "PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j",
		  "thumbnail_coordinate": {
		    "bottom_right_x": 640,
		    "bottom_right_y": 480,
		    "top_left_x": 0,
		    "top_left_y": 0
		  },
		  "youtube_url": "https://www.youtube.com/watch?v=0xl6Ufo4ZX0"
		}`, imageId)
	payload, err := cli.BuildRegisterPayload(payloadJson)
	if err != nil {
		return "", errors.Wrap(err, "error creating payload")
	}

	// send request
	res, err := endpoint(context.Background(), payload)
	if err != nil {
		return "", errors.Wrap(err, "error send data")
	}

	result, ok := res.(*artworks.RegisterResult)
	if !ok {
		return "", errors.New("invalid data")
	}

	fmt.Printf("Task id : %s\n", result.TaskID)
	return result.TaskID, nil
}

func watchStatus(client *cli.Client, taskId string) error {
	endpoint := client.RegisterTaskState()

	// create payload
	payload, err := cli.BuildRegisterTaskStatePayload(taskId)
	if err != nil {
		return errors.Wrap(err, "error creating payload")
	}

	// create stream
	data, err := endpoint(context.Background(), payload)
	if err != nil {
		return errors.Wrap(err, "error creating stream")
	}

	switch stream := data.(type) {
	case artworks.RegisterTaskStateClientStream:
		// result streaming
		for {
			p, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Task existed")
				break
			}
			if err != nil {
				return errors.Wrap(err, "error reading from stream")
			}
			prettyPrint(p)
		}
	}

	return nil
}

func main() {
	var (
		imageId = flag.String("imageid", "", "image id")
	)
	flag.Usage = usage
	flag.Parse()
	if len(*imageId) == 0 {
		usage()
		return
	}
	client := doHTTP("http", "localhost:8080", 100, true)

	fmt.Println("****************Register Ticket*************************")
	taskId, err := registerTicket(client, *imageId)
	if err != nil {
		fmt.Println(fmt.Errorf("error register ticket: %v", err))
		os.Exit(-1)
		return
	}

	fmt.Println("****************Watch status*************************")
	err = watchStatus(client, taskId)
	if err != nil {
		fmt.Println(fmt.Errorf("error watch status: %v", err))
		os.Exit(-1)
		return
	}
}

func prettyPrint(s interface{}) {
	m, _ := json.MarshalIndent(s, "", "    ")
	fmt.Println(string(m))
}

func usage() {
	fmt.Fprintf(os.Stderr, `%s is a command line client to simulate ui

Usage:
    %s [-imageid IMAGEID]

    -imageid IMAGEID:  image id

Commands:
%s
Additional help:
   None

Example:
%s -imageid 12345678
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
