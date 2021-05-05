package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/profile"
	"gonum.org/v1/gonum/mat"
	_ "gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	"github.com/pastelnetwork/gonode/common/errors"
	pfmt "github.com/pastelnetwork/gonode/common/fmt"
	"github.com/pastelnetwork/gonode/dupe-detection/pkg/dupedetection"

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
	strictness_factor    = 0.985
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

var dupe_detection_image_fingerprint_database_file_path string

func tryToFindLocalDatabaseFile() bool {
	if _, err := os.Stat(dupe_detection_image_fingerprint_database_file_path); os.IsNotExist(err) {
		return false
	}
	return true
}

func regenerate_empty_dupe_detection_image_fingerprint_database_func() error {
	defer pfmt.PrintExecutionTime(time.Now())
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
	fingerprints, err := dupedetection.ComputeImageDeepLearningFeatures(path_to_art_image_file)
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
	defer pfmt.PrintExecutionTime(time.Now())

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
	defer pfmt.PrintExecutionTime(time.Now())
	return dataFrame.RBind(dataFrameToJoinRows)
}

func get_all_image_fingerprints_from_dupe_detection_database_as_array() ([][]float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())

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
	defer pfmt.PrintExecutionTime(time.Now())
	fingerprints, err := dupedetection.ComputeImageDeepLearningFeatures(path_to_art_image_file)
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
	defer pfmt.PrintExecutionTime(time.Now())
	similarity_score_vector__pearson_all := make([]float64, len(final_combined_image_fingerprint_array))
	var err error
	for i, fingerprint := range final_combined_image_fingerprint_array {
		//similarity_score_vector__pearson_all[i] = wdm.Wdm(candidate_image_fingerprint, fingerprint, "pearson", []float64{})
		similarity_score_vector__pearson_all[i], err = stats.Pearson(candidate_image_fingerprint, fingerprint)
		if err != nil {
			return nil, err
		}
	}
	return similarity_score_vector__pearson_all, nil
}

func computeSpearmanForAllFingerprintPairs(candidate_image_fingerprint []float64, final_combined_image_fingerprint_array [][]float64) ([]float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())
	similarity_score_vector__spearman := make([]float64, len(final_combined_image_fingerprint_array))
	var err error
	for i, fingerprint := range final_combined_image_fingerprint_array {
		//similarity_score_vector__spearman[i] = wdm.Wdm(candidate_image_fingerprint, fingerprint, "spearman", []float64{})
		similarity_score_vector__spearman[i], err = Spearman2(candidate_image_fingerprint, fingerprint)
		if err != nil {
			return nil, err
		}
	}
	return similarity_score_vector__spearman, nil
}

func filterOutFingerprintsByThreshold(similarity_score_vector__pearson_all []float64, threshold float64, final_combined_image_fingerprint_array [][]float64) [][]float64 {
	defer pfmt.PrintExecutionTime(time.Now())
	var filteredByThresholdCombinedImageFingerprintArray [][]float64
	for i, valueToCheck := range similarity_score_vector__pearson_all {
		if !math.IsNaN(valueToCheck) && valueToCheck >= threshold {
			filteredByThresholdCombinedImageFingerprintArray = append(filteredByThresholdCombinedImageFingerprintArray, final_combined_image_fingerprint_array[i])
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

func compute_parallel_bootstrapped_kendalls_tau_func(x []float64, list_of_fingerprints_requiring_further_testing_2 [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())
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
	return robust_average_tau, robust_stdev_tau, nil
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

func compute_parallel_bootstrapped_randomized_dependence_func(x []float64, list_of_fingerprints_requiring_further_testing_3 [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64, error) {
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
	return robust_average_randomized_dependence, robust_stdev_randomized_dependence, nil
}

func compute_parallel_bootstrapped_bagged_hoeffdings_d_smaller_sample_size_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())
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
	return robust_average_D, robust_stdev_D, nil
}

func compute_parallel_bootstrapped_bagged_hoeffdings_d_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())
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
	return robust_average_D, robust_stdev_D, nil
}

func compute_parallel_bootstrapped_blomqvist_beta_func(x []float64, list_of_fingerprints_requiring_further_testing [][]float64, sample_size int, number_of_bootstraps int) ([]float64, []float64, error) {
	defer pfmt.PrintExecutionTime(time.Now())
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
	return robust_average_blomqvist, robust_stdev_blomqvist, nil
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
	defer pfmt.PrintExecutionTime(time.Now())
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
	pearson_max, _ := stats.Max(results)
	fmt.Printf("\nLength of computed pearson r vector: %v with max value %v", len(results), pearson_max)
}

func printKendallTauCalculationResults(similarityScore, similarityScoreStdev []float64) {
	stdev_as_pct_of_robust_avg__kendall := computeAverageRatioOfArrays(similarityScoreStdev, similarityScore)
	fmt.Printf("\nStandard Deviation as %% of Average Tau -- average across all fingerprints: %.2f%%", stdev_as_pct_of_robust_avg__kendall*100)
}

func printRDCCalculationResults(similarityScore, similarityScoreStdev []float64) {
	stdev_as_pct_of_robust_avg__randomized_dependence := computeAverageRatioOfArrays(similarityScoreStdev, similarityScore)
	fmt.Printf("\nStandard Deviation as %% of Average Randomized Dependence -- average across all fingerprints: %.2f%%", stdev_as_pct_of_robust_avg__randomized_dependence*100)
}

func printBlomqvistBetaCalculationResults(similarityScore, similarityScoreStdev []float64) {
	similarity_score_vector__blomqvist_average, _ := stats.Mean(similarityScore)
	fmt.Printf("\n Average for Blomqvist's beta: %.4f", strictness_factor*similarity_score_vector__blomqvist_average)
}

func measure_similarity_of_candidate_image_to_database_func(path_to_art_image_file string) (bool, error) {
	defer pfmt.PrintExecutionTime(time.Now())
	fmt.Printf("\nChecking if candidate image is a likely duplicate of a previously registered artwork:")
	fmt.Printf("\nRetrieving image fingerprints of previously registered images from local database...")

	pearson__dupe_threshold := 0.995
	spearman__dupe_threshold := 0.79
	kendall__dupe_threshold := 0.70
	randomized_dependence__dupe_threshold := 0.79
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

	totalFingerprintsCount := len(final_combined_image_fingerprint_array)

	fingerprintsForFurtherTesting, err := computeSimilarity(candidate_image_fingerprint, final_combined_image_fingerprint_array, computeData{
		title:                  "Pearson's R, which is fast to compute (We only perform the slower tests on the fingerprints that have a high R).",
		computeFunc:            computePearsonRForAllFingerprintPairs,
		printFunc:              printPearsonRCalculationResults,
		threshold:              strictness_factor * pearson__dupe_threshold,
		totalFingerprintsCount: totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                  "Spearman's Rho for selected fingerprints...",
		computeFunc:            computeSpearmanForAllFingerprintPairs,
		printFunc:              nil,
		threshold:              strictness_factor * spearman__dupe_threshold,
		totalFingerprintsCount: totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Bootstrapped Kendall's Tau for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: compute_parallel_bootstrapped_kendalls_tau_func,
		bootstrappedPrintFunc:   printKendallTauCalculationResults,
		sampleSize:              50,
		bootstrapsCount:         100,
		threshold:               strictness_factor * kendall__dupe_threshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Boostrapped Randomized Dependence Coefficient for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: compute_parallel_bootstrapped_randomized_dependence_func,
		bootstrappedPrintFunc:   printRDCCalculationResults,
		sampleSize:              50,
		bootstrapsCount:         100,
		threshold:               strictness_factor * randomized_dependence__dupe_threshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "bootstrapped Blomqvist's beta for selected fingerprints...",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: compute_parallel_bootstrapped_blomqvist_beta_func,
		bootstrappedPrintFunc:   printBlomqvistBetaCalculationResults,
		sampleSize:              100,
		bootstrapsCount:         100,
		threshold:               strictness_factor * randomized_blomqvist__dupe_threshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Hoeffding's D Round 1",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: compute_parallel_bootstrapped_bagged_hoeffdings_d_smaller_sample_size_func,
		bootstrappedPrintFunc:   nil,
		sampleSize:              20,
		bootstrapsCount:         50,
		threshold:               strictness_factor * hoeffding__dupe_threshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	fingerprintsForFurtherTesting, err = computeSimilarity(candidate_image_fingerprint, fingerprintsForFurtherTesting, computeData{
		title:                   "Hoeffding's D Round 2",
		computeFunc:             nil,
		printFunc:               nil,
		bootstrappedComputeFunc: compute_parallel_bootstrapped_bagged_hoeffdings_d_func,
		bootstrappedPrintFunc:   nil,
		sampleSize:              75,
		bootstrapsCount:         20,
		threshold:               strictness_factor * hoeffding_round2__dupe_threshold,
		totalFingerprintsCount:  totalFingerprintsCount,
	})
	if err != nil {
		return false, errors.New(err)
	}

	if len(fingerprintsForFurtherTesting) > 0 {
		fmt.Printf("\n\nWARNING! Art image file appears to be a duplicate!")
	} else {
		fmt.Printf("\n\nArt image file appears to be original! (i.e., not a duplicate of an existing image in the image fingerprint database)")
	}

	return len(fingerprintsForFurtherTesting) != 0, nil
}

func main() {
	defer pfmt.PrintExecutionTime(time.Now())

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
