package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"
	combinations "github.com/mxschmitt/golang-combinations"
	"gorm.io/driver/mysql"

	"github.com/c-bata/goptuna/rdb.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/pastelnetwork/gonode/dupe-detection/pkg/auprc"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"
)

const EvaluateNumberOfTimes = 1
const MinNumberOfCorrelationMethodsInChain = 4

// objective defines the objective of the study - find out the best aurpc value
func objective(trial goptuna.Trial) (float64, error) {
	var err error
	// Define the search space via Suggest APIs.
	config := dupedetection.NewComputeConfig()
	config.PearsonDupeThreshold, err = trial.SuggestFloat("Pearson", 0.5, 0.99999)
	if err != nil {
		return 0, err
	}
	config.SpearmanDupeThreshold, err = trial.SuggestFloat("Spearman", 0.5, 0.99999)
	if err != nil {
		return 0, err
	}
	config.KendallDupeThreshold, _ = trial.SuggestFloat("Kendall", 0.5, 0.99999)
	if err != nil {
		return 0, err
	}
	config.RandomizedDependenceDupeThreshold, _ = trial.SuggestFloat("RDC", 0.5, 0.99999)
	if err != nil {
		return 0, err
	}
	config.RandomizedBlomqvistDupeThreshold, _ = trial.SuggestFloat("Blomqvist", 0.5, 0.99999)
	if err != nil {
		return 0, err
	}
	config.HoeffdingDupeThreshold, _ = trial.SuggestFloat("HoeffdingD1", 0.1, 0.99999)
	if err != nil {
		return 0, err
	}
	config.HoeffdingRound2DupeThreshold, _ = trial.SuggestFloat("HoeffdingD2", 0.1, 0.99999)
	if err != nil {
		return 0, err
	}

	allCombinationsOfOrderedMethods := combinations.All(config.CorrelationMethodNameArray)
	var allCombinationsOfOrderedMethodsAsStrings []string
	for _, orderedMethods := range allCombinationsOfOrderedMethods {
		if len(orderedMethods) >= MinNumberOfCorrelationMethodsInChain {
			allCombinationsOfOrderedMethodsAsStrings = append(allCombinationsOfOrderedMethodsAsStrings, strings.Join(orderedMethods, " "))
		}
	}
	config.CorrelationMethodsOrder, err = trial.SuggestCategorical("CorrelationMethodsOrder", allCombinationsOfOrderedMethodsAsStrings)
	if err != nil {
		return 0, err
	}

	aurpc := auprc.MeasureAUPRC(config)
	return aurpc, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, _ := gorm.Open(mysql.Open("goptuna:password@tcp(localhost:3306)/goptuna?parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	storage := rdb.NewStorage(db)

	study, err := goptuna.CreateStudy(
		"dupe-detection-aurpc",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionRelativeSampler(cmaes.NewSampler()),
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMaximize),
		goptuna.StudyOptionLoadIfExists(true),
	)
	if err != nil {
		panic(err)
	}

	// Evaluate objective function specified number of times
	err = study.Optimize(objective, EvaluateNumberOfTimes)
	if err != nil {
		panic(err)
	}

	v, _ := study.GetBestValue()
	p, _ := study.GetBestParams()
	log.Printf("Best value=%f", v)
	for key, value := range p {
		log.Printf("%v=%v", key, value)
	}
}
