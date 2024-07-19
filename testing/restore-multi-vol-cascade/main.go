package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	json "github.com/json-iterator/go"
)

const (
	initialDelay = 1 * time.Second
	maxRetries   = 3
	burnAddress  = ""
	pastelID     = ""
	passphrase   = ""
	wnServerPort = ""
)

type result struct {
	ID      string
	Elapsed time.Duration
	Error   error
}

type uploadImageResponse struct {
	FileId                            string    `json:"file_id"`
	TotalEstimatedFee                 int       `json:"total_estimated_fee"`
	RequiredPreburnTransactionAmounts []float64 `json:"required_preburn_transaction_amounts"`
}

type payload struct {
	BurnTxids              []string `json:"burn_txids"`
	AppPastelid            string   `json:"app_pastelid"`
	MakePubliclyAccessible bool     `json:"make_publicly_accessible"`
}

type startResponse struct {
	TaskID string `json:"task_id"`
}

func main() {
	logFile, err := os.OpenFile("requests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	files := readFiles()

	count := 1
	taskIDs := make(map[string]time.Time)
	for fileName, filePath := range files {
		startReq := time.Now().UTC()

		uploadImageRes, err := doUploadImage("POST", filePath, fileName)
		if err != nil {
			fmt.Printf("Request to upload image failed:%v\n", err)
			logger.Printf("Request to upload image failed:%v\n", err)
		}
		logger.Printf("image uploaded:%d\n", count)

		logger.Printf("initiating pastel-client to burn-coins and get burn-txid for registration...")
		var burnTxIDs []string
		for _, amount := range uploadImageRes.RequiredPreburnTransactionAmounts {
			burnTxID, err := preBurnAmount(amount)
			if err != nil {
				logger.Printf("Request to pre burn amount failed:%v\n", err)
			}
			burnTxIDs = append(burnTxIDs, burnTxID)
		}
		logger.Printf("amounts pre-burned:%v, request-count:%d\n", burnTxIDs, count)

		payload := payload{BurnTxids: burnTxIDs, AppPastelid: pastelID, MakePubliclyAccessible: true}
		taskID, err := doCascadeRequest(payload, uploadImageRes.FileId, logger)
		if err != nil {
			logger.Printf("Request to cascade registration failed:%v\n", err)
		}
		logger.Printf("cascade tasks initiated:%s, request-count:%d\n", taskID, count)

		taskIDs[taskID] = startReq

		count++
		time.Sleep(60 * time.Second)

		err = stopServer(wnServerPort)
		if err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}

		pathToWN, err := getWalletNodePath()
		if err != nil {
			log.Fatalf("error finding path to WN: %v", err)
		}

		// Restart the server
		err = startServer(pathToWN)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}

	}
}

func doUploadImage(method, filePath, fileName string) (res uploadImageResponse, err error) {
	url := fmt.Sprintf("http://localhost:%s/openapi/cascade/v2/upload", wnServerPort)

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
		fmt.Println("curl resp -- err", resp, err)
		return res, err
	}

	if err := json.Unmarshal(resp, &res); err != nil {
		return res, nil
	}

	log.Printf("pre-burn-amounts: %v\n", res.RequiredPreburnTransactionAmounts)

	return res, nil
}

func preBurnAmount(amount float64) (string, error) {
	pastelCli, err := getPastelCliPath()
	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	cmd := exec.Command(pastelCli, "sendtoaddress", burnAddress, fmt.Sprint(amount))
	cmd.Dir = homeDir

	res, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	log.Printf("pre-burn-amount: %s\n", string(res))

	return strings.Replace(string(res), "\n", "", 1), nil
}

func doCascadeRequest(payload payload, taskID string, logger *log.Logger) (string, error) {
	url := fmt.Sprintf("http://localhost:%s/openapi/cascade/start/%s", wnServerPort, taskID)
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
	req.Header.Add("Authorization", passphrase)

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
		fmt.Printf("response: %s\n", string(body))
		return "", fmt.Errorf("received non-200 response code: %d %s", res.StatusCode, string(body))
	}

	var startResp startResponse
	err = json.Unmarshal(body, &startResp)
	if err != nil {
		return "", err
	}

	return startResp.TaskID, nil
}

func readFiles() map[string]string {
	filesInfo := make(map[string]string)
	dir := "./images" // Replace with your directory path

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}
	log.Printf("Total files: %d\n", len(files))
	for _, file := range files {
		filesInfo[file.Name()] = filepath.Join(".", "images", file.Name())
	}
	log.Printf("Total files after trunc: %d\n", len(filesInfo))

	return filesInfo
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

func getWalletNodePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Construct the path to the walletnode directory
	// Assuming the directory structure is:
	// gonode/
	// ├── walletnode/
	// └── testing/
	//     └── restore-script/
	walletnodePath := filepath.Join(cwd, "../../walletnode")

	absWalletnodePath, err := filepath.Abs(walletnodePath)
	if err != nil {
		return "", err
	}

	return absWalletnodePath, nil
}

func stopServer(port string) error {
	cmd := exec.Command("fuser", "-k", fmt.Sprintf("%s/tcp", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stop server on port %d: %v", port, err)
	}

	// Wait for a moment to ensure the process is terminated
	time.Sleep(3 * time.Second)

	return nil
}

func startServer(dir string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Optionally, wait for the server to start
	time.Sleep(5 * time.Second)

	return nil
}
