package tray

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

type TrayApi interface {
	SetTooltip(string)
	Quit()
	Run(func(), func())
	SetIcon()
	SetTitle(string)
	AddSeparator()
	AddMenuItem(string, string) *systray.MenuItem
	Disable(*systray.MenuItem)
}

type Tray struct{}

func (t *Tray) SetTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}
func (t *Tray) Quit() {
	systray.Quit()
}
func (t *Tray) Run(onReady func(), onExit func()) {
	systray.Run(onReady, onExit)
}
func (t *Tray) SetIcon() {
	systray.SetIcon(icon.Data)
}
func (t *Tray) SetTitle(title string) {
	systray.SetTitle(title)
}
func (t *Tray) AddSeparator() {
	systray.AddSeparator()
}
func (t *Tray) AddMenuItem(title string, tooltip string) *systray.MenuItem {
	return systray.AddMenuItem(title, tooltip)
}
func (t *Tray) Disable(title *systray.MenuItem) {
	title.Disable()
}
