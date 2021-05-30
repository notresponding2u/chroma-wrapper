package main

import (
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"time"
)

func main() {
	w, err := wrapper.New(
		"http://localhost:54235/razer/chromasdk",
		"L",
		"notresponding2u@gmail.com",
		"Heat map new",
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
	//err = w.Static()
	err = w.Custom()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
}
