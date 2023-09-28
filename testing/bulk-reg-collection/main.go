package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"

	json "github.com/json-iterator/go"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	initialDelay = 1 * time.Second
	maxRetries   = 5
	timeoutAfter = 1000
	NoOfRequests = 1
)

type collectionRegPayload struct {
	CollectionName                                 string   `json:"collection_name"`
	ItemType                                       string   `json:"item_type"`
	ListOfPastelidsOfAuthorizedContributors        []string `json:"list_of_pastelids_of_authorized_contributors"`
	MaxCollectionEntries                           int      `json:"max_collection_entries"`
	NoOfDaysToFinalizeCollection                   int      `json:"no_of_days_to_finalize_collection"`
	MaxPermittedOpenNsfwScore                      float64  `json:"max_permitted_open_nsfw_score"`
	MinimumSimilarityScoreToFirstEntryInCollection float64  `json:"minimum_similarity_score_to_first_entry_in_collection"`
	SpendableAddress                               string   `json:"spendable_address"`
	AppPastelid                                    string   `json:"app_pastelid"`
}

type collectionRegResponse struct {
	TaskID string `json:"task_id"`
}

type result struct {
	ID      string
	Elapsed time.Duration
	Error   error
}

func doCollectionRegRequest(payload collectionRegPayload) (string, error) {
	url := "http://localhost:18080/collection/register"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return "", err
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
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var resp collectionRegResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return resp.TaskID, nil
}

func doTaskState(taskID string, expectedValue string, logger *log.Logger) error {
	url := fmt.Sprintf("ws://127.0.0.1:18080/collection/%s/state", taskID)

	ctx := context.Background()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer c.Close()

	done := make(chan bool)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			resp := make(map[string]interface{})
			if err := json.Unmarshal(message, &resp); err != nil {
				continue
			}

			if val, ok := resp["status"]; ok {
				log.Printf("ws expected key: status - response: %s\n", val)
				if val == expectedValue {
					done <- true
				} else if val == "Task Rejected" {
					done <- false
				}

				logger.Printf("taskID:%s, block_count:%s  status - response: %s\n", taskID, getBlockCount(), val)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return errors.New("context cancelled")
	case <-time.After(timeoutAfter * time.Minute):
		return errors.New("timeout")
	case val := <-done:
		if val {
			return nil
		}

		return errors.New("task failed (Request Rejected), please see container logs")
	}
}

func getBlockCount() string {
	pastelCli, err := getPastelCliPath()
	if err != nil {
		return ""
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	cmd := exec.Command(pastelCli, "getblockcount")
	cmd.Dir = homeDir

	res, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.Replace(string(res), "\n", "", 1)
}

func main() {
	var mu sync.Mutex
	logFile, err := os.OpenFile("requests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	var results []result
	var wg sync.WaitGroup

	start := time.Now().UTC()
	count := 1
	taskIDs := make(map[string]time.Time)
	rand.Seed(time.Now().UTC().UnixNano())

	for count <= NoOfRequests {
		startReq := time.Now().UTC()

		req := collectionRegPayload{
			CollectionName:                                 "test-coll-sense-items",
			ItemType:                                       "sense",
			ListOfPastelidsOfAuthorizedContributors:        []string{"jXa6QiopivJLer8G65QsxwQmGELi1w6mbNXvrrYTvsddVE5BT57LtNCZ2SCmWStvLwWWTkuAFPsRREytgG62YX", "jXYjBDCtQk6c77DxTFVM28Cwuy6JkRGgbrhvES9paHZQEyg4ocD4a7GBs9XBk9na3fs7zcmcgZv77ugU4aoU8d"},
			MaxCollectionEntries:                           5,
			NoOfDaysToFinalizeCollection:                   rand.Intn(7) + 1,
			MaxPermittedOpenNsfwScore:                      rand.Float64(),
			MinimumSimilarityScoreToFirstEntryInCollection: rand.Float64(),
			SpendableAddress:                               "tPiW2QcxZNKXHEeYXn6kyyZbJuUsqzpThJC",
			AppPastelid:                                    "jXa6QiopivJLer8G65QsxwQmGELi1w6mbNXvrrYTvsddVE5BT57LtNCZ2SCmWStvLwWWTkuAFPsRREytgG62YX",
		}

		logger.Printf("payload for collection-name:%s, request:%d, payload: %v", req.CollectionName, count, req)

		taskID, err := doCollectionRegRequest(req)
		if err != nil {
			logger.Printf("Request to collection registration failed:%v\n", err)
		}
		logger.Printf("collection task initiated:%s, request-count:%d\n", taskID, count)

		taskIDs[taskID] = startReq
		count++
	}

	count = 1
	for taskID, startReq := range taskIDs {
		wg.Add(1)

		go func(count int, tID string) {
			defer wg.Done()

			logger.Printf("subscribing to task state:%s, request-count:%d\n", tID, count)
			if err = doTaskState(tID, "Task Completed", logger); err != nil {
				logger.Printf("Request to task state has been failed:%v\n", err)
			}

			results = appendResults(mu, results, result{
				ID:      fmt.Sprintf("request%d", count),
				Elapsed: time.Since(startReq),
				Error:   err,
			})
		}(count, taskID)

		count++
	}
	wg.Wait()

	totalElapsed := time.Since(start)

	var successes, failures int
	var totalReqTime time.Duration
	for _, result := range results {
		if result.Error != nil {
			failures++
			logger.Printf("Request %s failed after %v: %v\n", result.ID, result.Elapsed, result.Error)
		} else {
			successes++
			logger.Printf("Request %s succeeded after %v\n", result.ID, result.Elapsed)
		}
		totalReqTime += result.Elapsed
	}

	avgReqTime := totalReqTime / time.Duration(len(results))

	logger.Printf("Total time for all requests: %v\n", totalElapsed.String())
	logger.Printf("Average time per request: %v\n", avgReqTime)
	logger.Printf("Total successes: %d\n", successes)
	logger.Printf("Total failures: %d\n", failures)
}

func appendResults(mu sync.Mutex, results []result, result result) (res []result) {
	mu.Lock()
	defer mu.Unlock()

	results = append(results, result)

	return results
}

func getItemType(count int) string {
	if count == 0 {
		return "sense"
	}

	if count%2 == 0 {
		return "sense"
	}

	return "nft"
}

func getPastelCliPath() (path string, err error) {
	//create command
	findCmd := exec.Command("find", ".", "-print")
	grepCmd := exec.Command("grep", "-x", "./pastel/pastel-cli")

	findCmd.Dir, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	//make a pipe and set the input and output to reader and writer
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	findCmd.Stdout = writer
	grepCmd.Stdin = reader

	//cache the output of "grep" to memory
	grepCmd.Stdout = &buf

	//starting the commands
	findCmd.Start()
	grepCmd.Start()

	//waiting for commands to complete and close the reader & writer
	findCmd.Wait()
	writer.Close()

	grepCmd.Wait()
	reader.Close()

	pathWithEscapeCharacter := buf.String()
	return strings.Replace(pathWithEscapeCharacter, "\n", "", 1), nil
}
