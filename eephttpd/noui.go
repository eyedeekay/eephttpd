//go:build !gui
// +build !gui

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func UiMain() {

}

func runTray() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	exitCh := make(chan struct{})
	go func() {
		fmt.Println("eephttpd is running in the background. Press ^C to stop.")
		select {
		case <-signalCh:
			fmt.Println("received done, exiting in 500 milliseconds")
			time.Sleep(500 * time.Millisecond)
			close(exitCh)
		}
	}()

	<-exitCh
}
