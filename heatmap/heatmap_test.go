package heatmap

import (
	"encoding/json"
	"github.com/getlantern/systray"
	"github.com/notresponding2u/chroma-wrapper/heatmap/effect"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type trayMock struct{}

func (t *trayMock) SetTooltip(tooltip string)         {}
func (t *trayMock) Quit()                             {}
func (t *trayMock) Run(onReady func(), onExit func()) {}
func (t *trayMock) SetIcon()                          {}
func (t *trayMock) SetTitle(title string)             {}
func (t *trayMock) AddSeparator()                     {}
func (t *trayMock) AddMenuItem(title string, tooltip string) *systray.MenuItem {
	return &systray.MenuItem{}
}
func (t *trayMock) Disable(title *systray.MenuItem) {}

func GetMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			assert.Equal(t, request.Header.Get("Content-Type"), "application/json")

			res, err := json.Marshal(wrapper.SdkResponse{})
			if err != nil {
				t.Error(err)
			}

			writer.Write(res)
		}
	}))
}

func TestRemap(t *testing.T) {
	h := &heatmap{}
	h.grid = effect.BasicGrid()
	h.tray = &trayMock{}

	h.remap(key{
		X: 0,
		Y: 0,
	})

	assert.Equal(t, int64(0x0000FF), h.grid.Param[0][0])
}

func BenchmarkRemap(b *testing.B) {
	h := &heatmap{}
	h.grid = effect.BasicGrid()
	h.tray = &trayMock{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.remap(key{
			X: 0,
			Y: 0,
		})
	}
}

func TestHeatmap_Close(t *testing.T) {
	h := &heatmap{}
	h.wrapper = &wrapper.Wrapper{}
	h.wrapper.Client = &http.Client{}
	h.wrapper.KillChannel = make(chan bool, 1)
	h.grid = effect.BasicGrid()

	mock := GetMock(t)

	h.wrapper.Session.Uri = mock.URL

	go func() { h.Close("TestFile") }()

	msg := <-h.wrapper.KillChannel
	assert.Equal(t, true, msg)

	err := os.Remove("TestFile")
	if err != nil {
		t.Error(err)
	}
}
