package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	numRequests  = 5
	initialDelay = 1 * time.Second
	maxRetries   = 5
	timeoutAfter = 300
)

type result struct {
	ID      string
	Elapsed time.Duration
	Error   error
}

type uploadImageResponse struct {
	FileId                string    `json:"file_id"`
	ExpiresIn             time.Time `json:"expires_in"`
	TotalEstimatedFee     int       `json:"total_estimated_fee"`
	RequiredPreburnAmount int       `json:"required_preburn_amount"`
}

type payload struct {
	BurnTxid               string `json:"burn_txid"`
	AppPastelid            string `json:"app_pastelid"`
	MakePubliclyAccessible bool   `json:"make_publicly_accessible"`
}

type startResponse struct {
	TaskID string `json:"task_id"`
}

func doUploadImage(method, filePath, fileName string) (res uploadImageResponse, err error) {
	url := "http://localhost:18080/openapi/cascade/upload"

	app := "curl"
	arg0 := "--location"
	arg1 := "--request"
	arg4 := "--form"
	arg5 := fmt.Sprintf(`file=@"%s"`, filePath)
	arg6 := "--form"
	arg7 := fmt.Sprintf(`filename="%s"`, fileName)

	cmd := exec.Command(app, arg0, arg1, method, url, arg4, arg5, arg6, arg7)
	resp, err := cmd.Output()
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(resp, &res); err != nil {
		return res, nil
	}

	return res, nil
}

func preBurnAmount(amount int) (string, error) {
	pastelCli := "/home/btanveer/pastel/pastel-cli"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	cmd := exec.Command(pastelCli, "sendtoaddress", "tPpasteLBurnAddressXXXXXXXXXXX3wy7u", fmt.Sprint(amount))
	cmd.Dir = homeDir

	res, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	return strings.Replace(string(res), "\n", "", 1), nil
}

func getBlockCount() string {
	pastelCli := "/home/btanveer/pastel/pastel-cli"

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

func doCascadeRequest(payload payload, taskID string, logger *log.Logger) (string, error) {
	url := fmt.Sprintf("http://localhost:18080/openapi/cascade/start/%s", taskID)
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

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var startResp startResponse
	err = json.Unmarshal(body, &startResp)
	if err != nil {
		return "", err
	}

	return startResp.TaskID, nil
}

func doTaskState(taskID string, expectedValue string, logger *log.Logger) error {
	url := fmt.Sprintf("ws://127.0.0.1:18080/openapi/cascade/start/%s/state", taskID)

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

				logger.Printf("task_id:%s, block_count:%s  status - response: %s\n", taskID, getBlockCount(), val)
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

func readFiles() map[string]string {
	filesInfo := make(map[string]string)
	dir := "./images" // Replace with your directory path

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	for _, file := range files {
		filesInfo[file.Name()] = filepath.Join(".", "images", file.Name())
	}

	return filesInfo
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

	start := time.Now()

	files := readFiles()

	count := 1
	taskIDs := make(map[string]time.Time)
	for fileName, filePath := range files {
		startReq := time.Now()

		uploadImageRes, err := doUploadImage("POST", filePath, fileName)
		if err != nil {
			logger.Printf("Request to upload image failed:%v\n", err)
		}
		logger.Printf("image uploaded:%d\n", count)

		burnTxID, err := preBurnAmount(uploadImageRes.RequiredPreburnAmount)
		if err != nil {
			logger.Printf("Request to pre burn amount failed:%v\n", err)
		}
		logger.Printf("amount pre-burned:%s, request-count:%d\n", burnTxID, count)

		payload := payload{BurnTxid: burnTxID, AppPastelid: "jXa6QiopivJLer8G65QsxwQmGELi1w6mbNXvrrYTvsddVE5BT57LtNCZ2SCmWStvLwWWTkuAFPsRREytgG62YX"}
		taskID, err := doCascadeRequest(payload, uploadImageRes.FileId, logger)
		if err != nil {
			logger.Printf("Request to cascade registration failed:%v\n", err)
		}
		logger.Printf("cascade task initiated:%s, request-count:%d\n", taskID, count)

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

	avgReqTime := totalReqTime / time.Duration(successes)

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

func getPubliclyAccessible(count int) bool {
	if count == 0 {
		return true
	}

	return count%2 == 0
}
