package main

import (
	"github.com/notresponding2u/chroma-wrapper/wrapper"
)

func main() {
	w, err := wrapper.New(
		"http://localhost:54235/razer/chromasdk",
		"L",
		"notresponding2u@gmail.com",
		"Heat map",
		"Heatmap to be",
		[]string{wrapper.DeviceKeyboard},
	)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	err = w.Static()
	if err != nil {
		panic(err)
	}
	//time.Sleep(10 * time.Second)
}
