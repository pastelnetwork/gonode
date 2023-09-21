package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
)

var (
	defaultPath                   = configurer.DefaultPath()
	defaultPastelConfigFile       = filepath.Join(defaultPath, "pastel.conf")
	takeLastNBlocks         int32 = 5000
	maxTXIDsToDownload      int32 = 100
	pastelID                      = "jXZMSxS5w9UakpVMAs2vihcCVQ4fBrPsSriXmNqTq2nvK4awXvaP9hZJYL1eJ4o9y3jpvoGghVUQyvsU7Q64Jp"

	initialDelay = 1 * time.Second
	maxRetries   = 2

	/*mustTxIDs = []string{
		"0fc8340d7c5d18b0b8f98eea3ef546ab52c39e6796912c092ff8cf29360f115d",
		"3d6ddee1af1802b04a451543d54fd461fba75b5414e4d065733ebf1498126f1f",
		"f555d68c4e300353a9e30be930fc1d772687b039a861fea2b422da75671d0efd",
		"4047cdbfcd97836a16fc326d0e01d5d1ba99e6bf69fe0339144a8bb84b04bfbb",
		"0d74a8ef52ca917e879a1acd8b7991fa7e60b7c7599791f5bfb099357c63cca7",
		"3e0b6f63d9981e720e6954b4e9c18319f61c889ba10e739cadb73764b44a39e3",
		"93ae296f5f720707711d18181e25e76d3c81c76695a4b928a093e1a366e3b889",
		"14ee28f838d83af762febd1471d8d06f5f6cff018254523a890ef119d2c4e18e",
		"c45a0f033f1531e8c280104d6014b9b26e01a6e9fcf73a44daab62b38010fb27",
		"615bb44c984a7f59510cde05c7c72d9b06759928f56c69d8a185a526c00b5600",
		"2c9bfd03901053f43a2c39efd248c53a22c5c5ca46240953f525a7b5acf4324f",
		"15d98e7096903593d60f0d33287709dd6a387b6c01b3d88ce36a7b2ba99d7abc",
		"2526fb14e3649cc43f07367bc2d4d324a8aba1de807ee853b21673533d120d03",
		"7c0e14220b6e6dadbddcd2e9203867959d7ec68ed86345be58ed0c800faeb122",
		"abc7d04c0ee978ee3e284c7e0e9a87be3cc6dda4c9ea30b2f124b865bf409c09",
		"2b4ac1940c324516bef491a0d96955bb9fb96f27858fc621aa56aaac391edb49",
		"357f340c6b373ee880722591df1837c9bac1d83283e8ede7cadad5b9ec09f8c5",
		"3ffb9755695020bb4788948373b36b3144d5f0c3c885b4946f1522fb02e324f5",
		"574216a11d9c391018aac1be4e49973704376e9879cd08075bdfe6d325fb5820",
		"b5a991b540be614030c7c8cd4c227146c73486908be538c327df3c2c0880bd9b",
		"7d8341d196d7d248d431c31987ff639236654095a4c1eb9435187410958c4fe1",
		"9eda1409a1078411263ea6766a46367eb154d67521f08977de5f247240e99782",
		"15d98e7096903593d60f0d33287709dd6a387b6c01b3d88ce36a7b2ba99d7abc",
		"4de7cbc751b711ba0556f99a934645e164ea86071ebdd7718f0909aeedd785b2",
	}*/
)

func main() {
	ctx := context.Background()
	logFile, err := os.OpenFile("requests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	pastelClient, err := newPastelClient()
	if err != nil {
		logger.Println("unable to initialize pastel client")
	}
	logger.Println("Pastel-client has been initialised")

	regTxIDs, err := getRegTicketsTxIDs(ctx, pastelClient, logger)
	if err != nil {
		logger.Printf("error retrieving get reg tickets tx-ids:%s\n", err.Error())
	}

	logger.Printf("Total Tickets in the batch: %d\n", len(regTxIDs))

	startingTime := time.Now()
	logger.Printf("Batch starts at:%s\n", startingTime)

	var success int32
	var failure int32
	var wg sync.WaitGroup

	for _, txid := range regTxIDs {
		wg.Add(1)
		go func(txid string) {
			defer wg.Done()

			// Each goroutine waits for 100ms before starting
			time.Sleep(200 * time.Millisecond)

			logger.Printf("downloading:%s\n", txid)

			if err := doSenseDownload(txid); err != nil {
				logger.Printf("error downloading sense file: %s \n", err.Error())
				atomic.AddInt32(&failure, 1)
				return
			}

			atomic.AddInt32(&success, 1)
			logger.Printf("downloaded successfully:%s\n", txid)
		}(txid)
	}

	wg.Wait() // Wait for all goroutines to complete

	totalElapsed := time.Since(startingTime)
	logger.Printf("Total successful:%d\n", success)
	logger.Printf("Total failures:%d\n", failure)
	logger.Printf("Total Elapsed Time:%s\n", totalElapsed)

}

func newPastelClient() (pastel.Client, error) {
	var pastelConfigFile string
	cli.NewFlag("pastel-config-file", &pastelConfigFile).SetUsage("Set `path` to the pastel config file.").SetValue(defaultPastelConfigFile)

	var config = &pastel.Config{}

	if err := configurer.ParseFile(defaultPastelConfigFile, config); err != nil {
		return nil, errors.Errorf("unable to parse pastel.conf:%v\n", err)
	}
	pastelClient := pastel.NewClient(config, config.BurnAddress())

	return pastelClient, nil
}

func getSenseTicket(ctx context.Context, pastelClient pastel.Client, regTXID string) (*pastel.ActionRegTicket, *pastel.APISenseTicket, error) {
	regTicket, err := pastelClient.ActionRegTicket(ctx, regTXID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find action act by action reg txid: %w", err)
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
		return nil, nil, errors.New("not a sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode reg ticket: %w", err)
	}
	regTicket.ActionTicketData.ActionTicketData = *decTicket

	senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return &regTicket, senseTicket, nil
}

func getRegTicketsTxIDs(ctx context.Context, pastelClient pastel.Client, logger *log.Logger) ([]string, error) {
	currentBlockCount, err := pastelClient.GetBlockCount(ctx)
	if err != nil {
		errors.Errorf("unable to get block count")
	}
	logger.Printf("block count:%d\n", currentBlockCount)

	var regTxIDs []string

	actTickets, err := pastelClient.ActionActivationTicketsFromBlockHeight(ctx, uint64(currentBlockCount-takeLastNBlocks))
	if err != nil {
		logger.Println("get registered ticket - exit runTask")
		return nil, err
	}

	if len(actTickets) == 0 {
		return nil, errors.New("no tickets found")
	}

	sort.Slice(actTickets, func(i, j int) bool {
		return actTickets[i].Height < actTickets[j].Height
	})

	for j := 0; j < len(actTickets); j++ {
		regTicket, _, err := getSenseTicket(ctx, pastelClient, actTickets[j].ActTicketData.RegTXID)
		if err != nil {
			continue
		}

		regTxIDs = append(regTxIDs, regTicket.TXID)
	}

	log.Printf("Total Tickets in the batch: %d\n", len(regTxIDs))

	if len(regTxIDs) > int(maxTXIDsToDownload) {
		regTxIDs = regTxIDs[:maxTXIDsToDownload]
	}
	log.Printf("Got Tickets in the batch: %d\n", len(regTxIDs))

	//	regTxIDs = append(regTxIDs, mustTxIDs...)

	log.Printf("Total Tickets in the batch inlcuding must TxIDs: %d\n", len(regTxIDs))

	return regTxIDs, nil
}

func doSenseDownload(txID string) error {
	url := fmt.Sprintf("http://localhost:18080/openapi/sense/download?pid=%s&txid=%s", pastelID, txID)
	method := "GET"

	client := &http.Client{
		Timeout: 25 * time.Minute,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "passphrase")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading sense file: %w - txid: %s", err, txID)
	}

	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d - txid: %s", res.StatusCode, txID)
	}
	log.Printf("downloaded successfully: %s", txID)

	return nil
}
