package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/crypto/sha3"
)

// DiskStatus cotains info of disk storage
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// SafeErrStr returns err string
func SafeErrStr(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

// Sha3256hash returns SHA-256 hash of input message
func Sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// SafeString returns value of str ptr or empty string if ptr is nil
func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// IsContextErr checks if err is related to context
func IsContextErr(err error) bool {
	if err != nil {
		errStr := strings.ToLower(err.Error())
		return strings.Contains(errStr, "context") || strings.Contains(errStr, "ctx")
	}

	return false
}

// GetExternalIPAddress returns external IP address
func GetExternalIPAddress() (externalIP string, err error) {

	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// B64Encode base64 encodes
func B64Encode(in []byte) (out []byte) {
	out = make([]byte, base64.StdEncoding.EncodedLen(len(in)))
	base64.StdEncoding.Encode(out, in)

	return out
}

// B64Decode decode base64 input
func B64Decode(in []byte) (out []byte, err error) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(in)))
	n, err := base64.StdEncoding.Decode(out, in)
	if err != nil {
		return nil, errors.Errorf("b64 decode: %w", err)
	}

	return out[:n], nil
}

// EqualStrList tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func EqualStrList(a, b []string) error {
	if len(a) != len(b) {
		return errors.Errorf("unequal length: %d  & %d", len(a), len(b))
	}

	for i, v := range a {
		if v != b[i] {
			return errors.Errorf("index %d mismatch, vals: %s   &  %s", i, v, b[i])
		}
	}

	return nil
}

// GetHashFromString generate sha256 hash from a given string
func GetHashFromString(inputString string) string {
	h := sha3.New256()
	h.Write([]byte(inputString))
	return hex.EncodeToString(h.Sum(nil))
}

// XORBytes returns xor 2 same length bytes data
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

// BytesToInt convert input bytes data to big.Int value
func BytesToInt(inputBytes []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(inputBytes)
	return z
}

// ComputeXorDistanceBetweenTwoStrings func
func ComputeXorDistanceBetweenTwoStrings(string1 string, string2 string) uint64 {
	string1Hash := GetHashFromString(string1)
	string2Hash := GetHashFromString(string2)
	string1HashAsBytes := []byte(string1Hash)
	string2HashAsBytes := []byte(string2Hash)
	xorDistance, _ := XORBytes(string1HashAsBytes, string2HashAsBytes)
	xorDistanceAsInt := BytesToInt(xorDistance)
	xorDistanceAsString := fmt.Sprint(xorDistanceAsInt)
	if xorDistanceAsString == "0" {
		zeroAsUint64, _ := strconv.ParseUint("0", 10, 64)
		return zeroAsUint64
	}
	xorDistanceAsStringRescaled := fmt.Sprint(xorDistanceAsString[:len(xorDistanceAsString)-137])
	xorDistanceAsUint64, _ := strconv.ParseUint(xorDistanceAsStringRescaled, 10, 64)
	return xorDistanceAsUint64
}

// GetNClosestXORDistanceStringToAGivenComparisonString func
func GetNClosestXORDistanceStringToAGivenComparisonString(n int, comparisonString string, sliceOfComputingXORDistance []string, ignores ...string) []string {
	sliceOfXORDistance := make([]uint64, 0)
	XORDistanceToComputingStringMap := make(map[uint64]string)
	for _, currentComputing := range sliceOfComputingXORDistance {
		currentXORDistance := ComputeXorDistanceBetweenTwoStrings(currentComputing, comparisonString)
		for _, ignore := range ignores {
			if ignore == currentComputing {
				currentXORDistance = ^(uint64(0))
				break
			}
		}
		sliceOfXORDistance = append(sliceOfXORDistance, currentXORDistance)
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
