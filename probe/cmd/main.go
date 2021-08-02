package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/probe/configs"
	"github.com/pastelnetwork/gonode/probe/pkg/auprc"
	"github.com/pastelnetwork/gonode/probe/pkg/dupedetection"
	"github.com/pkg/profile"
)

const cacheFileName = "cached"

func main() {
	rootDirPtr := flag.String("rootDir", "", "a path to the directory with the test corpus of images.")
	configFile := flag.String("configFile", "", "a path to the directory with the configuration file.")
	numberOfImagesToValidate := flag.Int("imageCount", 0, "limits the number of dupes and original images to validate.")
	flag.Parse()

	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	memoizer := dupedetection.GetMemoizer()
	memoizer.Storage.LoadFile(cacheFileName)

	// TODO: Add conf to NewComputeConfig
	conf := configs.New()
	if *configFile != "" {
		if err := configurer.ParseFile(*configFile, conf); err != nil {
			panic(err)
		}
	}

	config := dupedetection.NewComputeConfig()
	config.RootDir = *rootDirPtr
	config.NumberOfImagesToValidate = *numberOfImagesToValidate

	if _, err := auprc.MeasureAUPRC(ctx, config); err != nil {
		log.Print(errors.ErrorStack(err))
		panic(err)
	}

	memoizer.Storage.SaveFile(cacheFileName)
}
