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

const cacheFileName = "cached"

func main() {
	rootDirPtr := flag.String("rootDir", "", "a path to the directory with the test corpus of images.")
	numberOfImagesToValidate := flag.Int("imageCount", 0, "limits the number of dupes and original images to validate.")
	flag.Parse()

	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	memoizer := dupedetection.GetMemoizer()
	memoizer.Storage.LoadFile(cacheFileName)

	config := dupedetection.NewComputeConfig()
	config.RootDir = *rootDirPtr
	config.NumberOfImagesToValidate = *numberOfImagesToValidate

	if _, err := auprc.MeasureAUPRC(config); err != nil {
		log.Print(errors.ErrorStack(err))
		panic(err)
	}

	memoizer.Storage.SaveFile(cacheFileName)
}
