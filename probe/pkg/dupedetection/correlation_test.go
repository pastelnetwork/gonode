package dupedetection

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/montanaflynn/stats"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/probe/wdm"
)

func readInputFileIntoArrays(inputFile string) ([][]float64, [][]float64, []float64, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, nil, nil, errors.New(err)
	}
	defer file.Close()

	var inputs1 [][]float64
	var inputs2 [][]float64
	var outputs []float64

	reader := bufio.NewReaderSize(file, 200000)
	for {
		newLine, isPrefix, err := reader.ReadLine()
		if isPrefix {
			return nil, nil, nil, errors.New(errors.Errorf("Reader bugger size is too short to contain whole line. Please increased buffer size."))
		}
		if err != nil {
			break
		}
		if len(newLine) > 0 {
			var input1 []float64
			values := strings.Split(string(newLine), " ")
			for i := range values {
				value, err := strconv.ParseFloat(values[i], 64)
				if err == nil {
					input1 = append(input1, value)
				}
			}
			inputs1 = append(inputs1, input1)
		}

		newLine, isPrefix, err = reader.ReadLine()
		if isPrefix {
			return nil, nil, nil, errors.New(errors.Errorf("Reader bugger size is too short to contain whole line. Please increased buffer size."))
		}
		if err != nil {
			break
		}
		if len(newLine) > 0 {
			var input2 []float64
			values := strings.Split(string(newLine), " ")
			for i := range values {
				value, err := strconv.ParseFloat(values[i], 64)
				if err == nil {
					input2 = append(input2, value)
				}
			}
			inputs2 = append(inputs2, input2)
		}

		newLine, isPrefix, err = reader.ReadLine()
		if isPrefix {
			return nil, nil, nil, errors.New(errors.Errorf("Reader bugger size is too short to contain whole line. Please increased buffer size."))
		}
		if err != nil {
			break
		}

		if len(newLine) > 0 {
			value, err := strconv.ParseFloat(string(newLine), 64)
			if err == nil {
				outputs = append(outputs, value)
			}
		}
	}
	return inputs1, inputs2, outputs, nil
}

func TestCorrelations(t *testing.T) {

	spearmanFloatingPointAccuracy := 0.000000000000001
	pearsonFloatingPointAccuracy := 0.000000000001
	kendallFloatingPointAccuracy := 0.000000000000001

	inputs1, inputs2, outputs, err := readInputFileIntoArrays(filepath.Join("testdata", "spearmanrho_data.txt"))
	if err != nil {
		t.Error(err)
	}

	if len(inputs1) != len(inputs2) || len(inputs1) != len(outputs) {
		t.Errorf("Inputs length doesn't match.")
	}

	fmt.Printf("\nNumber of inputs loaded: %v", len(inputs1))
	differentsCount := 0
	for i := range inputs1 {
		result, err := Spearman(inputs1[i], inputs2[i])
		if err != nil {
			t.Error(err)
		}
		if result != outputs[i] && spearmanFloatingPointAccuracy < math.Abs(result-outputs[i]) {
			differentsCount++
			t.Errorf("Spearman Rho Calculated correlation doesn't match the expected result:\n%v\n%v", result, outputs[i])
		}
	}
	fmt.Printf("\nSpearman Rho - Found calculation differences: %v from total number of %v with floating point accuracy %v", differentsCount, len(inputs1), spearmanFloatingPointAccuracy)

	inputs1, inputs2, outputs, err = readInputFileIntoArrays(filepath.Join("testdata", "pearsonr_data.txt"))
	if err != nil {
		t.Error(err)
	}

	if len(inputs1) != len(inputs2) || len(inputs1) != len(outputs) {
		t.Errorf("Inputs length doesn't match.")
	}

	fmt.Printf("\nNumber of inputs loaded: %v", len(inputs1))
	differentsCount = 0
	for i := range inputs1 {
		result, err := stats.Pearson(inputs1[i], inputs2[i])
		if err != nil {
			t.Error(err)
		}
		if result != outputs[i] && pearsonFloatingPointAccuracy < math.Abs(result-outputs[i]) {
			differentsCount++
			t.Errorf("Pearson R Calculated correlation doesn't match the expected result:\n%v\n%v", result, outputs[i])
		}
	}
	fmt.Printf("\nPearson R - Found calculation differences: %v from total number of %v with floating point accuracy %v", differentsCount, len(inputs1), pearsonFloatingPointAccuracy)

	inputs1, inputs2, outputs, err = readInputFileIntoArrays(filepath.Join("testdata", "kendall_data.txt"))
	if err != nil {
		t.Error(err)
	}

	if len(inputs1) != len(inputs2) || len(inputs1) != len(outputs) {
		t.Errorf("Inputs length doesn't match.")
	}

	fmt.Printf("\nNumber of inputs loaded: %v", len(inputs1))
	differentsCount = 0
	for i := range inputs1 {
		result := wdm.Wdm(inputs1[i], inputs2[i], "kendall")
		if result != outputs[i] && kendallFloatingPointAccuracy < math.Abs(result-outputs[i]) {
			differentsCount++
			t.Errorf("Kendall Calculated correlation doesn't match the expected result (diff is %v):\n%v\n%v", math.Abs(result-outputs[i]), result, outputs[i])
		}
	}
	fmt.Printf("\nKendall - Found calculation differences: %v from total number of %v with floating point accuracy %v", differentsCount, len(inputs1), kendallFloatingPointAccuracy)

	/*inputs1, inputs2, outputs, err = readInputFileIntoArrays("rdc_data.txt")
	if err != nil {
		t.Error(err)
	}

	if len(inputs1) != len(inputs2) || len(inputs1) != len(outputs) {
		t.Errorf("Inputs length doesn't match.")
	}

	fmt.Printf("\nNumber of inputs loaded: %v", len(inputs1))
	differentsCount = 0
	for i := range inputs1 {
		for repeat := 0; repeat < 10; repeat++ {
			result := ComputeRandomizedDependence(inputs1[i], inputs2[i])
			if result != outputs[i] {
				differentsCount++
				//t.Errorf("RDC Calculated correlation doesn't match the expected result (diff is %v):\n%v\n%v", math.Abs(result-outputs[i]), result, outputs[i])
				fmt.Printf("\n%v", result)
			}
		}
		fmt.Printf("\n--------------------------------")
	}
	fmt.Printf("\nRDC - Found calculation differences: %v from total number of %v with floating point accuracy %v", differentsCount, len(inputs1), kendallFloatingPointAccuracy)*/

}
