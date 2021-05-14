package main

import (
	"flag"
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
	rootDirPtr := flag.String("rootDir", "", "a path to the directory with the test corpus of images.")
	flag.Parse()

	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	config := dupedetection.NewComputeConfig()
	config.RootDir = *rootDirPtr

	if _, err := auprc.MeasureAUPRC(config); err != nil {
		log.Printf(errors.ErrorStack(err))
		panic(err)
	}
}
