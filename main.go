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
	//err = w.Static()

	time.Sleep(2 * time.Second)

	err = w.Custom()
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	//
	//fmt.Println("hook add...")
	//s := hook.Start()
	//defer hook.End()
	//
	//ct := false
	//for {
	//	i := <-s
	//
	//	if i.Kind == hook.KeyHold && i.Rawcode == 59 {
	//		ct = true
	//	}
	//
	//	if ct && i.Rawcode == 12 {
	//		break
	//	}
	//}

}
