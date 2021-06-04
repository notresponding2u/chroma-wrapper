package heatmap

import (
	"fmt"
	"github.com/notresponding2u/chroma-wrapper/wrapper"
	"testing"
)

func TestHeatUp(t *testing.T) {
	g := wrapper.BasicGrid()
	k := Key{
		X: 0,
		Y: 0,
	}
	i := 0
	fmt.Println(len(g.ColorMap))
	fmt.Println(int64(len(g.ColorMap)))
	for e, z := range g.ColorMap {
		fmt.Printf("%d  %X \n", e, z)
	}
	for {
		fmt.Printf("%d color: %X\n", i, g.Param[k.X][k.Y])
		i++
		HeatUp(k, g)
		//if g.Param[k.X][k.Y] <= 0x0000FF {
		//	break
		//}
		if i > 1030 {
			break
		}
	}

}
