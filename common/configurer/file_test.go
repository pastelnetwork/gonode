package configurer

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type configurerSuite struct {
	suite.Suite
}

func (s *configurerSuite) AfterTest(_, _ string) {
	s.TearDown()
}

func (s *configurerSuite) TearDown() {
	fileName := []string{
		"./testing.json",
		"./testing.yml",
	}

	//cleaning up the mockup file after test
	for _, name := range fileName {
		os.Remove(name)
	}

}

func TestConfigurer(t *testing.T) {
	suite.Run(t, new(configurerSuite))
}

func (s *configurerSuite) TestParseFile_Success() {
	var config interface{}

	fileName := "./testing.json"
	sampleTestData := struct {
		Name    string
		Key     int
		Element string
	}{"pastel", 1, "data"}

	testData, _ := json.Marshal(sampleTestData)

	//creating the mockup file to test
	MockupTestFile(fileName, testData)

	err := ParseFile(fileName, &config)
	//check whether the service can able to get the value after parsing file
	assert.Equal(s.T(), sampleTestData.Name, config.(map[string]interface{})["name"])
	assert.Equal(s.T(), float64(sampleTestData.Key), config.(map[string]interface{})["key"])
	assert.Equal(s.T(), sampleTestData.Element, config.(map[string]interface{})["element"], "Check the file is parsed and able to get the data")
	//expected no error
	assert.Equal(s.T(), nil, err, "Err should be nil as Process is succeed")
}

func (s *configurerSuite) TestParseYmlFile_Success() {
	var config interface{}

	fileName := "./testing.yml"
	sampleTestData := struct {
		APIVersion string `yaml:apiVersion`
		Kind       string `yaml:kind`
		Name       string `yaml:name`
		Namespace  string `yaml:namespace`
	}{"v1", "test", "testApp", "test"}

	//creating the mockup yaml file to test
	testData, _ := yaml.Marshal(&sampleTestData)
	MockupTestFile(fileName, testData)

	err := ParseFile(fileName, &config)

	assert.Equal(s.T(), nil, err, "Err should be nil as Process is succeed")
	assert.Equal(s.T(), sampleTestData.Name, config.(map[string]interface{})["name"])
	assert.Equal(s.T(), sampleTestData.Kind, config.(map[string]interface{})["kind"])

}


func (s *configurerSuite) TestParseFile_FileNotFoundErr() {
	var config interface{}
	fileName := "randomfile.json"

	err := ParseFile(fileName, &config)

	assert.NotNil(s.T(), err, "Err should not be nil as Config File not found")
	assert.Equal(s.T(), nil, config, "Config interface should be nil as unable to parse the file")
}


func MockupTestFile(fileName string, sampleTestData []byte) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}

	_, err = file.Write(sampleTestData)
	if err != nil {
		file.Close()
		return
	}
	err = file.Close()
	if err != nil {
		return
	}
}