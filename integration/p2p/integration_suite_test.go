package integration_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	it "github.com/pastelnetwork/gonode/integration/p2p"

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
	composeFilePath := filepath.Join(filepath.Dir("."), "docker-compose.yml")
	identifier := strings.ToLower(uuid.New().String())
	compose = tc.NewLocalDockerCompose([]string{composeFilePath}, identifier)

	Expect(compose.WithCommand([]string{"up", "-d", "--build"}).Invoke().Error).To(Succeed())

	// Backoff wait for api-server container to be available
	helper := it.NewItHelper()
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second
	b.InitialInterval = 200 * time.Millisecond

	pingSN1 := func() error {
		return helper.Ping(fmt.Sprintf("%v/%v", it.SN1BaseURI, "health"))
	}
	pingSN2 := func() error {
		return helper.Ping(fmt.Sprintf("%v/%v", it.SN2BaseURI, "health"))
	}
	pingSN3 := func() error {
		return helper.Ping(fmt.Sprintf("%v/%v", it.SN3BaseURI, "health"))
	}

	Expect(backoff.Retry(backoff.Operation(pingSN1), b)).To(Succeed())
	Expect(backoff.Retry(backoff.Operation(pingSN2), b)).To(Succeed())
	Expect(backoff.Retry(backoff.Operation(pingSN3), b)).To(Succeed())
	time.Sleep(10 * time.Second)
})

var _ = AfterSuite(func() {
	Expect(compose.Down().Error).NotTo(HaveOccurred(), "should do docker-compose down")
})
