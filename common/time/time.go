// Package time provides functions to work with execution time
package time

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

// Measure retrieves name of the caller function and prints it along with elapsed time
// Typical call of Measure should be deffered in the beggining of the function to measure
func Measure(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	pcFunc := runtime.FuncForPC(pc)
	funcNameOnly := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := funcNameOnly.ReplaceAllString(pcFunc.Name(), "$1")
	fmt.Printf("\n%s took %s\n", funcName, elapsed)
}
