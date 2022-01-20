package main_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tc "github.com/testcontainers/testcontainers-go"
)

var (
	compose *tc.LocalDockerCompose
)

func TestIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST_ENV") != "true" {
		log.Println("Skipping integration tests")
		return
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	composeFilePath := filepath.Join(filepath.Dir("."), "infra", "docker-compose.yml")
	identifier := strings.ToLower(uuid.New().String())
	compose = tc.NewLocalDockerCompose([]string{composeFilePath}, identifier)

	Expect(compose.WithCommand([]string{"up", "-d", "--build"}).Invoke().Error).To(Succeed())

	// Backoff wait for api-server container to be available
	helper := helper.NewItHelper()
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second
	b.InitialInterval = 200 * time.Millisecond

	pingServers := func() error {
		servers := append(it.SNServers, it.RQServers...)
		servers = append(servers, it.DDServers...)
		servers = append(servers, it.PasteldServers...)
		for _, addr := range servers {
			backoff.Retry(backoff.Operation(func() error {
				if err := helper.Ping(fmt.Sprintf("%v/%v", it.SN1BaseURI, "health")); err != nil {
					return fmt.Errorf("err reaching server: %s - err: %w", addr, err)
				}

				return nil
			}), b)
		}

		return nil
	}

	Expect(pingServers()).To(Succeed())
	time.Sleep(10 * time.Second)
})

var _ = AfterSuite(func() {
	Expect(compose.Down().Error).NotTo(HaveOccurred(), "should do docker-compose down")
})
