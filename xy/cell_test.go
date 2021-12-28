package xy

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func prettyClose(x, y float64) bool {
	return math.Abs(x-y) < 0.00001
}

func TestABC(t *testing.T) {
	cell := Encode((-112.321+180)/360, (33.123+90)/180)
	x, y := Decode(cell)
	fmt.Printf("%d %f %f '%s'\n", cell, x*360-180, y*180-90, cell.String())
}

func TestEncodeDecode(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	var count int
	for time.Since(start) < time.Millisecond*100 {
		for i := 0; i < 1000; i++ {
			x1 := rand.Float64()
			y1 := rand.Float64()

			switch i % 1000 {
			case 543:
				x1 = -0.0000001
			case 264:
				y1 = -0.0000001
			case 643:
				x1 = 1.0000001
			case 129:
				y1 = 1.0000001
			}
			cell := Encode(x1, y1)
			x2, y2 := Decode(cell)
			if !prettyClose(x1, x2) || !prettyClose(y1, y2) {
				t.Fatalf("\n"+
					"(%0.13f %0.13f)\n"+
					"(%0.13f %0.13f)\n",
					x1, y1, x2, y2)
			}
			str := cell.String()
			cell2 := CellFromString(str)
			if cell != cell2 {
				t.Fatalf("%s != %s", cell2, cell)
			}
			count++
		}
	}
}

func TestPerf(t *testing.T) {
	if os.Getenv("XY_PERF") != "1" {
		fmt.Printf("Use XY_PERF=1 to print performance details\n")
		return
	}
	rand.Seed(time.Now().UnixNano())
	N := 1_000_000
	xy := make([][2]float64, N*2)
	cells := make([]Cell, N)
	for i := 0; i < N; i++ {
		for j := 0; j < 2; j++ {
			xy[i][j] = rand.Float64()
		}
		cells[i] = Encode(xy[i][0], xy[i][1])
	}
	lotsa.Output = os.Stdout
	print("encode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Encode(xy[i][0], xy[i][1])
	})
	print("decode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Decode(cells[i])
	})
}

func TestQuad(t *testing.T) {
	if Encode(0.2, 0.2).Quad(0) != 0 {
		panic("!")
	}
	if Encode(0.2, 0.2).Quad(1) != 0 {
		panic("!")
	}
	if Encode(0.2, 0.6).Quad(1) != 1 {
		panic("!")
	}
	if Encode(0.6, 0.2).Quad(1) != 2 {
		panic("!")
	}
	if Encode(0.6, 0.6).Quad(1) != 3 {
		panic("!")
	}
}
