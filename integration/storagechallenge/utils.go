package storagechallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ory/dockertest/v3/docker"
)

func purgeDockerResources(client *docker.Client, containers ...*docker.Container) []error {
	var errs = make([]error, 0)
	for _, container := range containers {
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func createContainer(client *docker.Client, name string, image string, network *docker.Network, env []string, ipAddress string, noStart bool) (container *docker.Container, remove func(), err error) {
	remove = func() {}
	container, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
			Env:   env,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove:    true,
			RestartPolicy: docker.RestartPolicy{Name: "no"},
		},
		NetworkingConfig: &docker.NetworkingConfig{
			EndpointsConfig: map[string]*docker.EndpointConfig{
				network.Name: {
					Aliases:   []string{name},
					NetworkID: network.ID,
					IPAMConfig: &docker.EndpointIPAMConfig{
						IPv4Address: ipAddress,
					},
					IPAddress:   ipAddress,
					IPPrefixLen: 16,
					Gateway:     "192.168.100.1",
				},
			},
		},
	})
	if err != nil {
		return
	}

	remove = func() {
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		})
	}

	if noStart {
		return
	}

	if err = client.StartContainer(container.ID, nil); err != nil {
		return
	}

	remove = func() {
		client.StopContainer(container.ID, 30)
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		})
	}

	container, err = client.InspectContainer(container.ID)

	return
}

type mockSecClient struct{}

func (*mockSecClient) Sign(_ context.Context, _ []byte, _, _, _ string) ([]byte, error) {
	return []byte("mock signature"), nil
}

func (*mockSecClient) Verify(_ context.Context, _ []byte, signature, _, _ string) (bool, error) {
	return signature == "mock signature", nil
}

func getKeyList(httpClient *http.Client) ([]string, error) {
	res, err := httpClient.Get(fmt.Sprintf("http://%s:8088/store/keys", helperHost))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var ret []string
	return ret, json.NewDecoder(res.Body).Decode(&ret)
}

func getCurrentBlock(httpClient *http.Client) (int32, error) {
	res, err := httpClient.Get(fmt.Sprintf("http://%s:8088/getblockcount", helperHost))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	var ret int
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	blkStr := string(b)
	ret, err = strconv.Atoi(blkStr)
	if err != nil {
		return 0, err
	}
	return int32(ret), err
}

func prepareChallengeStatus(httpClient *http.Client, id, nodeID, blkSent string) error {
	formVal := url.Values{}
	formVal.Add("id", id)
	formVal.Add("node_id", nodeID)
	formVal.Add("sent_block", blkSent)
	_, err := httpClient.PostForm(fmt.Sprintf("http://%s:8088/sts/sent", helperHost), formVal)
	return err
}

// remove 20% of local keys
func removeKeys(httpClient *http.Client, nodeIP string) ([]string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:9090/p2p/0.2", nodeIP), nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var ret []string
	return ret, json.NewDecoder(res.Body).Decode(&ret)
}

type challengeStatusResult struct {
	// number of challenge send to a challenging masternode from challenger masternode
	NumberOfChallengeSent int `json:"sent"`
	// number of challenge which is challenging masternode respond to verifying masternode
	NumberOfChallengeRespondedTo int `json:"respond"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode successfully verified
	NumberOfChallengeSuccess int `json:"success"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode says incorrect hash response
	NumberOfChallengeFailed int `json:"failed"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode says the process take too long
	NumberOfChallengeTimeout int `json:"timeout"`
}
