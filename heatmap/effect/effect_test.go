package effect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicGridColorMap(t *testing.T) {
	g := BasicGrid()

	var color int64 = 0xFF0000
	for i := 0; i < 255; i++ {
		color += 0x000100
		assert.Equal(t, g.ColorMap[i], color)
	}
	color = 0xFFFF00
	for i := 255; i < 510; i++ {
		color -= 0x010000
		assert.Equal(t, g.ColorMap[i], color)
	}
	color = 0x00FF00
	for i := 510; i < 765; i++ {
		color += 0x000001
		assert.Equal(t, g.ColorMap[i], color)
	}
	color = 0x00FFFF
	for i := 765; i < 1021; i++ {
		color -= 0x000100
		assert.Equal(t, g.ColorMap[i], color)
	}
}

func TestGetBaseGrid(t *testing.T) {
	g := GetBaseGrid()

	for _, row := range g {
		for _, column := range row {
			assert.Equal(t, column, int64(0xFF0000))
		}
	}
}
