package main

import (
	"math/rand"
	"time"

	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/auprc"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"
	"github.com/pkg/profile"
)

func main() {
	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	auprc.MeasureAUPRC(dupedetection.NewComputeConfig())
}
