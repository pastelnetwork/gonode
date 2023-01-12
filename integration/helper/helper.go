package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HttpGet    = "GET"
	HttpPost   = "POST"
	HttpDelete = "DELETE"
	// wsExpectedResponseTimeout in seconds until the websocket recieves the expected response
	wsExpectedResponseTimeout = 200

	SN1BaseURI = "http://localhost:19090"
	// SN2BaseURI of SN2 Server
	SN2BaseURI = "http://localhost:19091"
	// SN3BaseURI of SN3 Server
	SN3BaseURI = "http://localhost:19092"
	// SN4BaseURI of SN4 Server
	SN4BaseURI = "http://localhost:19093"
	// SN5BaseURI of SN5 Server
	SN5BaseURI = "http://localhost:19094"
	// SN6BaseURI of SN6 Server
	SN6BaseURI = "http://localhost:19095"
	// SN7BaseURI of SN6 Server
	SN7BaseURI = "http://localhost:19096"
)

// ItHelper is used by integration tests to make requests to server
type ItHelper struct {
	client *http.Client
}

// NewItHelper returns instance of ItHelper with default http client
func NewItHelper() *ItHelper {
	return &ItHelper{
		client: &http.Client{Timeout: time.Duration(10 * time.Second)},
	}
}

func GetStoreURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "p2p")
}

func GetLocalStoreURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "local_p2p")
}

func GetStorageChallengeURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "storage/challenges")
}

func GetRetrieveURI(baseURI, key string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "p2p", key)
}

func GetTaskStatePath(taskID string) string {
	return fmt.Sprintf("nfts/register/%s/state", taskID)
}

// RequestRaw makes a http request onto uri with []byte payload
func (h *ItHelper) RequestRaw(method string, reqBody []byte, uri string, queryParams map[string]string) (resp []byte, status int, err error) {
	i := 0
	for key, val := range queryParams {
		if i == 0 {
			uri = uri + "?"
		} else {
			uri = uri + "&"
		}

		uri = uri + key + "=" + val
		i++
	}

	request, err := http.NewRequest(method, uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return resp, 0, err
	}

	request.Header.Set("Authorization", "passphrase")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("app_pastelid_passphrase", "passphrase")

	response, err := h.client.Do(request)
	if err != nil {
		return resp, 0, err
	}

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return resp, 0, err
	}

	return resp, response.StatusCode, nil
}

// Request makes a http request onto ItHelper.uri with ItHelper.token set for auth
func (h *ItHelper) Request(method string, payload interface{}, uri string, queryParams map[string]string) (resp []byte, status int, err error) {

	if payload != nil {
		reqBody, err := json.Marshal(payload)
		if err != nil {
			return resp, 0, err
		}
		return h.RequestRaw(method, reqBody, uri, queryParams)
	}

	return h.RequestRaw(method, nil, uri, queryParams)
}

// Ping checks the health of solo api server
func (h *ItHelper) Ping(uri string) error {
	request, err := http.NewRequest(HttpGet, uri, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %v", resp.StatusCode)
	}

	return nil
}

// PostForm makes a http request
func (h *ItHelper) HTTPCurlUploadFile(method, uri, file, filename string) (resp []byte, err error) {
	app := "curl"
	arg0 := "--location"
	arg1 := "--request"
	arg4 := "--form"
	arg5 := fmt.Sprintf(`file=@"%s"`, file)
	arg6 := "--form"
	arg7 := fmt.Sprintf(`filename="%s"`, filename)

	cmd := exec.Command(app, arg0, arg1, method, uri, arg4, arg5, arg6, arg7)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func DoNFTSearchWSReq(addr string, path string, expectedKey []string, expectedValue []string) error {
	ctx := context.Background()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := fmt.Sprintf("ws://%s/%s", addr, path)
	log.Printf("connecting %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
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

			count := 0
			for i := 0; i < len(expectedKey); i++ {
				resp := resp
				key := expectedKey[i]
				keysSplit := strings.Split(key, ".")
				if len(keysSplit) > 1 {
					resp = resp[keysSplit[0]].(map[string]interface{})
					key = keysSplit[1]
				}

				val, ok := resp[key]
				log.Printf("ws key: %s - expected response: %s - got response: %s\n",
					key, expectedValue[i], val)
				if ok {
					if val == expectedValue[i] {
						count++
						continue
					} else if val == "Request Rejected" {
						done <- false
					}
				}
			}

			if count == len(expectedKey) {
				done <- true
			} else {
				log.Printf("expected matches: %d - got matches: %d\n", len(expectedKey), count)
				done <- false
			}
		}
	}()

	select {
	case <-ctx.Done():
		return errors.New("context cancelled")
	case <-time.After(wsExpectedResponseTimeout * time.Second):
		return errors.New("timeout")
	case val := <-done:
		if val {
			return nil
		}

		return errors.New("task failed (Request Rejected), please see container logs")
	}

}

func DoWebSocketReq(addr, path, expectedKey, expectedValue string) error {
	ctx := context.Background()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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

			if val, ok := resp[expectedKey]; ok {
				log.Printf("ws expected key: %s - response: %s\n", expectedKey, val)
				if val == expectedValue {
					done <- true
				} else if val == "Request Rejected" {
					done <- false
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return errors.New("context cancelled")
	case <-time.After(wsExpectedResponseTimeout * time.Second):
		return errors.New("timeout")
	case val := <-done:
		if val {
			return nil
		}

		return errors.New("task failed (Request Rejected), please see container logs")
	}

}
