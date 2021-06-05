package heatmap

import (
	"encoding/json"
	"github.com/notresponding2u/chroma-wrapper/wrapper/effect"
	"io/ioutil"
	"os"
)

const FileAllTimeHeatMap = "./AllTimeHeatMap.json"

type Key struct {
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

func Remap(k Key, grid *effect.KeyboardGrid) {
	grid.MapCount[k.X][k.Y]++
	if grid.MaxKeyPresses < grid.MapCount[k.X][k.Y] {
		grid.MaxKeyPresses = grid.MapCount[k.X][k.Y]
	}
	for x, _ := range grid.Param {
		for y, _ := range grid.Param[x] {
			switch grid.MapCount[x][y] {
			case grid.MaxKeyPresses:
				grid.Param[x][y] = 0x0000FF
			case 0:
				grid.Param[x][y] = 0xFF0000
			default:
				//percentage := grid.MapCount[x][y] * 100 / grid.MaxKeyPresses * int64(len(grid.ColorMap)) / 100
				percentage := float64(grid.MapCount[x][y]) / float64(grid.MaxKeyPresses) * float64(len(grid.ColorMap))
				if int64(percentage) >= int64(len(grid.ColorMap)) {
					grid.Param[x][y] = 0x0000FF
				} else {
					grid.Param[x][y] = grid.ColorMap[int64(percentage)]
				}
			}
		}
	}
}

func Load() {

}

func SaveMap(e *effect.KeyboardGrid) error {
	if _, err := os.Stat(FileAllTimeHeatMap); os.IsNotExist(err) {
		err = save(e, FileAllTimeHeatMap)
		if err != nil {
			return err
		}
	} else {
		err = LoadFile(e, FileAllTimeHeatMap)
		if err != nil {
			return err
		}

		err = save(e, FileAllTimeHeatMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func save(e *effect.KeyboardGrid, file string) error {
	j, err := json.Marshal(e.MapCount)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, j, 0644)
	if err != nil {
		return err
	}
	return nil
}

func mergeHeatmaps(receiver *effect.KeyboardGrid, donor *effect.KeyboardGrid) {
	for x, _ := range receiver.MapCount {
		for y, _ := range receiver.MapCount[x] {
			receiver.MapCount[x][y] += donor.MapCount[x][y]
			if receiver.MapCount[x][y] > receiver.MaxKeyPresses {
				receiver.MaxKeyPresses = receiver.MapCount[x][y]
			}
		}
	}
}

func LoadFile(e *effect.KeyboardGrid, file string) error {
	g := &effect.KeyboardGrid{}

	j, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &g.MapCount)
	if err != nil {
		return err
	}

	mergeHeatmaps(e, g)

	return os.Remove(file)
}

func HeatUp(k Key, grid *effect.KeyboardGrid) {
	if grid.Param[k.X][k.Y] > 0x0000FF {
		grid.Param[k.X][k.Y] = grid.ColorMap[grid.MapCount[k.X][k.Y]]
		grid.MapCount[k.X][k.Y]++

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
	}
}

func NewMap() map[uint16]Key {
	m := make(map[uint16]Key)

	m[27] = Key{
		X: 0,
		Y: 1,
	}
	m[192] = Key{
		X: 1,
		Y: 1,
	}
	m[9] = Key{
		X: 2,
		Y: 1,
	}
	m[20] = Key{
		X: 3,
		Y: 1,
	}
	m[160] = Key{
		X: 4,
		Y: 1,
	}
	m[162] = Key{
		X: 5,
		Y: 1,
	}

	m[49] = Key{
		X: 1,
		Y: 2,
	}
	m[81] = Key{

		X: 2,
		Y: 2,
	}
	m[65] = Key{
		X: 3,
		Y: 2,
	}
	m[226] = Key{
		X: 4,
		Y: 2,
	}
	m[91] = Key{
		X: 5,
		Y: 2,
	}

	m[112] = Key{
		X: 0,
		Y: 3,
	}
	m[50] = Key{
		X: 1,
		Y: 3,
	}
	m[87] = Key{
		X: 2,
		Y: 3,
	}
	m[83] = Key{
		X: 3,
		Y: 3,
	}
	m[90] = Key{
		X: 4,
		Y: 3,
	}
	m[164] = Key{
		X: 5,
		Y: 3,
	}

	m[113] = Key{
		X: 0,
		Y: 4,
	}
	m[51] = Key{
		X: 1,
		Y: 4,
	}
	m[69] = Key{
		X: 2,
		Y: 4,
	}
	m[68] = Key{
		X: 3,
		Y: 4,
	}
	m[88] = Key{
		X: 4,
		Y: 4,
	}

	m[114] = Key{
		X: 0,
		Y: 5,
	}
	m[52] = Key{
		X: 1,
		Y: 5,
	}
	m[82] = Key{
		X: 2,
		Y: 5,
	}
	m[70] = Key{
		X: 3,
		Y: 5,
	}
	m[67] = Key{
		X: 4,
		Y: 5,
	}

	m[115] = Key{
		X: 0,
		Y: 6,
	}
	m[53] = Key{
		X: 1,
		Y: 6,
	}
	m[84] = Key{
		X: 2,
		Y: 6,
	}
	m[71] = Key{
		X: 3,
		Y: 6,
	}
	m[86] = Key{
		X: 4,
		Y: 6,
	}

	m[116] = Key{
		X: 0,
		Y: 7,
	}
	m[54] = Key{
		X: 1,
		Y: 7,
	}
	m[89] = Key{
		X: 2,
		Y: 7,
	}
	m[72] = Key{
		X: 3,
		Y: 7,
	}
	m[66] = Key{
		X: 4,
		Y: 7,
	}
	m[32] = Key{
		X: 5,
		Y: 7,
	}

	m[117] = Key{
		X: 0,
		Y: 8,
	}
	m[55] = Key{
		X: 1,
		Y: 8,
	}
	m[85] = Key{
		X: 2,
		Y: 8,
	}
	m[74] = Key{
		X: 3,
		Y: 8,
	}
	m[78] = Key{
		X: 4,
		Y: 8,
	}

	m[118] = Key{
		X: 0,
		Y: 9,
	}
	m[56] = Key{
		X: 1,
		Y: 9,
	}
	m[73] = Key{
		X: 2,
		Y: 9,
	}
	m[75] = Key{
		X: 3,
		Y: 9,
	}
	m[77] = Key{
		X: 4,
		Y: 9,
	}

	m[119] = Key{
		X: 0,
		Y: 10,
	}
	m[57] = Key{
		X: 1,
		Y: 10,
	}
	m[79] = Key{
		X: 2,
		Y: 10,
	}
	m[76] = Key{
		X: 3,
		Y: 10,
	}
	m[188] = Key{
		X: 4,
		Y: 10,
	}

	m[120] = Key{
		X: 0,
		Y: 11,
	}
	m[48] = Key{
		X: 1,
		Y: 11,
	}
	m[80] = Key{
		X: 2,
		Y: 11,
	}
	m[186] = Key{
		X: 3,
		Y: 11,
	}
	m[190] = Key{
		X: 4,
		Y: 11,
	}
	m[165] = Key{
		X: 5,
		Y: 11,
	}

	m[121] = Key{
		X: 0,
		Y: 12,
	}
	m[189] = Key{
		X: 1,
		Y: 12,
	}
	m[219] = Key{
		X: 2,
		Y: 12,
	}
	m[222] = Key{
		X: 3,
		Y: 12,
	}
	m[191] = Key{
		X: 4,
		Y: 12,
	}

	m[122] = Key{
		X: 0,
		Y: 13,
	}
	m[187] = Key{
		X: 1,
		Y: 13,
	}
	m[221] = Key{
		X: 2,
		Y: 13,
	}
	m[220] = Key{
		X: 3,
		Y: 13,
	}
	m[93] = Key{
		X: 5,
		Y: 13,
	}

	m[123] = Key{
		X: 0,
		Y: 14,
	}
	m[8] = Key{
		X: 1,
		Y: 14,
	}
	m[13] = Key{
		X: 3,
		Y: 14,
	}
	m[161] = Key{
		X: 4,
		Y: 14,
	}
	m[163] = Key{
		X: 5,
		Y: 14,
	}

	m[44] = Key{
		X: 0,
		Y: 15,
	}
	m[45] = Key{
		X: 1,
		Y: 15,
	}
	m[46] = Key{
		X: 2,
		Y: 15,
	}
	m[37] = Key{
		X: 5,
		Y: 15,
	}

	m[145] = Key{
		X: 0,
		Y: 16,
	}
	m[36] = Key{
		X: 1,
		Y: 16,
	}
	m[35] = Key{
		X: 2,
		Y: 16,
	}
	m[38] = Key{
		X: 4,
		Y: 16,
	}
	m[40] = Key{
		X: 5,
		Y: 16,
	}

	m[19] = Key{
		X: 0,
		Y: 17,
	}
	m[33] = Key{
		X: 1,
		Y: 17,
	}
	m[34] = Key{
		X: 2,
		Y: 17,
	}
	m[39] = Key{
		X: 5,
		Y: 17,
	}

	m[144] = Key{
		X: 1,
		Y: 18,
	}
	m[103] = Key{
		X: 2,
		Y: 18,
	}
	m[100] = Key{
		X: 3,
		Y: 18,
	}
	m[97] = Key{
		X: 4,
		Y: 18,
	}

	m[111] = Key{
		X: 1,
		Y: 19,
	}
	m[104] = Key{
		X: 2,
		Y: 19,
	}
	m[101] = Key{
		X: 3,
		Y: 19,
	}
	m[98] = Key{
		X: 4,
		Y: 19,
	}
	m[96] = Key{
		X: 5,
		Y: 19,
	}

	m[106] = Key{
		X: 1,
		Y: 20,
	}
	m[105] = Key{
		X: 2,
		Y: 20,
	}
	m[102] = Key{
		X: 3,
		Y: 20,
	}
	m[99] = Key{
		X: 4,
		Y: 20,
	}
	m[110] = Key{
		X: 5,
		Y: 20,
	}

	m[109] = Key{
		X: 1,
		Y: 21,
	}
	m[107] = Key{
		X: 2,
		Y: 21,
	}
	m[13] = Key{
		X: 4,
		Y: 21,
	}

	return m
}
