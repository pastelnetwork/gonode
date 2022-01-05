package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	HttpGet    = "GET"
	HttpPost   = "POST"
	HttpDelete = "DELETE"
)

// ItHelper is used by integration tests to make requests to server
type ItHelper struct {
	client *http.Client
	token  string
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

func GetRetrieveURI(baseURI, key string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "p2p", key)

}

// Request makes a http request onto ItHelper.uri with ItHelper.token set for auth
func (h *ItHelper) Request(method string, payload interface{}, uri string) (resp []byte, status int, err error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return resp, 0, err
	}

	request, err := http.NewRequest(method, uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return resp, 0, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", h.token))
	request.Header.Set("Content-Type", "application/json")

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
