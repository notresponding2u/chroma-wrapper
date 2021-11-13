package heatmap

import (
	"github.com/getlantern/systray"
	"github.com/notresponding2u/chroma-wrapper/heatmap/effect"
	"github.com/stretchr/testify/assert"
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

func TestRemap(t *testing.T) {
	// TODO: need to fix this, the systrayMock must be mocked.
	h := &heatmap{}
	h.grid = effect.BasicGrid()
	h.tray = &trayMock{}

	h.remap(key{
		X: 0,
		Y: 0,
	})

	assert.Equal(t, int64(0x0000FF), h.grid.Param[0][0])
}

//func BenchmarkRemap(b *testing.B) {
//	e := &effect.KeyboardGrid{}
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		Remap(key{
//			X: 5,
//			Y: 5,
//		}, e)
//		Remap(key{
//			X: 5,
//			Y: 5,
//		}, e)
//		Remap(key{
//			X: 2,
//			Y: 5,
//		}, e)
//	}
//}
