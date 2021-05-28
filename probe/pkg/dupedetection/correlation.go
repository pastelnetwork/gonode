package dupedetection

import (
	"math"

	onlinestats "github.com/dgryski/go-onlinestats"
	"github.com/kzahedi/goent/discrete"
	"github.com/montanaflynn/stats"
	"github.com/pastelnetwork/gonode/dupe-detection/wdm"
)

// Spearman calculates Spearman Rho correlation between arrays of input data
func Spearman(data1, data2 []float64) (float64, error) {
	r, _ := onlinestats.Spearman(data1, data2)
	return r, nil
}

// Pearson calculates Pearson R correlation between arrays of input data
func Pearson(data1, data2 []float64) (float64, error) {
	r, err := stats.Pearson(data1, data2)
	return r, err
}

// Kendall calculates Kendall Tau correlation between arrays of input data
func Kendall(data1, data2 []float64) (float64, error) {
	r := wdm.Wdm(data1, data2, "kendall")
	return r, nil
}

// HoeffdingD calculates HoeffdingD correlation between arrays of input data
func HoeffdingD(data1, data2 []float64) (float64, error) {
	r := wdm.Wdm(data1, data2, "hoeffding")
	return r, nil
}

// Blomqvist calculates Blomqvist Beta correlation between arrays of input data
func Blomqvist(data1, data2 []float64) (float64, error) {
	r := wdm.Wdm(data1, data2, "blomqvist")
	return r, nil
}

// MI calculates Mutual Information correlation between arrays of input data
func MI(data1, data2 []float64) (float64, error) {
	miInputPair := make([][]float64, 2)
	miInputPair[0] = data1
	miInputPair[1] = data2

	r := math.Pow(math.Abs(discrete.MutualInformationBase2(miInputPair)), 1.0/10.0)
	return r, nil
}
