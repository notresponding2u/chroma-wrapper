package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/notresponding2u/chroma-wrapper/heatmap"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
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

	go func() {
		systray.Run(onReady, onExit)
		evChan <- hook.Event{
			Kind:      hook.KeyUp,
			When:      time.Time{},
			Mask:      0,
			Reserved:  0,
			Keycode:   0,
			Rawcode:   123,
			Keychar:   0,
			Button:    0,
			Clicks:    0,
			X:         0,
			Y:         0,
			Amount:    0,
			Rotation:  0,
			Direction: 0,
		}
	}()

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
				err = heatmap.SaveMap(g)
				if err != nil {
					log.Fatal(err)
				}

				err = w.Close()
				if err != nil {
					log.Fatal(err)
				}

				break
			}
		}
	}
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit	F12", "Quit the whole app")
	mQuit.SetIcon(icon.Data)

	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {
}
