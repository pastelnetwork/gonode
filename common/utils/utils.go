package utils

import (
	"bytes"
	"container/heap"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/sync/semaphore"

	"github.com/klauspost/compress/zstd"
	"golang.org/x/crypto/sha3"
)

const (
	maxParallelHighCompressCalls = 5
	semaphoreWeight              = 1
	highCompressTimeout          = 30 * time.Minute
	highCompressionLevel         = 4
	iPCheckURL                   = "http://ipinfo.io/ip"
)

var sem = semaphore.NewWeighted(maxParallelHighCompressCalls)

// DiskStatus cotains info of disk storage
type DiskStatus struct {
	All  float64 `json:"all"`
	Used float64 `json:"used"`
	Free float64 `json:"free"`
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

// GetHashStringFromBytes generate sha256 hash string from a given byte array
func GetHashStringFromBytes(msg []byte) string {
	h := sha3.New256()
	if _, err := io.Copy(h, bytes.NewReader(msg)); err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

// SafeString returns value of str ptr or empty string if ptr is nil
func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// SafeInt returns value of int ptr or default int if ptr is nil
func SafeInt(ptr *int, def int) int {
	if ptr != nil {
		return *ptr
	}
	return def
}

// SafeFloat returns value of float ptr or default float if ptr is nil
func SafeFloat(ptr *float64, def float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return def
}

// SafeBool returns value of bool ptr or default bool if ptr is nil
func SafeBool(ptr *bool, def bool) bool {
	if ptr != nil {
		return *ptr
	}
	return def
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
	resp, err := http.Get(iPCheckURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ipString := normalizeString(string(body))
	if net.ParseIP(ipString) == nil {
		return "", errors.Errorf("invalid IP response from %s", iPCheckURL)
	}

	return ipString, nil
}

func normalizeString(s string) string {
	stringWithoutEscapes := strings.Replace(s, "\n", "", -1)
	stringWithoutEscapes = strings.Replace(stringWithoutEscapes, "\t", "", -1)
	stringWithoutEscapes = strings.Replace(stringWithoutEscapes, "\r", "", -1)

	return stringWithoutEscapes
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

// XORBytes returns the XOR of two same-length byte slices.
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

// ComputeXorDistanceBetweenTwoStrings computes the XOR distance between the hashes of two strings.
func ComputeXorDistanceBetweenTwoStrings(string1, string2 string) *big.Int {
	string1Hash := GetHashFromString(string1)
	string2Hash := GetHashFromString(string2)
	xorDistance, _ := XORBytes([]byte(string1Hash), []byte(string2Hash))
	return big.NewInt(0).SetBytes(xorDistance)
}

// GetNClosestXORDistanceStringToAGivenComparisonString finds the top N closest strings to the comparison string based on XOR distance.
func GetNClosestXORDistanceStringToAGivenComparisonString(n int, comparisonString string, sliceOfComputingXORDistance []string, ignores ...string) []string {
	h := &MaxHeap{}
	heap.Init(h)

	for _, currentComputing := range sliceOfComputingXORDistance {
		ignoreCurrent := false
		for _, ignore := range ignores {
			if ignore == currentComputing {
				ignoreCurrent = true
				break
			}
		}
		if !ignoreCurrent {
			currentXORDistance := ComputeXorDistanceBetweenTwoStrings(currentComputing, comparisonString)
			entry := DistanceEntry{distance: currentXORDistance, str: currentComputing}
			if h.Len() < n {
				heap.Push(h, entry)
			} else if entry.distance.Cmp((*h)[0].distance) < 0 {
				heap.Pop(h)
				heap.Push(h, entry)
			}
		}
	}

	result := make([]string, h.Len()) // Create the result slice with the actual size of the heap
	for i := h.Len() - 1; i >= 0; i-- {
		entry := heap.Pop(h).(DistanceEntry)
		result[i] = entry.str
	}

	return result
}

// GetFileSizeInMB returns size of the file in MB
func GetFileSizeInMB(data []byte) float64 {
	return math.Ceil(float64(len(data)) / (1024 * 1024))
}

// BytesToMB ...
func BytesToMB(bytes uint64) float64 {
	return float64(bytes) / 1048576.0
}

// StringInSlice checks if str is in slice
func StringInSlice(list []string, str string) bool {

	for _, v := range list {
		if strings.EqualFold(v, str) {
			return true
		}
	}

	return false
}

func hasActiveNetworkInterface() bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 {
			return true
		}
	}

	return false
}

func hasInternetConnectivity() bool {
	hosts := []string{"google.com:80", "microsoft.com:80", "amazon.com:80", "8.8.8.8:53", "1.1.1.1:53"}
	timeout := 1 * time.Second

	for _, host := range hosts {
		conn, err := net.DialTimeout("tcp", host, timeout)
		if err == nil {
			defer conn.Close()
			return true
		}
	}

	return false
}

// CheckInternetConnectivity checks if the device is connected to the internet
func CheckInternetConnectivity() bool {
	if hasActiveNetworkInterface() {
		return hasInternetConnectivity()
	}

	return false
}

// BytesIntToMB converts bytes to MB
func BytesIntToMB(b int) int {
	return int(math.Ceil(float64(b) / 1048576.0))
}

// Compress compresses the data
func Compress(data []byte, level int) ([]byte, error) {
	if level < 1 || level > 4 {
		return nil, fmt.Errorf("invalid compression level: %d - allowed levels are 1 - 4", level)
	}

	numCPU := runtime.NumCPU()
	// Create a buffer to store compressed data
	var compressedData bytes.Buffer

	// Create a new Zstd encoder with concurrency set to the number of CPU cores
	encoder, err := zstd.NewWriter(&compressedData, zstd.WithEncoderConcurrency(numCPU), zstd.WithEncoderLevel(zstd.EncoderLevel(level)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Zstd encoder: %v", err)
	}

	// Perform the compression
	_, err = io.Copy(encoder, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}

	// Close the encoder to flush any remaining data
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}

	return compressedData.Bytes(), nil
}

// Decompress decompresses the data
func Decompress(data []byte) ([]byte, error) {
	// Get the number of CPU cores available
	numCPU := runtime.NumCPU()

	// Create a new Zstd decoder with concurrency set to the number of CPU cores
	decoder, err := zstd.NewReader(bytes.NewReader(data), zstd.WithDecoderConcurrency(numCPU))
	if err != nil {
		return nil, fmt.Errorf("failed to create Zstd decoder: %v", err)
	}
	defer decoder.Close()

	// Perform the decompression
	decompressedData, err := io.ReadAll(decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	return decompressedData, nil
}

// RandomDuration returns a random duration between min and max
func RandomDuration(min, max int) time.Duration {
	if min > max {
		min, max = max, min
	}
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n) // read a random uint64
	randomMillisecond := min + int(n%(uint64(max-min+1)))
	return time.Duration(randomMillisecond) * time.Millisecond
}

// HighCompress compresses the data
func HighCompress(cctx context.Context, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(cctx, highCompressTimeout)
	defer cancel()

	// Acquire the semaphore. This will block if 5 other goroutines are already inside this function.
	if err := sem.Acquire(ctx, semaphoreWeight); err != nil {
		return nil, fmt.Errorf("failed to acquire semaphore: %v", err)
	}
	defer sem.Release(semaphoreWeight) // Ensure that the semaphore is always released

	numCPU := runtime.NumCPU()
	// Create a buffer to store compressed data
	var compressedData bytes.Buffer

	// Create a new Zstd encoder with concurrency set to the number of CPU cores
	encoder, err := zstd.NewWriter(&compressedData, zstd.WithEncoderConcurrency(numCPU), zstd.WithEncoderLevel(zstd.EncoderLevel(highCompressionLevel)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Zstd encoder: %v", err)
	}

	// Perform the compression
	_, err = io.Copy(encoder, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}

	// Close the encoder to flush any remaining data
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}

	return compressedData.Bytes(), nil
}
