// +build gui

package main

import (
	"github.com/andlabs/ui"
	"github.com/atotto/clipboard"
	"runtime"
)

var mainwin *ui.Window

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	copyButton := ui.NewButton("Copy I2P Address to Clipboard.")

	copyButton.OnClicked(func(*ui.Button) {
		err := clipboard.WriteAll("http://" + eepsite.Base32())
		if err != nil {
			ui.MsgBox(mainwin,
				"Error copying text to clipboard",
				err.Error())
		} else {
			ui.MsgBox(mainwin,
				"Copied to clipboard",
				"I2P Base32 address copied to clipboard: "+eepsite.Base32())
		}
	})

	hbox.Append(copyButton, false)

	vbox.Append(ui.NewLabel("This panel is for configuring eephttpd."), false)
	vbox.Append(ui.NewLabel("eephttpd is a simple file server for the I2P network, but with superpowers."), false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("Configuration")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	enterDirectory := ui.NewEntry()
	enterDirectory.SetText(*directory)

	enterGitURL := ui.NewEntry()
	enterGitURL.SetText(*giturl)

	entryForm.Append("Serve files from this directory", enterDirectory, false)
	entryForm.Append("Clone site from a git repository", enterGitURL, false)

	saveSettingsButton := ui.NewButton("Save and apply settings")

	saveSettingsButton.OnClicked(func(*ui.Button) {
		eepsite.ServeDir = enterDirectory.Text()
		eepsite.GitURL = enterGitURL.Text()

		if err := eepsite.Save(); err != nil {
			ui.MsgBox(mainwin,
				"Error saving config file",
				err.Error())
		}

		if err := eepsite.ResetGit(); err != nil {
			ui.MsgBox(mainwin,
				"Error cloning git repository",
				err.Error())
		}
	})
	entryForm.Append("Confirm", saveSettingsButton, false)

	return vbox
}

func UiMain() {
	gui = true
	ui.Main(IniEdit)
}

func IniEdit() {
	mainwin = ui.NewWindow("I2P Site Control Panel", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		/*if runtime.GOOS == "darwin" {
			return false
		}*/
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		/*if runtime.GOOS == "darwin" {
			return false
		}*/
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Basic Controls", makeBasicControlsPage())
	tab.SetMargined(0, true)
	mainwin.Show()
}
