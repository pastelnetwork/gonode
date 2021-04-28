package wdm

import (
	"github.com/pastelnetwork/gonode/dupe-detection/wdm/swig"
)

func Wdm(x, y []float64, method string, weights []float64) float64 {
	xVector := swig.NewDoubleVector(int64(len(x)))
	for i := range x {
		xVector.Set(i, x[i])
	}

	yVector := swig.NewDoubleVector(int64(len(y)))
	for i := range y {
		yVector.Set(i, y[i])
	}

	weightsVector := swig.NewDoubleVector(int64(len(weights)))
	for i := range weights {
		weightsVector.Set(i, weights[i])
	}

	return swig.Wdm(xVector, yVector, method, weightsVector)
}
