package utils

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
)

func XORBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}

func BytesToInt(input_bytes []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(input_bytes)
	return z
}

func ComputeXorDistanceBetweenTwoStrings(string1 string, string2 string) uint64 {
	string1Hash := GetHashFromString(string1)
	string2Hash := GetHashFromString(string2)
	string1HashAsBytes := []byte(string1Hash)
	string2HashAsBytes := []byte(string2Hash)
	xorDistance, _ := XORBytes(string1HashAsBytes, string2HashAsBytes)
	xorDistanceAsInt := BytesToInt(xorDistance)
	xorDistanceAsString := fmt.Sprint(xorDistanceAsInt)
	xorDistanceAsStringRescaled := fmt.Sprint(xorDistanceAsString[:len(xorDistanceAsString)-137])
	xorDistanceAsUint64, _ := strconv.ParseUint(xorDistanceAsStringRescaled, 10, 64)
	return xorDistanceAsUint64
}

func GetNClosestXORDistanceStringToAGivenComparisonString(n int, comparisonString string, sliceOfComputingXORDistance []string) []string {
	sliceOfXORDistance := make([]uint64, len(sliceOfComputingXORDistance))
	XORDistanceToComputingStringMap := make(map[uint64]string)
	for idx, currentComputing := range sliceOfComputingXORDistance {
		currentXORDistance := ComputeXorDistanceBetweenTwoStrings(currentComputing, comparisonString)
		sliceOfXORDistance[idx] = currentXORDistance
		XORDistanceToComputingStringMap[currentXORDistance] = currentComputing
	}
	sort.Slice(sliceOfXORDistance, func(i, j int) bool { return sliceOfXORDistance[i] < sliceOfXORDistance[j] })
	sliceOfTopNClosestString := make([]string, n)
	for ii, currentXORDistance := range sliceOfXORDistance {
		if ii < n {
			sliceOfTopNClosestString[ii] = XORDistanceToComputingStringMap[currentXORDistance]
		}
	}
	return sliceOfTopNClosestString
}
