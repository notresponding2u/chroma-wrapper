package wrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	SessionFile    = "Session.json"
	DeviceKeyboard = "keyboard"
)

type Wrapper struct {
	sync.Mutex
	url                string
	applicationContent string
	retryConnection    bool
	KillChannel        chan bool
	Client             *http.Client
	Session            connectionResponse
	application        app
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
	SessionId rune   `json:"sessionid"`
	Uri       string `json:"uri"`
}

type SdkResponse struct {
	Result int64  `json:"result"`
	Id     string `json:"id"`
}

func New(
	url string,
	authorName string,
	authorContact string,
	title string,
	description string,
	device []string,
) (*Wrapper, error) {
	w := &Wrapper{
		url:                url,
		applicationContent: "application/json",
		KillChannel:        make(chan bool, 1),
		application: app{
			Title:       title,
			Description: description,
			Author: author{
				Name:    authorName,
				Contact: authorContact,
			},
			DeviceSupported: device,
			Category:        "application",
		},
		retryConnection: true,
	}

	err := w.tryConnection()
	if err != nil {
		return nil, err
	}

	go w.heartbeat()

	time.Sleep(2 * time.Second)

	return w, nil
}

func (w *Wrapper) MakeKeyboardRequest(e interface{}) error {
	url := fmt.Sprintf("%s/keyboard", w.Session.Uri)

	return w.makeRequest(e, url, http.MethodPut)
}

func (w *Wrapper) Close() {
	err := w.makeRequest(nil, w.Session.Uri, http.MethodDelete)
	if err != nil {
		log.Fatal(err)
	}

	w.KillChannel <- true
}

func (w *Wrapper) makeRequest(e interface{}, url string, method string) error {
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
	if err != nil || res.StatusCode != 200 {
		err = w.tryConnection()
		if err != nil {
			log.Printf("Can't recoonect: %s", err.Error())
			return err
		}

		res, err = w.Client.Do(req)
		if err != nil {
			return err

		}
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

func (w *Wrapper) tryConnection() error {
	if w.retryConnection {
		w.Lock()
		w.retryConnection = false
		w.Unlock()

		err := w.openConnection()
		if err != nil {
			time.Sleep(5 * time.Second)
			err = w.openConnection()
			if err != nil {
				return err
			}
		}

		w.Lock()
		w.retryConnection = true
		w.Unlock()
	}

	return nil
}

func (w *Wrapper) checkIfStarted() error {
	if _, err := os.Stat(SessionFile); err == nil {
		s, err := ioutil.ReadFile(SessionFile)
		if err != nil {
			return err
		}

		err = json.Unmarshal(s, &w.Session)
		if err != nil {
			return err
		}

		w.Lock()
		w.retryConnection = false
		w.Unlock()

		err = w.pulse()
		if err == nil {
			// If heartbeat successful, then there is already running app.
			return errors.New("heatmap already running")
		}

		w.Lock()
		w.retryConnection = true
		w.Unlock()
	}
	return nil
}

func (w *Wrapper) openConnection() error {
	w.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(5) * time.Second,
	}

	err := w.checkIfStarted()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(w.application)
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

	err = json.Unmarshal(body, &w.Session)
	if err != nil {
		return err
	}

	err = saveSession(w.Session)
	if err != nil {
		return err
	}

	return nil
}

func saveSession(s connectionResponse) error {
	f, err := os.Create(SessionFile)
	if err != nil {
		return err
	}

	defer f.Close()

	j, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = f.Write(j)
	if err != nil {
		return err
	}

	return nil
}

func (w *Wrapper) pulse() error {
	url := fmt.Sprintf("%s/heartbeat", w.Session.Uri)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Cant create request: %s", err.Error()))
	}

	res, err := w.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		err = w.tryConnection()
		if err != nil {
			return errors.New(fmt.Sprintf("Missed heartbeat, can't recoonect: %s", err.Error()))
		}

		res, err = w.Client.Do(req)
		if err != nil {
			return errors.New(fmt.Sprintf("Missed heartbeat: %s", err.Error()))
		}
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", res.StatusCode))
	}

	return nil
}

func (w *Wrapper) heartbeat() {
	for {
		select {
		case <-w.KillChannel:
			return
		default:
			err := w.pulse()
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Second)
		}
	}
}
