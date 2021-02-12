// +build !gui

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func UiMain() {

}

func runTray() {
	c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})
	go func(ctx context.Context) {
		for {
			fmt.Println("eephttpd is running in the background. Press ^C to stop.")
			// Do something useful in a real usecase.
			// Here we just sleep for this example.
			time.Sleep(time.Minute * 10)

			select {
			case <-ctx.Done():
				fmt.Println("received done, exiting in 500 milliseconds")
				time.Sleep(500 * time.Millisecond)
				exitCh <- struct{}{}
				return
			default:
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
}
