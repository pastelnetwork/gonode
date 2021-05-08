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

// ① Define an objective function which returns a value you want to minimize.
func objective(trial goptuna.Trial) (float64, error) {
	// ② Define the search space via Suggest APIs.
	config := dupedetection.NewComputeConfig()
	config.PearsonDupeThreshold, _ = trial.SuggestFloat("PearsonDupeThreshold", 0.5, 0.99999)
	config.SpearmanDupeThreshold, _ = trial.SuggestFloat("SpearmanDupeThreshold", 0.5, 0.99999)

	aurpc := 1.0 - auprc.MeasureAUPRC(config)
	return aurpc, nil
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// ③ Create a study which manages each experiment.
	study, err := goptuna.CreateStudy(
		"goptuna-aurpc",
		goptuna.StudyOptionRelativeSampler(cmaes.NewSampler()))
	if err != nil {
		panic(err)
	}

	// ④ Evaluate your objective function.
	err = study.Optimize(objective, 25)
	if err != nil {
		panic(err)
	}

	// ⑤ Print the best evaluation parameters.
	v, _ := study.GetBestValue()
	p, _ := study.GetBestParams()
	log.Printf("Best value=%f (PearsonDupeThreshold=%f, SpearmanDupeThreshold=%f)",
		v, p["PearsonDupeThreshold"].(float64), p["SpearmanDupeThreshold"].(float64))
}
