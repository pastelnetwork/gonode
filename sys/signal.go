package sys

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// RegisterSignalHandler registers a handler of signals from the OS.
// Before using it you have to ensure a exit from each loop when <-ctx.Done() is releasing
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
func RegisterSignalHandler(cancel context.CancelFunc, sigs ...os.Signal) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)

	go func() {
		<-sigCh
		cancel()
	}()
}

// RegisterInterruptHandler registers a handler of interrupt signals from the OS.
// When signal os.Interrupt is coming, it informs the user about it and calls func `cancelApp`.
func RegisterInterruptHandler(cancelApp context.CancelFunc, eventFn func()) {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		RegisterSignalHandler(cancel, os.Interrupt)

		<-ctx.Done()
		fmt.Print("\r")
		eventFn()

		cancelApp()
	}()
}
