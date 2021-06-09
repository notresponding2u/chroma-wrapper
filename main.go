package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/notresponding2u/chroma-wrapper/heatmap"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"github.com/notresponding2u/chroma-wrapper/wrapper/effect"
	hook "github.com/robotn/gohook"
	"log"
	"time"
)

func main() {

	w, err := wrapper.New(
		"http://localhost:54235/razer/chromasdk",
		"L",
		"notresponding2u@gmail.com",
		"Heat map new1",
		"Heatmap to be",
		[]string{wrapper.DeviceKeyboard},
	)
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)

	err = w.Custom()
	if err != nil {
		panic(err)
	}

	h := heatmap.NewMap()
	g := wrapper.BasicGrid()

	evChan := hook.Start()
	defer hook.End()

	go startTray(g, w, evChan)

	for ev := range evChan {
		if ev.Kind == hook.KeyUp {
			if k, check := h[ev.Rawcode]; check {
				if ev.Rawcode == 13 {
					heatmap.Remap(heatmap.Key{
						X: 3,
						Y: 14,
					}, g)
					heatmap.Remap(heatmap.Key{
						X: 4,
						Y: 21,
					}, g)
				} else if ev.Rawcode == 122 {
					err = heatmap.LoadFile(g, heatmap.FileAllTimeHeatMap)
					if err != nil {
						log.Fatal(err)
					}

					heatmap.Remap(k, g)

					fmt.Printf("Map merged with all times.")
				} else if ev.Rawcode == 121 {
					err = heatmap.LoadFile(g, heatmap.FileAllTimeHeatMap)
					if err != nil {
						log.Fatal(err)
					}

					err = heatmap.SaveMap(g)
					if err != nil {
						log.Fatal(err)
					}

					g = wrapper.BasicGrid()

					heatmap.Remap(k, g)

					fmt.Println("Map saved and new loaded.")
				} else if ev.Rawcode == 120 {
					g = wrapper.BasicGrid()

					heatmap.Remap(k, g)

					fmt.Println("Map discarded")
				} else {
					heatmap.Remap(k, g)
				}
				err = w.MakeKeyboardRequest(&g)
				if err != nil {
					panic(err)
				}
			}
			if ev.Rawcode == 123 {
				// Quitting
				systray.Quit()

				break
			}

			if ev.Rawcode == 0 {
				// Quitting
				break
			}
		}
	}
}

func startTray(g *effect.KeyboardGrid, w *wrapper.Wrapper, evChan chan hook.Event) {
	systray.Run(onReady, onExit(g, w, evChan))
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Chroma heatmap")
	systray.SetTooltip("Chroma heatmap")
	mQuit := systray.AddMenuItem("Quit	F12", "Quit the whole app")
	mQuit.SetIcon(icon.Data)

	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
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
