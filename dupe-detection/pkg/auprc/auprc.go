package auprc

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pa-m/sklearn/metrics"
	"gonum.org/v1/gonum/mat"
	_ "gonum.org/v1/gonum/mat"

	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"

	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	_ "gorgonia.org/tensor"
)

const (
	cachedFingerprintsDB = "cachedFingerprints.sqlite"
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

	if len(results) > 30 {
		return results[:30], nil
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

func measureSimilarityOfCandidateImageToDatabase(imageFilePath string, config dupedetection.ComputeConfig) (int, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	fmt.Printf("\nChecking if candidate image is a likely duplicate of a previously registered artwork:")
	fmt.Printf("\nRetrieving image fingerprints of previously registered images from local database...")

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

	return dupedetection.MeasureImageSimilarity(candidateImageFingerprint, finalCombinedImageFingerprintArray, config)
}

func MeasureAUPRC(config dupedetection.ComputeConfig) (float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	rootPastelFolderPath := "./test_corpus"

	miscMasternodeFilesFolderPath := filepath.Join(rootPastelFolderPath, "misc_masternode_files")
	dupeDetectionImageFingerprintDatabaseFilePath = filepath.Join(rootPastelFolderPath, "dupe_detection_image_fingerprint_database.sqlite")
	pathToAllRegisteredWorksForDupeDetection := filepath.Join(rootPastelFolderPath, "dir_1")
	dupeDetectionTestImagesBaseFolderPath := filepath.Join(rootPastelFolderPath, "dir_2")
	nonDupeTestImagesBaseFolderPath := filepath.Join(rootPastelFolderPath, "dir_3")

	if _, err := os.Stat(miscMasternodeFilesFolderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(miscMasternodeFilesFolderPath, 0770); err != nil {
			return 0, errors.New(err)
		}
	}

	dbFound := tryToFindLocalDatabaseFile()
	if !dbFound {
		fmt.Printf("\nGenerating new image fingerprint database...")
		regenerateEmptyDupeDetectionImageFingerprintDatabase()
		err := addAllImagesInFolderToImageFingerprintDatabase(pathToAllRegisteredWorksForDupeDetection)
		if err != nil {
			return 0, errors.New(err)
		}
	} else {
		fmt.Printf("\nFound existing image fingerprint database.")
	}

	fmt.Printf("\n\nNow testing duplicate-detection scheme on known near-duplicate images:")
	nearDuplicates, err := getAllValidImageFilePathsInFolder(dupeDetectionTestImagesBaseFolderPath)
	if err != nil {
		if err != nil {
			return 0, errors.New(err)
		}
	}
	dupeCounter := 0
	var predictedY []float64
	for _, nearDupeFilePath := range nearDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Near Duplicate Image: %v", nearDupeFilePath)
		isLikelyDupe, err := measureSimilarityOfCandidateImageToDatabase(nearDupeFilePath, config)
		if err != nil {
			return 0, errors.New(err)
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
			return 0, errors.New(err)
		}
	}
	nondupeCounter := 0
	for _, nonDupeFilePath := range nonDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Non-Duplicate Test Image: %v", nonDupeFilePath)
		isLikelyDupe, err := measureSimilarityOfCandidateImageToDatabase(nonDupeFilePath, config)
		if err != nil {
			return 0, errors.New(err)
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

	if len(predictedY) == 0 {
		return 0, nil
	}

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
	return auprcMetric, nil
}
