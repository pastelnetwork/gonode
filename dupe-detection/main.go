package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/profile"
	"gonum.org/v1/gonum/mat"
	_ "gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"

	"github.com/disintegration/imaging"

	"github.com/pastelnetwork/gonode/common/errors"

	"encoding/binary"
	"encoding/hex"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"golang.org/x/crypto/sha3"

	"github.com/gonum/matrix/mat64"

	"github.com/montanaflynn/stats"

	_ "gorgonia.org/tensor"

	"github.com/pastelnetwork/gonode/dupe-detection/wdm"
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
		var current_image_file_path string
		var model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector []byte
		err = rows.Scan(&current_image_file_path, &model_1_image_fingerprint_vector, &model_2_image_fingerprint_vector, &model_3_image_fingerprint_vector, &model_4_image_fingerprint_vector, &model_5_image_fingerprint_vector, &model_6_image_fingerprint_vector, &model_7_image_fingerprint_vector)
		if err != nil {
			return nil, errors.New(err)
		}
		combined_image_fingerprint_vector := append(append(append(append(append(append(fromBytes(model_1_image_fingerprint_vector), fromBytes(model_2_image_fingerprint_vector)[:]...), fromBytes(model_3_image_fingerprint_vector)[:]...), fromBytes(model_4_image_fingerprint_vector)[:]...), fromBytes(model_5_image_fingerprint_vector)[:]...), fromBytes(model_6_image_fingerprint_vector)[:]...), fromBytes(model_7_image_fingerprint_vector)[:]...)
		return combined_image_fingerprint_vector, nil
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

		dupe_detection_image_fingerprint_database_creation_string := `
			CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file, model_1_image_fingerprint_vector array, model_2_image_fingerprint_vector array, model_3_image_fingerprint_vector array,
				model_4_image_fingerprint_vector array, model_5_image_fingerprint_vector array, model_6_image_fingerprint_vector array, model_7_image_fingerprint_vector array, datetime_fingerprint_added_to_database TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
				PRIMARY KEY (sha256_hash_of_art_image_file));
			`
		_, err = db.Exec(dupe_detection_image_fingerprint_database_creation_string)
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

	data_insertion_query_string := `
		INSERT OR REPLACE INTO image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file, path_to_art_image_file,
			model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector,
			model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	tx, err := db.Begin()
	if err != nil {
		return errors.New(err)
	}
	stmt, err := tx.Prepare(data_insertion_query_string)
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

func Measure(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	pcFunc := runtime.FuncForPC(pc)
	funcNameOnly := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := funcNameOnly.ReplaceAllString(pcFunc.Name(), "$1")
	fmt.Printf("\n%s took %s\n", funcName, elapsed)
}

var dupe_detection_image_fingerprint_database_file_path string

func tryToFindLocalDatabaseFile() bool {
	if _, err := os.Stat(dupe_detection_image_fingerprint_database_file_path); os.IsNotExist(err) {
		return false
	}
	return true
}

func regenerate_empty_dupe_detection_image_fingerprint_database_func() error {
	defer Measure(time.Now())
	os.Remove(dupe_detection_image_fingerprint_database_file_path)

	db, err := sql.Open("sqlite3", dupe_detection_image_fingerprint_database_file_path)
	if err != nil {
		return errors.New(err)
	}
	defer db.Close()

	dupe_detection_image_fingerprint_database_creation_string := `
	CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file, model_1_image_fingerprint_vector array, model_2_image_fingerprint_vector array, model_3_image_fingerprint_vector array,
		model_4_image_fingerprint_vector array, model_5_image_fingerprint_vector array, model_6_image_fingerprint_vector array, model_7_image_fingerprint_vector array, datetime_fingerprint_added_to_database TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		PRIMARY KEY (sha256_hash_of_art_image_file));
	`
	_, err = db.Exec(dupe_detection_image_fingerprint_database_creation_string)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func check_if_file_path_is_a_valid_image_func(path_to_file string) error {
	imageHeader, err := goimghdr.What(path_to_file)
	if err != nil {
		return err
	}

	if imageHeader == "gif" || imageHeader == "jpeg" || imageHeader == "png" || imageHeader == "bmp" {
		return nil
	}
	return errors.New("Image header is not supported.")
}

func get_all_valid_image_file_paths_in_folder_func(path_to_art_folder string) ([]string, error) {
	jpgMatches, err := filepath.Glob(filepath.Join(path_to_art_folder, "*.jpg"))
	if err != nil {
		return nil, errors.New(err)
	}

	jpegMatches, err := filepath.Glob(filepath.Join(path_to_art_folder, "*.jpeg"))
	if err != nil {
		return nil, errors.New(err)
	}

	pngMatches, err := filepath.Glob(filepath.Join(path_to_art_folder, "*.png"))
	if err != nil {
		return nil, errors.New(err)
	}

	bmpMatches, err := filepath.Glob(filepath.Join(path_to_art_folder, "*.bmp"))
	if err != nil {
		return nil, errors.New(err)
	}

	gifMatches, err := filepath.Glob(filepath.Join(path_to_art_folder, "*.gif"))
	if err != nil {
		return nil, errors.New(err)
	}

	allMatches := append(append(append(append(jpgMatches, jpegMatches...), pngMatches...), bmpMatches...), gifMatches...)
	var results []string
	for _, match := range allMatches {
		if err = check_if_file_path_is_a_valid_image_func(match); err == nil {
			results = append(results, match)
		}
	}
	return results, nil
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

var models = make(map[string]*tg.Model)

type compute struct {
	model string
	input string
}

var fingerprintSources = []compute{
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

func compute_image_deep_learning_features_func(path_to_art_image_file string) ([][]float64, error) {
	defer Measure(time.Now())

	m, err := loadImage(path_to_art_image_file, 224, 224)
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
		//fmt.Println(predictions)
		fingerprints[i] = fromFloat32To64(predictions)
	}

	return fingerprints, nil
}

func fromFloat32To64(input []float32) []float64 {
	output := make([]float64, len(input))
	for i, value := range input {
		output[i] = float64(value)
	}
	return output
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

func add_image_fingerprints_to_dupe_detection_database_func(path_to_art_image_file string) error {
	fingerprints, err := compute_image_deep_learning_features_func(path_to_art_image_file)
	if err != nil {
		return errors.New(err)
	}

	imageHash, err := getImageHashFromImageFilePath(path_to_art_image_file)
	if err != nil {
		return errors.New(err)
	}

	db, err := sql.Open("sqlite3", dupe_detection_image_fingerprint_database_file_path)
	if err != nil {
		return errors.New(err)
	}
	defer db.Close()

	data_insertion_query_string := `
		INSERT OR REPLACE INTO image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file, path_to_art_image_file,
			model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector,
			model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	tx, err := db.Begin()
	if err != nil {
		return errors.New(err)
	}
	stmt, err := tx.Prepare(data_insertion_query_string)
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(imageHash, path_to_art_image_file, toBytes(fingerprints[0]), toBytes(fingerprints[1]), toBytes(fingerprints[2]), toBytes(fingerprints[3]), toBytes(fingerprints[4]), toBytes(fingerprints[5]), toBytes(fingerprints[6]))
	if err != nil {
		return errors.New(err)
	}
	tx.Commit()

	return nil
}

func add_all_images_in_folder_to_image_fingerprint_database_func(path_to_art_folder string) error {
	valid_image_file_paths, err := get_all_valid_image_file_paths_in_folder_func(path_to_art_folder)
	if err != nil {
		return errors.New(err)
	}
	for _, current_image_file_path := range valid_image_file_paths {
		fmt.Printf("\nNow adding image file %v to image fingerprint database.", current_image_file_path)
		err = add_image_fingerprints_to_dupe_detection_database_func(current_image_file_path)
		if err != nil {
			return errors.New(err)
		}
	}
	return nil
}

func get_list_of_all_registered_image_file_hashes_func() ([]string, error) {
	db, err := sql.Open("sqlite3", dupe_detection_image_fingerprint_database_file_path)
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

func get_all_image_fingerprints_from_dupe_detection_database_as_dataframe_func() (*dataframe.DataFrame, error) {
	defer Measure(time.Now())

	hashes, err := get_list_of_all_registered_image_file_hashes_func()
	if err != nil {
		return nil, errors.New(err)
	}

	db, err := sql.Open("sqlite3", dupe_detection_image_fingerprint_database_file_path)
	if err != nil {
		return nil, errors.New(err)
	}
	defer db.Close()

	var list_of_combined_image_fingerprint_rows [][]float64
	combined_image_fingerprint_df := dataframe.New(
		series.New([]string{}, series.String, "0"),
		series.New([]string{}, series.String, "1"),
	)

	for _, current_image_file_hash := range hashes {
		selectQuery := `
			SELECT path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector,
				model_6_image_fingerprint_vector, model_7_image_fingerprint_vector FROM image_hash_to_image_fingerprint_table where sha256_hash_of_art_image_file = ? ORDER BY datetime_fingerprint_added_to_database DESC
		`
		rows, err := db.Query(selectQuery, current_image_file_hash)
		if err != nil {
			return nil, errors.New(err)
		}
		defer rows.Close()

		for rows.Next() {
			var current_image_file_path string
			var model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector []byte
			err = rows.Scan(&current_image_file_path, &model_1_image_fingerprint_vector, &model_2_image_fingerprint_vector, &model_3_image_fingerprint_vector, &model_4_image_fingerprint_vector, &model_5_image_fingerprint_vector, &model_6_image_fingerprint_vector, &model_7_image_fingerprint_vector)
			if err != nil {
				return nil, errors.New(err)
			}
			combined_image_fingerprint_vector := append(append(append(append(append(append(fromBytes(model_1_image_fingerprint_vector), fromBytes(model_2_image_fingerprint_vector)[:]...), fromBytes(model_3_image_fingerprint_vector)[:]...), fromBytes(model_4_image_fingerprint_vector)[:]...), fromBytes(model_5_image_fingerprint_vector)[:]...), fromBytes(model_6_image_fingerprint_vector)[:]...), fromBytes(model_7_image_fingerprint_vector)[:]...)
			list_of_combined_image_fingerprint_rows = append(list_of_combined_image_fingerprint_rows, combined_image_fingerprint_vector)
			current_combined_image_fingerprint_df_row := dataframe.LoadRecords(
				[][]string{
					{"0", "1"},
					{current_image_file_hash, current_image_file_path},
				},
			)

			combined_image_fingerprint_df = combined_image_fingerprint_df.RBind(current_combined_image_fingerprint_df_row)
		}
	}
	var combined_image_fingerprint_df_vectors dataframe.DataFrame
	for _, current_combined_image_fingerprint_vector := range list_of_combined_image_fingerprint_rows {

		current_combined_image_fingerprint_vector_gonum := mat64.NewDense(1, len(current_combined_image_fingerprint_vector), current_combined_image_fingerprint_vector)
		current_combined_image_fingerprint_vector_df := dataframe.LoadMatrix(current_combined_image_fingerprint_vector_gonum)

		if rows, columns := combined_image_fingerprint_df_vectors.Dims(); rows == 0 && columns == 0 {
			combined_image_fingerprint_df_vectors = current_combined_image_fingerprint_vector_df
		} else {
			//combined_image_fingerprint_df_vectors = combined_image_fingerprint_df_vectors.RBind(current_combined_image_fingerprint_vector_df)
			combined_image_fingerprint_df_vectors = bindRowsOfDataFrames(combined_image_fingerprint_df_vectors, current_combined_image_fingerprint_vector_df)
		}
	}

	return &combined_image_fingerprint_df_vectors, nil
}

func bindRowsOfDataFrames(dataFrame dataframe.DataFrame, dataFrameToJoinRows dataframe.DataFrame) dataframe.DataFrame {
	defer Measure(time.Now())
	return dataFrame.RBind(dataFrameToJoinRows)
}

func get_all_image_fingerprints_from_dupe_detection_database_as_array() ([][]float64, error) {
	defer Measure(time.Now())

	hashes, err := get_list_of_all_registered_image_file_hashes_func()
	if err != nil {
		return nil, errors.New(err)
	}

	db, err := sql.Open("sqlite3", dupe_detection_image_fingerprint_database_file_path)
	if err != nil {
		return nil, errors.New(err)
	}
	defer db.Close()

	var list_of_combined_image_fingerprint_rows [][]float64

	for _, current_image_file_hash := range hashes {
		selectQuery := `
			SELECT path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector,
				model_6_image_fingerprint_vector, model_7_image_fingerprint_vector FROM image_hash_to_image_fingerprint_table where sha256_hash_of_art_image_file = ? ORDER BY datetime_fingerprint_added_to_database DESC
		`
		rows, err := db.Query(selectQuery, current_image_file_hash)
		if err != nil {
			return nil, errors.New(err)
		}
		defer rows.Close()

		for rows.Next() {
			var current_image_file_path string
			var model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector []byte
			err = rows.Scan(&current_image_file_path, &model_1_image_fingerprint_vector, &model_2_image_fingerprint_vector, &model_3_image_fingerprint_vector, &model_4_image_fingerprint_vector, &model_5_image_fingerprint_vector, &model_6_image_fingerprint_vector, &model_7_image_fingerprint_vector)
			if err != nil {
				return nil, errors.New(err)
			}
			combined_image_fingerprint_vector := append(append(append(append(append(append(fromBytes(model_1_image_fingerprint_vector), fromBytes(model_2_image_fingerprint_vector)[:]...), fromBytes(model_3_image_fingerprint_vector)[:]...), fromBytes(model_4_image_fingerprint_vector)[:]...), fromBytes(model_5_image_fingerprint_vector)[:]...), fromBytes(model_6_image_fingerprint_vector)[:]...), fromBytes(model_7_image_fingerprint_vector)[:]...)
			list_of_combined_image_fingerprint_rows = append(list_of_combined_image_fingerprint_rows, combined_image_fingerprint_vector)
		}
	}
	return list_of_combined_image_fingerprint_rows, nil
}

func get_image_deep_learning_features_combined_vector_for_single_image_func(path_to_art_image_file string) ([]float64, error) {
	defer Measure(time.Now())
	fingerprints, err := compute_image_deep_learning_features_func(path_to_art_image_file)
	if err != nil {
		return nil, errors.New(err)
	}
	var combined_image_fingerprint_vector []float64
	for _, fingerprint := range fingerprints {
		combined_image_fingerprint_vector = append(combined_image_fingerprint_vector, fingerprint...)
	}
	cacheFingerprint(fingerprints, path_to_art_image_file)
	return combined_image_fingerprint_vector, err
}

func computePearsonRForAllFingerprintPairs(candidate_image_fingerprint []float64, final_combined_image_fingerprint_array [][]float64) ([]float64, error) {
	defer Measure(time.Now())
	var similarity_score_vector__pearson_all []float64
	for _, fingerprint := range final_combined_image_fingerprint_array {
		pearsonR, err := stats.Pearson(candidate_image_fingerprint, fingerprint)
		if err != nil {
			return nil, err
		}
		similarity_score_vector__pearson_all = append(similarity_score_vector__pearson_all, pearsonR)
	}
	return similarity_score_vector__pearson_all, nil
}

func computeSpearmanForAllFingerprintPairs(candidate_image_fingerprint []float64, final_combined_image_fingerprint_array [][]float64) ([]float64, error) {
	defer Measure(time.Now())
	var similarity_score_vector__spearman []float64
	for _, fingerprint := range final_combined_image_fingerprint_array {
		spearmanCorrelation, err := Spearman2(candidate_image_fingerprint, fingerprint)
		if err != nil {
			return nil, err
		}
		similarity_score_vector__spearman = append(similarity_score_vector__spearman, spearmanCorrelation)
	}
	return similarity_score_vector__spearman, nil
}

func filterOutFingerprintsByTreshhold(similarity_score_vector__pearson_all []float64, threshold float64, final_combined_image_fingerprint_array [][]float64) [][]float64 {
	defer Measure(time.Now())
	var filteredByTreshholdCombinedImageFingerprintArray [][]float64
	for i, valueToCheck := range similarity_score_vector__pearson_all {
		if !math.IsNaN(valueToCheck) && valueToCheck >= threshold {
			filteredByTreshholdCombinedImageFingerprintArray = append(filteredByTreshholdCombinedImageFingerprintArray, final_combined_image_fingerprint_array[i])
		}
	}
	return filteredByTreshholdCombinedImageFingerprintArray
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

func compute_average_and_stdev_of_25th_to_75th_percentile_func(input_vector []float64) (float64, float64) {
	percentile25, _ := stats.Percentile(input_vector, 25)
	percentile75, _ := stats.Percentile(input_vector, 75)

	trimmedVector := filterOutArrayValuesByRange(input_vector, percentile25, percentile75)
	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func compute_average_and_stdev_of_50th_to_90th_percentile_func(input_vector []float64) (float64, float64) {
	percentile50, _ := stats.Percentile(input_vector, 50)
	percentile90, _ := stats.Percentile(input_vector, 90)

	trimmedVector := filterOutArrayValuesByRange(input_vector, percentile50, percentile90)
	trimmedVectorAvg, _ := stats.Mean(trimmedVector)
	trimmedVectorStdev, _ := stats.StdDevS(trimmedVector)

	return trimmedVectorAvg, trimmedVectorStdev
}

func compute_parallel_bootstrapped_kendalls_tau_func(x []float64, list_of_fingerprints_requiring_further_testing_2 [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64) {
	defer Measure(time.Now())
	original_length_of_input := len(x)
	robust_average_tau := make([]float64, len(list_of_fingerprints_requiring_further_testing_2))
	robust_stdev_tau := make([]float64, len(list_of_fingerprints_requiring_further_testing_2))

	for fingerprintIdx, y := range list_of_fingerprints_requiring_further_testing_2 {
		list_of_bootstrap_sample_indices := make([][]int, number_of_bootstraps)
		x_bootstraps := make([][]float64, number_of_bootstraps)
		y_bootstraps := make([][]float64, number_of_bootstraps)
		bootstrapped_kendalltau_results := make([]float64, number_of_bootstraps)
		for i := 0; i < number_of_bootstraps; i++ {
			list_of_bootstrap_sample_indices[i] = randInts(0, original_length_of_input-1, sample_size)
		}
		for i, current_bootstrap_indices := range list_of_bootstrap_sample_indices {
			x_bootstraps[i] = arrayValuesFromIndexes(x, current_bootstrap_indices)
			y_bootstraps[i] = arrayValuesFromIndexes(y, current_bootstrap_indices)

			bootstrapped_kendalltau_results[i] = wdm.Wdm(x_bootstraps[i], y_bootstraps[i], "kendall", []float64{})
		}
		robust_average_tau[fingerprintIdx], robust_stdev_tau[fingerprintIdx] = compute_average_and_stdev_of_50th_to_90th_percentile_func(bootstrapped_kendalltau_results)
	}
	return robust_average_tau, robust_stdev_tau
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

func compute_randomized_dependence_func(x, y []float64) float64 {
	sinFunc := func(i, j int, v float64) float64 {
		return math.Sin(v)
	}
	s := 1.0 / 6.0
	k := 20
	cx := mat.NewDense(len(x), 2, rankArrayOrdinalCopulaTransformation(x))
	cy := mat.NewDense(len(y), 2, rankArrayOrdinalCopulaTransformation(y))

	Rx := mat.NewDense(2, k, randomLinearProjection(s, k))
	Ry := mat.NewDense(2, k, randomLinearProjection(s, k))

	var X, Y, stackedH mat.Dense
	X.Mul(cx, Rx)
	Y.Mul(cy, Ry)

	X.Apply(sinFunc, &X)
	Y.Apply(sinFunc, &Y)

	stackedH.Augment(&X, &Y)
	stackedT := stackedH.T()

	var CsymDense mat.SymDense

	stat.CovarianceMatrix(&CsymDense, stackedT, nil)
	C := mat.DenseCopyOf(&CsymDense)

	//fmt.Printf("\n%v", mat.Formatted(C, mat.Prefix(""), mat.Squeeze()))

	k0 := k
	lb := 1
	ub := k
	var maxEigVal float64
	for {
		maxEigVal = 0
		Cxx := C.Slice(0, k, 0, k).(*mat.Dense)
		Cyy := C.Slice(k0, k0+k, k0, k0+k).(*mat.Dense)
		Cxy := C.Slice(0, k, k0, k0+k).(*mat.Dense)
		Cyx := C.Slice(k0, k0+k, 0, k).(*mat.Dense)

		var CxxInversed, CyyInversed mat.Dense
		CxxInversed.Inverse(Cxx)
		CyyInversed.Inverse(Cyy)

		var CxxInversedMulCxy, CyyInversedMulCyx, resultMul mat.Dense
		CxxInversedMulCxy.Mul(&CxxInversed, Cxy)
		CyyInversedMulCyx.Mul(&CyyInversed, Cyx)

		resultMul.Mul(&CxxInversedMulCxy, &CyyInversedMulCyx)

		var eigs mat.Eigen
		eigs.Factorize(&resultMul, 0)

		continueLoop := false
		eigsVals := eigs.Values(nil)
		for _, eigVal := range eigsVals {
			realVal := real(eigVal)
			imageVal := imag(eigVal)
			if imageVal != 0.0 || 0 >= realVal || realVal >= 1 {
				ub -= 1
				k = (ub + lb) / 2
				continueLoop = true
				break
			}
			maxEigVal = math.Max(maxEigVal, realVal)
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
			k = (ub + lb) / 2
		}
	}
	sqrtResult := math.Sqrt(maxEigVal)
	result := math.Pow(sqrtResult, 12)

	return result
}

func compute_parallel_bootstrapped_randomized_dependence_func(x []float64, list_of_fingerprints_requiring_further_testing_3 [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64) {
	original_length_of_input := len(x)
	robust_average_randomized_dependence := make([]float64, len(list_of_fingerprints_requiring_further_testing_3))
	robust_stdev_randomized_dependence := make([]float64, len(list_of_fingerprints_requiring_further_testing_3))

	for fingerprintIdx, y := range list_of_fingerprints_requiring_further_testing_3 {
		list_of_bootstrap_sample_indices := make([][]int, number_of_bootstraps)
		x_bootstraps := make([][]float64, number_of_bootstraps)
		y_bootstraps := make([][]float64, number_of_bootstraps)
		bootstrapped_randomized_dependence_results := make([]float64, number_of_bootstraps)
		for i := 0; i < number_of_bootstraps; i++ {
			list_of_bootstrap_sample_indices[i] = randInts(0, original_length_of_input-1, sample_size)
		}
		for i, current_bootstrap_indices := range list_of_bootstrap_sample_indices {
			x_bootstraps[i] = arrayValuesFromIndexes(x, current_bootstrap_indices)
			y_bootstraps[i] = arrayValuesFromIndexes(y, current_bootstrap_indices)
			bootstrapped_randomized_dependence_results[i] = compute_randomized_dependence_func(x_bootstraps[i], y_bootstraps[i])
		}
		robust_average_randomized_dependence[fingerprintIdx], robust_stdev_randomized_dependence[fingerprintIdx] = compute_average_and_stdev_of_50th_to_90th_percentile_func(bootstrapped_randomized_dependence_results)
	}
	return robust_average_randomized_dependence, robust_stdev_randomized_dependence
}

func compute_parallel_bootstrapped_bagged_hoeffdings_d_smaller_sample_size_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) []float64 {
	defer Measure(time.Now())
	original_length_of_input := len(x)
	robust_average_D := make([]float64, len(list_of_fingerprints_requiring_further_testing))
	robust_stdev_D := make([]float64, len(list_of_fingerprints_requiring_further_testing))

	for fingerprintIdx, y := range list_of_fingerprints_requiring_further_testing {
		list_of_bootstrap_sample_indices := make([][]int, number_of_bootstraps)
		x_bootstraps := make([][]float64, number_of_bootstraps)
		y_bootstraps := make([][]float64, number_of_bootstraps)
		bootstrapped_hoeffdings_d_results := make([]float64, number_of_bootstraps)
		for i := 0; i < number_of_bootstraps; i++ {
			list_of_bootstrap_sample_indices[i] = randInts(0, original_length_of_input-1, sample_size)
		}
		for i, current_bootstrap_indices := range list_of_bootstrap_sample_indices {
			x_bootstraps[i] = arrayValuesFromIndexes(x, current_bootstrap_indices)
			y_bootstraps[i] = arrayValuesFromIndexes(y, current_bootstrap_indices)

			bootstrapped_hoeffdings_d_results[i] = wdm.Wdm(x_bootstraps[i], y_bootstraps[i], "hoeffding", []float64{})
		}
		robust_average_D[fingerprintIdx], robust_stdev_D[fingerprintIdx] = compute_average_and_stdev_of_25th_to_75th_percentile_func(bootstrapped_hoeffdings_d_results)
	}
	return robust_average_D
}

func compute_parallel_bootstrapped_bagged_hoeffdings_d_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) []float64 {
	defer Measure(time.Now())
	original_length_of_input := len(x)
	robust_average_D := make([]float64, len(list_of_fingerprints_requiring_further_testing))
	robust_stdev_D := make([]float64, len(list_of_fingerprints_requiring_further_testing))

	for fingerprintIdx, y := range list_of_fingerprints_requiring_further_testing {
		list_of_bootstrap_sample_indices := make([][]int, number_of_bootstraps)
		x_bootstraps := make([][]float64, number_of_bootstraps)
		y_bootstraps := make([][]float64, number_of_bootstraps)
		bootstrapped_hoeffdings_d_results := make([]float64, number_of_bootstraps)
		for i := 0; i < number_of_bootstraps; i++ {
			list_of_bootstrap_sample_indices[i] = randInts(0, original_length_of_input-1, sample_size)
		}
		for i, current_bootstrap_indices := range list_of_bootstrap_sample_indices {
			x_bootstraps[i] = arrayValuesFromIndexes(x, current_bootstrap_indices)
			y_bootstraps[i] = arrayValuesFromIndexes(y, current_bootstrap_indices)

			bootstrapped_hoeffdings_d_results[i] = wdm.Wdm(x_bootstraps[i], y_bootstraps[i], "hoeffding", []float64{})
		}
		robust_average_D[fingerprintIdx], robust_stdev_D[fingerprintIdx] = compute_average_and_stdev_of_50th_to_90th_percentile_func(bootstrapped_hoeffdings_d_results)
	}
	return robust_average_D
}

func compute_parallel_bootstrapped_blomqvist_beta_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) []float64 {
	defer Measure(time.Now())
	original_length_of_input := len(x)
	robust_average_blomqvist := make([]float64, len(list_of_fingerprints_requiring_further_testing))
	robust_stdev_blomqvist := make([]float64, len(list_of_fingerprints_requiring_further_testing))

	for fingerprintIdx, y := range list_of_fingerprints_requiring_further_testing {
		list_of_bootstrap_sample_indices := make([][]int, number_of_bootstraps)
		x_bootstraps := make([][]float64, number_of_bootstraps)
		y_bootstraps := make([][]float64, number_of_bootstraps)
		var bootstrapped_blomqvist_results []float64
		for i := 0; i < number_of_bootstraps; i++ {
			list_of_bootstrap_sample_indices[i] = randInts(0, original_length_of_input-1, sample_size)
		}
		for i, current_bootstrap_indices := range list_of_bootstrap_sample_indices {
			x_bootstraps[i] = arrayValuesFromIndexes(x, current_bootstrap_indices)
			y_bootstraps[i] = arrayValuesFromIndexes(y, current_bootstrap_indices)

			blomqvistValue := wdm.Wdm(x_bootstraps[i], y_bootstraps[i], "blomqvist", []float64{})
			if !math.IsNaN(blomqvistValue) {
				bootstrapped_blomqvist_results = append(bootstrapped_blomqvist_results, blomqvistValue)
			}
		}
		robust_average_blomqvist[fingerprintIdx], robust_stdev_blomqvist[fingerprintIdx] = compute_average_and_stdev_of_50th_to_90th_percentile_func(bootstrapped_blomqvist_results)
	}
	return robust_average_blomqvist
}

func measure_similarity_of_candidate_image_to_database_func(path_to_art_image_file string) (bool, error) {
	fmt.Printf("\nChecking if candidate image is a likely duplicate of a previously registered artwork:")
	fmt.Printf("\nRetrieving image fingerprints of previously registered images from local database...")

	pearson__dupe_threshold := 0.995
	spearman__dupe_threshold := 0.79
	kendall__dupe_threshold := 0.70
	strictness_factor := 0.985
	randomized_blomqvist__dupe_threshold := 0.7625
	hoeffding__dupe_threshold := 0.35
	hoeffding_round2__dupe_threshold := 0.23

	final_combined_image_fingerprint_array, err := get_all_image_fingerprints_from_dupe_detection_database_as_array()
	if err != nil {
		return false, errors.New(err)
	}
	number_of_previously_registered_images_to_compare := len(final_combined_image_fingerprint_array)
	length_of_each_image_fingerprint_vector := len(final_combined_image_fingerprint_array[0])
	fmt.Printf("\nComparing candidate image to the fingerprints of %v previously registered images. Each fingerprint consists of %v numbers.", number_of_previously_registered_images_to_compare, length_of_each_image_fingerprint_vector)
	fmt.Printf("\nComputing image fingerprint of candidate image...")

	candidate_image_fingerprint, err := fingerprintFromCache(path_to_art_image_file)
	if err != nil {
		candidate_image_fingerprint, err = get_image_deep_learning_features_combined_vector_for_single_image_func(path_to_art_image_file)
	}
	length_of_candidate_image_fingerprint := len(candidate_image_fingerprint)
	fmt.Printf("\nCandidate image fingerpint consists from %v numbers", length_of_candidate_image_fingerprint)
	if err != nil {
		return false, errors.New(err)
	}
	fmt.Printf("\nComputing Pearson's R, which is fast to compute (We only perform the slower tests on the fingerprints that have a high R).")
	similarity_score_vector__pearson_all, err := computePearsonRForAllFingerprintPairs(candidate_image_fingerprint, final_combined_image_fingerprint_array)
	if err != nil {
		return false, errors.New(err)
	}
	pearson_max, _ := stats.Max(similarity_score_vector__pearson_all)

	fmt.Printf("\nLength of computed pearson r vector: %v with max value %v", len(similarity_score_vector__pearson_all), pearson_max)

	list_of_fingerprints_requiring_further_testing_1 := filterOutFingerprintsByTreshhold(similarity_score_vector__pearson_all, strictness_factor*pearson__dupe_threshold, final_combined_image_fingerprint_array)
	percentage_of_fingerprints_requiring_further_testing_1 := float32(len(list_of_fingerprints_requiring_further_testing_1)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_1), percentage_of_fingerprints_requiring_further_testing_1*100)

	similarity_score_vector__spearman, err := computeSpearmanForAllFingerprintPairs(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_1)
	if err != nil {
		return false, errors.New(err)
	}
	list_of_fingerprints_requiring_further_testing_2 := filterOutFingerprintsByTreshhold(similarity_score_vector__spearman, strictness_factor*spearman__dupe_threshold, list_of_fingerprints_requiring_further_testing_1)
	percentage_of_fingerprints_requiring_further_testing_2 := float32(len(list_of_fingerprints_requiring_further_testing_2)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_2), percentage_of_fingerprints_requiring_further_testing_2*100)

	fmt.Printf("\nNow computing Bootstrapped Kendall's Tau for selected fingerprints...")
	sample_size__kendall := 50
	number_of_bootstraps__kendall := 100
	similarity_score_vector__kendall, similarity_score_vector__kendall__stdev := compute_parallel_bootstrapped_kendalls_tau_func(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_2, sample_size__kendall, number_of_bootstraps__kendall)
	stdev_as_pct_of_robust_avg__kendall := computeAverageRatioOfArrays(similarity_score_vector__kendall__stdev, similarity_score_vector__kendall)
	fmt.Printf("\nStandard Deviation as %% of Average Tau -- average across all fingerprints: %.2f%%", stdev_as_pct_of_robust_avg__kendall*100)

	list_of_fingerprints_requiring_further_testing_3 := filterOutFingerprintsByTreshhold(similarity_score_vector__kendall, strictness_factor*kendall__dupe_threshold, list_of_fingerprints_requiring_further_testing_2)
	percentage_of_fingerprints_requiring_further_testing_3 := float32(len(list_of_fingerprints_requiring_further_testing_3)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_3), percentage_of_fingerprints_requiring_further_testing_3*100)

	fmt.Printf("\nNow computing Boostrapped Randomized Dependence Coefficient for selected fingerprints...")
	sample_size__randomized_dep := 50
	number_of_bootstraps__randomized_dep := 100
	similarity_score_vector__randomized_dependence, similarity_score_vector__randomized_dependence__stdev := compute_parallel_bootstrapped_randomized_dependence_func(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_3, sample_size__randomized_dep, number_of_bootstraps__randomized_dep)
	stdev_as_pct_of_robust_avg__randomized_dependence := computeAverageRatioOfArrays(similarity_score_vector__randomized_dependence__stdev, similarity_score_vector__randomized_dependence)
	fmt.Printf("\nStandard Deviation as %% of Average Randomized Dependence -- average across all fingerprints: %.2f%%", stdev_as_pct_of_robust_avg__randomized_dependence*100)

	fmt.Printf("\nNow computing bootstrapped Blomqvist's beta for selected fingerprints...")
	sample_size_blomqvist := 100
	number_of_bootstraps_blomqvist := 100
	similarity_score_vector__blomqvist := compute_parallel_bootstrapped_blomqvist_beta_func(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_3, sample_size_blomqvist, number_of_bootstraps_blomqvist)
	similarity_score_vector__blomqvist_average, _ := stats.Mean(similarity_score_vector__blomqvist)
	fmt.Printf("\n Average for Blomqvist's beta: %.4f", strictness_factor*similarity_score_vector__blomqvist_average)

	list_of_fingerprints_requiring_further_testing_5 := filterOutFingerprintsByTreshhold(similarity_score_vector__blomqvist, strictness_factor*randomized_blomqvist__dupe_threshold, list_of_fingerprints_requiring_further_testing_3)
	percentage_of_fingerprints_requiring_further_testing_5 := float32(len(list_of_fingerprints_requiring_further_testing_5)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_5), percentage_of_fingerprints_requiring_further_testing_5*100)

	fmt.Printf("\nNow computing bootstrapped Hoeffding's D for selected fingerprints...")
	sample_size := 20
	number_of_bootstraps := 50
	fmt.Printf("\nHoeffding Round 1 | Sample Size: %v; Number of Bootstraps: %v", sample_size, number_of_bootstraps)
	similarity_score_vector__hoeffding := compute_parallel_bootstrapped_bagged_hoeffdings_d_smaller_sample_size_func(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_5, sample_size, number_of_bootstraps)

	list_of_fingerprints_requiring_further_testing_6 := filterOutFingerprintsByTreshhold(similarity_score_vector__hoeffding, strictness_factor*hoeffding__dupe_threshold, list_of_fingerprints_requiring_further_testing_5)
	percentage_of_fingerprints_requiring_further_testing_6 := float32(len(list_of_fingerprints_requiring_further_testing_6)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_6), percentage_of_fingerprints_requiring_further_testing_6*100)

	fmt.Printf("\nNow computing second round of bootstrapped Hoeffding's D for selected fingerprints using smaller sample size...")
	sample_size__round2 := 75
	number_of_bootstraps__round2 := 20
	fmt.Printf("\nHoeffding Round 2 | Sample Size: %v; Number of Bootstraps: %v", sample_size, number_of_bootstraps)
	similarity_score_vector__hoeffding_round2 := compute_parallel_bootstrapped_bagged_hoeffdings_d_func(candidate_image_fingerprint, list_of_fingerprints_requiring_further_testing_6, sample_size__round2, number_of_bootstraps__round2)

	list_of_fingerprints_requiring_further_testing_7 := filterOutFingerprintsByTreshhold(similarity_score_vector__hoeffding_round2, strictness_factor*hoeffding_round2__dupe_threshold, list_of_fingerprints_requiring_further_testing_6)
	percentage_of_fingerprints_requiring_further_testing_7 := float32(len(list_of_fingerprints_requiring_further_testing_7)) / float32(len(final_combined_image_fingerprint_array))
	fmt.Printf("\nSelected %v fingerprints for further testing(%.2f%% of the total registered fingerprints).", len(list_of_fingerprints_requiring_further_testing_7), percentage_of_fingerprints_requiring_further_testing_7*100)

	if len(list_of_fingerprints_requiring_further_testing_7) > 0 {
		fmt.Printf("\n\nWARNING! Art image file appears to be a duplicate!")
	} else {
		fmt.Printf("\n\nArt image file appears to be original! (i.e., not a duplicate of an existing image in the image fingerprint database)")
	}

	return len(list_of_fingerprints_requiring_further_testing_7) != 0, nil
}

func main() {
	defer Measure(time.Now())

	defer profile.Start(profile.ProfilePath(".")).Stop()

	root_pastel_folder_path := ""

	misc_masternode_files_folder_path := filepath.Join(root_pastel_folder_path, "misc_masternode_files")
	dupe_detection_image_fingerprint_database_file_path = filepath.Join(root_pastel_folder_path, "dupe_detection_image_fingerprint_database.sqlite")
	path_to_all_registered_works_for_dupe_detection := filepath.Join(root_pastel_folder_path, "Animecoin_All_Finished_Works")
	dupe_detection_test_images_base_folder_path := filepath.Join(root_pastel_folder_path, "dupe_detector_test_images")
	non_dupe_test_images_base_folder_path := filepath.Join(root_pastel_folder_path, "non_duplicate_test_images")

	if _, err := os.Stat(misc_masternode_files_folder_path); os.IsNotExist(err) {
		if err := os.MkdirAll(misc_masternode_files_folder_path, 0770); err != nil {
			panic(err)
		}
	}

	dbFound := tryToFindLocalDatabaseFile()
	if !dbFound {
		fmt.Printf("\nGenerating new image fingerprint database...")
		regenerate_empty_dupe_detection_image_fingerprint_database_func()
		err := add_all_images_in_folder_to_image_fingerprint_database_func(path_to_all_registered_works_for_dupe_detection)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	} else {
		fmt.Printf("\nFound existing image fingerprint database.")
	}

	fmt.Printf("\n\nNow testing duplicate-detection scheme on known near-duplicate images:")
	nearDuplicates, err := get_all_valid_image_file_paths_in_folder_func(dupe_detection_test_images_base_folder_path)
	if err != nil {
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}
	dupe_counter := 0
	for _, nearDupeFilePath := range nearDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Near Duplicate Image: %v", nearDupeFilePath)
		is_likely_dupe, err := measure_similarity_of_candidate_image_to_database_func(nearDupeFilePath)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
		if is_likely_dupe {
			dupe_counter++
		}
	}
	fmt.Printf("\n\n________________________________________________________________________________________________________________")
	fmt.Printf("\n________________________________________________________________________________________________________________")
	fmt.Printf("\nAccuracy Percentage in Detecting Near-Duplicate Images: %.2f %% from totally %v images", float32(dupe_counter)/float32(len(nearDuplicates))*100.0, len(nearDuplicates))

	fmt.Printf("\n\nNow testing duplicate-detection scheme on known non-duplicate images:")
	nonDuplicates, err := get_all_valid_image_file_paths_in_folder_func(non_dupe_test_images_base_folder_path)
	if err != nil {
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}
	nondupe_counter := 0
	for _, nonDupeFilePath := range nonDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Non-Duplicate Test Image: %v", nonDupeFilePath)
		is_likely_dupe, err := measure_similarity_of_candidate_image_to_database_func(nonDupeFilePath)
		if err != nil {
			fmt.Println(err.(*errors.Error).ErrorStack())
			panic(err)
		}
		if !is_likely_dupe {
			nondupe_counter++
		}
	}
	fmt.Printf("\n\n________________________________________________________________________________________________________________")
	fmt.Printf("\n________________________________________________________________________________________________________________")
	fmt.Printf("\nAccuracy Percentage in Detecting Non-Duplicate Images: %.2f %% from totally %v images", float32(nondupe_counter)/float32(len(nonDuplicates))*100.0, len(nonDuplicates))
}
