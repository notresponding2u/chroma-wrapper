package main

import (
	"fmt"
	"github.com/notresponding2u/chroma-wrapper/heatmap"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	hook "github.com/robotn/gohook"
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
	defer func() {
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(2 * time.Second)

	err = w.Custom()
	if err != nil {
		panic(err)
	}

	h := heatmap.NewMap()
	g := wrapper.BasicGrid()

	fmt.Println("hook start...")
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		//fmt.Println("hook: ", ev)
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
				} else {
					heatmap.Remap(k, g)
				}
				err = w.MakeKeyboardRequest(&g)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
