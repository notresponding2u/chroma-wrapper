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
	session            connectionResponse
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
}

type connectionResponse struct {
	SessionId int64  `json:"sessionid"`
	Uri       string `json:"uri"`
}

const DeviceKeyboard = "keyboard"

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
	go w.Heartbeat()
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
	return nil
}
func (w *wrapper) Heartbeat() {
	url := fmt.Sprintf("%s/heartbeat", w.session.Uri)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(errors.New(fmt.Sprintf("Status code %d", res.StatusCode)))
	}
	time.Sleep(5 * time.Second)
}

func (w *wrapper) Static() error {
	e := effect.Effect{
		Effect: effect.Custom,
		Param:  effect.Param{Color: 255},
	}
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	res, err := http.Post(w.url, w.applicationContent, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error effect change, code: %d", res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response effect.EffectResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Result != 0 {
		return errors.New(fmt.Sprintf("Wront response result code: %d", response.Result))
	}
	return nil
}
