// Package time provides internal functions to alter and measure execution time
package time

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"time"
)

// Measure retrieves name of the caller function and prints it along with elapsed time
// Should be deffered in the beginning of function to measure
func Measure(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	pcFunc := runtime.FuncForPC(pc)
	funcNameOnly := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := funcNameOnly.ReplaceAllString(pcFunc.Name(), "$1")
	fmt.Printf("\n%s took %s\n", funcName, elapsed)
}

// Sleep puts execution thread to sleep for pseudo-random number of milliseconds in range from 0 to 9
func Sleep() {
	r := rand.Intn(10)
	time.Sleep(time.Duration(r) * time.Millisecond)
}
