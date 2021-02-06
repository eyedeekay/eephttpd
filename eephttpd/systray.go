// +build gui

package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"os"
)

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("I2P Web Site")
	systray.SetTooltip("Administer your local I2P site.")
	if gui {
		mEdit := systray.AddMenuItem("Edit Config", "Change the configuration")
		go func() {
			<-mEdit.ClickedCh
			UiMain()
		}()
	}
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	//	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}

func runTray() {
	go systray.Run(onReady, onExit)
}
