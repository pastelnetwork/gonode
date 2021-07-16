// Package runtime provides functions to print various useful information
package runtime

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

// PrintExecutionTime retrieves name of the caller function and prints it along with elapsed time
// Typical call of PrintExecutionTime should be deffered in the beginning of the function to measure
func PrintExecutionTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	pcFunc := runtime.FuncForPC(pc)
	funcNameOnly := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := funcNameOnly.ReplaceAllString(pcFunc.Name(), "$1")
	fmt.Printf("\n%s took %s\n", funcName, elapsed)
}
