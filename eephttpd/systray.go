// +build gui

package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"os"
	"runtime"
	"time"
)

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("I2P Web Site")
	systray.SetTooltip("Administer your local I2P site.")
	var mEdit *systray.MenuItem
	if gui {
		mEdit = systray.AddMenuItem("Edit Config", "Change the configuration")

	}
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		go func() {
			<-mEdit.ClickedCh
			UiMain()
		}()

		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
			os.Exit(0)
		}()
		time.Sleep(time.Second)
	}
	// Sets the icon of a menu item. Only available on Mac and Windows.
	//	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}

func runTray() {
	if runtime.GOOS == "darwin" {
		gui = true
	}
	systray.Run(onReady, onExit)
}
