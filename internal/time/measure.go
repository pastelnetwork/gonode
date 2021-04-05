// Package time provides internal functions to measure execution time
package time

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"time"
)

// Measure retrieves name of the caller function and prints it along with elapsed time
// Should be deffered in the beggining of function to measure
func Measure(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	pcFunc := runtime.FuncForPC(pc)
	funcNameOnly := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := funcNameOnly.ReplaceAllString(pcFunc.Name(), "$1")
	log.Println(fmt.Sprintf("\n%s took %s", funcName, elapsed))
}
