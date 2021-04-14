package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"database/sql"

	"github.com/corona10/goimghdr"
	_ "github.com/mattn/go-sqlite3"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"

	"github.com/disintegration/imaging"

	"github.com/pastelnetwork/go-commons/errors"
)

var dupe_detection_image_fingerprint_database_file_path string

func get_list_of_all_registered_image_file_hashes_func() error {
	return errors.New("Not implemented")
}

func regenerate_empty_dupe_detection_image_fingerprint_database_func() error {
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
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	img = imaging.Resize(img, width, height, imaging.Linear)
	return img, nil
}

func compute_image_deep_learning_features_func(path_to_art_image_file string) error {
	m, err := loadImage(path_to_art_image_file, 224, 224)
	if err != nil {
		log.Fatal(err)
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

	model := tg.LoadModel("NASNetLarge.tf", []string{"serve"}, nil)

	fakeInput, _ := tf.NewTensor(inputTensor)
	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_input_6", 0): fakeInput,
	})

	predictions := results[0]
	fmt.Println(predictions.Value())
	return nil
}

func add_image_fingerprints_to_dupe_detection_database_func(path_to_art_image_file string) error {
	compute_image_deep_learning_features_func(path_to_art_image_file)
	return nil
}

/*global dupe_detection_image_fingerprint_database_file_path
  model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector,  sha256_hash_of_art_image_file, dupe_detection_model_1, dupe_detection_model_2, dupe_detection_model_3, dupe_detection_model_4, dupe_detection_model_5, dupe_detection_model_6, dupe_detection_model_7 = compute_image_deep_learning_features_func(path_to_art_image_file)
  conn = sqlite3.connect(dupe_detection_image_fingerprint_database_file_path)
  c = conn.cursor()
  data_insertion_query_string = """INSERT OR REPLACE INTO image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file, path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);"""
  c.execute(data_insertion_query_string, [sha256_hash_of_art_image_file, path_to_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector])
  conn.commit()
  conn.close()
  return  model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector, model_5_image_fingerprint_vector, model_6_image_fingerprint_vector, model_7_image_fingerprint_vector*/

func add_all_images_in_folder_to_image_fingerprint_database_func(path_to_art_folder string) error {
	valid_image_file_paths, err := get_all_valid_image_file_paths_in_folder_func(path_to_art_folder)
	if err != nil {
		return errors.New(err)
	}
	for _, current_image_file_path := range valid_image_file_paths {
		fmt.Printf("\nNow adding image file %v to image fingerprint database.", current_image_file_path)
		add_image_fingerprints_to_dupe_detection_database_func(current_image_file_path)
	}
	return nil
}

func main() {
	root_pastel_folder_path := ""

	misc_masternode_files_folder_path := filepath.Join(root_pastel_folder_path, "misc_masternode_files")
	dupe_detection_image_fingerprint_database_file_path = filepath.Join(root_pastel_folder_path, "dupe_detection_image_fingerprint_database.sqlite")
	path_to_all_registered_works_for_dupe_detection := filepath.Join(root_pastel_folder_path, "Animecoin_All_Finished_Works")
	/*dupe_detection_test_images_base_folder_path*/ _ = filepath.Join(root_pastel_folder_path, "dupe_detector_test_images")
	/*non_dupe_test_images_base_folder_path*/ _ = filepath.Join(root_pastel_folder_path, "non_duplicate_test_images")

	if _, err := os.Stat(misc_masternode_files_folder_path); os.IsNotExist(err) {
		if err := os.MkdirAll(misc_masternode_files_folder_path, 0770); err != nil {
			panic(err)
		}
	}

	err := get_list_of_all_registered_image_file_hashes_func()
	if err != nil {
		regenerate_empty_dupe_detection_image_fingerprint_database_func()
		add_all_images_in_folder_to_image_fingerprint_database_func(path_to_all_registered_works_for_dupe_detection)
	}
}
