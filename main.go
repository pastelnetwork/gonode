package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"

	"github.com/disintegration/imaging"

	"github.com/pastelnetwork/go-commons/errors"

	"encoding/binary"
	"encoding/hex"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"golang.org/x/crypto/sha3"

	"github.com/gonum/matrix/mat64"
)

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
		fmt.Println(predictions)
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

func bindRowsOfDataFrames(dataFrame dataframe.DataFrame, rBundDataDrame dataframe.DataFrame) dataframe.DataFrame {
	defer Measure(time.Now())
	return dataFrame.RBind(rBundDataDrame)
}

func measure_similarity_of_candidate_image_to_database_func(path_to_art_image_file string) (bool, error) {
	fmt.Printf("\nChecking if candidate image is a likely duplicate of a previously registered artwork:")
	fmt.Printf("\nRetrieving image fingerprints of previously registered images from local database...")

	final_combined_image_fingerprint_df, err := get_all_image_fingerprints_from_dupe_detection_database_as_dataframe_func()
	if err != nil {
		return false, errors.New(err)
	}
	rowCount, _ := final_combined_image_fingerprint_df.Dims()
	fmt.Printf("\nComparing candidate image to the fingerprints of %v previously registered images.", rowCount)
	return true, nil
}

func main() {
	root_pastel_folder_path := ""

	misc_masternode_files_folder_path := filepath.Join(root_pastel_folder_path, "misc_masternode_files")
	dupe_detection_image_fingerprint_database_file_path = filepath.Join(root_pastel_folder_path, "dupe_detection_image_fingerprint_database.sqlite")
	path_to_all_registered_works_for_dupe_detection := filepath.Join(root_pastel_folder_path, "Animecoin_All_Finished_Works")
	dupe_detection_test_images_base_folder_path := filepath.Join(root_pastel_folder_path, "dupe_detector_test_images")
	/*non_dupe_test_images_base_folder_path*/ _ = filepath.Join(root_pastel_folder_path, "non_duplicate_test_images")

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
	for _, nearDupeFilePath := range nearDuplicates {
		fmt.Printf("\n\n________________________________________________________________________________________________________________\n\n")
		fmt.Printf("\nCurrent Near Duplicate Image: %v", nearDupeFilePath)
		measure_similarity_of_candidate_image_to_database_func(nearDupeFilePath)
	}
}
