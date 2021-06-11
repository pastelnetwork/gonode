package test

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	//MobileNetV2 testing model
	MobileNetV2 = "MobileNetV2.tf"
	//NASNetMobile testing model
	NASNetMobile = "NASNetMobile.tf"

	//MobileNetV2Input input operation name of MobileNetV2 model
	MobileNetV2Input = "serving_default_input_1"
	//NASNetMobileInput input operation name of NASNetMobile model
	NASNetMobileInput = "serving_default_input_2"
)

var (
	testModels = []string{
		MobileNetV2,
		NASNetMobile,
	}
	testBaseDir string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	testBaseDir = regexp.MustCompile(`archive\.go`).ReplaceAllString(file, "")
}

//UnzipModels unzip testing models to temporary dir
func UnzipModels() (unzippedBaseDir string, unzippedModels []string, err error) {
	unzippedBaseDir, err = ioutil.TempDir("", "unzipped_models")
	if err != nil {
		return
	}
	var (
		modelBaseDir = filepath.Join(testBaseDir, "models")
		unzipped     []string
	)
	for _, m := range testModels {
		unzipped, err = unzip(filepath.Join(modelBaseDir, m)+".zip", unzippedBaseDir)
		if err != nil {
			return
		}
		unzippedModels = append(unzippedModels, unzipped...)
	}
	return
}

func unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
