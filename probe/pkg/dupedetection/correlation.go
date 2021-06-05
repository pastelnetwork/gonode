package dupedetection

import (
	"fmt"
	"math"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"

	onlinestats "github.com/dgryski/go-onlinestats"
	"github.com/kzahedi/goent/discrete"
	"github.com/montanaflynn/stats"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"
	"github.com/pastelnetwork/gonode/probe/wdm"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/mathext"
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

func tile(data mat.Matrix, row, column int) (*mat.Dense, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	if row != 1 && column != 1 {
		return nil, errors.New(errors.Errorf("One of the input dimensions should be equal to 1."))
	}

	if row == 1 && row == column {
		return nil, errors.New(errors.Errorf("Only one of the input dimensions should be equal to 1."))
	}

	inputRows, inputColumns := data.Dims()
	tiled := mat.NewDense(row*column, row*column, nil)
	if row != 1 {
		for i := 0; i < inputRows; i++ {
			var rawRowData []float64
			for j := 0; j < inputColumns; j++ {
				rawRowData = append(rawRowData, data.At(i, j))
			}
			for tiledRow := 0; tiledRow < row; tiledRow++ {
				tiled.SetRow(tiledRow, rawRowData)
			}
		}
	} else {
		for i := 0; i < inputColumns; i++ {
			var rawColumnData []float64
			for j := 0; j < inputRows; j++ {
				rawColumnData = append(rawColumnData, data.At(j, i))
			}
			for tiledColumn := 0; tiledColumn < column; tiledColumn++ {
				tiled.SetCol(tiledColumn, rawColumnData)
			}
		}
	}

	return tiled, nil
}

func identityMatrix(n int) *mat.Dense {
	defer pruntime.PrintExecutionTime(time.Now())
	d := make([]float64, n*n)
	for i := 0; i < n*n; i += n + 1 {
		d[i] = 1
	}
	return mat.NewDense(n, n, d)
}

func ones(r, c int) *mat.Dense {
	defer pruntime.PrintExecutionTime(time.Now())
	ones := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		ones[i] = 1
	}
	return mat.NewDense(r, c, ones)
}

func filterOutZeroes(input []float64) []float64 {
	defer pruntime.PrintExecutionTime(time.Now())
	var output []float64
	for _, value := range input {
		if value != 0 {
			output = append(output, value)
		}
	}
	return output
}

func rbfDot(pattern1, pattern2 *mat.Dense, deg float64) (*mat.Dense, error) {
	defer pruntime.PrintExecutionTime(time.Now())
	pattern1Vec := mat.NewVecDense(pattern1.RawMatrix().Rows, pattern1.RawMatrix().Data)
	var pattern1Mul mat.VecDense
	pattern1Mul.MulElemVec(pattern1Vec, pattern1Vec)

	pattern2Vec := mat.NewVecDense(pattern2.RawMatrix().Rows, pattern2.RawMatrix().Data)
	var pattern2Mul mat.VecDense
	pattern2Mul.MulElemVec(pattern2Vec, pattern2Vec)

	Q, err := tile(&pattern1Mul, 1, pattern1Mul.Len())
	if err != nil {
		return nil, errors.New(err)
	}

	R, err := tile(pattern2Mul.T(), pattern2Vec.Len(), 1)
	if err != nil {
		return nil, errors.New(err)
	}

	var QAddR, XpatternDot, H, scaledH mat.Dense
	QAddR.Add(Q, R)

	XpatternDot.Mul(pattern1Vec, pattern2Vec.T())
	XpatternDot.Scale(2, &XpatternDot)

	H.Sub(&QAddR, &XpatternDot)

	scaledH.Scale(-1.0/2.0/math.Pow(deg, 2), &H)

	rawData := scaledH.RawMatrix().Data
	for i, value := range rawData {
		rawData[i] = math.Exp(value)
	}

	return &scaledH, nil
}

func diag(i, j int, v float64) float64 {
	if i == j {
		return v
	}
	return 0
}

func addScalar1(_, _ int, v float64) float64 {
	return v + 1.0
}

func elemPow2(_, _ int, v float64) float64 {
	return math.Pow(v, 2.0)
}

// HSIC Hilbert-Schmidt Independence Criterion between arrays of input data
func HSIC(data1, data2 []float64) (float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	n := len(data1)

	// ----- width of X -----
	X := mat.NewDense(len(data1), 1, data1)

	Xmed := mat.DenseCopyOf(X)

	var XMul mat.Dense

	XMul.MulElem(X, X)

	G := mat.DenseCopyOf(&XMul)

	Q, err := tile(G, 1, n)
	if err != nil {
		return 0, errors.New(err)
	}

	R, err := tile(G.T(), n, 1)
	if err != nil {
		return 0, errors.New(err)
	}

	var QAddR, XmedDot, dists, subDists mat.Dense
	QAddR.Add(Q, R)

	XmedDot.Mul(Xmed, Xmed.T())
	XmedDot.Scale(2, &XmedDot)

	dists.Sub(&QAddR, &XmedDot)

	triDists := mat.NewTriDense(dists.RawMatrix().Rows, mat.Lower, dists.RawMatrix().Data)
	subDists.Sub(&dists, triDists)

	finalDists := mat.NewDense(int(math.Pow(float64(n), 2)), 1, subDists.RawMatrix().Data)

	finalDistsNoZeroes := filterOutZeroes(finalDists.RawMatrix().Data)

	median, err := stats.Median(finalDistsNoZeroes)
	if err != nil {
		return 0, errors.New(err)
	}

	widthX := math.Sqrt(0.5 * median)

	// ----- width of Y -----
	Y := mat.NewDense(len(data2), 1, data2)

	Ymed := mat.DenseCopyOf(Y)

	var YMul mat.Dense

	YMul.MulElem(Y, Y)

	G = mat.DenseCopyOf(&YMul)

	Q, err = tile(G, 1, n)
	if err != nil {
		return 0, errors.New(err)
	}

	R, err = tile(G.T(), n, 1)
	if err != nil {
		return 0, errors.New(err)
	}

	var YmedDot mat.Dense
	QAddR.Add(Q, R)

	YmedDot.Mul(Ymed, Ymed.T())
	YmedDot.Scale(2, &YmedDot)

	dists.Sub(&QAddR, &YmedDot)

	triDists = mat.NewTriDense(dists.RawMatrix().Rows, mat.Lower, dists.RawMatrix().Data)
	subDists.Sub(&dists, triDists)

	finalDists = mat.NewDense(int(math.Pow(float64(n), 2)), 1, subDists.RawMatrix().Data)

	finalDistsNoZeroes = filterOutZeroes(finalDists.RawMatrix().Data)

	median, err = stats.Median(finalDistsNoZeroes)
	if err != nil {
		return 0, errors.New(err)
	}

	widthY := math.Sqrt(0.5 * median)

	bone := ones(n, 1)

	identityMat := identityMatrix(n)

	nOne := ones(n, n)
	nOne.Scale(1.0/float64(n), nOne)
	var H mat.Dense
	H.Sub(identityMat, nOne)

	K, err := rbfDot(X, X, widthX)
	if err != nil {
		return 0, errors.New(err)
	}

	L, err := rbfDot(Y, Y, widthY)
	if err != nil {
		return 0, errors.New(err)
	}

	var Kc, Lc mat.Dense
	Kc.Mul(&H, K)
	Kc.Mul(&Kc, &H)

	Lc.Mul(&H, L)
	Lc.Mul(&Lc, &H)

	var KcTMulLc mat.Dense
	KcTMulLc.MulElem((&Kc).T(), &Lc)

	testStat := mat.Sum(&KcTMulLc) / float64(n)

	var KcMulLc mat.Dense
	KcMulLc.MulElem(&Kc, &Lc)

	rawData := KcMulLc.RawMatrix().Data
	for i, value := range rawData {
		rawData[i] = math.Pow(1.0/6.0*value, 2.0)
	}

	varHSIC := (mat.Sum(&KcMulLc) - mat.Trace(&KcMulLc)) / float64(n) / (float64(n) - 1.0)

	varHSIC = varHSIC * 72.0 * (float64(n) - 4.0) * (float64(n) - 5.0) / float64(n) / (float64(n) - 1.0) / (float64(n) - 2.0) / (float64(n) - 3.0)

	var diagK, diagL mat.Dense
	diagK.Apply(diag, K)
	K.Sub(K, &diagK)

	diagL.Apply(diag, L)
	L.Sub(L, &diagL)

	var dotBoneK, muX mat.Dense
	dotBoneK.Mul(bone.T(), K)
	muX.Mul(&dotBoneK, bone)
	muX.Scale(1.0/float64(n)/(float64(n)-1.0), &muX)

	var dotBoneL, muY mat.Dense
	dotBoneL.Mul(bone.T(), L)
	muY.Mul(&dotBoneL, bone)

	muY.Scale(1.0/float64(n)/(float64(n)-1.0), &muY)

	var mHSIC mat.Dense
	mHSIC.MulElem(&muX, &muY)
	mHSIC.Apply(addScalar1, &mHSIC)
	mHSIC.Sub(&mHSIC, &muX)
	mHSIC.Sub(&mHSIC, &muY)
	mHSIC.Scale(1.0/float64(n), &mHSIC)

	var al mat.Dense
	al.Apply(elemPow2, &mHSIC)
	al.Scale(1.0/varHSIC, &al)

	bet := mat.DenseCopyOf(&mHSIC)
	bet.Set(0, 0, varHSIC*float64(n)/bet.At(0, 0))

	thresh := mathext.GammaIncRegInv(al.At(0, 0), 0.95) * bet.At(0, 0)
	fmt.Printf("\ntestStat=%v", testStat)
	fmt.Printf("\nthresh=%v", thresh)
	result := 0.0
	if testStat > thresh {
		result = 1.0
	}

	return result, nil
}
