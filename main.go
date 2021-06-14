package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/notresponding2u/chroma-wrapper/heatmap"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"github.com/notresponding2u/chroma-wrapper/wrapper/effect"
	hook "github.com/robotn/gohook"
	"log"
)

func main() {
	w, err := wrapper.New(
		"http://localhost:54235/razer/chromasdk",
		"L",
		"notresponding2u@gmail.com",
		"Heat map",
		"Heatmap",
		[]string{wrapper.DeviceKeyboard},
	)
	if err != nil {
		panic(err)
	}

	g := wrapper.BasicGrid()

	evChan := hook.Start()
	defer hook.End()

	go startTray(g, w, evChan)

	err = heatmap.Listen(evChan, g, w)
	if err != nil {
		log.Fatal(err)
	}
}

func startTray(g *effect.KeyboardGrid, w *wrapper.Wrapper, evChan chan hook.Event) {
	systray.Run(onReady(evChan), onExit(g, w, evChan))
}

func onReady(evChan chan hook.Event) func() {
	return func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Chroma heatmap")
		systray.SetTooltip("Chroma heatmap")

		mDiscard := systray.AddMenuItem("Discard	F9  ", "Discard current and start new")
		mSaveAndNew := systray.AddMenuItem("Save and new	F10", "Save into all time and start new heatmap")
		mMergeAndLoad := systray.AddMenuItem("Load all time	F11", "Load all time heatmap and merge to current")
		mQuit := systray.AddMenuItem("Quit	F12", "Quit the whole app")

		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mMergeAndLoad.ClickedCh:
				evChan <- hook.Event{
					Kind:    hook.KeyUp,
					Rawcode: 122,
				}
			case <-mSaveAndNew.ClickedCh:
				evChan <- hook.Event{
					Kind:    hook.KeyUp,
					Rawcode: 121,
				}
			case <-mDiscard.ClickedCh:
				evChan <- hook.Event{
					Kind:    hook.KeyUp,
					Rawcode: 120,
				}
			}
		}
	}
}

func onExit(g *effect.KeyboardGrid, w *wrapper.Wrapper, evChan chan hook.Event) func() {
	return func() {
		err := heatmap.SaveMap(g)
		if err != nil {
			log.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			log.Fatal(err)
		}

		evChan <- hook.Event{
			Kind: hook.KeyUp,
		}
	}
}
