package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"

	"github.com/pastelnetwork/gonode/dupe-detection/pkg/auprc"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"
)

const EvaluateNumberOfTimes = 3

// objective defines the objective of the study - find out the best aurpc value
func objective(trial goptuna.Trial) (float64, error) {
	// Define the search space via Suggest APIs.
	config := dupedetection.NewComputeConfig()
	config.PearsonDupeThreshold, _ = trial.SuggestFloat("PearsonDupeThreshold", 0.5, 0.99999)
	config.SpearmanDupeThreshold, _ = trial.SuggestFloat("SpearmanDupeThreshold", 0.5, 0.99999)
	config.KendallDupeThreshold, _ = trial.SuggestFloat("KendallDupeThreshold", 0.5, 0.99999)
	config.RandomizedDependenceDupeThreshold, _ = trial.SuggestFloat("RandomizedDependenceDupeThreshold", 0.5, 0.99999)
	config.RandomizedBlomqvistDupeThreshold, _ = trial.SuggestFloat("RandomizedBlomqvistDupeThreshold", 0.5, 0.99999)
	config.HoeffdingDupeThreshold, _ = trial.SuggestFloat("HoeffdingDupeThreshold", 0.1, 0.99999)
	config.HoeffdingRound2DupeThreshold, _ = trial.SuggestFloat("HoeffdingRound2DupeThreshold", 0.1, 0.99999)

	aurpc := 1.0 - auprc.MeasureAUPRC(config)
	return aurpc, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	study, err := goptuna.CreateStudy(
		"goptuna-aurpc",
		goptuna.StudyOptionRelativeSampler(cmaes.NewSampler()))
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
	log.Printf("Best value=%f \nPearsonDupeThreshold=%f,\nSpearmanDupeThreshold=%f\nKendallDupeThreshold=%f\nRandomizedDependenceDupeThreshold=%f\nRandomizedBlomqvistDupeThreshold=%f\nHoeffdingDupeThreshold=%f\nHoeffdingRound2DupeThreshold=%f",
		v,
		p["PearsonDupeThreshold"].(float64),
		p["SpearmanDupeThreshold"].(float64),
		p["KendallDupeThreshold"].(float64),
		p["RandomizedDependenceDupeThreshold"].(float64),
		p["RandomizedBlomqvistDupeThreshold"].(float64),
		p["HoeffdingDupeThreshold"].(float64),
		p["HoeffdingRound2DupeThreshold"].(float64))
}
