package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pa-m/sklearn/metrics"
	"github.com/pkg/profile"
	"gonum.org/v1/gonum/mat"
	_ "gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"

	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"

	"github.com/montanaflynn/stats"

	_ "gorgonia.org/tensor"

	"github.com/pastelnetwork/gonode/dupe-detection/wdm"
)

const (
	cachedFingerprintsDB = "cachedFingerprints.sqlite"
	strictnessFactor     = 0.985
)

func fingerprintFromCache(filePath string) ([]float64, error) {
	if _, err := os.Stat(cachedFingerprintsDB); os.IsNotExist(err) {
		return nil, errors.New(errors.Errorf("Cache database is not found."))
	}
	db, err := sql.Open("sqlite3", cachedFingerprintsDB)
	if err != nil {
		return nil, errors.New(err)
	}
	defer db.Close()

	imageHash, err := getImageHashFromImageFilePath(filePath)
	if err != nil {
		return nil, errors.New(err)
	}

	selectQuery := `
			SELECT path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector,
				model_6_image_fingerprint_vector, model_7_image_fingerprint_vector FROM image_hash_to_image_fingerprint_table where sha256_hash_of_art_image_file = ? ORDER BY datetime_fingerprint_added_to_database DESC
		`
	rows, err := db.Query(selectQuery, imageHash)
	if err != nil {
		return nil, errors.New(err)
	}
	defer rows.Close()

	for rows.Next() {
		var currentImageFilePath string
		var model1ImageFingerprintVector, model2ImageFingerprintVector, model3ImageFingerprintVector, model4ImageFingerprintVector, model5ImageFingerprintVector, model6ImageFingerprintVector, model7ImageFingerprintVector []byte
		err = rows.Scan(&currentImageFilePath, &model1ImageFingerprintVector, &model2ImageFingerprintVector, &model3ImageFingerprintVector, &model4ImageFingerprintVector, &model5ImageFingerprintVector, &model6ImageFingerprintVector, &model7ImageFingerprintVector)
		if err != nil {
			return nil, errors.New(err)
		}
		combinedImageFingerprintVector := append(append(append(append(append(append(fromBytes(model1ImageFingerprintVector), fromBytes(model2ImageFingerprintVector)[:]...), fromBytes(model3ImageFingerprintVector)[:]...), fromBytes(model4ImageFingerprintVector)[:]...), fromBytes(model5ImageFingerprintVector)[:]...), fromBytes(model6ImageFingerprintVector)[:]...), fromBytes(model7ImageFingerprintVector)[:]...)
		return combinedImageFingerprintVector, nil
	}
	return nil, errors.New(errors.Errorf("Fingerprint is not found"))
}

func cacheFingerprint(fingerprints [][]float64, filePath string) error {
	if _, err := os.Stat(cachedFingerprintsDB); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", cachedFingerprintsDB)
		if err != nil {
			return errors.New(err)
		}
		defer db.Close()

		dupeDetectionImageFingerprintDatabaseCreationQuery := `
			CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file, model_1_image_fingerprint_vector array, model_2_image_fingerprint_vector array, model_3_image_fingerprint_vector array,
				model_4_image_fingerprint_vector array, model_5_image_fingerprint_vector array, model_6_image_fingerprint_vector array, model_7_image_fingerprint_vector array, datetime_fingerprint_added_to_database TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
				PRIMARY KEY (sha256_hash_of_art_image_file));
			`
		_, err = db.Exec(dupeDetectionImageFingerprintDatabaseCreationQuery)
		if err != nil {
			return errors.New(err)
		}
	}

	imageHash, err := getImageHashFromImageFilePath(filePath)
	if err != nil {
		return errors.New(err)
	}

	db, err := sql.Open("sqlite3", cachedFingerprintsDB)
	if err != nil {
		return errors.New(err)
	}
	defer db.Close()

	dataInsertionQuery := `
		INSERT OR REPLACE INTO image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file, path_to_art_image_file,
			model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector,
			model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	tx, err := db.Begin()
	if err != nil {
		return errors.New(err)
	}
	stmt, err := tx.Prepare(dataInsertionQuery)
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(imageHash, filePath, toBytes(fingerprints[0]), toBytes(fingerprints[1]), toBytes(fingerprints[2]), toBytes(fingerprints[3]), toBytes(fingerprints[4]), toBytes(fingerprints[5]), toBytes(fingerprints[6]))
	if err != nil {
		return errors.New(err)
	}
	tx.Commit()

	return nil
}

var dupeDetectionImageFingerprintDatabaseFilePath string

func tryToFindLocalDatabaseFile() bool {
	if _, err := os.Stat(dupeDetectionImageFingerprintDatabaseFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func regenerateEmptyDupeDetectionImageFingerprintDatabase() error {
	defer pruntime.PrintExecutionTime(time.Now())
	os.Remove(dupeDetectionImageFingerprintDatabaseFilePath)

	db, err := sql.Open("sqlite3", dupeDetectionImageFingerprintDatabaseFilePath)
	if err != nil {
		return errors.New(err)
	}
	defer db.Close()

	dupeDetectionImageFingerprintDatabaseCreationQuery := `
	CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file, model_1_image_fingerprint_vector array, model_2_image_fingerprint_vector array, model_3_image_fingerprint_vector array,
		model_4_image_fingerprint_vector array, model_5_image_fingerprint_vector array, model_6_image_fingerprint_vector array, model_7_image_fingerprint_vector array, datetime_fingerprint_added_to_database TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		PRIMARY KEY (sha256_hash_of_art_image_file));
	`
	_, err = db.Exec(dupeDetectionImageFingerprintDatabaseCreationQuery)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func checkIfFilePathIsAValidImage(filePath string) error {
	imageHeader, err := goimghdr.What(filePath)
	if err != nil {
		return err
	}

	if imageHeader == "gif" || imageHeader == "jpeg" || imageHeader == "png" || imageHeader == "bmp" {
		return nil
	}
	return errors.New("Image header is not supported.")
}

func getAllValidImageFilePathsInFolder(artFolderPath string) ([]string, error) {
	jpgMatches, err := filepath.Glob(filepath.Join(artFolderPath, "*.jpg"))
	if err != nil {
		return nil, errors.New(err)
	}

	jpegMatches, err := filepath.Glob(filepath.Join(artFolderPath, "*.jpeg"))
	if err != nil {
		return nil, errors.New(err)
	}

	pngMatches, err := filepath.Glob(filepath.Join(artFolderPath, "*.png"))
	if err != nil {
		return nil, errors.New(err)
	}

	bmpMatches, err := filepath.Glob(filepath.Join(artFolderPath, "*.bmp"))
	if err != nil {
		return nil, errors.New(err)
	}

	gifMatches, err := filepath.Glob(filepath.Join(artFolderPath, "*.gif"))
	if err != nil {
		return nil, errors.New(err)
	}

	allMatches := append(append(append(append(jpgMatches, jpegMatches...), pngMatches...), bmpMatches...), gifMatches...)
	var results []string
	for _, match := range allMatches {
		if err = checkIfFilePathIsAValidImage(match); err == nil {
			results = append(results, match)
		}
	}
	return results, nil
}

func getImageHashFromImageFilePath(sampleImageFilePath string) (string, error) {
	f, err := os.Open(sampleImageFilePath)
	if err != nil {
		return "", errors.New(err)
	}

	defer f.Close()
	hash := sha3.New256()
	if _, err := io.Copy(hash, f); err != nil {
		return "", errors.New(err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func toBytes(data []float64) []byte {
	output := new(bytes.Buffer)
	_ = binary.Write(output, binary.LittleEndian, data)
	return output.Bytes()
}

func fromBytes(data []byte) []float64 {
	output := make([]float64, len(data)/8)
	for i := range output {
		bits := binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
		output[i] = math.Float64frombits(bits)
	}
	return output
}

func addImageFingerprintsToDupeDetectionDatabase(imageFilePath string) error {
	fingerprints, err := dupedetection.ComputeImageDeepLearningFeatures(imageFilePath)
	if err != nil {
		return errors.New(err)
	}

	imageHash, err := getImageHashFromImageFilePath(imageFilePath)
	if err != nil {
		return errors.New(err)
	}

	db, err := sql.Open("sqlite3", dupeDetectionImageFingerprintDatabaseFilePath)
	if err != nil {
		return errors.New(err)
	}
	defer db.Close()

	dataInsertionQuery := `
		INSERT OR REPLACE INTO image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file, path_to_art_image_file,
			model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector,
			model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	tx, err := db.Begin()
	if err != nil {
		return errors.New(err)
	}
	stmt, err := tx.Prepare(dataInsertionQuery)
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(imageHash, imageFilePath, toBytes(fingerprints[0]), toBytes(fingerprints[1]), toBytes(fingerprints[2]), toBytes(fingerprints[3]), toBytes(fingerprints[4]), toBytes(fingerprints[5]), toBytes(fingerprints[6]))
	if err != nil {
		return errors.New(err)
	}
	tx.Commit()

	return nil
}

func addAllImagesInFolderToImageFingerprintDatabase(artFolderPath string) error {
	validImageFilePaths, err := getAllValidImageFilePathsInFolder(artFolderPath)
	if err != nil {
		return errors.New(err)
	}
	for _, currentImageFilePath := range validImageFilePaths {
		fmt.Printf("\nNow adding image file %v to image fingerprint database.", currentImageFilePath)
		err = addImageFingerprintsToDupeDetectionDatabase(currentImageFilePath)
		if err != nil {
			return errors.New(err)
		}
	}
	return nil
}

func getListOfAllRegisteredImageFileHashes() ([]string, error) {
	db, err := sql.Open("sqlite3", dupeDetectionImageFingerprintDatabaseFilePath)
	if err != nil {
		return nil, errors.New(err)
	}
	defer db.Close()

	selectQuery := "SELECT sha256_hash_of_art_image_file FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC"
	rows, err := db.Query(selectQuery)
	if err != nil {
		return nil, errors.New(err)
	}
	defer rows.Close()

	var hashes []string
	for rows.Next() {
		var imageHash string
		err = rows.Scan(&imageHash)
		if err != nil {
			return nil, errors.New(err)
		}
		hashes = append(hashes, imageHash)
	}
	return hashes, nil
}

func getAllImageFingerprintsFromDupeDetectionDatabaseAsArray() ([][]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	hashes, err := getListOfAllRegisteredImageFileHashes()
	if err != nil {
		return nil, errors.New(err)
	}

	db, err := sql.Open("sqlite3", dupeDetectionImageFingerprintDatabaseFilePath)
	if err != nil {
		return nil, errors.New(err)
	}
	defer db.Close()

	var arrayOfCombinedImageFingerprintRows [][]float64

	for _, currentImageFileHash := range hashes {
		selectQuery := `
			SELECT path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector,
				model_6_image_fingerprint_vector, model_7_image_fingerprint_vector FROM image_hash_to_image_fingerprint_table where sha256_hash_of_art_image_file = ? ORDER BY datetime_fingerprint_added_to_database DESC
		`
		rows, err := db.Query(selectQuery, currentImageFileHash)
		if err != nil {
			return nil, errors.New(err)
		}
		defer rows.Close()

		for rows.Next() {
			var currentImageFilePath string
			var model1ImageFingerprintVector, model2ImageFingerprintVector, model3ImageFingerprintVector, model4ImageFingerprintVector, model5ImageFingerprintVector, model6ImageFingerprintVector, model7ImageFingerprintVector []byte
			err = rows.Scan(&currentImageFilePath, &model1ImageFingerprintVector, &model2ImageFingerprintVector, &model3ImageFingerprintVector, &model4ImageFingerprintVector, &model5ImageFingerprintVector, &model6ImageFingerprintVector, &model7ImageFingerprintVector)
			if err != nil {
				return nil, errors.New(err)
			}
			combinedImageFingerprintVector := append(append(append(append(append(append(fromBytes(model1ImageFingerprintVector), fromBytes(model2ImageFingerprintVector)[:]...), fromBytes(model3ImageFingerprintVector)[:]...), fromBytes(model4ImageFingerprintVector)[:]...), fromBytes(model5ImageFingerprintVector)[:]...), fromBytes(model6ImageFingerprintVector)[:]...), fromBytes(model7ImageFingerprintVector)[:]...)
			arrayOfCombinedImageFingerprintRows = append(arrayOfCombinedImageFingerprintRows, combinedImageFingerprintVector)
		}
	}
	return arrayOfCombinedImageFingerprintRows, nil
}

func getImageDeepLearningFeaturesCombinedVectorForSingleImage(artImageFilePath string) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	fingerprints, err := dupedetection.ComputeImageDeepLearningFeatures(artImageFilePath)
	if err != nil {
		return nil, errors.New(err)
	}
	var combinedImageFingerprintVector []float64
	for _, fingerprint := range fingerprints {
		combinedImageFingerprintVector = append(combinedImageFingerprintVector, fingerprint...)
	}
	cacheFingerprint(fingerprints, artImageFilePath)
	return combinedImageFingerprintVector, err
}

func computePearsonRForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorPearsonAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			//similarityScoreVectorPearsonAll[currentIndex] = wdm.Wdm(candidateImageFingerprint, currentFingerprint, "pearson", []float64{})
			similarityScoreVectorPearsonAll[currentIndex], err = stats.Pearson(candidateImageFingerprint, currentFingerprint)
			if err != nil {
				return err
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

func computeSpearmanForAllFingerprintPairs(candidateImageFingerprint []float64, finalCombinedImageFingerprintArray [][]float64) ([]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	similarityScoreVectorSpearmanAll := make([]float64, len(finalCombinedImageFingerprintArray))
	var err error
	g, _ := errgroup.WithContext(context.Background())
	for i, fingerprint := range finalCombinedImageFingerprintArray {
		currentIndex := i
		currentFingerprint := fingerprint
		g.Go(func() error {
			//similarityScoreVectorSpearmanAll[currentIndex] = wdm.Wdm(candidateImageFingerprint, currentFingerprint, "spearman", []float64{})
			similarityScoreVectorSpearmanAll[currentIndex], err = Spearman2(candidateImageFingerprint, currentFingerprint)
			if err != nil {
				return err
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

func filterOutFingerprintsByThreshold(similarityScoreVector []float64, threshold float64, combinedImageFingerprintArray [][]float64) [][]float64 {
	defer pruntime.PrintExecutionTime(time.Now())
	var filteredByThresholdCombinedImageFingerprintArray [][]float64
	for i, valueToCheck := range similarityScoreVector {
		if !math.IsNaN(valueToCheck) && valueToCheck >= threshold {
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

func computeAverageAndStdevOf25thTo75thPercentile(inputVector []float64) (float64, float64) {
	percentile25, _ := stats.Percentile(inputVector, 25)
	percentile75, _ := stats.Percentile(inputVector, 75)

	trimmedVector := filterOutArrayValuesByRange(inputVector, percentile25, percentile75)
	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func computeAverageAndStdevOf50thTo90thPercentile(inputVector []float64) (float64, float64) {
	percentile50, _ := stats.Percentile(inputVector, 50)
	percentile90, _ := stats.Percentile(inputVector, 90)

	trimmedVector := filterOutArrayValuesByRange(inputVector, percentile50, percentile90)
	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func computeParallelBootstrappedKendallsTau(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int) ([]float64, []float64, error) {
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

				bootstrappedKendallTauResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "kendall", []float64{})
			}
			robustAverageTau[fingerprintIdx], robustStdevTau[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedKendallTauResults)
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
				ranks[outputIdx] += 1
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

func computeRandomizedDependence(x, y []float64) float64 {
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
		if k == 4.5 {
			//fmt.Printf("\n%v", k)
		}
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
				ub -= 1
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

func computeParallelBootstrappedRandomizedDependence(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int) ([]float64, []float64, error) {
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

				bootstrappedRandomizedDependenceResults[currentIdx] = computeRandomizedDependence(xBootstraps[currentIdx], yBootstraps[currentIdx])

			}

			robustAverageRandomizedDependence[fingerprintIdx], robustStdevRandomizedDependence[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedRandomizedDependenceResults)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}
	return robustAverageRandomizedDependence, robustStdevRandomizedDependence, nil
}

func computeParallelBootstrappedBaggedHoeffdingsDSmallerSampleSize(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	for fingerprintIdx, y := range arrayOfFingerprintsRequiringFurtherTesting {
		arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
		xBootstraps := make([][]float64, numberOfBootstraps)
		yBootstraps := make([][]float64, numberOfBootstraps)
		bootstrappedHoeffdingsDResults := make([]float64, numberOfBootstraps)
		for i := 0; i < numberOfBootstraps; i++ {
			arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
		}
		for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
			xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
			yBootstraps[i] = arrayValuesFromIndexes(y, currentBootstrapIndices)

			bootstrappedHoeffdingsDResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "hoeffding", []float64{})
		}
		robustAverageD[fingerprintIdx], robustStdevD[fingerprintIdx] = computeAverageAndStdevOf25thTo75thPercentile(bootstrappedHoeffdingsDResults)
	}
	return robustAverageD, robustStdevD, nil
}

func computeParallelBootstrappedBaggedHoeffdingsD(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int) ([]float64, []float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	originalLengthOfInput := len(x)
	robustAverageD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))
	robustStdevD := make([]float64, len(arrayOfFingerprintsRequiringFurtherTesting))

	for fingerprintIdx, y := range arrayOfFingerprintsRequiringFurtherTesting {
		arrayOfBootstrapSampleIndices := make([][]int, numberOfBootstraps)
		xBootstraps := make([][]float64, numberOfBootstraps)
		yBootstraps := make([][]float64, numberOfBootstraps)
		bootstrappedHoeffdingsDResults := make([]float64, numberOfBootstraps)
		for i := 0; i < numberOfBootstraps; i++ {
			arrayOfBootstrapSampleIndices[i] = randInts(0, originalLengthOfInput-1, sampleSize)
		}
		for i, currentBootstrapIndices := range arrayOfBootstrapSampleIndices {
			xBootstraps[i] = arrayValuesFromIndexes(x, currentBootstrapIndices)
			yBootstraps[i] = arrayValuesFromIndexes(y, currentBootstrapIndices)

			bootstrappedHoeffdingsDResults[i] = wdm.Wdm(xBootstraps[i], yBootstraps[i], "hoeffding", []float64{})
		}
		robustAverageD[fingerprintIdx], robustStdevD[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedHoeffdingsDResults)
	}
	return robustAverageD, robustStdevD, nil
}

func computeParallelBootstrappedBlomqvistBeta(x []float64, arrayOfFingerprintsRequiringFurtherTesting [][]float64, sampleSize int, numberOfBootstraps int) ([]float64, []float64, error) {
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

				blomqvistValue := wdm.Wdm(xBootstraps[i], yBootstraps[i], "blomqvist", []float64{})
				if !math.IsNaN(blomqvistValue) {
					bootstrappedBlomqvistResults = append(bootstrappedBlomqvistResults, blomqvistValue)
				}
			}
			robustAverageBlomqvist[fingerprintIdx], robustStdevBlomqvist[fingerprintIdx] = computeAverageAndStdevOf50thTo90thPercentile(bootstrappedBlomqvistResults)
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
type correlation func([]float64, [][]float64) ([]float64, error)
type bootstrappedCorrelation func([]float64, [][]float64, int, int) ([]float64, []float64, error)

type computeData struct {
	title                       string
	computeFunc                 correlation
	printFunc                   printInfo
	bootstrappedComputeFunc     bootstrappedCorrelation
	bootstrappedPrintFunc       bootstrappedPrintInfo
	sampleSize, bootstrapsCount int
	threshold                   float64
	totalFingerprintsCount      int
}

func computeSimilarity(candidateImageFingerprint []float64, fingerprints [][]float64, data computeData) ([][]float64, error) {
	if len(fingerprints) == 0 {
		return [][]float64{}, nil
	}
	defer pruntime.PrintExecutionTime(time.Now())
	fmt.Printf("\n\n-->Computing %v", data.title)

	var similarityScore, similarityScoreStdev []float64
	var err error

	if data.computeFunc != nil {
		similarityScore, err = data.computeFunc(candidateImageFingerprint, fingerprints)
		if err != nil {
			return nil, errors.New(err)
		}
		if data.printFunc != nil {
			data.printFunc(similarityScore)
		}
	}

	if data.bootstrappedComputeFunc != nil {
		similarityScore, similarityScoreStdev, err = data.bootstrappedComputeFunc(candidateImageFingerprint, fingerprints, data.sampleSize, data.bootstrapsCount)
		if err != nil {
			return nil, errors.New(err)
		}

		if data.bootstrappedPrintFunc != nil {
			data.bootstrappedPrintFunc(similarityScore, similarityScoreStdev)
		}
	}

	fingerprintsRequiringFurtherTesting := filterOutFingerprintsByThreshold(similarityScore, data.threshold, fingerprints)
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
	similarityScoreVectorBlomqvistAverage, _ := stats.Mean(similarityScore)
	fmt.Printf("\n Average for Blomqvist's beta: %.4f", strictnessFactor*similarityScoreVectorBlomqvistAverage)
}

func measureSimilarityOfCandidateImageToDatabase(imageFilePath string) (int, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	fmt.Printf("\nChecking if candidate image is a likely duplicate of a previously registered artwork:")
	fmt.Printf("\nRetrieving image fingerprints of previously registered images from local database...")

	pearsonDupeThreshold := 0.995
	spearmanDupeThreshold := 0.79
	kendallDupeThreshold := 0.70
	randomizedDependenceDupeThreshold := 0.79
	randomizedBlomqvistDupeThreshold := 0.7625
	hoeffdingDupeThreshold := 0.35
	hoeffdingRound2DupeThreshold := 0.23

	finalCombinedImageFingerprintArray, err := getAllImageFingerprintsFromDupeDetectionDatabaseAsArray()
	if err != nil {
		return 0, errors.New(err)
	}
	numberOfPreviouslyRegisteredImagesToCompare := len(finalCombinedImageFingerprintArray)
	lengthOfEachImageFingerprintVector := len(finalCombinedImageFingerprintArray[0])
	fmt.Printf("\nComparing candidate image to the fingerprints of %v previously registered images. Each fingerprint consists of %v numbers.", numberOfPreviouslyRegisteredImagesToCompare, lengthOfEachImageFingerprintVector)
	fmt.Printf("\nComputing image fingerprint of candidate image...")

	candidateImageFingerprint, err := fingerprintFromCache(imageFilePath)
	if err != nil {
		candidateImageFingerprint, err = getImageDeepLearningFeaturesCombinedVectorForSingleImage(imageFilePath)
	}
	lengthOfCandidateImageFingerprint := len(candidateImageFingerprint)
	fmt.Printf("\nCandidate image fingerpint consists from %v numbers", lengthOfCandidateImageFingerprint)
	if err != nil {
		return 0, errors.New(err)
	}

	totalFingerprintsCount := len(finalCombinedImageFingerprintArray)

	fingerprintsForFurtherTesting, err := computeSimilarity(candidateImageFingerprint, finalCombinedImageFingerprintArray, computeData{
		title:                  "Pearson's R, which is fast to compute (We only perform the slower tests on the fingerprints that have a high R).",
		computeFunc:            computePearsonRForAllFingerprintPairs,
		printFunc:              printPearsonRCalculationResults,
		threshold:              strictnessFactor * pearsonDupeThreshold,
		totalFingerprintsCount: totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                  "Spearman's Rho for selected fingerprints...",
		computeFunc:            computeSpearmanForAllFingerprintPairs,
		printFunc:              nil,
		threshold:              strictnessFactor * spearmanDupeThreshold,
		totalFingerprintsCount: totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Bootstrapped Kendall's Tau for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: computeParallelBootstrappedKendallsTau,
		bootstrappedPrintFunc:   printKendallTauCalculationResults,
		sampleSize:              50,
		bootstrapsCount:         100,
		threshold:               strictnessFactor * kendallDupeThreshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Boostrapped Randomized Dependence Coefficient for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: computeParallelBootstrappedRandomizedDependence,
		bootstrappedPrintFunc:   printRDCCalculationResults,
		sampleSize:              50,
		bootstrapsCount:         100,
		threshold:               strictnessFactor * randomizedDependenceDupeThreshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "bootstrapped Blomqvist's beta for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: computeParallelBootstrappedBlomqvistBeta,
		bootstrappedPrintFunc:   printBlomqvistBetaCalculationResults,
		sampleSize:              100,
		bootstrapsCount:         100,
		threshold:               strictnessFactor * randomizedBlomqvistDupeThreshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Hoeffding's D Round 1",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: computeParallelBootstrappedBaggedHoeffdingsDSmallerSampleSize,
		bootstrappedPrintFunc:   nil,
		sampleSize:              20,
		bootstrapsCount:         50,
		threshold:               strictnessFactor * hoeffdingDupeThreshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidateImageFingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Hoeffding's D Round 2",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: computeParallelBootstrappedBaggedHoeffdingsD,
		bootstrappedPrintFunc:   nil,
		sampleSize:              75,
		bootstrapsCount:         20,
		threshold:               strictnessFactor * hoeffdingRound2DupeThreshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return 0, errors.New(err)
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

func main() {
	defer pruntime.PrintExecutionTime(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	rand.Seed(time.Now().UnixNano())

	rootPastelFolderPath := ""

	miscMasternodeFilesFolderPath := filepath.Join(rootPastelFolderPath, "misc_masternode_files")
	dupeDetectionImageFingerprintDatabaseFilePath = filepath.Join(rootPastelFolderPath, "dupe_detection_image_fingerprint_database.sqlite")
	pathToAllRegisteredWorksForDupeDetection := filepath.Join(rootPastelFolderPath, "Animecoin_All_Finished_Works")
	dupeDetectionTestImagesBaseFolderPath := filepath.Join(rootPastelFolderPath, "dupe_detector_test_images")
	nonDupeTestImagesBaseFolderPath := filepath.Join(rootPastelFolderPath, "non_duplicate_test_images")

	if _, err := os.Stat(miscMasternodeFilesFolderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(miscMasternodeFilesFolderPath, 0770); err != nil {
			panic(err)
		}
	}

	dbFound := tryToFindLocalDatabaseFile()
	if !dbFound {
		fmt.Printf("\nGenerating new image fingerprint database...")
		regenerateEmptyDupeDetectionImageFingerprintDatabase()
		err := addAllImagesInFolderToImageFingerprintDatabase(pathToAllRegisteredWorksForDupeDetection)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	} else {
		fmt.Printf("\nFound existing image fingerprint database.")
	}

	fmt.Printf("\n\nNow testing duplicate-detection scheme on known near-duplicate images:")
	nearDuplicates, err := getAllValidImageFilePathsInFolder(dupeDetectionTestImagesBaseFolderPath)
	if err != nil {
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}
	dupeCounter := 0
	var predictedY []float64
	for _, nearDupeFilePath := range nearDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Near Duplicate Image: %v", nearDupeFilePath)
		isLikelyDupe, err := measureSimilarityOfCandidateImageToDatabase(nearDupeFilePath)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
		dupeCounter += isLikelyDupe
		predictedY = append(predictedY, float64(isLikelyDupe))
	}
	fmt.Printf("\n\n________________________________________________________________________________________________________________")
	fmt.Printf("\n________________________________________________________________________________________________________________")
	fmt.Printf("\nAccuracy Percentage in Detecting Near-Duplicate Images: %.2f %% from totally %v images", float32(dupeCounter)/float32(len(nearDuplicates))*100.0, len(nearDuplicates))

	fmt.Printf("\n\nNow testing duplicate-detection scheme on known non-duplicate images:")
	nonDuplicates, err := getAllValidImageFilePathsInFolder(nonDupeTestImagesBaseFolderPath)
	if err != nil {
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}
	nondupeCounter := 0
	for _, nonDupeFilePath := range nonDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Non-Duplicate Test Image: %v", nonDupeFilePath)
		isLikelyDupe, err := measureSimilarityOfCandidateImageToDatabase(nonDupeFilePath)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}

		if isLikelyDupe == 0 {
			nondupeCounter += 1
			predictedY = append(predictedY, 1.0)
		} else {
			predictedY = append(predictedY, 0.0)
		}
	}
	fmt.Printf("\n\n________________________________________________________________________________________________________________")
	fmt.Printf("\n________________________________________________________________________________________________________________")
	fmt.Printf("\nAccuracy Percentage in Detecting Non-Duplicate Images: %.2f %% from totally %v images", float32(nondupeCounter)/float32(len(nonDuplicates))*100.0, len(nonDuplicates))

	fmt.Printf("\n\n\n_______________________________Summary:_______________________________\n\n")
	fmt.Printf("\nAccuracy Percentage in Detecting Near-Duplicate Images: %.2f %% from totally %v images", float32(dupeCounter)/float32(len(nearDuplicates))*100.0, len(nearDuplicates))
	fmt.Printf("\nAccuracy Percentage in Detecting Non-Duplicate Images: %.2f %% from totally %v images\n", float32(nondupeCounter)/float32(len(nonDuplicates))*100.0, len(nonDuplicates))

	actualY := make([]float64, len(predictedY))
	for i := range actualY {
		actualY[i] = 1.0
	}

	Ytrue := mat.NewDense(len(predictedY), 1, predictedY)
	Yscores := mat.NewDense(len(actualY), 1, actualY)
	precision, recall, _ := metrics.PrecisionRecallCurve(Ytrue, Yscores, 1, nil)
	sort.Float64s(recall)
	auprcMetric := metrics.AUC(recall, precision)
	fmt.Printf("\nAcross all near-duplicate and non-duplicate test images, precision is %v and the Area Under the Precision-Recall Curve (AUPRC) is %.3f\n", precision, auprcMetric)
}
