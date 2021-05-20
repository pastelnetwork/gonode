package dupedetection

import (
	onlinestats "github.com/dgryski/go-onlinestats"
)

// Spearman calculates Spearman Rho correlation between arrays of input data
func Spearman(data1, data2 []float64) (float64, error) {
	r, _ := onlinestats.Spearman(data1, data2)
	return r, nil
}
