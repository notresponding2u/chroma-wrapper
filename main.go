package main

import (
	"github.com/notresponding2u/chroma-wrapper/heatmap"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"log"
)

func main() {
	w, err := wrapper.New(
		"http://localhost:54235/razer/chromasdk",
		"L",
		"notresponding2u@gmail.com",
		"Heat map",
		"heatmap",
		[]string{wrapper.DeviceKeyboard},
	)
	if err != nil {
		panic(err)
	}

	h, err := heatmap.New(w)
	if err != nil {
		log.Fatal(err)
	}

	defer h.Close()

	err = h.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
