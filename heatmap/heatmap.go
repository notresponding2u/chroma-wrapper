package heatmap

import (
	"encoding/json"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"github.com/notresponding2u/chroma-wrapper/wrapper/effect"
	hook "github.com/robotn/gohook"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const FileAllTimeHeatMap = "./AllTime.json"

type key struct {
	X int64
	Y int64
}

/**
LOGIC:
FF0000
FFFF00
00FF00
00FFFF
0000FF
*/

type callback func(k key) error

type heatmap struct {
	sync.Mutex
	isKeyHeld   bool
	timeoutChan chan bool
	grid        *effect.KeyboardGrid
	evChan      chan hook.Event
	wrapper     *wrapper.Wrapper
	sigc        chan os.Signal
}

func New(w *wrapper.Wrapper) (*heatmap, error) {
	h := &heatmap{}
	h.grid = effect.BasicGrid()
	err := w.MakeKeyboardRequest(h.grid)
	if err != nil {
		return nil, err
	}
	h.evChan = hook.Start()
	h.wrapper = w
	go h.startTray()
	return h, nil
}

func (h *heatmap) Close() {
	err := h.saveMap()
	if err != nil {
		log.Fatal(err)
	}
	defer hook.End()
}

func (h *heatmap) remap(k key) {
	h.grid.MapCount[k.X][k.Y]++
	if h.grid.MaxKeyPresses < h.grid.MapCount[k.X][k.Y] {
		h.grid.MaxKeyPresses = h.grid.MapCount[k.X][k.Y]
		systray.SetTooltip(fmt.Sprintf("Max count %d", h.grid.MaxKeyPresses))
	}
	for x, _ := range h.grid.Param {
		for y, _ := range h.grid.Param[x] {
			switch h.grid.MapCount[x][y] {
			case 0:
				h.grid.Param[x][y] = 0xFF0000
			case h.grid.MaxKeyPresses:
				h.grid.Param[x][y] = 0x0000FF
			default:
				percentage := float64(h.grid.MapCount[x][y]) / float64(h.grid.MaxKeyPresses) * float64(len(h.grid.ColorMap))
				if int64(percentage) < int64(len(h.grid.ColorMap)) {
					h.grid.Param[x][y] = h.grid.ColorMap[int64(percentage)]
				} else {
					h.grid.Param[x][y] = 0x0000FF
				}
			}
		}
	}
}

func (h *heatmap) saveMap() error {
	if _, err := os.Stat(FileAllTimeHeatMap); os.IsNotExist(err) {
		return h.save(FileAllTimeHeatMap)
	} else {
		err = h.loadFile(FileAllTimeHeatMap)
		if err != nil {
			return err
		}

		return h.save(FileAllTimeHeatMap)
	}
}

func (h *heatmap) save(file string) error {
	j, err := json.Marshal(h.grid.MapCount)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, j, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (h *heatmap) mergeHeatmaps(donor *effect.KeyboardGrid) {
	for x, _ := range h.grid.MapCount {
		for y, _ := range h.grid.MapCount[x] {
			h.grid.MapCount[x][y] += donor.MapCount[x][y]
			if h.grid.MapCount[x][y] > h.grid.MaxKeyPresses {
				h.grid.MaxKeyPresses = h.grid.MapCount[x][y]
			}
		}
	}
}

func (h *heatmap) Listen() error {
	m := newMap()

	h.sigc = make(chan os.Signal, 1)
	signal.Notify(h.sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	h.timeoutChan = make(chan bool)

	var lastEvent uint16
	var lastEventKind uint8

	defer systray.Quit()
	for {
		select {
		case ev := <-h.evChan:
			if k, check := m[ev.Rawcode]; check {
				if (ev.Kind == hook.KeyHold || ev.Kind == hook.KeyUp) && (lastEvent != ev.Rawcode || lastEventKind != ev.Kind) {
					lastEvent = ev.Rawcode
					lastEventKind = ev.Kind

					if ev.Kind == hook.KeyUp {
						if h.isKeyHeld {
							go func() {
								h.timeoutChan <- true
							}()
						}

					}

					if ev.Kind == hook.KeyHold {
						switch ev.Rawcode {
						case 122:
							// Load all times
							if !h.isKeyHeld {
								h.isKeyHeld = true
								go func() {
									err := h.processCallback(func(k key) error {
										err := h.loadFile(FileAllTimeHeatMap)
										if err != nil {
											return err
										}

										h.remap(k)

										return nil
									}, k, true)
									if err != nil {
										log.Fatal(err)
									}
								}()
							}
						case 121:
							// Save and new
							if !h.isKeyHeld {
								h.isKeyHeld = true
								go func() {
									err := h.processCallback(func(k key) error {
										err := h.saveMap()
										if err != nil {
											return err
										}

										h.grid = effect.BasicGrid()
										return nil
									}, k, true)
									if err != nil {
										log.Fatal(err)
									}
								}()
							}
						case 120:
							// Discard
							if !h.isKeyHeld {
								h.isKeyHeld = true
								go func() {
									err := h.processCallback(func(k key) error {
										h.grid = effect.BasicGrid()
										h.remap(k)
										return nil
									}, k, true)
									if err != nil {
										log.Fatal(err)
									}
								}()
							}
						case 123:
							// Quitting
							if !h.isKeyHeld {
								h.isKeyHeld = true
								go func() {
									err := h.processCallback(func(k key) error {
										h.sigc <- syscall.SIGQUIT
										return nil
									}, k, false)
									if err != nil {
										log.Fatal(err)
									}
								}()
							}
						case 13:
							h.remap(key{
								X: 3,
								Y: 14,
							})
							h.remap(key{
								X: 4,
								Y: 21,
							})

							err := h.wrapper.MakeKeyboardRequest(&h.grid)
							if err != nil {
								return err
							}
						default:
							h.remap(k)

							err := h.wrapper.MakeKeyboardRequest(&h.grid)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		case <-h.sigc:
			return nil
		}
	}
}

func (h *heatmap) processCallback(f callback, k key, shouldRefresh bool) error {
	var err error
	select {
	case <-h.timeoutChan:
		break
	case <-time.After(time.Second):
		h.Lock()
		err = f(k)
		h.Unlock()
		if shouldRefresh {
			err := h.wrapper.MakeKeyboardRequest(h.grid)
			if err != nil {
				return err
			}
		}

		break
	}
	h.Lock()
	h.isKeyHeld = false
	h.Unlock()

	return err
}

func (h *heatmap) loadFile(file string) error {
	if _, err := os.Stat(FileAllTimeHeatMap); os.IsNotExist(err) {
		return nil
	} else {
		g := &effect.KeyboardGrid{}

		j, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		err = json.Unmarshal(j, &g.MapCount)
		if err != nil {
			return err
		}

		h.mergeHeatmaps(g)

		return os.Remove(file)
	}
}

//func (h *heatmap) HeatUp(k key) {
//	if h.grid.Param[k.X][k.Y] > 0x0000FF {
//		h.grid.Param[k.X][k.Y] = h.grid.ColorMap[h.grid.MapCount[k.X][k.Y]]
//		h.grid.MapCount[k.X][k.Y]++

// So sad that I don't want this q.Q
//switch {
//case grid.Param[k.X][k.Y]&0xFF0000 == 0xFF0000 && grid.Param[k.X][k.Y] != 0xFFFF00: //	From blue to blue/green
//	fmt.Println("more green")
//	grid.Param[k.X][k.Y] += 0x000100
//case (grid.Param[k.X][k.Y]&0x00FF00 == 0x00FF00 || grid.Param[k.X][k.Y] == 0xFFFF00) && grid.Param[k.X][k.Y] > 0x00FFFF: // From blue/green to green
//	fmt.Println("less blue")
//	grid.Param[k.X][k.Y] -= 0x010000
//case (grid.Param[k.X][k.Y] < 0x00FFFF || grid.Param[k.X][k.Y] == 0x00FF00) && grid.Param[k.X][k.Y]&0x0000FF != 0x0000FF: //	From green to green/red
//	fmt.Println("more red")
//	grid.Param[k.X][k.Y] += 0x000001
//case grid.Param[k.X][k.Y] <= 0x0FFFF && grid.Param[k.X][k.Y] > 0x0000FF: //	From green/red to red
//	fmt.Println("less green")
//	grid.Param[k.X][k.Y] -= 0x000100
//}
//	}
//}

func (h *heatmap) startTray() {
	systray.Run(h.onReady(), h.onExit())
}

func (h *heatmap) onReady() func() {
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
				h.sigc <- syscall.SIGQUIT
			case <-mMergeAndLoad.ClickedCh:
				h.evChan <- hook.Event{
					Kind:    hook.KeyHold,
					Rawcode: 122,
				}
			case <-mSaveAndNew.ClickedCh:
				h.evChan <- hook.Event{
					Kind:    hook.KeyHold,
					Rawcode: 121,
				}
			case <-mDiscard.ClickedCh:
				h.evChan <- hook.Event{
					Kind:    hook.KeyHold,
					Rawcode: 120,
				}
			}
		}
	}
}

func (h *heatmap) onExit() func() {
	return func() {
		err := h.saveMap()
		if err != nil {
			log.Fatal(err)
		}

		h.sigc <- syscall.SIGQUIT
	}
}

func newMap() map[uint16]key {
	m := make(map[uint16]key)

	m[27] = key{
		X: 0,
		Y: 1,
	}
	m[192] = key{
		X: 1,
		Y: 1,
	}
	m[9] = key{
		X: 2,
		Y: 1,
	}
	m[20] = key{
		X: 3,
		Y: 1,
	}
	m[160] = key{
		X: 4,
		Y: 1,
	}
	m[162] = key{
		X: 5,
		Y: 1,
	}

	m[49] = key{
		X: 1,
		Y: 2,
	}
	m[81] = key{

		X: 2,
		Y: 2,
	}
	m[65] = key{
		X: 3,
		Y: 2,
	}
	m[226] = key{
		X: 4,
		Y: 2,
	}
	m[91] = key{
		X: 5,
		Y: 2,
	}

	m[112] = key{
		X: 0,
		Y: 3,
	}
	m[50] = key{
		X: 1,
		Y: 3,
	}
	m[87] = key{
		X: 2,
		Y: 3,
	}
	m[83] = key{
		X: 3,
		Y: 3,
	}
	m[90] = key{
		X: 4,
		Y: 3,
	}
	m[164] = key{
		X: 5,
		Y: 3,
	}

	m[113] = key{
		X: 0,
		Y: 4,
	}
	m[51] = key{
		X: 1,
		Y: 4,
	}
	m[69] = key{
		X: 2,
		Y: 4,
	}
	m[68] = key{
		X: 3,
		Y: 4,
	}
	m[88] = key{
		X: 4,
		Y: 4,
	}

	m[114] = key{
		X: 0,
		Y: 5,
	}
	m[52] = key{
		X: 1,
		Y: 5,
	}
	m[82] = key{
		X: 2,
		Y: 5,
	}
	m[70] = key{
		X: 3,
		Y: 5,
	}
	m[67] = key{
		X: 4,
		Y: 5,
	}

	m[115] = key{
		X: 0,
		Y: 6,
	}
	m[53] = key{
		X: 1,
		Y: 6,
	}
	m[84] = key{
		X: 2,
		Y: 6,
	}
	m[71] = key{
		X: 3,
		Y: 6,
	}
	m[86] = key{
		X: 4,
		Y: 6,
	}

	m[116] = key{
		X: 0,
		Y: 7,
	}
	m[54] = key{
		X: 1,
		Y: 7,
	}
	m[89] = key{
		X: 2,
		Y: 7,
	}
	m[72] = key{
		X: 3,
		Y: 7,
	}
	m[66] = key{
		X: 4,
		Y: 7,
	}
	m[32] = key{
		X: 5,
		Y: 7,
	}

	m[117] = key{
		X: 0,
		Y: 8,
	}
	m[55] = key{
		X: 1,
		Y: 8,
	}
	m[85] = key{
		X: 2,
		Y: 8,
	}
	m[74] = key{
		X: 3,
		Y: 8,
	}
	m[78] = key{
		X: 4,
		Y: 8,
	}

	m[118] = key{
		X: 0,
		Y: 9,
	}
	m[56] = key{
		X: 1,
		Y: 9,
	}
	m[73] = key{
		X: 2,
		Y: 9,
	}
	m[75] = key{
		X: 3,
		Y: 9,
	}
	m[77] = key{
		X: 4,
		Y: 9,
	}

	m[119] = key{
		X: 0,
		Y: 10,
	}
	m[57] = key{
		X: 1,
		Y: 10,
	}
	m[79] = key{
		X: 2,
		Y: 10,
	}
	m[76] = key{
		X: 3,
		Y: 10,
	}
	m[188] = key{
		X: 4,
		Y: 10,
	}

	m[120] = key{
		X: 0,
		Y: 11,
	}
	m[48] = key{
		X: 1,
		Y: 11,
	}
	m[80] = key{
		X: 2,
		Y: 11,
	}
	m[186] = key{
		X: 3,
		Y: 11,
	}
	m[190] = key{
		X: 4,
		Y: 11,
	}
	m[165] = key{
		X: 5,
		Y: 11,
	}

	m[121] = key{
		X: 0,
		Y: 12,
	}
	m[189] = key{
		X: 1,
		Y: 12,
	}
	m[219] = key{
		X: 2,
		Y: 12,
	}
	m[222] = key{
		X: 3,
		Y: 12,
	}
	m[191] = key{
		X: 4,
		Y: 12,
	}

	m[122] = key{
		X: 0,
		Y: 13,
	}
	m[187] = key{
		X: 1,
		Y: 13,
	}
	m[221] = key{
		X: 2,
		Y: 13,
	}
	m[220] = key{
		X: 3,
		Y: 13,
	}
	m[93] = key{
		X: 5,
		Y: 13,
	}

	m[123] = key{
		X: 0,
		Y: 14,
	}
	m[8] = key{
		X: 1,
		Y: 14,
	}
	m[13] = key{
		X: 3,
		Y: 14,
	}
	m[161] = key{
		X: 4,
		Y: 14,
	}
	m[163] = key{
		X: 5,
		Y: 14,
	}

	m[44] = key{
		X: 0,
		Y: 15,
	}
	m[45] = key{
		X: 1,
		Y: 15,
	}
	m[46] = key{
		X: 2,
		Y: 15,
	}
	m[37] = key{
		X: 5,
		Y: 15,
	}

	m[145] = key{
		X: 0,
		Y: 16,
	}
	m[36] = key{
		X: 1,
		Y: 16,
	}
	m[35] = key{
		X: 2,
		Y: 16,
	}
	m[38] = key{
		X: 4,
		Y: 16,
	}
	m[40] = key{
		X: 5,
		Y: 16,
	}

	m[19] = key{
		X: 0,
		Y: 17,
	}
	m[33] = key{
		X: 1,
		Y: 17,
	}
	m[34] = key{
		X: 2,
		Y: 17,
	}
	m[39] = key{
		X: 5,
		Y: 17,
	}

	m[144] = key{
		X: 1,
		Y: 18,
	}
	m[103] = key{
		X: 2,
		Y: 18,
	}
	m[100] = key{
		X: 3,
		Y: 18,
	}
	m[97] = key{
		X: 4,
		Y: 18,
	}

	m[111] = key{
		X: 1,
		Y: 19,
	}
	m[104] = key{
		X: 2,
		Y: 19,
	}
	m[101] = key{
		X: 3,
		Y: 19,
	}
	m[98] = key{
		X: 4,
		Y: 19,
	}
	m[96] = key{
		X: 5,
		Y: 19,
	}

	m[106] = key{
		X: 1,
		Y: 20,
	}
	m[105] = key{
		X: 2,
		Y: 20,
	}
	m[102] = key{
		X: 3,
		Y: 20,
	}
	m[99] = key{
		X: 4,
		Y: 20,
	}
	m[110] = key{
		X: 5,
		Y: 20,
	}

	m[109] = key{
		X: 1,
		Y: 21,
	}
	m[107] = key{
		X: 2,
		Y: 21,
	}
	m[13] = key{
		X: 4,
		Y: 21,
	}

	return m
}
