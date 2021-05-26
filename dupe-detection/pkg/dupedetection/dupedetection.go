// Package dupedetection provides functions to compute dupe detection fingerprints for specific image
package dupedetection

import (
	"context"
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kofalt/go-memoize"
	"github.com/montanaflynn/stats"
	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/dupe-detection/wdm"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	"github.com/disintegration/imaging"
	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
)

const strictnessFactor = 0.985

// InterfaceTypeError indicates unexpected variable type is returned by invoked correlation calculation func
var InterfaceTypeError = errors.Errorf("Calculation function returned value of unexpected type.")

var models = make(map[string]*tg.Model)

type modelData struct {
	model string
	input string
}

var fingerprintSources = []modelData{
	{
		model: "EfficientNetB7.tf",
		input: "serving_default_input_1",
	},
	{
		model: "EfficientNetB6.tf",
		input: "serving_default_input_2",
	},
	{
		model: "InceptionResNetV2.tf",
		input: "serving_default_input_3",
	},
	{
		model: "DenseNet201.tf",
		input: "serving_default_input_4",
	},
	{
		model: "InceptionV3.tf",
		input: "serving_default_input_5",
	},
	{
		model: "NASNetLarge.tf",
		input: "serving_default_input_6",
	},
	{
		model: "ResNet152V2.tf",
		input: "serving_default_input_7",
	},
}

func tgModel(path string) *tg.Model {
	m, ok := models[path]
	if !ok {
		m = tg.LoadModel(path, []string{"serve"}, nil)
		models[path] = m
	}
	return m
}

func fromFloat32To64(input []float32) []float64 {
	output := make([]float64, len(input))
	for i, value := range input {
		output[i] = float64(value)
	}
	return output
}

func loadImage(imagePath string, width int, height int) (image.Image, error) {
	reader, err := os.Open(imagePath)
	if err != nil {
		return nil, errors.New(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.New(err)
	}

	img = imaging.Resize(img, width, height, imaging.Linear)

	return img, nil
}

// ComputeImageDeepLearningFeatures computes dupe detection fingerprints for image with imagePath
func ComputeImageDeepLearningFeatures(imagePath string) ([][]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	m, err := loadImage(imagePath, 224, 224)
	if err != nil {
		return nil, errors.New(err)
	}

	bounds := m.Bounds()

	var inputTensor [1][224][224][3]float32

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// height = y and width = x
			inputTensor[0][y][x][0] = float32(r >> 8)
			inputTensor[0][y][x][1] = float32(g >> 8)
			inputTensor[0][y][x][2] = float32(b >> 8)
		}
	}

	fingerprints := make([][]float64, len(fingerprintSources))
	for i, source := range fingerprintSources {
		model := tgModel(source.model)

		fakeInput, _ := tf.NewTensor(inputTensor)
		results := model.Exec([]tf.Output{
			model.Op("StatefulPartitionedCall", 0),
		}, map[tf.Output]*tf.Tensor{
			model.Op(source.input, 0): fakeInput,
		})

		predictions := results[0].Value().([][]float32)[0]
		fingerprints[i] = fromFloat32To64(predictions)
	}

	return fingerprints, nil
}

func computeMIForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorMI := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "MI", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorMI[currentIndex], err = memoizeCorrelationFuncCall(MI, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorMI[currentIndex], err = MI(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorMI, nil
}

var memoizer *memoize.Memoizer

// GetMemoizer returns memoizer used to cache correlations calculations results
func GetMemoizer() *memoize.Memoizer {
	if memoizer == nil {
		memoizer = memoize.NewMemoizer(cache.NoExpiration, cache.NoExpiration)
	}
	return memoizer
}

type funcCallToMemoize func([]float64, []float64) (float64, error)

func memoizeCorrelationFuncCall(f funcCallToMemoize, data1, data2 []float64, memoKey string) (float64, error) {
	pearsonCall := func() (interface{}, error) {
		return f(data1, data2)
	}

	result, err, _ := memoizer.Memoize(memoKey, pearsonCall)
	if err != nil {
		return 0, err
	}

	value, ok := result.(float64)
	if !ok {
		return 0, errors.New(InterfaceTypeError)
	}
	return value, nil
}

func computePearsonRForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorPearsonAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "PearsonR", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorPearsonAll[currentIndex], err = memoizeCorrelationFuncCall(Pearson, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorPearsonAll[currentIndex], err = Pearson(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorPearsonAll, nil
}

func computeSpearmanForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorSpearmanAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "Spearman", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorSpearmanAll[currentIndex], err = memoizeCorrelationFuncCall(Spearman, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorSpearmanAll[currentIndex], err = Spearman(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorSpearmanAll, nil
}

func computeKendallForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorKendallAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "Kendall", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorKendallAll[currentIndex], err = memoizeCorrelationFuncCall(Kendall, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorKendallAll[currentIndex], err = Kendall(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorKendallAll, nil
}

func computeHoeffdingDForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorHoeffdingBetaAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "Hoeffding", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorHoeffdingBetaAll[currentIndex], err = memoizeCorrelationFuncCall(HoeffdingD, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorHoeffdingBetaAll[currentIndex], err = HoeffdingD(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorHoeffdingBetaAll, nil
}

func computeBlomqvistBetaForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64, memoizationData MemoizationImageData, _ ComputeConfig) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorBlomqvistBetaAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			if memoizer != nil {
				//currentImageSHA256 := memoizationData.SHA256HashOfFetchedImages[currentIndex]
				memoKey := fmt.Sprintf("%v:%v:%v", "Blomqvist", currentIndex, memoizationData.SHA256HashOfCurrentImage)
				similarityScoreVectorBlomqvistBetaAll[currentIndex], err = memoizeCorrelationFuncCall(Blomqvist, candidateImageFingerprint, currentFingerprint, memoKey)
				if err != nil {
					return err
				}
			} else {
				similarityScoreVectorBlomqvistBetaAll[currentIndex], err = Blomqvist(candidateImageFingerprint, currentFingerprint)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return similarityScoreVectorBlomqvistBetaAll, nil
}

func filterOutFingerprintsByThreshold(similarityScoreVector []float64, threshold float64, cutAboveThreshold bool, combinedImageFingerprintArray [][]float64) [][]float64 {
	defer pruntime.PrintExecutionTime(time.Now())
	var filteredByThresholdCombinedImageFingerprintArray [][]float64
	for i, valueToCheck := range similarityScoreVector {
		if !math.IsNaN(valueToCheck) && ((cutAboveThreshold && valueToCheck >= threshold) || (!cutAboveThreshold && valueToCheck < threshold)) {
			filteredByThresholdCombinedImageFingerprintArray = append(filteredByThresholdCombinedImageFingerprintArray, combinedImageFingerprintArray[i])
		}
	}
	return filteredByThresholdCombinedImageFingerprintArray
}

func randInts(min int, max int, size int) []int {
	output := make([]int, size)
	for i := range output {
		output[i] = rand.Intn(max-min) + min
	}
	return output
}

func arrayValuesFromIndexes(input []float64, indexes []int) []float64 {
	output := make([]float64, len(indexes))
	for i := range output {
		output[i] = input[indexes[i]]
	}
	return output
}

func filterOutArrayValuesByRange(input []float64, min, max float64) []float64 {
	var output []float64
	for i := range input {
		if input[i] >= min && input[i] <= max {
			output = append(output, input[i])
		}
	}
	return output
}

func computeAverageAndStdevOf25thTo75thPercentile(inputVector []float64, filterByPercentile bool) (float64, float64) {
	trimmedVector := inputVector
	if filterByPercentile {
		percentile25, _ := stats.Percentile(inputVector, 25)
		percentile75, _ := stats.Percentile(inputVector, 75)

		trimmedVector = filterOutArrayValuesByRange(inputVector, percentile25, percentile75)
	}

	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func computeAverageAndStdevOf50thTo90thPercentile(inputVector []float64, filterByPercentile bool) (float64, float64) {
	trimmedVector := inputVector
	if filterByPercentile {
		percentile50, _ := stats.Percentile(inputVector, 50)
		percentile90, _ := stats.Percentile(inputVector, 90)

		trimmedVector = filterOutArrayValuesByRange(inputVector, percentile50, percentile90)
	}

	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func computeParallelBootstrappedKendallsTau(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int, config ComputeConfig) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageTau := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevTau := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	g, _ := errgroup.WithContext(context.Background())
	for i, y := range arrayOfFingerprintsRequiringFurtherTesting {
		fingerprintIdx := i
		currentFingerprint := y
		g.Go(func() error {
			arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
			xBootstraps := make([][]float64, numberOfBootstraps)
			yBootstraps := make([][]float64, numberOfBootstraps)
			bootstrappedKendallTauResults := make([]float64, numberOfBootstraps)
			for i := 0; i < numberOfBootstraps; i++ {
				arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
			}
			for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
				xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
				yBootstraps[i] = arrayValuesFromIndexes(currentFingerprint, currentBootstrapIndices)

				bootstrappedKendallTauResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "kendall")
			}
			robustAverageTau[fingerprintIdx], robustStdevTau[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedKendallTauResults, config.TrimByPercentile)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}
	return robustAverageTau, robustStdevTau, nil
}

func computeAverageRatioOfArrays(numerator []float64, denominator []float64) float64 {
	ratio := make([]float64, len(numerator))
	for i := range ratio {
		ratio[i] = numerator[i] / denominator[i]
	}
	averageRatio, _ := stats.Mean(ratio)
	return averageRatio
}

func rankArrayOrdinalCopulaTransformation(input []float64) []float64 {
	ranks := make([]float64, len(input)*2)
	outputIdx := 0
	for i := range input {
		ranks[outputIdx] = 1
		for j := range input {
			if (i != j) && (input[j] <= input[i]) {
				ranks[outputIdx]++
			}
		}
		ranks[outputIdx] /= float64(len(input))
		outputIdx++
		ranks[outputIdx] = 1
		outputIdx++
	}
	return ranks
}

func randomLinearProjection(s float64, size int) []float64 {
	output := make([]float64, size*2)
	for i := range output {
		output[i] = s / 2 * rand.NormFloat64()
	}
	return output
}

// ComputeRandomizedDependence computes RDC correlation between input arrays of data
func ComputeRandomizedDependence(x, y []float64) float64 {
	sinFunc := func(i, j int, v float64) float64 {
		return math.Sin(v)
	}
	s := 1.0 / 6.0
	k := 20.0
	cx := mat.NewDense(len(x), 2, rankArrayOrdinalCopulaTransformation(x))
	cy := mat.NewDense(len(y), 2, rankArrayOrdinalCopulaTransformation(y))

	Rx := mat.NewDense(2, int(k), randomLinearProjection(s, int(k)))
	Ry := mat.NewDense(2, int(k), randomLinearProjection(s, int(k)))

	var X, Y, stackedH mat.Dense
	X.Mul(cx, Rx)
	Y.Mul(cy, Ry)

	X.Apply(sinFunc, &X)
	Y.Apply(sinFunc, &Y)

	stackedH.Augment(&X, &Y)
	//stackedT := mat.DenseCopyOf(stackedH.T())

	var CsymDense mat.SymDense

	stat.CovarianceMatrix(&CsymDense, &stackedH, nil)
	C := mat.DenseCopyOf(&CsymDense)

	//fmt.Printf("\n%v", mat.Formatted(C, mat.Prefix(""), mat.Squeeze()))

	k0 := int(k)
	lb := 1.0
	ub := k
	counter := 0
	var maxEigVal float64
	for {
		maxEigVal = 0
		counter++
		//fmt.Printf("\n%v", k)
		Cxx := mat.DenseCopyOf(C.Slice(0, int(k), 0, int(k)))
		Cyy := C.Slice(k0, k0+int(k), k0, k0+int(k)).(*mat.Dense)
		Cxy := C.Slice(0, int(k), k0, k0+int(k)).(*mat.Dense)
		Cyx := C.Slice(k0, k0+int(k), 0, int(k)).(*mat.Dense)

		var CxxInversed, CyyInversed mat.Dense
		CxxInversed.Inverse(Cxx)
		CyyInversed.Inverse(Cyy)

		var CxxInversedMulCxy, CyyInversedMulCyx, resultMul mat.Dense
		CxxInversedMulCxy.Mul(&CxxInversed, Cxy)
		CyyInversedMulCyx.Mul(&CyyInversed, Cyx)

		resultMul.Mul(&CxxInversedMulCxy, &CyyInversedMulCyx)

		//fmt.Printf("\nresultMul: %v", mat.Formatted(&resultMul, mat.Prefix(""), mat.Squeeze()))

		var eigs mat.Eigen
		eigs.Factorize(&resultMul, 0)

		continueLoop := false
		eigsVals := eigs.Values(nil)
		for _, eigVal := range eigsVals {
			realVal := real(eigVal)
			imageVal := imag(eigVal)
			if !(imageVal == 0.0 && 0 <= realVal && realVal <= 1) {
				ub--
				k = (ub + lb) / 2.0
				continueLoop = true
				maxEigVal = 0
				break
			}
			//maxEigVal = math.Max(maxEigVal, realVal)
			if maxEigVal < realVal {
				maxEigVal = realVal
				//fmt.Printf("\n%v", maxEigVal)
			}
		}

		if continueLoop {
			continue
		}

		if lb == ub {
			break
		}

		lb = k

		if ub == lb+1 {
			k = ub
		} else {
			k = (ub + lb) / 2.0
		}
	}
	sqrtResult := math.Sqrt(maxEigVal)
	result := math.Pow(sqrtResult, 12)

	return result
}

func computeParallelBootstrappedRandomizedDependence(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int, config ComputeConfig) ([]float64, []float64, error) {
	originalLengthOfInput := len(x)
	robustAverageRandomizedDependence := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevRandomizedDependence := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	g, _ := errgroup.WithContext(context.Background())
	for i, y := range arrayOfFingerprintsRequiringFurtherTesting {
		fingerprintIdx := i
		currentFingerprint := y

		g.Go(func() error {
			arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
			xBootstraps := make([][]float64, numberOfBootstraps)
			yBootstraps := make([][]float64, numberOfBootstraps)
			bootstrappedRandomizedDependenceResults := make([]float64, numberOfBootstraps)
			for i := 0; i < numberOfBootstraps; i++ {
				arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
			}

			for currentIdx, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
				xBootstraps[currentIdx] = arrayValuesFromIndexes(x, currentBootstrapIndices)
				yBootstraps[currentIdx] = arrayValuesFromIndexes(currentFingerprint, currentBootstrapIndices)

				/*var rdcProbes []float64
				const probeCount = 3
				for probe := 0; probe < probeCount; probe++ {
					rdcProbes = append(rdcProbes, computeRandomizedDependence(xBootstraps[current], yBootstraps[current]))
				}
				var err error
				bootstrappedRandomizedDependenceResults[current], err = stats.Trimean(rdcProbes)
				if err != nil {
					return errors.New(err)
				}*/

				bootstrappedRandomizedDependenceResults[currentIdx] = ComputeRandomizedDependence(xBootstraps[currentIdx], yBootstraps[currentIdx])

			}

			robustAverageRandomizedDependence[fingerprintIdx], robustStdevRandomizedDependence[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedRandomizedDependenceResults, config.TrimByPercentile)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}
	return robustAverageRandomizedDependence, robustStdevRandomizedDependence, nil
}

func computeParallelBootstrappedBaggedHoeffdingsDSmallerSampleSize(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int, config ComputeConfig) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	g, _ := errgroup.WithContext(context.Background())
	for i, y := range arrayOfFingerprintsRequiringFurtherTesting {
		fingerprintIdx := i
		currentFingerprint := y

		g.Go(func() error {
			arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
			xBootstraps := make([][]float64, numberOfBootstraps)
			yBootstraps := make([][]float64, numberOfBootstraps)
			bootstrappedHoeffdingsDResults := make([]float64, numberOfBootstraps)
			for i := 0; i < numberOfBootstraps; i++ {
				arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
			}
			for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
				xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
				yBootstraps[i] = arrayValuesFromIndexes(currentFingerprint, currentBootstrapIndices)

				bootstrappedHoeffdingsDResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "hoeffding")
			}
			robustAverageD[fingerprintIdx], robustStdevD[fingerprintIdx] = computeAverageAndStdevOf25thTo75thPercentile(bootstrappedHoeffdingsDResults, config.TrimByPercentile)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}
	return robustAverageD, robustStdevD, nil
}

func computeParallelBootstrappedBaggedHoeffdingsD(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int, config ComputeConfig) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	g, _ := errgroup.WithContext(context.Background())
	for i, y := range arrayOfFingerprintsRequiringFurtherTesting {
		fingerprintIdx := i
		currentFingerprint := y

		g.Go(func() error {
			arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
			xBootstraps := make([][]float64, numberOfBootstraps)
			yBootstraps := make([][]float64, numberOfBootstraps)
			bootstrappedHoeffdingsDResults := make([]float64, numberOfBootstraps)
			for i := 0; i < numberOfBootstraps; i++ {
				arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
			}
			for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
				xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
				yBootstraps[i] = arrayValuesFromIndexes(currentFingerprint, currentBootstrapIndices)

				bootstrappedHoeffdingsDResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "hoeffding")
			}
			robustAverageD[fingerprintIdx], robustStdevD[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedHoeffdingsDResults, config.TrimByPercentile)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}
	return robustAverageD, robustStdevD, nil
}

func computeParallelBootstrappedBlomqvistBeta(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int, config ComputeConfig) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageBlomqvist := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevBlomqvist := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	g, _ := errgroup.WithContext(context.Background())
	for i, y := range arrayOfFingerprintsRequiringFurtherTesting {
		fingerprintIdx := i
		currentFingerprint := y
		g.Go(func() error {
			arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
			xBootstraps := make([][]float64, numberOfBootstraps)
			yBootstraps := make([][]float64, numberOfBootstraps)
			var bootstrappedBlomqvistResults []float64
			for i := 0; i < numberOfBootstraps; i++ {
				arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
			}
			for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
				xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
				yBootstraps[i] = arrayValuesFromIndexes(currentFingerprint, currentBootstrapIndices)

				blomqvistValue := wdm.Wdm(xBootstraps[i], yBootstraps[i], "blomqvist")
				if !math.IsNaN(blomqvistValue) {
					bootstrappedBlomqvistResults = append(bootstrappedBlomqvistResults, blomqvistValue)
				}
			}
			robustAverageBlomqvist[fingerprintIdx], robustStdevBlomqvist[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedBlomqvistResults, config.TrimByPercentile)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, nil, errors.New(err)
	}
	return robustAverageBlomqvist, robustStdevBlomqvist, nil
}

type printInfo func([]float64)
type bootstrappedPrintInfo func([]float64, []float64)
type correlation func([]float64, [][]float64, MemoizationImageData, ComputeConfig) ([]float64, error)
type bootstrappedCorrelation func([]float64, [][]float64, int, int, ComputeConfig) ([]float64, []float64, error)

type computeData struct {
	name                        string
	title                       string
	computeFunc                 correlation
	printFunc                   printInfo
	bootstrappedComputeFunc     bootstrappedCorrelation
	bootstrappedPrintFunc       bootstrappedPrintInfo
	sampleSize, bootstrapsCount int
	threshold                   float64
	cutAboveThreshold           bool
	totalFingerprintsCount      int
	config                      ComputeConfig
}

func computeSimilarity(candidateImageFingerprint []float64, fingerprints [][]float64, memoizationData MemoizationImageData, data computeData) ([][]float64, error) {
	if len(fingerprints) == 0 {
		return [][]float64{}, nil
	}
	defer pruntime.PrintExecutionTime(time.Now())
	fmt.Printf("\n\n-->Computing %v", data.title)

	var similarityScore, similarityScoreStdev []float64
	var err error

	if data.computeFunc != nil {
		similarityScore, err = data.computeFunc(candidateImageFingerprint, fingerprints, memoizationData, data.config)
		if err != nil {
			return nil, errors.New(err)
		}
		if data.printFunc != nil {
			data.printFunc(similarityScore)
		}
	}

	if data.bootstrappedComputeFunc != nil {
		similarityScore, similarityScoreStdev, err = data.bootstrappedComputeFunc(candidateImageFingerprint, fingerprints, data.sampleSize, data.bootstrapsCount, data.config)
		if err != nil {
			return nil, errors.New(err)
		}

		if data.bootstrappedPrintFunc != nil {
			data.bootstrappedPrintFunc(similarityScore, similarityScoreStdev)
		}
	}

	fingerprintsRequiringFurtherTesting := filterOutFingerprintsByThreshold(similarityScore, data.threshold, data.cutAboveThreshold, fingerprints)
	percentageOfFingerprintsRequiringFurtherTesting := float32(len(fingerprintsRequiringFurtherTesting)) / float32(data.totalFingerprintsCount)
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(fingerprintsRequiringFurtherTesting), percentageOfFingerprintsRequiringFurtherTesting*100)
	return fingerprintsRequiringFurtherTesting, nil
}

func printPearsonRCalculationResults(results []float64) {
	pearsonMax, _ := stats.Max(results)
	fmt.Printf("\nLength of computed pearson r vector: %v with max value %v", len(results), pearsonMax)
}

func printKendallTauCalculationResults(similarityScore, similarityScoreStdev []float64) {
	stdevAsPctOfRobustAvgKendall := computeAverageRatioOfArrays(similarityScoreStdev, similarityScore)
	fmt.Printf("\nStandard Deviation as %% of Average Tau -- average across all fingerprints: %.2f%%", stdevAsPctOfRobustAvgKendall*100)
}

func printRDCCalculationResults(similarityScore, similarityScoreStdev []float64) {
	stdevAsPctOfRobustAvgRandomizedDependence := computeAverageRatioOfArrays(similarityScoreStdev, similarityScore)
	fmt.Printf("\nStandard Deviation as %% of Average Randomized Dependence -- average across all fingerprints: %.2f%%", stdevAsPctOfRobustAvgRandomizedDependence*100)
}

func printBlomqvistBetaCalculationResults(similarityScore, similarityScoreStdev []float64) {
	_ = similarityScoreStdev
	similarityScoreVectorBlomqvistAverage, _ := stats.Mean(similarityScore)
	fmt.Printf("\n Average for Blomqvist's beta: %.4f", strictnessFactor*similarityScoreVectorBlomqvistAverage)
}

// MemoizationImageData provides images data for memoization of correlation calculation methods results
type MemoizationImageData struct {
	SHA256HashOfFetchedImages []string
	SHA256HashOfCurrentImage  string
}

// ComputeConfig contains configurable parameters to calculate AUPRC of image similariy measurement
type ComputeConfig struct {
	CorrelationMethodNameArray        []string
	StableOrderOfCorrelationMethods   []string
	UnstableOrderOfCorrelationMethods []string
	CorrelationMethodsOrder           string
	PearsonDupeThreshold              float64
	SpearmanDupeThreshold             float64
	KendallDupeThreshold              float64
	RandomizedDependenceDupeThreshold float64
	BlomqvistDupeThreshold            float64
	HoeffdingDupeThreshold            float64
	HoeffdingRound1DupeThreshold      float64
	HoeffdingRound2DupeThreshold      float64
	MIThreshold                       float64

	RootDir                  string
	NumberOfImagesToValidate int
	TrimByPercentile         bool
}

// NewComputeConfig retirieves new ComputeConfig with default values
func NewComputeConfig() ComputeConfig {
	config := ComputeConfig{}

	config.TrimByPercentile = false

	config.PearsonDupeThreshold = 0.995
	config.SpearmanDupeThreshold = 0.79
	config.KendallDupeThreshold = 0.70
	config.RandomizedDependenceDupeThreshold = 0.79
	config.BlomqvistDupeThreshold = 0.723
	config.HoeffdingDupeThreshold = 0.45
	config.HoeffdingRound1DupeThreshold = 0.35
	config.HoeffdingRound2DupeThreshold = 0.23
	config.MIThreshold = 5.28

	config.CorrelationMethodNameArray = []string{
		//"MI",
		"PearsonR",
		"SpearmanRho",
		"KendallTau",
		"HoeffdingD",
		//"BootstrappedKendallTau",
		//"BootstrappedRDC",
		"BlomqvistBeta",
		//"BootstrappedBlomqvistBeta",

		//"HoeffdingDRound1",
		//"HoeffdingDRound2",
	}
	config.CorrelationMethodsOrder = strings.Join(config.CorrelationMethodNameArray, " ")

	config.StableOrderOfCorrelationMethods = []string{
		"PearsonR",
		"SpearmanRho",
	}

	config.UnstableOrderOfCorrelationMethods = []string{
		"MI",
		"KendallTau",
		"HoeffdingD",
		"BlomqvistBeta",
	}

	return config
}

// MeasureImageSimilarity calculates similarity between candidateImageFingerprint and each value in fingerprintsArrayToCompareWith
func MeasureImageSimilarity(candidateImageFingerprint []float64, fingerprintsArrayToCompareWith [][]float64, memoizationData MemoizationImageData, config ComputeConfig) (int, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	var err error
	totalFingerprintsCount := len(fingerprintsArrayToCompareWith)

	fmt.Printf("\nConfigured CorrelationMethodsOrder: %v", config.CorrelationMethodsOrder)

	var orderedComputeBlocks []computeData
	orderedCorrelationMethods := strings.Split(config.CorrelationMethodsOrder, " ")
	for _, computeBlockName := range orderedCorrelationMethods {
		var nextComputeBlock computeData
		switch {
		case computeBlockName == "PearsonR":
			nextComputeBlock = computeData{
				name:                   "PearsonR",
				title:                  "Pearson's R, which is fast to compute (We only perform the slower tests on the fingerprints that have a high R).",
				computeFunc:            computePearsonRForAllFingerprintPairs,
				printFunc:              printPearsonRCalculationResults,
				threshold:              strictnessFactor * config.PearsonDupeThreshold,
				cutAboveThreshold:      true,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		case computeBlockName == "SpearmanRho":
			nextComputeBlock = computeData{
				name:                   "SpearmanRho",
				title:                  "Spearman's Rho for selected fingerprints...",
				computeFunc:            computeSpearmanForAllFingerprintPairs,
				printFunc:              nil,
				threshold:              strictnessFactor * config.SpearmanDupeThreshold,
				cutAboveThreshold:      true,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		case computeBlockName == "KendallTau":
			nextComputeBlock = computeData{
				name:                   "KendallTau",
				title:                  "Kendall's Tau for selected fingerprints...",
				computeFunc:            computeKendallForAllFingerprintPairs,
				printFunc:              nil,
				threshold:              strictnessFactor * config.KendallDupeThreshold,
				cutAboveThreshold:      true,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		case computeBlockName == "BootstrappedKendallTau":
			nextComputeBlock = computeData{
				name:                    "BootstrappedKendallTau",
				title:                   "Bootstrapped Kendall's Tau for selected fingerprints...",
				computeFunc:             nil,
				printFunc:               nil,
				bootstrappedComputeFunc: computeParallelBootstrappedKendallsTau,
				bootstrappedPrintFunc:   printKendallTauCalculationResults,
				sampleSize:              50,
				bootstrapsCount:         100,
				threshold:               strictnessFactor * config.KendallDupeThreshold,
				cutAboveThreshold:       true,
				totalFingerprintsCount:  totalFingerprintsCount,
				config:                  config,
			}
		case computeBlockName == "BootstrappedRDC":
			nextComputeBlock = computeData{
				name:                    "BootstrappedRDC",
				title:                   "Boostrapped Randomized Dependence Coefficient for selected fingerprints...",
				computeFunc:             nil,
				printFunc:               nil,
				bootstrappedComputeFunc: computeParallelBootstrappedRandomizedDependence,
				bootstrappedPrintFunc:   printRDCCalculationResults,
				sampleSize:              50,
				bootstrapsCount:         100,
				threshold:               strictnessFactor * config.RandomizedDependenceDupeThreshold,
				cutAboveThreshold:       true,
				totalFingerprintsCount:  totalFingerprintsCount,
				config:                  config,
			}
		case computeBlockName == "BlomqvistBeta":
			nextComputeBlock = computeData{
				name:                   "BlomqvistBeta",
				title:                  "BlomqvistBeta for selected fingerprints...",
				computeFunc:            computeBlomqvistBetaForAllFingerprintPairs,
				printFunc:              nil,
				threshold:              strictnessFactor * config.BlomqvistDupeThreshold,
				cutAboveThreshold:      true,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		case computeBlockName == "BootstrappedBlomqvistBeta":
			nextComputeBlock = computeData{
				name:                    "BootstrappedBlomqvistBeta",
				title:                   "Bootstrapped Blomqvist's beta for selected fingerprints...",
				computeFunc:             nil,
				printFunc:               nil,
				bootstrappedComputeFunc: computeParallelBootstrappedBlomqvistBeta,
				bootstrappedPrintFunc:   printBlomqvistBetaCalculationResults,
				sampleSize:              100,
				bootstrapsCount:         100,
				threshold:               strictnessFactor * config.BlomqvistDupeThreshold,
				cutAboveThreshold:       true,
				totalFingerprintsCount:  totalFingerprintsCount,
				config:                  config,
			}
		case computeBlockName == "HoeffdingD":
			nextComputeBlock = computeData{
				name:                   "HoeffdingD",
				title:                  "HoeffdingD for selected fingerprints...",
				computeFunc:            computeHoeffdingDForAllFingerprintPairs,
				printFunc:              nil,
				threshold:              strictnessFactor * config.HoeffdingDupeThreshold,
				cutAboveThreshold:      true,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		case computeBlockName == "HoeffdingDRound1":
			nextComputeBlock = computeData{
				name:                    "HoeffdingDRound1",
				title:                   "Hoeffding's D Round 1",
				computeFunc:             nil,
				printFunc:               nil,
				bootstrappedComputeFunc: computeParallelBootstrappedBaggedHoeffdingsDSmallerSampleSize,
				bootstrappedPrintFunc:   nil,
				sampleSize:              20,
				bootstrapsCount:         50,
				threshold:               strictnessFactor * config.HoeffdingRound1DupeThreshold,
				cutAboveThreshold:       true,
				totalFingerprintsCount:  totalFingerprintsCount,
				config:                  config,
			}
		case computeBlockName == "HoeffdingDRound2":
			nextComputeBlock = computeData{
				name:                    "HoeffdingDRound2",
				title:                   "Hoeffding's D Round 2",
				computeFunc:             nil,
				printFunc:               nil,
				bootstrappedComputeFunc: computeParallelBootstrappedBaggedHoeffdingsD,
				bootstrappedPrintFunc:   nil,
				sampleSize:              75,
				bootstrapsCount:         20,
				threshold:               strictnessFactor * config.HoeffdingRound2DupeThreshold,
				cutAboveThreshold:       true,
				totalFingerprintsCount:  totalFingerprintsCount,
				config:                  config,
			}
		case computeBlockName == "MI":
			nextComputeBlock = computeData{
				name:                   "MI",
				title:                  "Mutual Information",
				computeFunc:            computeMIForAllFingerprintPairs,
				printFunc:              nil,
				threshold:              strictnessFactor * config.MIThreshold,
				cutAboveThreshold:      false,
				totalFingerprintsCount: totalFingerprintsCount,
				config:                 config,
			}
		default:
			return 0, errors.New(errors.Errorf("Unrecognized similarity computation block name \"%v\"", computeBlockName))
		}
		orderedComputeBlocks = append(orderedComputeBlocks, nextComputeBlock)
	}

	fingerprintsForFurtherTesting := fingerprintsArrayToCompareWith
	for _, computeBlock := range orderedComputeBlocks {

		fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, memoizationData, computeBlock)
		if err != nil {
			return 0, errors.New(err)
		}
		if len(fingerprintsForFurtherTesting) == 0 {
			break
		}
	}

	if len(fingerprintsForFurtherTesting) > 0 {
		fmt.Printf("\n\nWARNING! Art image file appears to be a duplicate!")
	} else {
		fmt.Printf("\n\nArt image file appears to be original! (i.e., not a duplicate of an existing image in the image fingerprint database)")
	}

	result := 1
	if len(fingerprintsForFurtherTesting) == 0 {
		result = 0
	}
	return result, nil
}
