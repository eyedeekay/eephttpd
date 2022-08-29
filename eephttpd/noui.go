//go:build !gui
// +build !gui

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func UiMain() {

}

func runTray() {
	exitCh := make(chan struct{})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		fmt.Println("eephttpd is running in the background. Press ^C to stop.")
		select {
		case <-signalCh:
			fmt.Println("received done, exiting...")
			close(exitCh)
		}
	}()

	<-exitCh
}
