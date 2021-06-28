package effect

const Static = "CHROMA_STATIC"
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

type Identifier struct {
	Id string `json:"id"`
}
type List struct {
	Ids []string `json:"ids"`
}

//type Number struct {
//	Number int64
//	Value  int64
//}

type KeyboardGrid struct {
	ColorMap      [1021]int64  `json:"-"`
	MapCount      [6][22]int64 `json:"-"`
	MaxKeyPresses int64        `json:"-"`
	Effect        string       `json:"effect"`
	Param         [6][22]int64 `json:"param"`
}

func BasicGrid() *KeyboardGrid {
	e := &KeyboardGrid{
		Effect: Custom,
		Param:  GetKeyboardStruct(),
	}

	setColorMap(e)

	return e
}

func GetKeyboardStruct() [KeyboardMaxRows][KeyboardMaxColumns]int64 {
	var grid KeyboardGrid
	for i, _ := range grid.Param {
		for y, _ := range grid.Param[i] {
			grid.Param[i][y] = 0xFF0000
		}
	}
	return grid.Param
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
