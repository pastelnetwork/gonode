package sys

import (
	"context"
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
