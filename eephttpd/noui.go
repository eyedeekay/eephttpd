//go:build !gui
// +build !gui

package main

import (
	"os"
	"os/signal"
)

func UiMain() {

}

func runTray() {
	//c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			os.Exit(0)
			return
		}
	}()
	<-exitCh
}
