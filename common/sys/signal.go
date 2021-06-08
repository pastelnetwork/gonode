package sys

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// RegisterSignalInterceptor registers a signal interceptor from the OS.
// You have to ensure an exit from each loop by using ctx.Done() to break loop
// Example:
// func HelloWorld(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 		}
// 		.....
// 	}
// }
func RegisterSignalInterceptor(cancel context.CancelFunc, sigs ...os.Signal) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)

	go func() {
		<-sigCh
		cancel()
	}()
}

// RegisterInterruptHandler registers a handler of interrupt signal from the OS.
// When signal os.Interrupt is coming, it informs the user about it by calling `notifyFn`.
func RegisterInterruptHandler(notifyFn func()) {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		RegisterSignalInterceptor(cancel, os.Interrupt)

		<-ctx.Done()
		fmt.Print("\r")

		notifyFn()
	}()
}
