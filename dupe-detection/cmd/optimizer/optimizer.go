package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/c-bata/goptuna"
	"gorm.io/driver/mysql"

	"github.com/c-bata/goptuna/cmaes"
	"github.com/c-bata/goptuna/rdb.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/auprc"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"
)

var evaluateNumberOfTimes = 500
var rootDir = ""
var numberOfImagesToValidate = 0

// objective defines the objective of the study - find out the best aurpc value
func objective(trial goptuna.Trial) (float64, error) {
	var err error
	// Define the search space via Suggest APIs.
	config := dupedetection.NewComputeConfig()

	config.RootDir = rootDir
	err = trial.SetUserAttr("RootDir", config.RootDir)
	if err != nil {
		return 0, errors.New(err)
	}

	config.NumberOfImagesToValidate = numberOfImagesToValidate
	err = trial.SetUserAttr("MaxImageCountToEvaluate", fmt.Sprintf("%v", config.NumberOfImagesToValidate*2))
	if err != nil {
		return 0, errors.New(err)
	}

	config.MIThreshold, err = trial.SuggestFloat("MIThreshold", 5.2, 5.4)
	if err != nil {
		return 0, errors.New(err)
	}

	config.PearsonDupeThreshold, err = trial.SuggestFloat("Pearson", 0.99, 0.99999)
	if err != nil {
		return 0, errors.New(err)
	}
	config.SpearmanDupeThreshold, err = trial.SuggestFloat("Spearman", 0.75, 0.85)
	if err != nil {
		return 0, errors.New(err)
	}
	config.KendallDupeThreshold, _ = trial.SuggestFloat("Kendall", 0.68, 0.72)
	if err != nil {
		return 0, errors.New(err)
	}
	/*config.RandomizedDependenceDupeThreshold, _ = trial.SuggestFloat("RDC", 0.5, 0.99999)
	if err != nil {
		return 0, errors.New(err)
	}*/
	config.HoeffdingDupeThreshold, _ = trial.SuggestFloat("Hoeffding", 0.2, 0.6)
	if err != nil {
		return 0, errors.New(err)
	}
	config.BlomqvistDupeThreshold, _ = trial.SuggestFloat("Blomqvist", 0.6, 0.8)
	if err != nil {
		return 0, errors.New(err)
	}
	/*config.HoeffdingDupeThreshold, _ = trial.SuggestFloat("HoeffdingD1", 0.1, 0.99999)
	if err != nil {
		return 0, errors.New(err)
	}
	config.HoeffdingRound2DupeThreshold, _ = trial.SuggestFloat("HoeffdingD2", 0.1, 0.99999)
	if err != nil {
		return 0, errors.New(err)
	}

	allCombinationsOfUnstableMethods := combinations.All(config.UnstableOrderOfCorrelationMethods)
	var allOrderedCombinationsOfUnstableMethodsAsStrings []string
	for _, combination := range allCombinationsOfUnstableMethods {
		permutator := permutation.New(permutation.StringSlice(combination))
		for permutator.Next() {
			fmt.Println(combination)
			allOrderedCombinationsOfUnstableMethodsAsStrings = append(allOrderedCombinationsOfUnstableMethodsAsStrings, strings.Join(combination, " "))
		}
	}*/

	/*correlationMethodIndex, err := trial.SuggestStepInt("CorrelationMethodsOrderIndex", 0, len(allOrderedCombinationsOfUnstableMethodsAsStrings)-1, 1)
	if err != nil {
		return 0, errors.New(err)
	}
	correlationMethodsOrder := append(config.StableOrderOfCorrelationMethods, allOrderedCombinationsOfUnstableMethodsAsStrings[correlationMethodIndex])
	config.CorrelationMethodsOrder = strings.Join(correlationMethodsOrder, " ")
	if err != nil {
		return 0, errors.New(err)
	}*/

	//config.CorrelationMethodsOrder = "MI PearsonR SpearmanRho BootstrappedKendallTau BootstrappedBlomqvistBeta HoeffdingDRound1 HoeffdingDRound2"
	config.CorrelationMethodsOrder = "MI PearsonR SpearmanRho KendallTau HoeffdingD BlomqvistBeta"

	err = trial.SetUserAttr("CorrelationMethodsOrder", config.CorrelationMethodsOrder)
	if err != nil {
		return 0, errors.New(err)
	}

	aurpc, err := auprc.MeasureAUPRC(config)
	if err != nil {
		return 0, errors.New(err)
	}
	return aurpc, nil
}

func runStudy(studyName string) error {
	db, err := gorm.Open(mysql.Open("goptuna:password@tcp(localhost:3306)/goptuna?parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return errors.New(err)
	}

	storage := rdb.NewStorage(db)

	// Creates new study or loads available
	study, err := goptuna.CreateStudy(
		studyName,
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionRelativeSampler(cmaes.NewSampler()),
		//goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMaximize),
		goptuna.StudyOptionLoadIfExists(true),
	)
	if err != nil {
		return errors.New(err)
	}

	// Evaluate objective function specified number of times
	err = study.Optimize(objective, evaluateNumberOfTimes)
	if err != nil {
		return errors.New(err)
	}

	v, err := study.GetBestValue()
	if err != nil {
		return errors.New(err)
	}
	p, err := study.GetBestParams()
	if err != nil {
		return errors.New(err)
	}

	log.Printf("Best value=%f", v)
	for key, value := range p {
		log.Printf("%v=%v", key, value)
	}
	return nil
}

func main() {
	rootDirPtr := flag.String("rootDir", "", "a path to the directory with the test corpus of images.")
	goptunaStudyNamePtr := flag.String("studyName", "dupe-detection-aurpc", "a name of the Goptuna study to create or continue available.")
	numberOfImagesToValidatePtr := flag.Int("imageCount", 0, "limits the number of dupes and original images to validate.")
	evaluateNumberOfTimesPtr := flag.Int("runCount", 0, "defines the number of times goptuna will evaluate optimization objective.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	rootDir = *rootDirPtr
	numberOfImagesToValidate = *numberOfImagesToValidatePtr
	evaluateNumberOfTimes = *evaluateNumberOfTimesPtr

	if err := runStudy(*goptunaStudyNamePtr); err != nil {
		log.Printf(errors.ErrorStack(err))
		panic(err)
	}
}
