package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
)

var (
	defaultPath                   = configurer.DefaultPath()
	defaultPastelConfigFile       = filepath.Join(defaultPath, "pastel.conf")
	noOfDownloadReqs        int32 = 1000
	pastelID                      = "jXa6QiopivJLer8G65QsxwQmGELi1w6mbNXvrrYTvsddVE5BT57LtNCZ2SCmWStvLwWWTkuAFPsRREytgG62YX"

	initialDelay = 1 * time.Second
	maxRetries   = 5
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

	regTxIDs, err := getCascadeRegTicketsTxIDs(ctx, pastelClient, logger)
	if err != nil {
		logger.Printf("error retrieving get reg tickets tx-ids:%s\n", err.Error())
	}

	logger.Printf("Total Tickets in the batch: %d\n", len(regTxIDs))

	startingTime := time.Now().UTC()
	logger.Printf("Batch starts at:%s\n", startingTime)
	var success = 0
	var failure = 0
	for _, txid := range regTxIDs {
		logger.Printf("downloading:%s\n", txid)

		if err := doCascadeDownload(txid); err != nil {
			failure++
			continue
		}
		success++

		logger.Printf("downloaded successfully:%s\n", txid)
	}

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

func getCascadeTicket(ctx context.Context, pastelClient pastel.Client, regTXID string) (*pastel.ActionRegTicket, *pastel.APICascadeTicket, error) {
	regTicket, err := pastelClient.ActionRegTicket(ctx, regTXID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find action act by action reg txid: %w", err)
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeCascade {
		return nil, nil, errors.New("not a sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode reg ticket: %w", err)
	}
	regTicket.ActionTicketData.ActionTicketData = *decTicket

	cascadeTicket, err := regTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return &regTicket, cascadeTicket, nil
}

func getCascadeRegTicketsTxIDs(ctx context.Context, pastelClient pastel.Client, logger *log.Logger) ([]string, error) {
	currentBlockCount, err := pastelClient.GetBlockCount(ctx)
	if err != nil {
		errors.Errorf("unable to get block count")
	}
	logger.Printf("block count:%d\n", currentBlockCount)

	var regTxIDs []string
	for i := currentBlockCount; i > currentBlockCount-noOfDownloadReqs; i-- {
		actTickets, err := pastelClient.ActionActivationTicketsFromBlockHeight(ctx, uint64(currentBlockCount))
		if err != nil {
			logger.Println("get registered ticket - exit runTask")
			continue
		}

		if len(actTickets) == 0 {
			continue
		}

		sort.Slice(actTickets, func(i, j int) bool {
			return actTickets[i].Height < actTickets[j].Height
		})

		for j := 0; j < len(actTickets); j++ {
			regTicket, _, err := getCascadeTicket(ctx, pastelClient, actTickets[j].ActTicketData.RegTXID)
			if err != nil {
				continue
			}

			regTxIDs = append(regTxIDs, regTicket.TXID)
		}
	}

	return regTxIDs, nil
}

func doCascadeDownload(txID string) error {
	url := fmt.Sprintf("http://localhost:18080/openapi/cascade/download?pid=%s&txid=%s", pastelID, txID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "passphrase")

	var res *http.Response
	delay := initialDelay
	for retries := 0; retries < maxRetries; retries++ {
		res, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(delay)
		delay *= 2
	}
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	return nil
}
