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
	List               effect.List
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

type SdkResponse struct {
	Result int64  `json:"result"`
	Id     string `json:"id"`
}

type SdkResults struct {
	Results []SdkResponse
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
	return nil
}

func (w *wrapper) heartbeat() {
	for {
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
		fmt.Println("Beep")
		time.Sleep(5 * time.Second)
	}
}

func (w *wrapper) Close() error {
	req, err := http.NewRequest(http.MethodDelete, w.session.Uri, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Close status code: %d", res.StatusCode))
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
		return errors.New(fmt.Sprintf("Error closing connection, response result: %d", response.Result))
	}
	return nil
}

func (w *wrapper) Static() error {
	e := effect.Effect{
		Effect: effect.Static,
		Param:  effect.Param{Color: 200},
	}
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/keyboard", w.session.Uri)
	res, err := http.Post(url, w.applicationContent, bytes.NewBuffer(payload))
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
	var response SdkResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Result != 0 {
		return errors.New(fmt.Sprintf("Status code: %d", response.Result))
	}
	err = w.setEffect(response)
	return err
}

func (w *wrapper) setEffect(ef SdkResponse) error {
	w.List.Ids = append(w.List.Ids, ef.Id)
	fmt.Println(ef.Id)
	e := effect.Identifier{Id: ef.Id}
	url := fmt.Sprintf("%s/effect", w.session.Uri)
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
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
		return errors.New(fmt.Sprintf("Bad result code on effect apply: %d", response.Result))
	}
	return nil
}

func (w *wrapper) DeleteEffects() error {
	url := fmt.Sprintf("%s/effect", w.session.Uri)
	payload, err := json.Marshal(w.List)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var results SdkResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return err
	}
	for _, sdkRes := range results.Results {
		if sdkRes.Result != 0 {
			return errors.New(fmt.Sprintf("Wrong result code on deleting effect: %d", sdkRes.Result))
		}
	}
	return nil
}
