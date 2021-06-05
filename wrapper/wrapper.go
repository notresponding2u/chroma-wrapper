package wrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/notresponding2u/chroma-wrapper/wrapper/effect"
	"io/ioutil"
	"net/http"
	"time"
)

type wrapper struct {
	url                string
	applicationContent string
	dead               chan bool
	session            connectionResponse
	List               effect.List
	Client             *http.Client
}

type author struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type app struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Author          author   `json:"author"`
	DeviceSupported []string `json:"device_supported"`
	Category        string   `json:"category"`
	alive           chan bool
}

type connectionResponse struct {
	SessionId rune   `json:"sessionid"`
	Uri       string `json:"uri"`
}

type SdkResponse struct {
	Result int64  `json:"result"`
	Id     string `json:"id"`
}

type SdkResults struct {
	Results []SdkResponse
}

const (
	DeviceKeyboard     = "keyboard"
	KeyboardMaxRows    = 6
	KeyboardMaxColumns = 22
)

//const DeviceMouse = "mouse"
//const DeviceHeadset = "headset"
//const DeviceMousepad = "mousepad"
//const DeviceKeypad = "keypad"
//const DeviceChromalink = "chromalink"

// New
// device must be one of the constants of the package.
// Only DeviceKeyboard supported right now.
func New(
	url string,
	authorName string,
	authorContact string,
	title string,
	description string,
	device []string,
) (*wrapper, error) {
	w := &wrapper{
		url:                url,
		applicationContent: "application/json",
		dead:               make(chan bool),
	}
	a := app{
		Title:       title,
		Description: description,
		Author: author{
			Name:    authorName,
			Contact: authorContact,
		},
		DeviceSupported: device,
		Category:        "application",
	}
	err := w.openConnection(a)
	if err != nil {
		return nil, err
	}
	go w.heartbeat()
	return w, nil
}

func (w *wrapper) openConnection(a app) error {
	payload, err := json.Marshal(a)
	if err != nil {
		return err
	}

	res, err := http.Post(w.url, w.applicationContent, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &w.session)
	if err != nil {
		return err
	}

	fmt.Printf("Session %q", w.session)

	w.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(5) * time.Second,
	}

	return nil
}

func (w *wrapper) heartbeat() {
	for {
		select {
		case <-w.dead:
			return
		default:
			url := fmt.Sprintf("%s/heartbeat", w.session.Uri)

			req, err := http.NewRequest(http.MethodPut, url, nil)
			if err != nil {
				panic(err)
			}

			res, err := w.Client.Do(req)
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()
			if res.StatusCode != 200 {
				panic(errors.New(fmt.Sprintf("Status code %d", res.StatusCode)))
			}

			time.Sleep(time.Second)
		}
	}
}

func (w *wrapper) deleteEffects() error {
	url := fmt.Sprintf("%s/effect", w.session.Uri)
	return w.makeRequest(w.List, url, http.MethodDelete)
}

func (w *wrapper) Close() error {
	w.dead <- false
	time.Sleep(2 * time.Second)
	return w.makeRequest(nil, w.session.Uri, http.MethodDelete)
	//err = w.deleteEffects()
}

func (w *wrapper) Static() error {
	e := &effect.Effect{
		Effect: effect.Static,
		Param:  effect.Param{Color: 200},
	}
	return w.MakeKeyboardRequest(e)
}

func (w *wrapper) MakeKeyboardRequest(e interface{}) error {
	url := fmt.Sprintf("%s/keyboard", w.session.Uri)
	return w.makeRequest(e, url, http.MethodPut)
}

func (w *wrapper) makeRequest(e interface{}, url string, method string) error {
	payload, err := json.Marshal(&e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := w.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error, httpcode: %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response SdkResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Result != 0 {
		return errors.New(fmt.Sprintf("Status code: %d", response.Result))
	}

	return err

}

func (w *wrapper) setEffect(ef SdkResponse) error {
	e := &effect.Identifier{Id: ef.Id}
	url := fmt.Sprintf("%s/effect", w.session.Uri)

	return w.makeRequest(e, url, http.MethodPost)
}

func GetKeyboardStruct() [KeyboardMaxRows][KeyboardMaxColumns]int64 {
	var grid effect.KeyboardGrid
	for i, _ := range grid.Param {
		for y, _ := range grid.Param[i] {
			grid.Param[i][y] = 0xFF0000
		}
	}
	return grid.Param
}

func (w *wrapper) Custom() error {
	e := &effect.KeyboardGrid{
		Effect: effect.Custom,
		Param:  GetKeyboardStruct(),
	}
	return w.MakeKeyboardRequest(e)
}

func BasicGrid() *effect.KeyboardGrid {
	e := &effect.KeyboardGrid{
		Effect: effect.Custom,
		Param:  GetKeyboardStruct(),
	}

	setColorMap(e)

	return e
}

func setColorMap(e *effect.KeyboardGrid) {
	var color int64 = 0xFF0000
	for i := 0; i < 255; i++ {
		color += 0x000100
		e.ColorMap[i] = color
	}
	color = 0xFFFF00
	for i := 255; i < 510; i++ {
		color -= 0x010000
		e.ColorMap[i] = color
	}
	color = 0x00FF00
	for i := 510; i < 765; i++ {
		color += 0x000001
		e.ColorMap[i] = color
	}
	color = 0x00FFFF
	for i := 765; i < 1021; i++ {
		color -= 0x000100
		e.ColorMap[i] = color
	}
}
