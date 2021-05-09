package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/auprc"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"
	"github.com/pkg/profile"
)

func main() {
	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	if _, err := auprc.MeasureAUPRC(dupedetection.NewComputeConfig()); err != nil {
		log.Printf(errors.ErrorStack(err))
		panic(err)
	}
}
