package effect

const Custom = "CHROMA_CUSTOM"
const KeyboardMaxRows = 6
const KeyboardMaxColumns = 22

type Effect struct {
	Effect string `json:"effect"`
	Param  Param  `json:"param"`
}

type Param struct {
	Color int64 `json:"color"`
}

type KeyboardGrid struct {
	ColorMap        [1021]int64  `json:"-"`
	MapCount        [6][22]int64 `json:"-"`
	MaxKeyPresses   int64        `json:"-"`
	TotalKeyPresses int64        `json:"-"`
	Effect          string       `json:"effect"`
	Param           [6][22]int64 `json:"param"`
}

func BasicGrid() *KeyboardGrid {
	e := &KeyboardGrid{
		Effect: Custom,
		Param:  GetBaseGrid(),
	}

	setColorMap(e)

	return e
}

func GetBaseGrid() [KeyboardMaxRows][KeyboardMaxColumns]int64 {
	var g [KeyboardMaxRows][KeyboardMaxColumns]int64
	for i, _ := range g {
		for y, _ := range g[i] {
			g[i][y] = 0xFF0000
		}
	}

	return g
}

func setColorMap(e *KeyboardGrid) {
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
